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

func main() {
	startDirectory, dirSizeLimit, sortType := getLogInput()
	err := checkLogInput(startDirectory, dirSizeLimit, sortType)
	if err != nil {
		log.Fatal(err)
	}

	dirSizes := getMapDirSizes(startDirectory)
	sortedDirSizes := sortDirSizes(dirSizes, sortType)
	dirSizesLargerLimit := getDirsLargerLimit(sortedDirSizes, dirSizeLimit)
	printDirsToLog(sortedDirSizes)

	err = writeDirSizesToFile("D:\\goolangProjects\\SecondTask\\result.txt", dirSizesLargerLimit)
	if err != nil {
		log.Fatal(err)
	}
}

func getLogInput() (string, int64, string) {
	var startDirectory string
	var dirSizeLimit int64
	var sortType = "ASC"

	flag.StringVar(&startDirectory, "pathToDirectory", "", "Path to directory")
	flag.Int64Var(&dirSizeLimit, "limitOfDirSize", 0, "Limit of directory size in bytes")
	flag.StringVar(&sortType, "typeOfSort", "", "Type of sort ASC/DESK")

	flag.Parse()

	fmt.Printf("%s - path to directory\n", startDirectory)
	fmt.Printf("%d - limit of directory size in bytes\n", dirSizeLimit)
	fmt.Printf("%s - type of sort\n", sortType)

	return startDirectory, dirSizeLimit, strings.ToTitle(sortType)
}

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

func getMapDirSizes(startDirectory string) map[string]int64 {
	dirSizes := make(map[string]int64)

	filepath.Walk(startDirectory, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() && path != "." {
			var size int64
			filepath.Walk(path, func(subPath string, subInfo os.FileInfo, subErr error) error {
				if !subInfo.IsDir() {
					size += subInfo.Size()
				}
				return nil
			})
			dirSizes[path] = size
		}
		return nil
	})
	return dirSizes
}

func sortDirSizes(dirSizes map[string]int64, sortType string) map[string]int64 {
	type kvForSort struct {
		Key   string
		Value int64
	}
	var sortedDirSizesSplit []kvForSort
	var sortedDirSizesMap = make(map[string]int64, len(dirSizes))

	for key, value := range dirSizes {
		sortedDirSizesSplit = append(sortedDirSizesSplit, kvForSort{key, value})
	}

	if sortType == "ASC" {
		sort.Slice(sortedDirSizesSplit, func(i, j int) bool {
			return sortedDirSizesSplit[i].Value < sortedDirSizesSplit[j].Value
		})
	} else if sortType == "DESK" {
		sort.Slice(sortedDirSizesSplit, func(i, j int) bool {
			return sortedDirSizesSplit[i].Value > sortedDirSizesSplit[j].Value
		})
	}

	for _, keyValue := range sortedDirSizesSplit {
		sortedDirSizesMap[keyValue.Key] = keyValue.Value
	}
	return sortedDirSizesMap
}

func printDirsToLog(dirSizes map[string]int64) {
	for key, value := range dirSizes {
		fmt.Printf("%s : %d bytes \n", key, value)
	}
}

func getDirsLargerLimit(dirSizes map[string]int64, dirSizeLimit int64) map[string]int64 {
	largDirs := make(map[string]int64)

	for key, value := range dirSizes {
		if value > dirSizeLimit {
			largDirs[key] = value
		}
	}

	return largDirs
}

func writeDirSizesToFile(pathToFile string, dirSizes map[string]int64) error {
	file, err := os.Create(pathToFile)
	if err != nil {
		return fmt.Errorf("Error by trying to write file by path : %s \n %s \n", pathToFile, err)
	}

	for key, value := range dirSizes {
		_, err := file.WriteString(fmt.Sprintf("%s : %d \n", key, value))
		if err != nil {
			return fmt.Errorf("Error by trying to write file by path : %s \n %s \n", pathToFile, err)
		}
	}
	return nil
}
