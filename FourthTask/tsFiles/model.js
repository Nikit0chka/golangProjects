var __awaiter = (this && this.__awaiter) || function (thisArg, _arguments, P, generator) {
    function adopt(value) { return value instanceof P ? value : new P(function (resolve) { resolve(value); }); }
    return new (P || (P = Promise))(function (resolve, reject) {
        function fulfilled(value) { try { step(generator.next(value)); } catch (e) { reject(e); } }
        function rejected(value) { try { step(generator["throw"](value)); } catch (e) { reject(e); } }
        function step(result) { result.done ? resolve(result.value) : adopt(result.value).then(fulfilled, rejected); }
        step((generator = generator.apply(thisArg, _arguments || [])).next());
    });
};
//url сервера
const url = "./";
//структура json ответа
class JsonResponse {
    constructor(name, fileOrder, path, size, type) {
        this.name = name;
        this.fileOrder = fileOrder;
        this.path = path;
        this.size = size;
        this.type = type;
    }
}
// структура json запроса
class JsonRequest {
    constructor(path, sortType) {
        this.path = path;
        this.sortType = sortType;
    }
}
// sendJsonRequest отправляет и обрабатывает json к серверу и от сервера
export function sendJsonRequest(path, sortType) {
    return __awaiter(this, void 0, void 0, function* () {
        const data = new JsonRequest(path, sortType);
        const response = yield fetch(url, {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify(data)
        });
        const json = yield response.json();
        return json.map((item) => new JsonResponse(item.name, item.fileOrder, item.path, item.size, item.type));
    });
}
