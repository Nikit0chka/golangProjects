import { sendJsonRequest } from "./model";
import {renderElements} from "./pageRender";
import {addEventHandlerById} from "./addEventsOnElems";

//виды сортировок
enum SortType{
    ASC = "ASC",
    DESK = "DESK"
}

let render = new renderElements()

//текущий путь и сортировка
let currentPath: string = 'D:\\';
var currentSort: SortType = SortType.ASC;

//buttonCLickHandler обработчик события нажатия на кнопку
async function fileListClickHandler (this: HTMLElement): Promise<void> {
    const startTime = new Date();

    const clickedItem = event.target as HTMLElement;
    const liElement = clickedItem.closest("li")

    const newPath = liElement.getAttribute("path")
    const data = await sendJsonRequest(newPath, currentSort);

    render.addFilesToHtml(data);
    render.addRootToHtml(newPath);

    currentPath = newPath;

    render.addTimerToHtml(new Date().getTime() - startTime.getTime())
}

//buttonCLickHandler обработчик события нажатия на кнопку
async function rootClickHandler(this: HTMLElement): Promise<void> {
    const startTime = new Date();

    const clickedItem = event.target as HTMLElement;

    const newPath = clickedItem.getAttribute("path")
    const data = await sendJsonRequest(newPath, currentSort);

    render.addFilesToHtml(data);
    render.addRootToHtml(newPath);

    currentPath = newPath;

    render.addTimerToHtml(new Date().getTime() - startTime.getTime())
}

//buttonCLickHandler обработчик события нажатия на кнопку
async function buttonClickHandler(): Promise<void> {
    const startTime = new Date();

    const newSort = currentSort === SortType.ASC ? SortType.DESK : SortType.ASC;

    const data = await sendJsonRequest(currentPath, newSort);

    render.addFilesToHtml(data);
    render.addRootToHtml(currentPath);

    currentSort = newSort;

    render.addTimerToHtml(new Date().getTime() - startTime.getTime())
}

//buttonCLickHandler обработчик события нажатия на кнопку
async function parentDirButtonClickHandler(): Promise<void> {
    const startTime = new Date();

    const lastSlashIndex = currentPath.lastIndexOf('\\')
    let newPath = currentPath.slice(0, lastSlashIndex)

    if (newPath[0] != "D")
        newPath = "D"
    const data = await sendJsonRequest(newPath, currentSort);

    render.addFilesToHtml(data);
    render.addRootToHtml(newPath);

    currentPath = newPath;

    render.addTimerToHtml(new Date().getTime() - startTime.getTime())
}

//
async function directoriesButtonClickHandler(): Promise<void>{
    const startTime = new Date();

    const data = await sendJsonRequest(currentPath, currentSort);

    render.addFilesToHtml(data);
    render.addRootToHtml(currentPath);
    render.addTimerToHtml(new Date().getTime() - startTime.getTime())
}

//
async function statisticButtonClickHandler():Promise<void>{

}

//события при загрузке окна
window.onload = async function(): Promise<void> {
    const startTime = new Date();
    const data = await sendJsonRequest(currentPath, currentSort);

    addEventHandlerById("sortButt", "click", buttonClickHandler)
    addEventHandlerById("folder-list", "click", fileListClickHandler)
    addEventHandlerById("currentDir", "click", rootClickHandler)
    addEventHandlerById("parentButt", "click", parentDirButtonClickHandler)
    addEventHandlerById("directoriesButt", "click", directoriesButtonClickHandler)


    render.addFilesToHtml(data);
    render.addRootToHtml(currentPath);

    render.addTimerToHtml(new Date().getTime() - startTime.getTime())
}
