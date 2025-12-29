const express = require("express");
const app = express();
app.use(express.json());

// Lifecycle routes
const lifecycleRoutes = require("./routes/lifecycles");
app.use("/lifecycles", lifecycleRoutes);

// Projects routes (new)
const projectsRoutes = require("./routes/projects");
app.use("/projects", projectsRoutes);

// Start server
app.listen(3000, () => console.log("Server running on port 3000"));
