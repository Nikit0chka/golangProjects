//getCurrentPath возвращает путь к текущей директории
function getCurrentPath(){
  const roots = document.querySelectorAll(".root")
  let maxRootLen = ""

  for (let i = 0; i < roots.length; i++)
    if (roots[i].getAttribute("path") > maxRootLen)
      maxRootLen = roots[i].getAttribute("path")

  return maxRootLen
}

// sendJsonAndUpdateHtml отправляет json и по результатам ответа изменяет html
async function sendJsonAndUpdateHtml(path){
  const start= new Date().getTime()
  let jsonRequest = new JsonRequest(path, currentSort)

  await sendJsonRequest(jsonRequest).then(data => {addRootToHtml(data[0].path); addFilesToHtml(data) }).catch(error => console.error(error))
  const end = new Date().getTime()
  addTimerToHtml(end - start)
}

//события при загрузке окна
window.onload = function (){
  addEventOnButton()
  sendJsonAndUpdateHtml("/")
}