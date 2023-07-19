package packagesSecTask

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
)

// PathSize структура похожая на map[path]size
type PathSize struct {
	Path string
	Size int64
}

// Типы доступных сортировок
const (
	ascType  = "ASC"
	deskType = "DESK"
)

// Реализация метода String использующий интерфейс stringer, для форматированного вывода
func (p PathSize) String() string {
	return fmt.Sprintf("%s : %f mb\n", p.Path, float32(p.Size)/1048576.0)
}

// GetLogInput запрашивает и возвращает ввод с консоли директории, лимит размера директории, тип сортировки
func GetLogInput() (string, int64, string) {
	var startDirectory string
	var dirSizeLimitBits int64
	var sortType string

	flag.StringVar(&startDirectory, "pathToDirectory", "", "Path to directory")
	flag.Int64Var(&dirSizeLimitBits, "limitOfDirSize", 0, "Limit of directory size in bytes")
	flag.StringVar(&sortType, "typeOfSort", ascType, "Type of sort ASC/DESK")

	flag.Parse()

	fmt.Printf("%s - path to directory\n", startDirectory)
	fmt.Printf("%d - limit of directory size in mb\n", dirSizeLimitBits)
	fmt.Printf("%s - type of sort\n", sortType)

	return startDirectory, dirSizeLimitBits * 1048576.0, strings.ToTitle(sortType)
}

// CheckInput проверяет ввод с консоли
func CheckInput(startDirectory string, dirSizeLimit int64, sortType string) error {
	if _, err := os.Stat(startDirectory); os.IsNotExist(err) {
		return fmt.Errorf("directory by path : %s is not exist", startDirectory)
	}
	if dirSizeLimit < 0 {
		return fmt.Errorf("directory size limit can not be less than 0 : %d", dirSizeLimit)
	}
	if strings.ToUpper(sortType) != ascType && strings.ToUpper(sortType) != deskType {
		return fmt.Errorf("sort type can be ASC or DESK! %s", sortType)
	}
	return nil
}

// GetDirSizes проходит по всем вложенным директориям и возвращает массив PathSize всех директорий и их размер
func GetDirSizes(startDirectory string) ([]PathSize, error) {
	var dirSizes []PathSize
	var dirSizesMutex sync.Mutex

	waitGroup := sync.WaitGroup{}

	err := filepath.Walk(startDirectory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() || path == startDirectory {
			return nil
		}

		waitGroup.Add(1)

		go func(path string) {
			defer waitGroup.Done()

			var size int64
			err = filepath.Walk(path, func(subPath string, subInfo os.FileInfo, subErr error) error {
				if subErr != nil {
					return subErr
				}
				if subInfo == nil {
					return nil
				}
				if !subInfo.IsDir() {
					size += subInfo.Size()
				}
				return nil
			})
			if err != nil {
				log.Printf("error while walking directory %q: %v\n", path, err)
				return
			}

			dirSizesMutex.Lock()
			dirSizes = append(dirSizes, PathSize{path, size})
			dirSizesMutex.Unlock()
		}(path)

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("error by reading directory : %s", err)
	}
	waitGroup.Wait()

	return dirSizes, nil
}

// SortDirSizes сортирует размер директорий в зависимости от типа сортировки ASC/DESK
func SortDirSizes(dirSizes []PathSize, sortType string) ([]PathSize, error) {
	sortedPathSizes := make([]PathSize, len(dirSizes))
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

// getDirsLargerLimit выводит размер директорий
func PrintDirsToLog(dirSizes []PathSize) {
	for _, value := range dirSizes {
		fmt.Print(value)
	}
}

// GetDirsLargerLimit возврашает массив PathSize размеров директорий, которые больше лимита размера директории
func GetDirsLargerLimit(dirSizes []PathSize, dirSizeLimit int64) []PathSize {
	var largDirs []PathSize

	for _, value := range dirSizes {
		if value.Size > dirSizeLimit {
			largDirs = append(largDirs, PathSize{value.Path, value.Size})
		}
	}

	return largDirs
}

// WriteDirSizesToFile записывает массив PathSize директорий и их размеров в файл
func WriteDirSizesToFile(fileName string, dirSizes []PathSize) error {
	file, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf("error by trying to write file by path : %s \n%e", fileName, err)
	}
	defer file.Close()

	for _, value := range dirSizes {
		_, err := file.WriteString(fmt.Sprint(value))
		if err != nil {
			return fmt.Errorf("error by trying to write file by path : %s \n%e", fileName, err)
		}
	}
	return nil
}
