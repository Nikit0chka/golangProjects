//структура json ответа
class JsonResponse{
    constructor(name, dirId, path, size, type) {
        this.name = name;
        this.dirId = dirId
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
const url = "http://localhost:1144/dir-sizes";

//типы сортировок
const asc = "ASC"
const desk = "DESK"

//addEventsOnDirRoot добавляет обработчиков на нажатие по корню
function addEventsOnDirRoot() {
    const roots = document.querySelectorAll(".root-path");
    roots.forEach(root => {
        root.addEventListener("click", function (){
            let path = root.getAttribute("path")
            sendJsonAndUpdateHtml(path)
        });
    });
}

//addEventsToFolders добавляет обработчиков на нажатие по папкам
function addEventsToFolders(){
    const folders = document.querySelectorAll(".folder-list li a");
    folders.forEach(folder => {
        folder.addEventListener("click", function() {
            let path = folder.getAttribute("path");
            sendJsonAndUpdateHtml(path)
        });
    });
}

//addEventOnButton добавляет обработчиков на нажатие по кнопке сортировки
 function addEventOnButton() {
     const button = document.getElementById("sortButt");
     button.addEventListener("click", function () {
         button.getAttribute("currentSort") === asc ? button.setAttribute("currentSort", desk) : button.setAttribute("currentSort", asc)
         let currentPath = getCurrentPath()
         sendJsonAndUpdateHtml(currentPath);
     });
}

//getCurrentSortType возвращает текущий тип сортировки с кнопки
function getCurrentSortType(){
    const button = document.getElementById("sortButt")
    return button.getAttribute("currentSort") === "ASC"? asc:desk
}

//getCurrentPath возвращает путь к текущей директории
function getCurrentPath(){
    const roots = document.querySelectorAll(".root-path")
    let maxRootLen = ""

    for (let i = 0; i < roots.length; i++)
        if (roots[i].getAttribute("path") > maxRootLen)
            maxRootLen = roots[i].getAttribute("path")

    return maxRootLen
}

// sendJsonRequest отправляет и обрабатывает json к серверу и от сервера
async function sendJsonRequest(jsonRequest) {
    const data = {
        path: jsonRequest.path,
        sortType: jsonRequest.sortType
    };
    const response = await fetch(url, {
        method: "POST",
        headers: {
            "Content-Type": "application/json"
        },
        body: JSON.stringify(data)
    });
    const json = await response.json();
    return json.map(item => new JsonResponse(item.name,  item.dirId, item.path, item.size, item.type));
}

//addRootToHtml добавляет путь от корня в html
function addRootToHtml(path){
    let dirs = path.split('/');
    let root = document.getElementById("currentDir")
    root.innerHTML = "";
    let currentPath = "";

    for (let i = 0; i < dirs.length-1; i++)
    {
        let span = document.createElement("span");
        currentPath += dirs[i] + '/';
        span.setAttribute("path", currentPath);
        span.className = "root-path"
        span.appendChild(document.createTextNode(dirs[i]))
        span.appendChild(document.createTextNode("/"))
        root.appendChild(span);
    }
    addEventsOnDirRoot()
}

//addFilesToHtml добавляет файлы в html
function addFilesToHtml(folders) {
    let folderList = document.querySelector(".folder-list");
    folderList.innerHTML = "";
    for (let i = 0; i < folders.length; i++)
    {
        for (let j = 0; j < folders.length; j++)
        {
            if (folders[j].dirId === i)
            {
                let folder = folders[i];
                let li = document.createElement("li");
                let a = document.createElement("a");
                let img = document.createElement("img");
                folder.type === "DIR" ? img.src = "dirImage.png" : img.src = "fileImg.jpg"
                let span = document.createElement("span");
                img.className = "image"
                li.appendChild(a);
                a.setAttribute("name", folder.name);
                a.setAttribute("path", folder.path)
                a.appendChild(img)
                a.appendChild(document.createTextNode(folder.name));
                a.appendChild(span);
                span.classList.add("folder-size");
                span.appendChild(document.createTextNode(folder.size + "  mb"));
                folderList.appendChild(li);
                break
            }
        }    
    }
    addEventsToFolders()
}

//addTimerToHtml выводит результат работы таймера
function addTimerToHtml(result){
    let timer = document.getElementById("timer")
    timer.innerHTML = ""

    let span = document.createElement("span")
    timer.appendChild(span)
    span.appendChild(document.createTextNode(result))
}

// sendJsonAndUpdateHtml отправляет json и по результатам ответа изменяет html
function sendJsonAndUpdateHtml(path){
    const start= new Date().getTime();

    let currentSortType = getCurrentSortType()
    let jsonRequest = new JsonRequest(path, currentSortType)

    sendJsonRequest(jsonRequest).then(data => {addRootToHtml(data[0].path); addFilesToHtml(data)}).catch(error => console.error(error))
    const end = new Date().getTime();

    addTimerToHtml(end - start)
}

//события при загрузке окна
window.onload = function (){
    addEventOnButton()
    sendJsonAndUpdateHtml("/home")
}
