package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
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

// Println выводит в консоль отформатированные значения PathSize
// Так и не понял, о каком переопределении функции вы говорили
// поэтому решил, что просто сделаю метод для структуры выводящий ее данные
func (s PathSize) Println() {
	fmt.Printf("%s : %f mb\n", s.Path, float32(s.Size)/1048576.0)
}

// Типы доступных сортировок
const (
	ascType  = "ASC"
	deskType = "DESK"
)

func main() {
	timer := time.Now()
	//Получаем начальную директорию, лимит размера директории в битай, тип сортировки
	startDirectory, dirSizeLimitbits, sortType := getLogInput()
	//Проверяем результат ввода
	err := checkLogInput(startDirectory, dirSizeLimitbits, sortType)

	if err != nil {
		log.Fatal(err)
	}

	//Получаем []PathSize с директориями и их размерами
	dirSizes, err := getDirSizes(startDirectory)
	if err != nil {
		log.Fatal(err)
	}
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
	fmt.Printf("Time of going : %s\n", time.Since(timer))
	createServ(sortedDirSizes)
}

func createServ(dirSizes []PathSize) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		for _, value := range dirSizes {
			fmt.Fprintf(w, "%s : %f mb\n", value.Path, float32(value.Size)/1048576.0)
		}
	})

	http.ListenAndServe(":8080", nil)
}

// getLogInput запрашивает и возвращает ввод с консоли директории, лимит размера директории, тип сортировки
func getLogInput() (string, int64, string) {
	var startDirectory string
	var dirSizeLimitbits int64
	var sortType string

	flag.StringVar(&startDirectory, "pathToDirectory", "", "Path to directory")
	flag.Int64Var(&dirSizeLimitbits, "limitOfDirSize", 0, "Limit of directory size in bytes")
	flag.StringVar(&sortType, "typeOfSort", ascType, "Type of sort ASC/DESK")

	flag.Parse()

	fmt.Printf("%s - path to directory\n", startDirectory)
	fmt.Printf("%d - limit of directory size in mb\n", dirSizeLimitbits)
	fmt.Printf("%s - type of sort\n", sortType)

	return startDirectory, dirSizeLimitbits * 1048576.0, strings.ToTitle(sortType)
}

// checkLogInput проверяет ввод с консоли
func checkLogInput(startDirectory string, dirSizeLimit int64, sortType string) error {
	if _, err := os.Stat(startDirectory); os.IsNotExist(err) {
		return fmt.Errorf("Directory by path : %s is not exist!\n", startDirectory)
	}
	if dirSizeLimit < 0 {
		return fmt.Errorf("Directory size limit can not be less than 0 : %d\n", dirSizeLimit)
	}
	if sortType != ascType && sortType != deskType {
		return fmt.Errorf("Sort type can be ASC or DESK! %s\n", sortType)
	}
	return nil
}

// getDirSizes проходит по всем вложенным директориям и возвращает массив PathSize всех директорий и их размер
func getDirSizes(startDirectory string) ([]PathSize, error) {
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
				log.Printf("Error while walking directory %q: %v\n", path, err)
				return
			}

			dirSizesMutex.Lock()
			dirSizes = append(dirSizes, PathSize{path, size})
			dirSizesMutex.Unlock()
		}(path)

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("Error by reading directory : %s \n", err)
	}
	waitGroup.Wait()

	return dirSizes, nil
}

// sortDirSizes сортирует размер директорий в зависимости от типа сортировки ASC/DESK
func sortDirSizes(dirSizes []PathSize, sortType string) ([]PathSize, error) {
	sortedPathSizes := make([]PathSize, len(dirSizes))
	copy(sortedPathSizes, dirSizes)

	switch sortType {
	case ascType:
		sort.Slice(sortedPathSizes, func(i, j int) bool {
			return sortedPathSizes[i].Size < sortedPathSizes[j].Size
		})
	case deskType:
		sort.Slice(sortedPathSizes, func(i, j int) bool {
			return sortedPathSizes[i].Size > sortedPathSizes[j].Size
		})
	default:
		return nil, fmt.Errorf("Error by trying sort dirSizes \n")
	}

	return sortedPathSizes, nil
}

// getDirsLargerLimit выводит размер директорий
func printDirsToLog(dirSizes []PathSize) {
	for _, value := range dirSizes {
		value.Println()
	}
}

// getDirsLargerLimit возврашает массив PathSize размеров директорий, которые больше лимита размера директории
func getDirsLargerLimit(dirSizes []PathSize, dirSizeLimit int64) []PathSize {
	var largDirs []PathSize

	for _, value := range dirSizes {
		if value.Size > dirSizeLimit {
			largDirs = append(largDirs, PathSize{value.Path, value.Size})
		}
	}

	return largDirs
}

// writeDirSizesToFile записывает массив PathSize директорий и их размеров в файл
func writeDirSizesToFile(fileName string, dirSizes []PathSize) error {
	file, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf("Error by trying to write file by path : %s \n %s \n", fileName, err)
	}
	defer file.Close()

	for _, value := range dirSizes {
		_, err := file.WriteString(fmt.Sprintf("%s : %f mb\n", value.Path, float32(value.Size)/1048576.0))
		if err != nil {
			return fmt.Errorf("Error by trying to write file by path : %s \n %s \n", fileName, err)
		}
	}
	return nil
}
