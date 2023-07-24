package pack

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
)

// PathSize структура похожая на map[path]size
type PathSize struct {
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

// GetFiles считывает размеры файлов и директорий
func GetFiles(path string) ([]PathSize, error) {
	var pathSizes []PathSize
	dirs, err := ioutil.ReadDir(path)

	if err != nil {
		return nil, fmt.Errorf("error by geting dir size: %s", path)
	}

	for _, file := range dirs {
		filePath, err := filepath.Abs(path)
		if err != nil {
			return nil, fmt.Errorf("error by reading full path : %s", err)
		}

		var fileTyp string
		if file.IsDir() {
			fileTyp = dirType
		} else {
			fileTyp = fileType
		}
		pathSizes = append(pathSizes, PathSize{fmt.Sprintf("%s/%s", filePath, file.Name()), float32(file.Size()) / 8388608.0, file.Name(), fileTyp})
	}
	return pathSizes, nil
}

// addFileSizes считает размер файлов в директории
func addFileSizes(directory string, pathSizes []PathSize) error {
	files, err := ioutil.ReadDir(".")
	if err != nil {
		return fmt.Errorf("error reading file : %s", err)
	}

	for _, file := range files {
		pathSizes = append(pathSizes, PathSize{fmt.Sprintf("%s/%s", directory, file.Name()), float32(file.Size()), file.Name(), fileType})
	}
	return nil
}

// dirSize считает размер текущей директории
func dirSize(path string, wg *sync.WaitGroup, sizes *[]PathSize) error {
	defer wg.Done()
	var size int64
	err := filepath.Walk(path, func(filePath string, f os.FileInfo, err error) error {
		if !f.IsDir() {
			size += f.Size()
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("error reading file : %s", err)
	}
	name := strings.Split(path, "/")[len(strings.Split(path, "/"))-1]
	*sizes = append(*sizes, PathSize{fmt.Sprintf("%s/%s", path, name), float32(size), name, dirType})
	return nil
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
