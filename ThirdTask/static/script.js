//структура json ответа
class JsonResponse{
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

//url сервера
const url = "./"
//типы сортировок
const asc = "ASC"
const desk = "DESK"

var currentSort = asc

//addEventsOnDirRoot добавляет обработчиков на нажатие по корню
function addEventsOnDirRoot() {
  const roots = document.querySelectorAll(".root")
  roots.forEach(root => {
    root.addEventListener("click", function (){
      let path = root.getAttribute("path")
      sendJsonAndUpdateHtml(path)
    })
  })
}

//addEventsToFolders добавляет обработчиков на нажатие по папкам
function addEventsToFolders(){
  const folders = document.querySelectorAll(".folder-list li")
  folders.forEach(folder => {
    folder.addEventListener("click", function() {
      let path = folder.getAttribute("path")
      sendJsonAndUpdateHtml(path)
    })
  })
}

//addEventOnButton добавляет обработчиков на нажатие по кнопке сортировки
function addEventOnButton() {
  const button = document.getElementById("sortButt")
  button.addEventListener("click", function () {
    let currentPath = getCurrentPath();
    currentSort === asc ? currentSort = desk : currentSort = asc
    sendJsonAndUpdateHtml(currentPath)
  })
}

//getCurrentPath возвращает путь к текущей директории
function getCurrentPath(){
  const roots = document.querySelectorAll(".root")
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
  }
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

//addRootToHtml добавляет путь от корня в html
function addRootToHtml(path){
// ничего не менять если путь пустой
  if (path.length === 0)
    return

  let dirs = path.split('/')
  let htmlPath = document.getElementById("currentDir")
  htmlPath.innerHTML = ""

//всегда добавляем путь к корню
  let root = document.createElement("a")
  let currentPath = '/'
  root.setAttribute("path", currentPath)
  root.className = "root"
  root.appendChild(document.createTextNode("/"))
  htmlPath.appendChild(root)
  currentPath = ''

  for (let i = 1; i < dirs.length-1; i++)
  {
    if (dirs[i].length === 0)
      continue

    let root = document.createElement("a")
    currentPath += '/' + dirs[i]
    root.setAttribute("path", currentPath)
    root.className = "root"
    root.appendChild(document.createTextNode(dirs[i]))
    root.appendChild(document.createTextNode("/"))
    htmlPath.appendChild(root)
  }
  addEventsOnDirRoot()

}

//addFilesToHtml добавляет файлы в html
function addFilesToHtml(folders) {
  let folderList = document.querySelector(".folder-list")
  folderList.innerHTML = ""
  for (let i = 0; i < folders.length; i++)
  {
//ничего не добавлять, если имя нового файла пустое
    if (folders[i].name.length === 0)
      continue
    for (let j = 0; j < folders.length; j++)
    {
      if (folders[j].fileOrder === i)
      {
        let folder = folders[i]
        let fileSpace = document.createElement("li")
        let fileNameSpace = document.createElement("li")
        let fileIcon = document.createElement("img")
        folder.type === "DIR" ? fileIcon.src = "/static/dirImage.png" : fileIcon.src = "/static/fileImg.jpg"
        let folderSize = document.createElement("span")
        fileIcon.className = "file-icon"
        fileSpace.className = "file-space"
        fileSpace.appendChild(fileNameSpace)
        fileSpace.setAttribute("name", folder.name)
        fileSpace.setAttribute("path", folder.path)
        fileNameSpace.appendChild(fileIcon)
        fileNameSpace.appendChild(document.createTextNode(folder.name))
        fileSpace.appendChild(folderSize)
        folderSize.className = "folder-size"
        folderSize.appendChild(document.createTextNode(folder.size + " mb"))
        folderList.appendChild(fileSpace)

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
  span.appendChild(document.createTextNode(result  + "ms"))
}

// sendJsonAndUpdateHtml отправляет json и по результатам ответа изменяет html
async function sendJsonAndUpdateHtml(path){
  const start= new Date().getTime()
  let jsonRequest = new JsonRequest(path, currentSort)

  await sendJsonRequest(jsonRequest).then(data => {addRootToHtml(data[0].path); addFilesToHtml(data) }).catch(error => console.error(error))
  const end = new Date().getTime()
  addTimerToHtml(end - start)

}

//события при загрузке окна
window.onload = function (){
  addEventOnButton()
  sendJsonAndUpdateHtml("/")
}