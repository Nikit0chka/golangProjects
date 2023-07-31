export function addEventHandlerById(elementId: string, event:string, handler:() => void){
    const element = document.getElementById(elementId)
    element.addEventListener(event, handler)
}
