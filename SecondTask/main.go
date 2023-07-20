package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"
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

func main() {
	//Запуск таймера
	timer := time.Now()
	//Получаем начальную директорию, лимит размера директории в битах, тип сортировки
	startDirectory, dirSizeLimitbits, sortType := getLogInput()
	//Проверяем результат ввода
	if err := checkLogInput(startDirectory, dirSizeLimitbits, sortType); err != nil {
		log.Fatal(err)
	}

	//Получаем []PathSize с директориями и их размерами
	dirSizes := getDirSizes(startDirectory)

	//Получаем отсортированный []PathSize с директориями и их размерами
	sortedDirSizes, err := sortDirSizes(dirSizes, sortType)
	if err != nil {
		log.Fatal(err)
	}
	//Получаем []PathSize с превышающими лимит битов директориями
	dirSizesLargerLimit := getDirsLargerLimit(sortedDirSizes, dirSizeLimitbits)
	printDirsToLog(sortedDirSizes)

	//Записываем директории в файл
	err = writeDirSizesToFile("result.txt", dirSizesLargerLimit)
	if err != nil {
		log.Fatal(err)
	}
	//Вывод времени работы
	fmt.Printf("Time of going : %s\n", time.Since(timer))

}

// getLogInput запрашивает и возвращает ввод с консоли директории, лимит размера директории, тип сортировки
func getLogInput() (string, int64, string) {
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

// checkLogInput проверяет ввод с консоли
func checkLogInput(startDirectory string, dirSizeLimit int64, sortType string) error {
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

// getDirSizes считает размер каждой поддиректории
func getDirSizes(path string) []PathSize {
	var pathSizes []PathSize
	var wg sync.WaitGroup

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
	return pathSizes
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

// sortDirSizes сортирует размер директорий в зависимости от типа сортировки ASC/DESK
func sortDirSizes(dirSizes []PathSize, sortType string) ([]PathSize, error) {
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
func printDirsToLog(dirSizes []PathSize) {
	for _, value := range dirSizes {
		fmt.Print(value)
	}
}

// getDirsLargerLimit возвращает массив PathSize размеров директорий, которые больше лимита размера директории
func getDirsLargerLimit(dirSizes []PathSize, dirSizeLimit int64) []PathSize {
	var largeDirs []PathSize

	for _, value := range dirSizes {
		if value.Size > dirSizeLimit {
			largeDirs = append(largeDirs, PathSize{value.Path, value.Size})
		}
	}

	return largeDirs
}

// writeDirSizesToFile записывает массив PathSize директорий и их размеров в файл
func writeDirSizesToFile(fileName string, dirSizes []PathSize) error {
	file, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf("error by trying to write file by path : %s \n %s", fileName, err)
	}
	defer file.Close()

	for _, value := range dirSizes {
		_, err := file.WriteString(fmt.Sprint(value))
		if err != nil {
			return fmt.Errorf("error by trying to write file by path : %s \n %s", fileName, err)
		}
	}
	return nil
}
