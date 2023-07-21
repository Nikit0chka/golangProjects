package pack

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

// CheckInput проверяет ввод с консоли
func CheckInput(startDirectory string, sortType string) error {
	if _, err := os.Stat(startDirectory); os.IsNotExist(err) {
		return fmt.Errorf("directory by path : %s is not exist", startDirectory)
	}
	if strings.ToUpper(sortType) != ascType && strings.ToUpper(sortType) != deskType {
		return fmt.Errorf("sort type can be ASC or DESK! %s", sortType)
	}
	return nil
}

// PrintDirsToLog выводит размер директорий
func PrintDirsToLog(dirSizes []PathSize) {
	for _, value := range dirSizes {
		fmt.Print(value)
	}
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
