package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	pcg "mainModule/pack"
	"net/http"
	"os"
	"strings"
)

type RequestBody struct {
	Error string `json:"error"`
}

func main() {
	config, err := getConfigSettings("C:\\Users\\voron\\OneDrive\\Рабочий стол\\golangProjects\\ThirdTask\\config")
	if err != nil {
		log.Fatal(err)
	}
	http.Handle("/static/", http.StripPrefix("/static/",)
	hostServ(config)
}

// hostServ запускает сервер HTTP на порту указанном в конфиге и слушает его
func hostServ(configSettings map[string]string) {
	domainName := configSettings["domain"]
	port := configSettings["port"]

	http.HandleFunc(domainName, fileWorkHandler)

	fmt.Println("Server is listening...")
	err := http.ListenAndServe(port, fileWorkHandler)

	if err != nil {
		log.Fatal(fmt.Sprintf("Error by trying listen serv : %s", err))
	}
}

// fileWorkHandler обработчик, принимающий из url начальную директорию, лимит размера директории, тип сортировки и работает с директориями
func fileWorkHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Print("sdf")

	reqBody := RequestBody  {Error: "asds"}

	w.Header().Set("Content-Type", "application/json")

	sjsj, _ := json.Marshal(reqBody)
	w.Write(sjsj)
}

func getRequestData(r *http.Request) ([]byte, error) {
	if r.Method != "POST" {
		return nil, fmt.Errorf("allowed only POST methods")
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading request body")
	}

	return body, nil
}

func getDirectories(startDirectory string, sizeLimitMb int64, sortType string) ([]pcg.PathSize, error) {
	//Проверяем данные
	if err := pcg.CheckInput(startDirectory, sizeLimitMb, sortType); err != nil {
		return nil, err
	}

	//Получаем директории - размеры
	dirSizes, err := pcg.GetDirSizes(startDirectory)
	if err != nil {
		return nil, err
	}

	//Сортируем директории - размеры
	dirSizes, err = pcg.SortDirSizes(dirSizes, sortType)
	if err != nil {
		return nil, err
	}

	return dirSizes, nil
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
		parts := strings.Split(line, " ")
		host := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		configMap[host] = value
	}

	return configMap, nil
}
