function askAI() {
  const input = document.getElementById("userInput").value;
  const error = document.getElementById("error");
  const answer = document.getElementById("answer");

  error.textContent = "";
  answer.textContent = "";

  // AI BUG FIX
  if (input.trim() === "") {
    error.textContent = "❌ Сұрақ бос болмауы керек";
    return;
  }

  answer.textContent = "🤖 AI жауабы: " + input;
}
