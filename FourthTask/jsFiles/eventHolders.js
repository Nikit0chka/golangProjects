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
const asc = "ASC";
const desk = "DESK";
//текущий путь и сортировка
let currentPath = '/';
let currentSort = asc;
const firstPathIndex = 0;
//buttonCLickHandler обработчик события нажатия на кнопку
const folderClickHandler = function () {
    return __awaiter(this, void 0, void 0, function* () {
        let startTime = new Date();
        let data = yield sendJsonRequest(this.getAttribute("path"), currentSort);
        addFilesToHtml(data, folderClickHandler);
        addRootToHtml(data[firstPathIndex].path, rootClickHandler);
        currentPath = data[firstPathIndex].path;
        addTimerToHtml(new Date().getTime() - startTime.getTime() / 1000);
    });
};
//buttonCLickHandler обработчик события нажатия на кнопку
function rootClickHandler() {
    return __awaiter(this, void 0, void 0, function* () {
        let startTime = new Date();
        let data = yield sendJsonRequest(this.getAttribute("path"), currentSort);
        addFilesToHtml(data, folderClickHandler);
        addRootToHtml(data[firstPathIndex].path, rootClickHandler);
        currentPath = data[firstPathIndex].path;
        addTimerToHtml(new Date().getTime() - startTime.getTime() / 1000);
    });
}
//buttonCLickHandler обработчик события нажатия на кнопку
function buttonClickHandler() {
    return __awaiter(this, void 0, void 0, function* () {
        let startTime = new Date();
        let newSort = currentSort === asc ? desk : asc;
        let data = yield sendJsonRequest(currentPath, newSort);
        addFilesToHtml(data, folderClickHandler);
        addRootToHtml(data[firstPathIndex].path, rootClickHandler);
        currentSort = newSort;
        addTimerToHtml(new Date().getTime() - startTime.getTime() / 1000);
    });
}
//события при загрузке окна
window.onload = function () {
    return __awaiter(this, void 0, void 0, function* () {
        let startTime = new Date();
        let data = yield sendJsonRequest("/", currentSort);
        addFilesToHtml(data, folderClickHandler);
        addRootToHtml(data[firstPathIndex].path, rootClickHandler);
        addTimerToHtml(new Date().getTime() - startTime.getTime() / 1000);
    });
};
