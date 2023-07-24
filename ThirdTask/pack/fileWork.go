package pack

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"sort"
	"strings"
)

// PathSize структура похожая на map[path]size
type PathSize struct {
	Path string
	Size int64
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

// GetDirSizes считает размер каждой директории
func GetDirSizes(path string) ([]PathSize, error) {
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
		pathSizes = append(pathSizes, PathSize{fmt.Sprintf("%s\\%s", filePath, file.Name()), file.Size() / 8388608, file.Name(), fileTyp})
	}
	return pathSizes, nil
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
