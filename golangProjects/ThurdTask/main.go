package main

import (
	"fmt"
	"net/http"
	"strconv"

	pcg "golangProjects/pack"
)

func main() {
	hostServ()
}

// hostServ запускает сервер HTTP на порту 8080 и назначает обработчиков для путей
func hostServ() {
	http.HandleFunc("/dir-size", fileWorkHandler)

	fmt.Println("Server is listening...")
	http.ListenAndServe(":8080", nil)
}

// fileWorkHandler обработчик, принимающий из url начальную директорию, лимит размера директории, тип сортировки и работает с директориями
func fileWorkHandler(w http.ResponseWriter, r *http.Request) {
	//Получаем данные из url
	startDirectory, sizeLimitMb, sortType, err := getUrlInput(r)
	if err != nil {
		fmt.Fprintf(w, "Error by trying get input! %e\n", err)
	}

	//Проверяем данные
	if err := pcg.CheckInput(startDirectory, sizeLimitMb, sortType); err != nil {
		fmt.Fprintf(w, "Error by checking input! %e\n", err)
		return
	}

	//Получаем директории - размеры
	dirSizes, err := pcg.GetDirSizes(startDirectory)
	if err != nil {
		fmt.Fprintf(w, "Error by geting directories %e\n", err)
		return
	}

	//Сортируем директории - размеры
	dirSizes, err = pcg.SortDirSizes(dirSizes, sortType)
	if err != nil {
		fmt.Fprintf(w, "Error by sorting directories! %e\n", err)
		return
	}

	//Получаем директории больше лимита
	dirSizeLargerLimit := pcg.GetDirsLargerLimit(dirSizes, sizeLimitMb)
	if err := pcg.WriteDirSizesToFile("result.txt", dirSizeLargerLimit); err != nil {
		fmt.Fprintf(w, "Eror by writing in file! %e\n", err)
		return
	}
	fmt.Fprint(w, "Done!")
}

func getUrlInput(r *http.Request) (string, int64, string, error) {
	vars := r.URL.Query()
	startDirectory := vars.Get("root")
	SizeLimitBytes, err := strconv.ParseInt(vars.Get("limit"), 10, 64)

	if err != nil {
		return "", 0, "", fmt.Errorf("eror by trying convert limit to int64! %e", err)
	}
	sortType := vars.Get("sort")

	return startDirectory, SizeLimitBytes * 1048576.0, sortType, nil
}
