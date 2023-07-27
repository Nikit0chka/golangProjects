(()=>{"use strict";var e=function(e,t,n,i){return new(n||(n=Promise))((function(c,o){function a(e){try{r(i.next(e))}catch(e){o(e)}}function d(e){try{r(i.throw(e))}catch(e){o(e)}}function r(e){var t;e.done?c(e.value):(t=e.value,t instanceof n?t:new n((function(e){e(t)}))).then(a,d)}r((i=i.apply(e,t||[])).next())}))};const t="./";class n{constructor(e,t,n,i,c){this.name=e,this.fileOrder=t,this.path=n,this.size=i,this.type=c}}class i{constructor(e,t){this.path=e,this.sortType=t}}function c(c,o){return e(this,void 0,void 0,(function*(){const e=new i(c,o),a=yield fetch(t,{method:"POST",headers:{"Content-Type":"application/json"},body:JSON.stringify(e)});return(yield a.json()).map((e=>new n(e.name,e.fileOrder,e.path,e.size,e.type)))}))}var o;function a(e){const t=document.querySelector(".folder-list");t.innerHTML="";for(let n=0;n<e.length;n++)for(let i=0;i<e.length;i++)if(e[i].fileOrder===n){const i=e[n],c=document.createElement("li"),a=document.createElement("div"),d=document.createElement("img");i.type===o.dir||void 0===i.type?d.src="/static/dirImage.png":d.src="/static/fileImg.jpg";const r=document.createElement("span");d.className="file-icon",c.className="file-space",c.setAttribute("name",i.name),c.setAttribute("path",i.path),a.appendChild(d),a.appendChild(document.createTextNode(i.name)),c.appendChild(a),c.appendChild(r),r.className="folder-size",r.appendChild(document.createTextNode(i.size+" mb")),t.appendChild(c);break}}function d(e){const t=e.split("/"),n=document.getElementById("currentDir");n.innerHTML="";const i=document.createElement("a");let c="/";i.setAttribute("path",c),i.className="root",i.appendChild(document.createTextNode("start:")),n.appendChild(i),c="";for(let e=1;e<t.length;e++){const i=document.createElement("a");c+="/"+t[e],i.setAttribute("path",c),i.className="root",i.appendChild(document.createTextNode(t[e])),i.appendChild(document.createTextNode("/")),n.appendChild(i)}}function r(e){const t=document.getElementById("timer");t.innerHTML="";const n=document.createElement("span");t.appendChild(n),n.appendChild(document.createTextNode(e+"ms"))}!function(e){e.file="FILE",e.dir="DIR"}(o||(o={}));var s,l=function(e,t,n,i){return new(n||(n=Promise))((function(c,o){function a(e){try{r(i.next(e))}catch(e){o(e)}}function d(e){try{r(i.throw(e))}catch(e){o(e)}}function r(e){var t;e.done?c(e.value):(t=e.value,t instanceof n?t:new n((function(e){e(t)}))).then(a,d)}r((i=i.apply(e,t||[])).next())}))};!function(e){e.asc="ASC",e.desk="DESK"}(s||(s={}));let u="/";var m=s.asc;const p=function(){return l(this,void 0,void 0,(function*(){const e=new Date,t=event.target.closest("li").getAttribute("path");a(yield c(t,m)),d(t),u=t,r((new Date).getTime()-e.getTime())}))};function h(){return l(this,void 0,void 0,(function*(){const e=new Date,t=event.target.getAttribute("path");a(yield c(t,m)),d(t),u=t,r((new Date).getTime()-e.getTime())}))}function f(){return l(this,void 0,void 0,(function*(){const e=new Date,t=m===s.asc?s.desk:s.asc;a(yield c(u,t)),d(u),m=t,r((new Date).getTime()-e.getTime())}))}function g(){return l(this,void 0,void 0,(function*(){const e=new Date,t=u.lastIndexOf("/");let n=u.slice(0,t);"/"!=n[0]&&(n="/");const i=yield c(n,m);console.log(n),a(i),d(n),u=n,r((new Date).getTime()-e.getTime())}))}window.onload=function(){return l(this,void 0,void 0,(function*(){const e=new Date,t=yield c("/",m);document.getElementById("sortButt").addEventListener("click",f),document.getElementById("folder-list").addEventListener("click",p),document.getElementById("currentDir").addEventListener("click",h),document.getElementById("parentButt").addEventListener("click",g),a(t),d(u),r((new Date).getTime()-e.getTime())}))}})();