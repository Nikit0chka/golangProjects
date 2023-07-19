package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"
)

func main() {
	pathToFile, pathToDirectory := getLogInput()
	if !isInputValid(pathToFile, pathToDirectory) {
		log.Fatal("Check input!")
	}

	timer := time.Now()
	workWithFiles(pathToFile, pathToDirectory)

	fmt.Printf("Time of going : %s\n", time.Since(timer))
}

// workWithFiles вызывает методы чтения url из файла, отправки http запросаЯ и записи результата в файл
func workWithFiles(pathToFile, pathToDirectory string) {
	urls, err := getFileStrContent(pathToFile)
	if err != nil {
		log.Fatal(err)
	}

	var wg sync.WaitGroup
	for _, urlName := range urls {
		wg.Add(1)

		go func(url string) {
			defer wg.Done()

			domain, err := getDomainName(url)
			if err != nil {
				fmt.Print(err)
				return
			}

			httpRequestinfo, err := getHttpRequestInfo(url)
			if err != nil {
				fmt.Print(err)
				return
			}

			err = writeToFile(pathToDirectory, domain, httpRequestinfo)
			if err != nil {
				fmt.Print(err)
				return
			}
		}(urlName)
	}
	wg.Wait()
}

// getDomainName возвращает имя домейна url из строки
func getDomainName(urlName string) (string, error) {
	parsedUrl, err := url.Parse(urlName)
	if err != nil {
		return "", fmt.Errorf("Error parsing URL: %s", err)
	}
	return parsedUrl.Hostname(), nil
}

// getLogInput обробатывает ввод в консоль
func getLogInput() (string, string) {
	var pathToFile string
	var pathToDirectory string

	flag.StringVar(&pathToFile, "pathToFile", "", "Path to file")
	flag.StringVar(&pathToDirectory, "pathToDirectory", "", "Path to directory")

	flag.Parse()

	fmt.Printf("%s - path to file \n", pathToFile)
	fmt.Printf("%s - path to directory\n", pathToDirectory)
	return pathToFile, pathToDirectory
}

// isInputValid проверяет введеные данные
func isInputValid(pathToFile, pathToDirectory string) bool {
	return strings.Contains(pathToFile, ".txt") && !strings.Contains(pathToDirectory, ".txt")
}

// getFileStrContent читатет файл по пути и возвращает массив строк
func getFileStrContent(pathToFile string) ([]string, error) {
	content, err := ioutil.ReadFile(pathToFile)

	if err != nil {
		return nil, fmt.Errorf("Error of reading file by path : %s \n", err)
	}

	lines := strings.Split(string(content), "\n")
	return lines, nil
}

// writeToFile записывает массив байтов в файл с именем по пути
func writeToFile(pathToFile, fileName string, content []byte) error {
	file, err := os.OpenFile(fmt.Sprintf("%s/%s", pathToFile, fileName), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("Erorr by trying open file by path : %s\n %s \n", pathToFile, err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	_, err = writer.Write(content)

	if err != nil {
		return fmt.Errorf("Error by opening file by path : %s\n %s \n", pathToFile, err)
	}
	err = writer.Flush()
	if err != nil {
		return fmt.Errorf("Error by closing file by path : %s\n %s \n", pathToFile, err)
	}
	return nil
}

// getHttpRequestInfo возвращает массив байтов из http запроса к сайту по url
func getHttpRequestInfo(siteUrl string) ([]byte, error) {
	resp, err := http.Get(strings.TrimRight(siteUrl, "\r"))
	if err != nil {
		return nil, fmt.Errorf("Eror by trying do http request to : %s\n %s \n", siteUrl, err)
	}
	defer resp.Body.Close()

	bytesFromBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Eror by trying do http request to : %s\n %s \n", siteUrl, err)
	}

	return bytesFromBody, nil
}
