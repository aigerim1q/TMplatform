const { Pool } = require("pg");

const pool = new Pool({
  user: "postgres",
  host: "localhost",
  database: "silence_ai",
  password: "AikonNokia1987@",
  port: 5432,
});

module.exports = pool;
