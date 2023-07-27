import { sendJsonRequest } from "./model";
import { addRootToHtml, addTimerToHtml, addFilesToHtml } from "./draw";

//виды сортировок
const asc: string = "ASC";
const desk: string = "DESK";

//текущий путь и сортировка
let currentPath: string = '\\';
let currentSort: string = asc;

const firstPathIndex: number = 0;

//buttonCLickHandler обработчик события нажатия на кнопку
const folderClickHandler = async function(this: HTMLElement): Promise<void> {
    let startTime = new Date();
    let data = await sendJsonRequest(this.getAttribute("path")!, currentSort);
    addFilesToHtml(data, folderClickHandler);
    addRootToHtml(data[firstPathIndex].path, rootClickHandler);
    currentPath = data[firstPathIndex].path;
    addTimerToHtml(new Date().getTime() - startTime.getTime()/ 1000);
};

//buttonCLickHandler обработчик события нажатия на кнопку
async function rootClickHandler(this: HTMLElement): Promise<void> {
    let startTime = new Date();
    let data = await sendJsonRequest(this.getAttribute("path")!, currentSort);
    addFilesToHtml(data, folderClickHandler);
    addRootToHtml(data[firstPathIndex].path, rootClickHandler);
    currentPath = data[firstPathIndex].path;
    addTimerToHtml(new Date().getTime() - startTime.getTime()/ 1000);
}

//buttonCLickHandler обработчик события нажатия на кнопку
async function buttonClickHandler(): Promise<void> {
    let startTime = new Date();
    let newSort = currentSort === asc ? desk : asc;
    let data = await sendJsonRequest(currentPath, newSort);
    addFilesToHtml(data, folderClickHandler);
    addRootToHtml(data[firstPathIndex].path, rootClickHandler);
    currentSort = newSort;
    addTimerToHtml(new Date().getTime() - startTime.getTime()/ 1000);
}

//события при загрузке окна
window.onload = async function(): Promise<void> {
    let startTime = new Date();
    let data = await sendJsonRequest("\\", currentSort);

    addFilesToHtml(data, folderClickHandler);
    addRootToHtml(data[firstPathIndex].path, rootClickHandler);
    addTimerToHtml(new Date().getTime() - startTime.getTime()/ 1000);
};