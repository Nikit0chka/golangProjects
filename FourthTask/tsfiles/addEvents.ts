// addEventsOnDirRoot добавляет обработчиков на нажатие по корню
function addEventsOnDirRoot() {
    const roots = document.querySelectorAll(".root");
    roots.forEach((root) => {
        root.addEventListener("click", function () {
            const path = root.getAttribute("path");
            sendJsonAndUpdateHtml(path);
        });
    });
}

// addEventsToFolders добавляет обработчиков на нажатие по папкам
function addEventsToFolders() {
    const folders = document.querySelectorAll(".folder-list li");
    folders.forEach((folder) => {
        folder.addEventListener("click", function () {
            const path = folder.getAttribute("path");
            sendJsonAndUpdateHtml(path);
        });
    });
}

// addEventOnButton добавляет обработчиков на нажатие по кнопке сортировки
function addEventOnButton() {
    const button = document.getElementById("sortButt");
    button.addEventListener("click", function () {
        const currentPath = getCurrentPath();
        currentSort === asc ? (currentSort = desk) : (currentSort = asc);
        sendJsonAndUpdateHtml(currentPath);
    });
}

// getCurrentPath возвращает путь к текущей директории
function getCurrentPath() {
    const roots = document.querySelectorAll(".root");
    let maxRootLen = "";

    for (let i = 0; i < roots.length; i++) if (roots[i].getAttribute("path") > maxRootLen) maxRootLen = roots[i].getAttribute("path");

    return maxRootLen;
}

// события при загрузке окна
window.onload = function () {
    addEventOnButton();
    sendJsonAndUpdateHtml("/");
};