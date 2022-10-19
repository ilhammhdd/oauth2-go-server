const xhttp = new XMLHttpRequest();
xhttp.onload = function () {
  window.self.document.write(this.response);
};
xhttp.open("GET", "192.168.8.8:6464/register?client_id=48a51453-032b-4fe9-b0aa-5fb57e0058e5");
xhttp.setRequestHeader("Authorization", "Bearer EEhZRk2co1gwLWNJc-zkbIFav27Q-rC_AMPziZXFMC4");
xhttp.send();

fetch("http://192.168.8.8:6464/register", {
  method: "GET",
  headers: {
    "Authorization": "Bearer EEhZRk2co1gwLWNJc-zkbIFav27Q-rC_AMPziZXFMC4"
  }
}).then(html => {
  const win = window.open(``, `_blank`);
  win.document.body.innerHTML = html;
  win.focus();
});