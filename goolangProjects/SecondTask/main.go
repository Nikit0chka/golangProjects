package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// PathSize структура похожая на map[path]size
type PathSize struct {
	Path string
	Size int64
}

func main() {
	startDirectory, dirSizeLimit, sortType := getLogInput()
	err := checkLogInput(startDirectory, dirSizeLimit, sortType)
	if err != nil {
		fmt.Print(err)
	}

	dirSizes, err := getMapDirSizes(startDirectory)
	if err != nil {
		log.Fatal(err)
	}
	sortedDirSizes, err := sortDirSizes(dirSizes, sortType)
	if err != nil {
		log.Fatal(err)
	}
	dirSizesLargerLimit := getDirsLargerLimit(sortedDirSizes, dirSizeLimit)
	printDirsToLog(sortedDirSizes)

	err = writeDirSizesToFile("D:\\goolangProjects\\SecondTask\\result.txt", dirSizesLargerLimit)
	if err != nil {
		log.Fatal(err)
	}
}

// getLogInput запрашивает и возвращает ввод с консоли директории, лимит размера директории, тип сортировки
func getLogInput() (string, int64, string) {
	var startDirectory string
	var dirSizeLimit int64
	var sortType string

	flag.StringVar(&startDirectory, "pathToDirectory", "", "Path to directory")
	flag.Int64Var(&dirSizeLimit, "limitOfDirSize", 0, "Limit of directory size in bytes")
	flag.StringVar(&sortType, "typeOfSort", "ASC", "Type of sort ASC/DESK")

	flag.Parse()

	fmt.Printf("%s - path to directory\n", startDirectory)
	fmt.Printf("%d - limit of directory size in bytes\n", dirSizeLimit)
	fmt.Printf("%s - type of sort\n", sortType)

	return startDirectory, dirSizeLimit, strings.ToTitle(sortType)
}

// checkLogInput проверяет ввод с консоли
func checkLogInput(startDirectory string, dirSizeLimit int64, sortType string) error {
	if strings.Contains(startDirectory, ".txt") {
		return fmt.Errorf("You must input directory, not path to file : %s\n", startDirectory)
	}
	if dirSizeLimit < 0 {
		return fmt.Errorf("Directory size limit can not be less than 0 : %d\n", dirSizeLimit)
	}
	if sortType != "ASC" && sortType != "DESK" {
		return fmt.Errorf("Sort type can be ASC or DESK %s\n", sortType)
	}
	return nil
}

// getMapDirSizes проходит по всем вложенным директориям и возвращает массив всех директорий
func getMapDirSizes(startDirectory string) ([]PathSize, error) {
	var dirSizes []PathSize

	err := filepath.Walk(startDirectory, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			return nil
		}
		if path == startDirectory {
			return nil
		}

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
			return err
		}
		dirSizes = append(dirSizes, PathSize{path, size})
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("Error by reading directory : %s \n", err)
	}
	return dirSizes, nil
}

// sortDirSizes сортирует размер директорий в зависимости от типа сортировки ASC/DESK
func sortDirSizes(dirSizes []PathSize, sortType string) ([]PathSize, error) {
	var sortedPathSizes []PathSize

	for _, value := range dirSizes {
		sortedPathSizes = append(sortedPathSizes, PathSize{value.Path, value.Size})
	}

	if sortType == "ASC" {
		sort.Slice(sortedPathSizes, func(i, j int) bool {
			return sortedPathSizes[i].Size < sortedPathSizes[j].Size
		})
	} else if sortType == "DESK" {
		sort.Slice(sortedPathSizes, func(i, j int) bool {
			return sortedPathSizes[i].Size > sortedPathSizes[j].Size
		})
	} else {
		return nil, fmt.Errorf("Error by trying sort dirSizes \n")
	}

	return sortedPathSizes, nil
}

// getDirsLargerLimit выводит размер директорий
func printDirsToLog(dirSizes []PathSize) {
	for _, value := range dirSizes {
		fmt.Printf("%s : %d bytes \n", value.Path, value.Size)
	}
}

// getDirsLargerLimit возврашает массив размеров директорий, которые больше лимита размера директории
func getDirsLargerLimit(dirSizes []PathSize, dirSizeLimit int64) []PathSize {
	var largDirs []PathSize

	for _, value := range dirSizes {
		if value.Size > dirSizeLimit {
			largDirs = append(largDirs, PathSize{value.Path, value.Size})
		}
	}

	return largDirs
}

// writeDirSizesToFile записывает массив размеров директорий в файл
func writeDirSizesToFile(pathToFile string, dirSizes []PathSize) error {
	file, err := os.Create(pathToFile)
	if err != nil {
		return fmt.Errorf("Error by trying to write file by path : %s \n %s \n", pathToFile, err)
	}
	defer file.Close()

	for _, value := range dirSizes {
		_, err := file.WriteString(fmt.Sprintf("%s : %d \n", value.Path, value.Size))
		if err != nil {
			return fmt.Errorf("Error by trying to write file by path : %s \n %s \n", pathToFile, err)
		}
	}
	return nil
}
