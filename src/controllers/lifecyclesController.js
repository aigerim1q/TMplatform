const pool = require("../db");
const fs = require("fs");
const { callAI } = require("../services/aiService");

// Upload a lifecycle file
exports.uploadLifecycle = async (req, res) => {
  try {
    const userId = req.body.user_id;         // assuming user_id sent in body
    const filePath = req.file.path;

    const result = await pool.query(
      "INSERT INTO lifecycles (user_id, file_url) VALUES ($1, $2) RETURNING *",
      [userId, filePath]
    );

    res.json(result.rows[0]);
  } catch (err) {
    console.error(err);
    res.status(500).send("Server error");
  }
};

// Get all lifecycle files for a specific user
exports.getMyLifecycles = async (req, res) => {
  try {
    const { user_id } = req.params;

    const result = await pool.query(
      "SELECT * FROM lifecycles WHERE user_id = $1",
      [user_id]
    );

    res.json(result.rows);
  } catch (err) {
    console.error(err);
    res.status(500).send("Server error");
  }
};

// Parse uploaded file with AI (Step 6: enforce JSON output)
exports.parseLifecycleWithAI = async (req, res) => {
  try {
    const filePath = req.file.path;

    // Read file content (adjust for PDF/image/PPT as needed)
    const fileContent = fs.readFileSync(filePath, "utf-8");

    const prompt = `
You are given a project lifecycle document.
Extract project stages, tasks, and subtasks.
Return ONLY valid JSON in this format:
{
  "stages": [
    {
      "name": "...",
      "tasks": [
        {
          "name": "...",
          "subtasks": ["...", "..."]
        }
      ]
    }
  ]
}
`;

    const aiResult = await callAI(prompt, fileContent);

    res.json(aiResult); // returns structured JSON from AI
  } catch (err) {
    console.error(err);
    res.status(500).send("AI parsing failed");
  }
};
