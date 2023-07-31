//FileTypes типы файлов
var FileTypes;
(function (FileTypes) {
    FileTypes["FILE"] = "FILE";
    FileTypes["DIRECTORY"] = "DIR";
})(FileTypes || (FileTypes = {}));
export class renderElements {
    //addFilesToHtml добавляет файлы в html
    addFilesToHtml(folders) {
        const folderList = document.querySelector(".folder-list");
        folderList.innerHTML = "";
        for (let i = 0; i < folders.length; i++) {
            for (let j = 0; j < folders.length; j++) {
                if (folders[j].fileOrder === i) {
                    const folder = folders[i];
                    const fileSpace = document.createElement("li");
                    const fileNameSpace = document.createElement("div");
                    const fileIcon = document.createElement("img");
                    folder.type === FileTypes.DIRECTORY || folder.type === undefined ? fileIcon.src = "/static/dirImage.png" : fileIcon.src = "/static/fileImg.jpg";
                    const folderSize = document.createElement("span");
                    fileIcon.className = "file-icon";
                    fileSpace.className = "file-space";
                    fileSpace.setAttribute("name", folder.name);
                    fileSpace.setAttribute("path", folder.path);
                    fileNameSpace.appendChild(fileIcon);
                    fileNameSpace.appendChild(document.createTextNode(folder.name));
                    fileSpace.appendChild(fileNameSpace);
                    fileSpace.appendChild(folderSize);
                    folderSize.className = "folder-size";
                    folderSize.appendChild(document.createTextNode(folder.size + " mb"));
                    folderList.appendChild(fileSpace);
                    break;
                }
            }
        }
    }
    //addRootToHtml добавляет путь от корня в html
    addRootToHtml(path) {
        const dirs = path.split('\\');
        const htmlPath = document.getElementById("currentDir");
        htmlPath.innerHTML = "";
        //всегда добавляем путь к корню
        const root = document.createElement("a");
        let currentPath = 'D:';
        root.setAttribute("path", currentPath);
        root.className = "root";
        root.appendChild(document.createTextNode("start\\"));
        htmlPath.appendChild(root);
        for (let i = 1; i < dirs.length; i++) {
            if (dirs[i] === "")
                continue;
            const root = document.createElement("a");
            currentPath += '\\' + dirs[i];
            root.setAttribute("path", currentPath);
            root.className = "root";
            root.appendChild(document.createTextNode(dirs[i]));
            root.appendChild(document.createTextNode("\\"));
            htmlPath.appendChild(root);
        }
    }
    //addTimerToHtml выводит результат работы таймера
    addTimerToHtml(result) {
        const timer = document.getElementById("timer");
        timer.innerHTML = "";
        const span = document.createElement("span");
        timer.appendChild(span);
        span.appendChild(document.createTextNode(result + "ms"));
    }
}
