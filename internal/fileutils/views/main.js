const msg = document.getElementById("msg");
const sendMsg = document.getElementById("send-msg");
const output = document.getElementById("output");

const websocket = new WebSocket("ws://localhost:3000/echo");
websocket.onopen = function () {
  const message = "Conexion establecida";
  console.log(message);
  // websocket.send(message);
};
websocket.onmessage = function (event) {
  const message = event.data;
  if (message) {
    output.innerText += `${message}\n`;
  }
};

sendMsg.addEventListener("click", function () {
  if (websocket.readyState === websocket.CLOSED)
    return window.alert("Connection is clossed");
  const message = msg.value;
  if (message) {
    websocket.send(message);
    msg.value = "";
    return;
  }
});