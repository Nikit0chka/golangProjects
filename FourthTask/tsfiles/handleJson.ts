import {addRootToHtml, addFilesToHtml, addTimerToHtml} from "./editHtml";
// типы сортировок
export const asc = "ASC";
export const desk = "DESK";

export let currentSort = asc;

// url сервера
export const url = "./";

// структура json ответа
class JsonResponse {
    public name: string;
    public fileOrder: number;
    public path: string;
    public size: number;
    public type: string;

    constructor(name: string, fileOrder: number, path: string, size: number, type: string) {
        this.name = name;
        this.fileOrder = fileOrder;
        this.path = path;
        this.size = size;
        this.type = type;
    }
}

// структура json запроса
class JsonRequest {
    public path: string;
    public sortType: string;

    constructor(path: string, sortType: string) {
        this.path = path;
        this.sortType = sortType;
    }
}

// sendJsonRequest отправляет и обрабатывает json к серверу и от сервера
async function sendJsonRequest(jsonRequest: JsonRequest) {
    const data = {
        path: jsonRequest.path,
        sortType: jsonRequest.sortType,
    };
    const response = await fetch(url, {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
        },
        body: JSON.stringify(data),
    });
    const json = await response.json();
    return json.map((item: any) => new JsonResponse(item.name, item.fileOrder, item.path, item.size, item.type));
}

// sendJsonAndUpdateHtml отправляет json и по результатам ответа изменяет html
export async function sendJsonAndUpdateHtml(path: string) {
    const start = new Date().getTime();

    let jsonRequest = new JsonRequest(path, currentSort);

    await sendJsonRequest(jsonRequest)
        .then((data) => {
            addRootToHtml(data[0].path);
            addFilesToHtml(data);
        })
        .catch((error) => console.error(error));
    const end = new Date().getTime();
    addTimerToHtml(end - start);
}