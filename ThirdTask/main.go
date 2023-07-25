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

	flwrk "mainModule/fileWork"
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

	hostServ(config["port"], config["domain"], config["ip"])
}

// hostServ запускает сервер по ip и порту указанных в конфиге
func hostServ(port, domain, ip string) {
	http.HandleFunc(domain, fileWorkHandler)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	fmt.Println("Server is listening...")

	err := http.ListenAndServe(ip+port, nil)
	if err != nil {
		log.Fatal(fmt.Sprintf("Error by trying listen serv : %s", err))
	}
}

// fileWorkHandler обработчик, принимающий и возвращающий json, работает с директориями
func fileWorkHandler(w http.ResponseWriter, r *http.Request) {
	//декодируем запрос
	requestJson, err := decodeRequest(r)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error by decoding json %s", err), 1)
	}

	//получаем директории
	directories, err := getDirectories(requestJson.Path, requestJson.SortType)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error by reading directories %s", err), 2)
	}

	//создаем ответный json
	responseJson, err := createRequest(directories)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error by creating json %s", err), 3)
	}

	//возвращаем json
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(responseJson)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error by responding json %s", err), 4)
	}
}

// createRequest создает json для ответа
func createRequest(directories []flwrk.FileInfo) ([]ResponseJson, error) {
	var bodies []ResponseJson

	if len(directories) == 0 {
		bodies = append(bodies, ResponseJson{"", 0, "", 0, ""})
		return bodies, nil
	}

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
func getDirectories(startDirectory string, sortType string) ([]flwrk.FileInfo, error) {
	//Проверяем данные
	if err := flwrk.CheckInput(startDirectory, sortType); err != nil {
		return nil, err
	}

	//Получаем директории - размеры
	dirSizes, err := flwrk.GetDirSizes(startDirectory)
	if err != nil {
		return nil, err
	}

	//Сортируем директории - размеры
	dirSizes, err = flwrk.SortDirSizes(dirSizes, sortType)
	if err != nil {
		return nil, err
	}
	return dirSizes, nil
}

// getConfigSettings разбивает файл с конфигом
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
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		configMap[key] = value
	}
	return configMap, nil
}
