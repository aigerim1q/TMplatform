// controllers/projectsController.js
const pool = require("../db");
const fs = require("fs");
const { callAI } = require("../services/aiService");

// Helper functions for DB inserts
async function createProject(name, lifecycleId) {
  const result = await pool.query(
    "INSERT INTO projects (name, lifecycle_id) VALUES ($1, $2) RETURNING *",
    [name, lifecycleId]
  );
  return result.rows[0];
}

async function createStage(name, projectId) {
  const result = await pool.query(
    "INSERT INTO stages (name, project_id) VALUES ($1, $2) RETURNING *",
    [name, projectId]
  );
  return result.rows[0];
}

async function createTask(name, stageId) {
  const result = await pool.query(
    "INSERT INTO tasks (name, stage_id) VALUES ($1, $2) RETURNING *",
    [name, stageId]
  );
  return result.rows[0];
}

async function createSubtask(name, taskId) {
  await pool.query(
    "INSERT INTO subtasks (name, task_id) VALUES ($1, $2)",
    [name, taskId]
  );
}

// Main bridge API
exports.createProjectFromLifecycle = async (req, res) => {
  try {
    const { lifecycle_id } = req.params;

    // 1️⃣ Load lifecycle file path from DB
    const lifecycleResult = await pool.query(
      "SELECT * FROM lifecycles WHERE id = $1",
      [lifecycle_id]
    );

    if (lifecycleResult.rows.length === 0) {
      return res.status(404).send("Lifecycle not found");
    }

    const filePath = lifecycleResult.rows[0].file_url;

    // 2️⃣ Read file content
    const fileContent = fs.readFileSync(filePath, "utf-8");

    // 3️⃣ Send to AI
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

    // 4️⃣ Nested loops: create project, stages, tasks, subtasks
    const project = await createProject("New Project", lifecycle_id);

    for (const stage of aiResult.stages) {
      const stageRecord = await createStage(stage.name, project.id);

      for (const task of stage.tasks) {
        const taskRecord = await createTask(task.name, stageRecord.id);

        for (const subtask of task.subtasks) {
          await createSubtask(subtask, taskRecord.id);
        }
      }
    }

    res.json({ message: "Project created successfully", project });
  } catch (err) {
    console.error(err);
    res.status(500).send("Project creation failed");
  }
};
