package fileWork

import (
	"context"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
)

// FileInfo структура для работы с файлами и директориями
type FileInfo struct {
	Name      string
	FileOrder int
	Path      string
	Size      float32
	Type      string
}

// Типы файлов
const (
	dirType  = "DIR"
	fileType = "FILE"
)

// Типы доступных сортировок
const (
	ascType  = "ASC"
	deskType = "DESK"
)

func GetFileInfos(path string) ([]FileInfo, error) {
	// Создаем контекст с отменой
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Список объектов в текущей директории
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, fmt.Errorf("error reading dir : %s", err)
	}

	// Канал результата
	resCh := make(chan FileInfo)

	// Канал для ошибок
	errCh := make(chan error)

	// Создаем WaitGroup
	var wg sync.WaitGroup

	// Обойти все объекты в текущей директории
	for _, file := range files {
		if file.IsDir() {
			// Увеличиваем счетчик WaitGroup
			wg.Add(1)

			// Начать выполнение горутины
			go func(dirName string) {
				defer wg.Done()

				// Подсчитать размер директории
				dirPath, err := filepath.Abs(filepath.Join(path, dirName))
				if err != nil {
					errCh <- fmt.Errorf("error getting absolute path: %s", err)
					return
				}
				size, err := dirSize(dirPath)
				if err != nil {
					errCh <- fmt.Errorf("error getting directory size: %s", err)
					return
				}

				// Отправить результат в канал
				resCh <- FileInfo{dirName, 0, dirPath, roundTo3Digit(size), dirType}
			}(file.Name())
		}
	}

	// Закрыть каналы после завершения работы всех горутин
	go func() {
		wg.Wait()
		close(resCh)
		close(errCh)
	}()

	// Считать локальные файлы
	pathSizes, err := addFiles(path, nil)
	if err != nil {
		return nil, fmt.Errorf("error adding local files: %s", err)
	}

	result := make([]FileInfo, 0)
	for {
		select {
		case res, ok := <-resCh:
			if !ok {
				// Канал закрыт, все горутины завершились
				return append(result, pathSizes...), nil
			}
			result = append(result, res)
		case err := <-errCh:
			// Выводим сообщение об ошибке и продолжаем работу
			fmt.Printf("error: %s\n", err)
		case <-ctx.Done():
			// Контекст отменен, завершаем работу
			return nil, ctx.Err()
		}
	}
}

// dirSize считает размер одной директории
func dirSize(path string) (int64, error) {
	var size int64

	err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			size += info.Size()
		}
		return nil
	})

	if err != nil {
		return 0, err
	}

	return size, nil
}

// roundTo3Digit округляет число до 3 знаков после запятой
func roundTo3Digit(value int64) float32 {
	return float32(math.Round(float64(value)/float64(1024*1024)*1000) / 1000)
}

// addFiles добавляет в слайс файлы
func addFiles(directory string, pathSizes []FileInfo) ([]FileInfo, error) {
	files, err := ioutil.ReadDir(directory)
	if err != nil {
		return nil, fmt.Errorf("error reading file : %s", err)
	}

	for _, file := range files {
		if !file.IsDir() {
			pathSizes = append(pathSizes, FileInfo{file.Name(), 0, fmt.Sprintf("%s\\%s", directory, file.Name()), roundTo3Digit(file.Size()), fileType})
		}
	}
	return pathSizes, nil
}

// SortDirSizes сортирует размер директорий в зависимости от типа сортировки ASC/DESK
func SortDirSizes(dirSizes []FileInfo, sortType string) ([]FileInfo, error) {
	sortedPathSizes := make([]FileInfo, len(dirSizes))
	copy(sortedPathSizes, dirSizes)

	switch strings.ToUpper(sortType) {
	case ascType:
		sort.Slice(sortedPathSizes, func(i, j int) bool {
			return sortedPathSizes[i].Size < sortedPathSizes[j].Size
		})
	case deskType:
		sort.Slice(sortedPathSizes, func(i, j int) bool {
			return sortedPathSizes[i].Size > sortedPathSizes[j].Size
		})
	default:
		return nil, fmt.Errorf("error by trying sort dirSizes")
	}

	for i := range sortedPathSizes {
		sortedPathSizes[i].FileOrder = i
	}

	return sortedPathSizes, nil
}

// CheckInput проверяет ввод данных для работы с файлами
func CheckInput(startDirectory string, sortType string) error {
	if _, err := os.Stat(startDirectory); os.IsNotExist(err) {
		return fmt.Errorf("directory by path : %s is not exist", startDirectory)
	}
	if strings.ToUpper(sortType) != ascType && strings.ToUpper(sortType) != deskType {
		return fmt.Errorf("sort type can be ASC or DESK! %s", sortType)
	}

	return nil
}
