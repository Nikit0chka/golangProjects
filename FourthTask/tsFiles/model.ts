//url сервера
const url: string = "./"

//структура json ответа
export class JsonResponse {
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
export async function sendJsonRequest(path: string, sortType: string): Promise<JsonResponse[]> {
    const data = new JsonRequest(path, sortType);
    const response = await fetch(url, {
        method: "POST",
        headers: {
            "Content-Type": "application/json"
        },
        body: JSON.stringify(data)
    });
    const json = await response.json();
    return json.map((item: any) => new JsonResponse(item.name, item.fileOrder, item.path, item.size, item.type));
}