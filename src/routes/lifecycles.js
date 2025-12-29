const express = require("express");
const router = express.Router();
const upload = require("../services/uploadService");
const { 
  uploadLifecycle, 
  getMyLifecycles, 
  parseLifecycleWithAI 
} = require("../controllers/lifecyclesController");

// Upload a file
router.post("/upload", upload.single("file"), uploadLifecycle);

// List all uploaded files for a user
router.get("/:user_id", getMyLifecycles);

// Parse an uploaded file with AI
router.post("/parse", upload.single("file"), parseLifecycleWithAI);

module.exports = router;
