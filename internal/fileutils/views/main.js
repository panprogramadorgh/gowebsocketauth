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

sendMsg.addEventListener("click", sendMessage);
msg.addEventListener("keydown", (event) => {
  if (event.key === "Enter") sendMessage()
})

function sendMessage() {
  if (websocket.readyState === websocket.CLOSED)
    return window.alert("Connection is clossed. Reload the page to re-connect");
  const message = msg.value;
  if (message) {
    // Clearing the screen remainds on the client
    if (message === "/clear") {
      output.innerHTML = ""
    } else {
      websocket.send(message);
    }
    msg.value = "";
  }
}