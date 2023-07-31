package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	flWrk "mainModule/fileWork"
	"net/http"
	"time"
)

// RequestJson структура json запроса
type RequestJson struct {
	Path     string `json:"path"`
	SortType string `json:"sortType"`
}

type PostJson struct {
	Path        string  `json:"path"`
	Size        float32 `json:"size"`
	Seconds     int     `json:"seconds"`
	DateOfEntry string  `json:"dateOfEntry"`
}

type GetJson struct {
	Id          int       `json:"id"`
	Path        string    `json:"path"`
	Size        float32   `json:"size"`
	Seconds     int       `json:"seconds"`
	DateOfEntry time.Time `json:"dateOfEntry"`
}

func main() {
	hostServ()
}

// hostServ запускает сервер
func hostServ() {
	http.HandleFunc("/", fileWorkHandler)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	fmt.Println("Server is listening...")

	err := http.ListenAndServe(":8081", nil)
	if err != nil {
		log.Fatal(fmt.Sprintf("Error by trying listen serv : %s", err))
	}
}

// fileWorkHandler обработчик, принимающий и возвращающий json, работает с директориями
func fileWorkHandler(w http.ResponseWriter, r *http.Request) {
	//декодируем запрос
	requestJson, err := decodeRequest(r)
	if err != nil {
		http.ServeFile(w, r, "static/dir-sizes")
		return
	}

	//получаем директории
	directories, err := getDirectories(requestJson.Path, requestJson.SortType)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error by reading directories %s", err), 2)
	}

	err = sendPostToApache(PostJson{Path: "requestJson.Path", Size: 12.32, Seconds: 1, DateOfEntry: "2023-07-31"})
	if err != nil {
		fmt.Println(err)
	}

	//возвращаем json
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(directories)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error by responding json %s", err), 4)
	}
}

func sendGetToApache() ([]GetJson, error) {
	resp, err := http.Get("http://192.168.1.31:8083/")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result []GetJson
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}

	fmt.Println(resp)
	fmt.Println(result)
	return result, nil
}

// send
func sendPostToApache(postJson PostJson) error {
	jsonBody, err := json.Marshal(postJson)
	fmt.Println("slo")
	if err != nil {
		return err
	}

	// Создаем HTTP-запрос
	resp, err := http.Post("http://192.168.1.31:8083", "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	// Обрабатываем ответ
	fmt.Println("Response Status:", resp.Status)
	return nil
}

// decodeRequest декодирует запрос в RequestJson
func decodeRequest(r *http.Request) (RequestJson, error) {
	if r.Method != "POST" {
		return RequestJson{}, fmt.Errorf("only POST method allowed")
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
func getDirectories(startDirectory string, sortType string) ([]flWrk.FileInfo, error) {
	//Проверяем данные
	if err := flWrk.CheckInput(startDirectory, sortType); err != nil {
		return nil, err
	}

	//Получаем директории - размеры
	dirSizes, err := flWrk.GetFileInfos(startDirectory)
	if err != nil {
		return dirSizes, err
	}

	//Сортируем директории - размеры
	dirSizes, err = flWrk.SortDirSizes(dirSizes, sortType)
	if err != nil {
		return nil, err
	}
	return dirSizes, nil
}
