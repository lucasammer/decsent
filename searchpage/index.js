const express = require("express");
const app = express();
const rateLimit = require("express-rate-limit");

const limiter = rateLimit({
  windowMs: 15 * 60 * 1000,
  max: 100,
  standardHeaders: "draft-7",
  legacyHeaders: false,
});

app.use(express.static(__dirname + "/static"));

app.get("/", (_req, res) => {
  res.sendFile(__dirname + "/html/home.html");
});

app.listen(3000, () => {
  console.log("Running on port 3000! :3");
});
