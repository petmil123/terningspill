const ws = new WebSocket("ws://localhost:8080/ws");

let state = "toConnect";

ws.onopen = () => {
  console.log("WebSocket connection established");
};

ws.onmessage = (event) => {
  const message = JSON.parse(event.data);
  console.log("Message received:", message)
  if (state === "toConnect" && message.action === "connected") {
    state = "connected";
    console.log("Connected to the server");
    document.getElementById("status").textContent = "Waiting...";
  } else if (state === "connected" && message.action === "start") {
    state = "inGame";
    console.log("Game started");
    document.getElementsByClassName("lobby")[0].classList.add("hidden");
    document.getElementsByClassName("game-ui")[0].classList.remove("hidden");
    document.getElementById("status").textContent = "Game Started!";
  }
};
ws.onclose = () => {
    console.log("WebSocket connection closed");
};

ws.onerror = (error) => {
  console.error("WebSocket error:", error);
};
