//структура json ответа
class JsonResponse{
    constructor(name, path, size, type) {
        this.name = name;
        this.path = path;
        this.size = size;
        this.type = type;
    }
}
// структура json запроса
class JsonRequest{
    constructor(path, sortType) {
        this.path = path;
        this.sortType = sortType;
    }
}

//url домена
const url = '/dir-sizes';

//addEventsToFolders добавляет обработчиков на нажатие
function addEventsToFolders(){
    const folders = document.querySelectorAll('.folder-list li a');
    folders.forEach(folder => {
        folder.addEventListener('click', function() {
            sendJsonRequest(new JsonRequest(folder.getAttribute("path"), "ASC")).then(data => addFilesToHtml(data)).catch(error => console.error(error));
        });
    });
}

// sendJsonRequest отправляет и обрабатывает json к серверу и от сервера
async function sendJsonRequest(jsonRequest) {
    const data = {
        path: jsonRequest.path,
        sortType: jsonRequest.sortType
    };
    console.log(jsonRequest.path)

    const response = await fetch(url, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(data)
    });

    const json = await response.json();
    const result = json.map(item => new JsonResponse(item.name, item.path, item.size, item.type));
    return result;
}

//addFilesToHtml добавляет файлы в html
function addFilesToHtml(folders) {
    let folderList = document.querySelector(".folder-list");
    folderList.innerHTML = "";
    for (let i = 0; i < folders.length; i++) {
        let folder = folders[i];
        let li = document.createElement("li");
        let a = document.createElement("a");
        let img = document.createElement("img");
        folder.type == "DIR" ? img.src = "dirImage.png" : img.src = "fileImg.jpg"
        let span = document.createElement("span");
        img.className = "image"
        li.appendChild(a);
        a.setAttribute("name", folder.name);
        a.setAttribute("path", folder.path)
        a.appendChild(img)
        a.appendChild(document.createTextNode(folder.name));
        a.appendChild(span);
        span.classList.add("folder-size");
        if (folder.size != 0)
            span.appendChild(document.createTextNode(folder.size + "  mb"));
        folderList.appendChild(li);
    }
    addEventsToFolders()
}

window.onload = function (){
    let jsonRequest = new JsonRequest("D:\\", "ASC");

    sendJsonRequest(jsonRequest).then(data => addFilesToHtml(data)).catch(error => console.error(error));
}
