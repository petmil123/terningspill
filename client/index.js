const ws = new WebSocket("ws://localhost:8080/ws");

const testMessage = {
  action: "test",
  data: "This is a test message"
}

ws.onopen = () => {
  console.log("WebSocket connection established");
  ws.send(JSON.stringify(testMessage));
};

ws.onmessage = (event) => {
  console.log("Message from server:", event.data);
};
ws.onclose = () => {
    console.log("WebSocket connection closed");
};

ws.onerror = (error) => {
  console.error("WebSocket error:", error);
};
