//url сервера
const url = "./"

//структура json ответа
export class JsonResponse{
    constructor(name, fileOrder, path, size, type) {
        this.name = name
        this.fileOrder = fileOrder
        this.path = path
        this.size = size
        this.type = type
    }
}

// структура json запроса
class JsonRequest{
    constructor(path, sortType) {
        this.path = path
        this.sortType = sortType
    }
}

// sendJsonRequest отправляет и обрабатывает json к серверу и от сервера
export async function sendJsonRequest(path, sortType) {
    const data = new JsonRequest(path, sortType);
    const response = await fetch(url, {
        method: "POST",
        headers: {
            "Content-Type": "application/json"
        },
        body: JSON.stringify(data)
    })
    const json = await response.json()
    return json.map(item => new JsonResponse(item.name, item.fileOrder, item.path, item.size, item.type))
}