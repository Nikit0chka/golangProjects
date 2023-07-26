import {sendJsonRequest, currentSort, currentPath, asc, desk} from "./model";

const anyPath = 0

function folderClickHandler(){
    let data = sendJsonRequest(this.getAttribute("path"), currentSort)
    addFilesToHtml(data);
    addRootToHtml(data);
    currentPath = data[anyPath].path
}

function rootClickHandler(){
    let data = sendJsonRequest(this.getAttribute("path"), currentSort)
    addFilesToHtml(data);
    addRootToHtml(data);
    currentPath = data[anyPath].path
}

function buttonClickHandler(){
    let newSort = currentSort === asc?desk:asc
    let data = sendJsonRequest(currentPath, newSort)
    addFilesToHtml(data)
    addRootToHtml(data)
    currentSort = newSort
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
                fileSpace.addEventListener("click", folderClickHandler)
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
    root.addEventListener("click", rootClickHandler)
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
        root.addEventListener("click", rootClickHandler)
        root.appendChild(document.createTextNode(dirs[i]))
        root.appendChild(document.createTextNode("/"))
        htmlPath.appendChild(root)
    }
}

//addTimerToHtml выводит результат работы таймера
function addTimerToHtml(result){
    let timer = document.getElementById("timer")
    timer.innerHTML = ""
    let span = document.createElement("span")
    timer.appendChild(span)
    span.appendChild(document.createTextNode(result  + "ms"))
}