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
    state = "started";
    console.log("Game started");
    document.getElementsByClassName("lobby")[0].classList.add("hidden");
    document.getElementsByClassName("game-ui")[0].classList.remove("hidden");
    document.getElementById("status").textContent = "Game Started!";
  } else if (state === "started" && message.action === "gameState") {
    console.log("Game state updated:", message.data);
    // Update the game UI based on the new game state
    updateBoard(message.data);
  }
};
ws.onclose = () => {
    console.log("WebSocket connection closed");
};

ws.onerror = (error) => {
  console.error("WebSocket error:", error);
};

function updateBoard(gameState) {
  // Extract the different variables based on the gameState string
  const [playerCoveredRaw, opponentCoveredRaw, isYourTurnRaw, turnPhase, chosenField, diceRoll] = gameState.split(",");
  const playerCovered = playerCoveredRaw.slice(1, -1).split(" ").map(Number);
  const opponentCovered = opponentCoveredRaw.slice(1, -1).split(" ").map(Number);
  const isYourTurn = isYourTurnRaw === "true";
  console.log("Updating board with:", {playerCovered, opponentCovered, isYourTurn, turnPhase, chosenField, diceRoll });
  // Cover the fields based on covered

  // How do we do the covered?
  // If its not your turn, disable all controls.
  if (!isYourTurn) {
    disableButtons();
    disableDie();
    disableBoard();
  } else {
    switch (turnPhase) {
      case "toThrow":
        enableDie();
        disableBoard();
        disableButtons();
        break;
      case "toCover":
        disableDie();
        enableBoard();
        disableButtons();
        setDie(diceRoll);
        break;
      case "toCall":
        disableDie();
        disableBoard();
        enableButtons();
        setCall(chosenField);
        break;
  }
}
}

function disableDie() {
  const die = document.getElementById("ownDice");
  die.removeEventListener("click", throwDie);
  setDie(die, 0);
}
function setDie(die, value) {
  if (value < 1 || value > 6 ) {
    die.src = "dices/none.png";
    die.alt = "Dice showing none";
  } else {
    die.src = `dices/${value}.png`;
    die.alt = `Dice showing ${value}`;
  }
}

function enableDie() {
  const die = document.getElementById("ownDice");
  die.addEventListener("click", throwDie);
  die.classList.add("clickable")
  setDie(die, Math.ceil(Math.random() * 6));
}

function disableBoard() {
  document.querySelectorAll(".ownRow .field").forEach(field => field.classList.add("disabled"));
}

function enableBoard() {
  document.querySelectorAll(".ownRow .field").forEach(field => field.classList.remove("disabled"));
}

function enableButtons() {
  document.querySelectorAll("button").forEach(btn => btn.disabled = false);
}
function disableButtons() {
  document.querySelectorAll("button").forEach(btn => btn.disabled = true);
}

function setCall(number) {
  document.getElementById("opponentCall").textContent = number != 0 ? number : "";
}

function throwDie() {
  console.log("Die thrown");
  ws.send(JSON.stringify({ action: "throwDie", data: "" }));
}