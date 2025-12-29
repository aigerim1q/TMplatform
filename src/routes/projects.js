// routes/projects.js
const express = require("express");
const router = express.Router();
const { createProjectFromLifecycle } = require("../controllers/projectsController");

router.post("/from-lifecycle/:lifecycle_id", createProjectFromLifecycle);

module.exports = router;
