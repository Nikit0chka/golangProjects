package pack

import (
	"fmt"
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

// GetDirSizes считает размер каждой поддиректории
func GetDirSizes(path string) ([]PathSize, error) {
	var pathSizes []PathSize
	var wg sync.WaitGroup

	_, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	wg.Add(1)
	go dirSize(path, &wg, &pathSizes)
	filepath.Walk(path, func(filePath string, f os.FileInfo, err error) error {
		if f.IsDir() {
			wg.Add(1)
			go dirSize(filePath, &wg, &pathSizes)
		}
		return nil
	})
	wg.Wait()
	return pathSizes, nil
}

// dirSize считает размер текущей директории
func dirSize(path string, wg *sync.WaitGroup, sizes *[]PathSize) {
	defer wg.Done()
	var size int64
	filepath.Walk(path, func(filePath string, f os.FileInfo, err error) error {
		if !f.IsDir() {
			size += f.Size()
		}
		return nil
	})
	*sizes = append(*sizes, PathSize{Path: path, Size: size})
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

// GetDirsLargerLimit возвращает массив PathSize размеров директорий, которые больше лимита размера директории
func GetDirsLargerLimit(dirSizes []PathSize, dirSizeLimit int64) []PathSize {
	var largeDirs []PathSize

	for _, value := range dirSizes {
		if value.Size > dirSizeLimit {
			largeDirs = append(largeDirs, PathSize{value.Path, value.Size})
		}
	}

	return largeDirs
}

// WriteDirSizesToFile записывает массив PathSize директорий и их размеров в файл
func WriteDirSizesToFile(fileName string, dirSizes []PathSize) error {
	file, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf("error by trying to write file by path : %s \n%s", fileName, err)
	}
	defer file.Close()

	for _, value := range dirSizes {
		_, err := file.WriteString(fmt.Sprint(value))
		if err != nil {
			return fmt.Errorf("error by trying to write file by path : %s \n%s", fileName, err)
		}
	}
	return nil
}
