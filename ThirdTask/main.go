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
	errorId  int    `json:"errorId"`
	path     string `json:"path"`
	size     int64  `json:"size"`
	sortType string `json:"sortType"`
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

	http.Handle("/", http.FileServer(http.Dir("static")))
	http.HandleFunc("/"+domainName, fileWorkHandler)
	fmt.Println("Server is listening...")
	err := http.ListenAndServe(port, nil)

	if err != nil {
		log.Fatal(fmt.Sprintf("Error by trying listen serv : %s", err))
	}
}

// fileWorkHandler обработчик, принимающий json работает с директориями
func fileWorkHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	http.ServeFile(w, r, "static/dir-sizes")

	var errors []error
	//декодируем запрос
	decodedRequest, err := decodeRequest(r)
	if err != nil {
		errors = append(errors, err)
	}

	//получаем директории
	directories, err := getDirectories(decodedRequest.path, decodedRequest.sortType)
	if err != nil {
		errors = append(errors, err)
	}

	//создаем ответный json
	jsonRequest, err := createRequest(directories)
	if err != nil {
		errors = append(errors, err)
	}
	fmt.Println(json.Marshal(jsonRequest))
	asdf := RequestBody{1, "valuePath", 3, "sdf"}

	json.NewEncoder(w).Encode(asdf)
}

// createRequest создает json для ответа
func createRequest(directories []pcg.PathSize) ([]byte, error) {
	var bodies []RequestBody
	for _, value := range directories {
		bodies = append(bodies, RequestBody{1, value.Path, int64(value.Size), "sdf"})
	}

	jsonRequest, err := json.Marshal(bodies)
	if err != nil {
		return nil, fmt.Errorf("error by creatin json body %s", err)
	}

	return jsonRequest, nil
}

// decodeRequest декодирует запрос в RequestBody
func decodeRequest(r *http.Request) (RequestBody, error) {
	if r.Method != "POST" {
		return RequestBody{}, fmt.Errorf("allowed only POST methods")
	}
	var decodedRequest RequestBody

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return RequestBody{}, fmt.Errorf("error reading request body %s", err)
	}

	err = json.Unmarshal(body, &decodedRequest)
	if err != nil {
		return RequestBody{}, fmt.Errorf("error reading request body %s", err)
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
