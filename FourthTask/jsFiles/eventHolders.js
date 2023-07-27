var __awaiter = (this && this.__awaiter) || function (thisArg, _arguments, P, generator) {
    function adopt(value) { return value instanceof P ? value : new P(function (resolve) { resolve(value); }); }
    return new (P || (P = Promise))(function (resolve, reject) {
        function fulfilled(value) { try { step(generator.next(value)); } catch (e) { reject(e); } }
        function rejected(value) { try { step(generator["throw"](value)); } catch (e) { reject(e); } }
        function step(result) { result.done ? resolve(result.value) : adopt(result.value).then(fulfilled, rejected); }
        step((generator = generator.apply(thisArg, _arguments || [])).next());
    });
};
import { sendJsonRequest } from "./model";
import { addRootToHtml, addTimerToHtml, addFilesToHtml } from "./draw";
//виды сортировок
var SortType;
(function (SortType) {
    SortType["asc"] = "ASC";
    SortType["desk"] = "DESK";
})(SortType || (SortType = {}));
//текущий путь и сортировка
let currentPath = '/';
var currentSort = SortType.asc;
//buttonCLickHandler обработчик события нажатия на кнопку
const fileListClickHandler = function () {
    return __awaiter(this, void 0, void 0, function* () {
        const startTime = new Date();
        const clickedItem = event.target;
        const liElement = clickedItem.closest("li");
        const path = liElement.getAttribute("path");
        const data = yield sendJsonRequest(path, currentSort);
        addFilesToHtml(data);
        addRootToHtml(path);
        currentPath = path;
        addTimerToHtml(new Date().getTime() - startTime.getTime());
    });
};
//buttonCLickHandler обработчик события нажатия на кнопку
function rootClickHandler() {
    return __awaiter(this, void 0, void 0, function* () {
        const startTime = new Date();
        const clickedItem = event.target;
        const newPath = clickedItem.getAttribute("path");
        const data = yield sendJsonRequest(newPath, currentSort);
        addFilesToHtml(data);
        addRootToHtml(newPath);
        currentPath = newPath;
        addTimerToHtml(new Date().getTime() - startTime.getTime());
    });
}
//buttonCLickHandler обработчик события нажатия на кнопку
function buttonClickHandler() {
    return __awaiter(this, void 0, void 0, function* () {
        const startTime = new Date();
        const newSort = currentSort === SortType.asc ? SortType.desk : SortType.asc;
        const data = yield sendJsonRequest(currentPath, newSort);
        addFilesToHtml(data);
        addRootToHtml(currentPath);
        currentSort = newSort;
        addTimerToHtml(new Date().getTime() - startTime.getTime());
    });
}
//buttonCLickHandler обработчик события нажатия на кнопку
function parentDirButtonClickHandler() {
    return __awaiter(this, void 0, void 0, function* () {
        const startTime = new Date();
        const lastSlashIndex = currentPath.lastIndexOf('/');
        let newPath = currentPath.slice(0, lastSlashIndex);
        if (newPath[0] != "/")
            newPath = "/";
        const data = yield sendJsonRequest(newPath, currentSort);
        console.log(newPath);
        addFilesToHtml(data);
        addRootToHtml(newPath);
        currentPath = newPath;
        addTimerToHtml(new Date().getTime() - startTime.getTime());
    });
}
//addEventHandlerOnSortButton добавляет обработчик события нажатия на кнопку сортировки
function addEventHandlerOnSortButton() {
    const sortButton = document.getElementById("sortButt");
    sortButton.addEventListener("click", buttonClickHandler);
}
// addEventHandlerOnFileList добавляет обработчик события нажатия на список директорий
function addEventHandlerOnFileList() {
    const fileList = document.getElementById("folder-list");
    fileList.addEventListener("click", fileListClickHandler);
}
// addEventHandlerOnRoot добавляет обработчик события нажатия на путь
function addEventHandlerOnRoot() {
    const root = document.getElementById("currentDir");
    root.addEventListener("click", rootClickHandler);
}
function addEventHandlerOnParentDirButton() {
    const parentDirButton = document.getElementById("parentButt");
    parentDirButton.addEventListener("click", parentDirButtonClickHandler);
}
//события при загрузке окна
window.onload = function () {
    return __awaiter(this, void 0, void 0, function* () {
        const startTime = new Date();
        const data = yield sendJsonRequest("/", currentSort);
        addEventHandlerOnSortButton();
        addEventHandlerOnFileList();
        addEventHandlerOnRoot();
        addEventHandlerOnParentDirButton();
        addFilesToHtml(data);
        addRootToHtml(currentPath);
        addTimerToHtml(new Date().getTime() - startTime.getTime());
    });
};
