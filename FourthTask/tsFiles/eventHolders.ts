import { sendJsonRequest } from "./model";
import { addRootToHtml, addTimerToHtml, addFilesToHtml } from "./draw";

//виды сортировок
const asc: string = "ASC";
const desk: string = "DESK";

//текущий путь и сортировка
let currentPath: string = '/';
let currentSort: string = asc;

//индекс для получения первого пути
const firstPathIndex: number = 0;

//buttonCLickHandler обработчик события нажатия на кнопку
const fileListClickHandler = async function(this: HTMLElement): Promise<void> {
    let startTime = new Date();
    const clickedItem = event.target as HTMLElement;

    let liElement = clickedItem.closest("li")
    let path = liElement.getAttribute("path")
    let data = await sendJsonRequest(path, currentSort);

    addFilesToHtml(data);
    addRootToHtml(path);

    currentPath = path;

    addTimerToHtml(new Date().getTime() - startTime.getTime());
};

//buttonCLickHandler обработчик события нажатия на кнопку
async function rootClickHandler(this: HTMLElement): Promise<void> {
    let startTime = new Date();
    const clickedItem = event.target as HTMLElement;

    let path = clickedItem.getAttribute("path")
    let data = await sendJsonRequest(path, currentSort);

    addFilesToHtml(data);
    addRootToHtml(path);

    currentPath = path;

    addTimerToHtml(new Date().getTime() - startTime.getTime());
}

//buttonCLickHandler обработчик события нажатия на кнопку
async function buttonClickHandler(): Promise<void> {
    let startTime = new Date();
    let newSort = currentSort === asc ? desk : asc;

    let data = await sendJsonRequest(currentPath, newSort);

    addFilesToHtml(data);
    addRootToHtml(currentPath);

    currentSort = newSort;

    addTimerToHtml(new Date().getTime() - startTime.getTime());
}

function addEventHandlerOnSortButton(){
    let sortButton = document.getElementById("sortButt");
    sortButton.addEventListener("click", buttonClickHandler);
}

function addEventHandlerOnFileList(){
    let fileList = document.getElementById("folder-list")
    fileList.addEventListener("click", fileListClickHandler)
}

function addEventHandlerOnRoot(){
    let root = document.getElementById("currentDir")
    root.addEventListener("click", rootClickHandler)
}

//события при загрузке окна
window.onload = async function(): Promise<void> {
    let startTime = new Date();
    let data = await sendJsonRequest("/", currentSort);

    addEventHandlerOnSortButton();
    addEventHandlerOnFileList();
    addEventHandlerOnRoot();

    addFilesToHtml(data);
    addRootToHtml(currentPath);

    addTimerToHtml(new Date().getTime() - startTime.getTime());
};