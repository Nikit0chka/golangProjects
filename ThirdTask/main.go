package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	pcg "mainModule/pack"
)

type ErrorContent struct {
	ErId    string
	Content string
}

func main() {
	config, err := getConfigSettings("/home/vorontsov/Desktop/golangProjects/ThirdTask/config")
	if err != nil {
		log.Fatal(err)
	}

	hostServ(config)
}

// hostServ запускает сервер HTTP на порту указанном в конфиге и слушает его
func hostServ(configSettings map[string]string) {
	domainName := configSettings["domain"]
	port := configSettings["port"]
	fmt.Printf("%s %s", domainName, port)
	http.HandleFunc(fmt.Sprintf("/%s", domainName), fileWorkHandler)

	fmt.Println("Server is listening...")
	err := http.ListenAndServe(fmt.Sprintf(":%s", port), http.FileServer(http.Dir("static")))

	if err != nil {
		log.Fatal(fmt.Sprintf("Error by trying listen serv : %s", err))
	}
}

// fileWorkHandler обработчик, принимающий из url начальную директорию, лимит размера директории, тип сортировки и работает с директориями
func fileWorkHandler(w http.ResponseWriter, r *http.Request) {
	var errors []ErrorContent

	//Получаем данные из url
	startDirectory, sizeLimitMb, sortType, err := getUrlInput(r)
	if err != nil {
		errors = append(errors, ErrorContent{"0", fmt.Sprint(err)})
	}

	//Проверяем данные
	if err := pcg.CheckInput(startDirectory, sizeLimitMb, sortType); err != nil {
		errors = append(errors, ErrorContent{"1", fmt.Sprint(err)})
	}

	//Получаем директории - размеры
	dirSizes, err := pcg.GetDirSizes(startDirectory)
	if err != nil {
		errors = append(errors, ErrorContent{"2", fmt.Sprint(err)})
	}

	//Сортируем директории - размеры
	dirSizes, err = pcg.SortDirSizes(dirSizes, sortType)
	if err != nil {
		errors = append(errors, ErrorContent{"3", fmt.Sprint(err)})
	}

	//Получаем директории больше лимита
	dirSizeLargerLimit := pcg.GetDirsLargerLimit(dirSizes, sizeLimitMb)
	if err := pcg.WriteDirSizesToFile("result.txt", dirSizeLargerLimit); err != nil {
		errors = append(errors, ErrorContent{"4", fmt.Sprint(err)})
	}

	jsonBytes, err := json.Marshal(errors)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(jsonBytes)
}

// getUrlInput получает ввод с url сайта
func getUrlInput(r *http.Request) (string, int64, string, error) {
	vars := r.URL.Query()
	startDirectory := vars.Get("root")
	SizeLimitBytes, err := strconv.ParseInt(vars.Get("limit"), 10, 64)

	if err != nil {
		return "", 0, "", fmt.Errorf("eror by trying convert limit to int64! %s", err)
	}
	sortType := vars.Get("sort")

	return startDirectory, SizeLimitBytes * 1048576.0, sortType, nil
}

// getConfigSettings получает файл с конфигом
func getConfigSettings(pathToFile string) (map[string]string, error) {
	configMap := make(map[string]string)

	file, err := os.Open(pathToFile)
	if err != nil {
		return nil, fmt.Errorf("error by trying open config file by path %s : %s", pathToFile, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, ":")
		host := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		configMap[host] = value
	}

	return configMap, nil
}
