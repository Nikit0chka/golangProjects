import {sendJsonRequest} from "./model.js";
import {addRootToHtml, addTimerToHtml, addFilesToHtml} from "./draw.js";

//виды сортировок
const asc = "ASC"
const desk = "DESK"

//текущий путь и сортировка
var currentPath = '\\'
var currentSort = asc

//
const firstPathIndex = 0

//buttonCLickHandler обработчик события нажатия на кнопку
const folderClickHandler = async function (){
    let startTime = new Date();
    let data = await sendJsonRequest(this.getAttribute("path"), currentSort)
    addFilesToHtml(data, folderClickHandler);
    addRootToHtml(data[firstPathIndex].path, rootClickHandler);
    currentPath = data[firstPathIndex].path
    addTimerToHtml(new Date() - startTime)
};

//buttonCLickHandler обработчик события нажатия на кнопку
async function rootClickHandler(){
    let startTime = new Date();
    let data = await sendJsonRequest(this.getAttribute("path"), currentSort)
    addFilesToHtml(data, folderClickHandler);
    addRootToHtml(data[firstPathIndex].path, rootClickHandler);
    currentPath = data[firstPathIndex].path
    addTimerToHtml(new Date() - startTime)
}

//buttonCLickHandler обработчик события нажатия на кнопку
async function buttonClickHandler(){
    let startTime = new Date();
    let newSort = currentSort === asc?desk:asc
    let data =  await sendJsonRequest(currentPath, newSort)
    addFilesToHtml(data, folderClickHandler)
    addRootToHtml(data, rootClickHandler)
    currentSort = newSort
    addTimerToHtml(new Date() - startTime)
}

//события при загрузке окна
window.onload = async function () {
    let startTime = new Date()
    let data = await sendJsonRequest("\\", currentSort)

    addFilesToHtml(data, folderClickHandler)
    addRootToHtml(data[firstPathIndex].path, rootClickHandler)
    addTimerToHtml(new Date() - startTime)
}