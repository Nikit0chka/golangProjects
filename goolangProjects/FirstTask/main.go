package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

func main() {
	startProg()
}

func getLogInput() (string, string) {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()

	isInputValid(scanner.Text())

	return strings.Split(scanner.Text(), " ")[0], strings.Split(scanner.Text(), " ")[1]
}

func startProg() {
	fmt.Println("Input path to file and path to directory : ")
	pathToFile, pathToDirectory := getLogInput()
	timer := time.Now()

	WorkWithFiles(pathToFile, pathToDirectory)

	elapsed := time.Since(timer)
	fmt.Printf("Time of going : %s", elapsed)
}

func isInputValid(inputString string) bool {
	if len(strings.Split(inputString, " ")) != 2 {
		log.Fatal("Check input!")
	}

	return true
}

func WorkWithFiles(pathToFile string, pathToDirectory string) {
	urls := readFile(pathToFile)

	for _, url := range urls {
		writeFile(pathToDirectory, string(getHttpRequestInfo(url)))
	}
}

func readFile(pathToFile string) []string {
	if isFileExist(pathToFile) {
		content, err := ioutil.ReadFile(pathToFile)
		if err != nil {
			fmt.Printf("Eror by reading file by path : %s", pathToFile)
		}

		lines := strings.Split(string(content), "\n")
		return lines
	}
	return nil
}

func isFileExist(pathToFile string) bool {
	_, err := os.Stat(pathToFile)
	if os.IsNotExist(err) {
		log.Fatal("File by path" + pathToFile + " is not exists!")
	}

	return true
}

func writeFile(pathToFile string, content string) {
	file, err := os.OpenFile(pathToFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)

	for _, value := range []byte(content) {
		_, err = writer.WriteString(fmt.Sprint(int(value)))
		if err != nil {
			fmt.Println(err)
		}
	}
	// _, err = writer.WriteString(content)

	writer.Flush()
}

func getHttpRequestInfo(siteUrl string) string {
	resp, err := http.Get(strings.TrimRight(siteUrl, "\r"))
	if err != nil {
		fmt.Printf("Eror by try do http request to : %s %s", siteUrl, err)
	}

	bytesFromBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Eror by try do http request to : %s %s", siteUrl, err)
	}

	defer resp.Body.Close()

	return string(bytesFromBody)
}
