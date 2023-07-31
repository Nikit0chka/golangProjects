export function addEventHandlerById(elementId, event, handler) {
    const element = document.getElementById(elementId);
    element.addEventListener(event, handler);
}
