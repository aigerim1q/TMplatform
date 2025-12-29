// services/aiService.js
const axios = require("axios"); // for HTTP requests to AI API

async function callAI(prompt, fileContent) {
  try {
    // Example for Qwen API, replace with real endpoint and auth
    const response = await axios.post(
      "https://api.qwen.com/v1/coder", // example endpoint
      {
        prompt: prompt,
        content: fileContent
      },
      {
        headers: {
          "Authorization": `Bearer YOUR_QWEN_API_KEY`,
          "Content-Type": "application/json"
        }
      }
    );

    // AI returns structured JSON
    return response.data;
  } catch (err) {
    console.error("AI call failed:", err);
    throw err;
  }
}

module.exports = { callAI };
