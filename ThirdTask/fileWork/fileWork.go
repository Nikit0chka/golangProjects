package fileWork

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
)

// FileInfo структура для работы с файлами и директориями
type FileInfo struct {
	Path string
	Size float32
	Name string
	Type string
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

// GetDirSizes считает размер каждой поддиректории
func GetDirSizes(path string) ([]FileInfo, error) {
	// Список объектов в текущей директории
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, fmt.Errorf("error reading dir : %s", err)
	}
	wg := sync.WaitGroup{}

	// Канал результата
	resCh := make(chan FileInfo)

	// Канал для ошибок
	errCh := make(chan error)

	// Обойти все объекты в текущей директории
	for _, file := range files {
		if file.IsDir() {
			// Добавить горутину в WaitGroup
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
				resCh <- FileInfo{Path: dirPath, Size: float32(size) / (1024 * 1024), Name: dirName, Type: dirType}
			}(file.Name())
		}
	}

	// Закрыть канал результатов, когда все горутины завершатся
	go func() {
		wg.Wait()
		close(resCh)
		close(errCh)
	}()

	result := make([]FileInfo, 0)
	for res := range resCh {
		result = append(result, res)
	}

	for err := range errCh {
		return nil, err
	}

	//добавляем к результату файлы
	result, err = addFileSizes(path, result)
	if err != nil {
		return nil, fmt.Errorf("error trying get files from dir: %s", err)
	}

	return result, nil
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

// addFileSizes возвращает файлы и их размер
func addFileSizes(directory string, pathSizes []FileInfo) ([]FileInfo, error) {
	files, err := ioutil.ReadDir(directory)
	if err != nil {
		return nil, fmt.Errorf("error reading file : %s", err)
	}

	for _, file := range files {
		if !file.IsDir() {
			pathSizes = append(pathSizes, FileInfo{fmt.Sprintf("%s/%s", directory, file.Name()), float32(file.Size()), file.Name(), fileType})
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

	return sortedPathSizes, nil
}

// CheckInput проверяет ввод данынх для поиска директорий
func CheckInput(startDirectory string, sortType string) error {
	if _, err := os.Stat(startDirectory); os.IsNotExist(err) {
		return fmt.Errorf("directory by path : %s is not exist", startDirectory)
	}
	if strings.ToUpper(sortType) != ascType && strings.ToUpper(sortType) != deskType {
		return fmt.Errorf("sort type can be ASC or DESK! %s", sortType)
	}

	return nil
}
