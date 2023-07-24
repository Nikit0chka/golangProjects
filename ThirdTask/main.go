package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	pcg "mainModule/pack"
)

// ResponseJson структура json ответа
type ResponseJson struct {
	Name   string  `json:"name"`
	DireId int     `json:"dirId"`
	Path   string  `json:"path"`
	Size   float32 `json:"size"`
	Type   string  `json:"type"`
}

// RequestJson структура json запроса
type RequestJson struct {
	Path     string `json:"path"`
	SortType string `json:"sortType"`
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
	domain := configSettings["domain"]
	port := configSettings["port"]

	http.HandleFunc(domain, fileWorkHandler)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	fmt.Println("Server is listening...")

	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatal(fmt.Sprintf("Error by trying listen serv : %s", err))
	}
}

// fileWorkHandler обработчик, принимающий и возвращающий json,работает с директориями
func fileWorkHandler(w http.ResponseWriter, r *http.Request) {
	//декодируем запрос
	requestJson, err := decodeRequest(r)
	if err != nil {
		//http.Error(w, fmt.Sprintf("Error by decoding json %s", err), 1)
		fmt.Println(err)
		return
	}

	//получаем директории
	directories, err := getDirectories(requestJson.Path, requestJson.SortType)
	if err != nil {
		fmt.Println(err)
		http.Error(w, fmt.Sprintf("Error by reading directories %s", err), 2)
		return
	}

	//создаем ответный json
	responseJson, err := createRequest(directories)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error by creating json %s", err), 3)
		return
	}

	//возвращаем json
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(responseJson)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error by responding json %s", err), 4)
		return
	}
}

// createRequest создает json для ответа
func createRequest(directories []pcg.PathSize) ([]ResponseJson, error) {
	var bodies []ResponseJson
	for i, value := range directories {
		bodies = append(bodies, ResponseJson{value.Name, i, value.Path, value.Size, value.Type})
	}

	return bodies, nil
}

// decodeRequest декодирует запрос в RequestJson
func decodeRequest(r *http.Request) (RequestJson, error) {
	if r.Method != "POST" {
		return RequestJson{}, fmt.Errorf("only POST request allowed")
	}
	var decodedRequest RequestJson

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return RequestJson{}, fmt.Errorf("error reading request body %s", err)
	}

	err = json.Unmarshal(body, &decodedRequest)
	if err != nil {
		return RequestJson{}, fmt.Errorf("error reading request body %s", err)
	}

	return decodedRequest, nil
}

// getDirectories возвращает отсортированные локальные директории
func getDirectories(startDirectory string, sortType string) ([]pcg.PathSize, error) {
	//Проверяем данные
	if err := pcg.CheckInput(startDirectory, sortType); err != nil {
		return nil, err
	}

	//Получаем директории - размеры
	dirSizes, err := pcg.GetFiles(startDirectory)
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
