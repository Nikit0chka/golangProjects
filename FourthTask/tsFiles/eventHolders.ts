import { sendJsonRequest } from "./model";
import { addRootToHtml, addTimerToHtml, addFilesToHtml } from "./draw";

//виды сортировок
enum SortType{
    asc = "ASC",
    desk = "DESK"
}

//текущий путь и сортировка
let currentPath: string = '/';
var currentSort: SortType = SortType.asc;

//buttonCLickHandler обработчик события нажатия на кнопку
const fileListClickHandler = async function(this: HTMLElement): Promise<void> {
    const startTime = new Date();

    const clickedItem = event.target as HTMLElement;
    const liElement = clickedItem.closest("li")

    const path = liElement.getAttribute("path")
    const data = await sendJsonRequest(path, currentSort);

    addFilesToHtml(data);
    addRootToHtml(path);

    currentPath = path;

    addTimerToHtml(new Date().getTime() - startTime.getTime());
}

//buttonCLickHandler обработчик события нажатия на кнопку
async function rootClickHandler(this: HTMLElement): Promise<void> {
    const startTime = new Date();

    const clickedItem = event.target as HTMLElement;

    const newPath = clickedItem.getAttribute("path")
    const data = await sendJsonRequest(newPath, currentSort);

    addFilesToHtml(data);
    addRootToHtml(newPath);

    currentPath = newPath;

    addTimerToHtml(new Date().getTime() - startTime.getTime());
}

//buttonCLickHandler обработчик события нажатия на кнопку
async function buttonClickHandler(): Promise<void> {
    const startTime = new Date();

    const newSort = currentSort === SortType.asc ? SortType.desk : SortType.asc;

    const data = await sendJsonRequest(currentPath, newSort);

    addFilesToHtml(data);
    addRootToHtml(currentPath);

    currentSort = newSort;

    addTimerToHtml(new Date().getTime() - startTime.getTime());
}

//buttonCLickHandler обработчик события нажатия на кнопку
async function parentDirButtonClickHandler(): Promise<void> {
    const startTime = new Date();

    const lastSlashIndex = currentPath.lastIndexOf('/')
    let newPath = currentPath.slice(0, lastSlashIndex)

    if (newPath[0] != "/")
        newPath = "/"
    const data = await sendJsonRequest(newPath, currentSort);

    console.log(newPath)

    addFilesToHtml(data);
    addRootToHtml(newPath);

    currentPath = newPath;

    addTimerToHtml(new Date().getTime() - startTime.getTime());
}

//addEventHandlerOnSortButton добавляет обработчик события нажатия на кнопку сортировки
function addEventHandlerOnSortButton(){
    const sortButton = document.getElementById("sortButt");
    sortButton.addEventListener("click", buttonClickHandler);
}

// addEventHandlerOnFileList добавляет обработчик события нажатия на список директорий
function addEventHandlerOnFileList(){
    const fileList = document.getElementById("folder-list")
    fileList.addEventListener("click", fileListClickHandler)
}

// addEventHandlerOnRoot добавляет обработчик события нажатия на путь
function addEventHandlerOnRoot(){
    const root = document.getElementById("currentDir")
    root.addEventListener("click", rootClickHandler)
}

function addEventHandlerOnParentDirButton(){
    const parentDirButton = document.getElementById("parentButt")
    parentDirButton.addEventListener("click", parentDirButtonClickHandler)
}

//события при загрузке окна
window.onload = async function(): Promise<void> {
    const startTime = new Date();
    const data = await sendJsonRequest("/", currentSort);

    addEventHandlerOnSortButton();
    addEventHandlerOnFileList();
    addEventHandlerOnRoot();
    addEventHandlerOnParentDirButton()

    addFilesToHtml(data);
    addRootToHtml(currentPath);

    addTimerToHtml(new Date().getTime() - startTime.getTime());
}