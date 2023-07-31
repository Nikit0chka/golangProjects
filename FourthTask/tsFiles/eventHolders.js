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
import { renderElements } from "./pageRender";
import { addEventHandlerById } from "./addEventsOnElems";
//виды сортировок
var SortType;
(function (SortType) {
    SortType["ASC"] = "ASC";
    SortType["DESK"] = "DESK";
})(SortType || (SortType = {}));
let render = new renderElements();
//текущий путь и сортировка
let currentPath = 'D:\\';
var currentSort = SortType.ASC;
//buttonCLickHandler обработчик события нажатия на кнопку
function fileListClickHandler() {
    return __awaiter(this, void 0, void 0, function* () {
        const startTime = new Date();
        const clickedItem = event.target;
        const liElement = clickedItem.closest("li");
        const newPath = liElement.getAttribute("path");
        const data = yield sendJsonRequest(newPath, currentSort);
        render.addFilesToHtml(data);
        render.addRootToHtml(newPath);
        currentPath = newPath;
        render.addTimerToHtml(new Date().getTime() - startTime.getTime());
    });
}
//buttonCLickHandler обработчик события нажатия на кнопку
function rootClickHandler() {
    return __awaiter(this, void 0, void 0, function* () {
        const startTime = new Date();
        const clickedItem = event.target;
        const newPath = clickedItem.getAttribute("path");
        const data = yield sendJsonRequest(newPath, currentSort);
        render.addFilesToHtml(data);
        render.addRootToHtml(newPath);
        currentPath = newPath;
        render.addTimerToHtml(new Date().getTime() - startTime.getTime());
    });
}
//buttonCLickHandler обработчик события нажатия на кнопку
function buttonClickHandler() {
    return __awaiter(this, void 0, void 0, function* () {
        const startTime = new Date();
        const newSort = currentSort === SortType.ASC ? SortType.DESK : SortType.ASC;
        const data = yield sendJsonRequest(currentPath, newSort);
        render.addFilesToHtml(data);
        render.addRootToHtml(currentPath);
        currentSort = newSort;
        render.addTimerToHtml(new Date().getTime() - startTime.getTime());
    });
}
//buttonCLickHandler обработчик события нажатия на кнопку
function parentDirButtonClickHandler() {
    return __awaiter(this, void 0, void 0, function* () {
        const startTime = new Date();
        const lastSlashIndex = currentPath.lastIndexOf('\\');
        let newPath = currentPath.slice(0, lastSlashIndex);
        if (newPath[0] != "D")
            newPath = "D";
        const data = yield sendJsonRequest(newPath, currentSort);
        render.addFilesToHtml(data);
        render.addRootToHtml(newPath);
        currentPath = newPath;
        render.addTimerToHtml(new Date().getTime() - startTime.getTime());
    });
}
//
function directoriesButtonClickHandler() {
    return __awaiter(this, void 0, void 0, function* () {
        const startTime = new Date();
        const data = yield sendJsonRequest(currentPath, currentSort);
        render.addFilesToHtml(data);
        render.addRootToHtml(currentPath);
        render.addTimerToHtml(new Date().getTime() - startTime.getTime());
    });
}
//
function statisticButtonClickHandler() {
    return __awaiter(this, void 0, void 0, function* () {
    });
}
//события при загрузке окна
window.onload = function () {
    return __awaiter(this, void 0, void 0, function* () {
        const startTime = new Date();
        const data = yield sendJsonRequest(currentPath, currentSort);
        addEventHandlerById("sortButt", "click", buttonClickHandler);
        addEventHandlerById("folder-list", "click", fileListClickHandler);
        addEventHandlerById("currentDir", "click", rootClickHandler);
        addEventHandlerById("parentButt", "click", parentDirButtonClickHandler);
        addEventHandlerById("directoriesButt", "click", directoriesButtonClickHandler);
        render.addFilesToHtml(data);
        render.addRootToHtml(currentPath);
        render.addTimerToHtml(new Date().getTime() - startTime.getTime());
    });
};
