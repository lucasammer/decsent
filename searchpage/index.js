const express = require("express");
const app = express();
const rateLimit = require("express-rate-limit");

const limiter = rateLimit({
  windowMs: 15 * 60 * 1000,
  max: 100,
  standardHeaders: "draft-7",
  legacyHeaders: false,
});

app.use(limiter);

app.use(express.static(__dirname + "/static"));

app.use((req, res, next) => {
  let rand = Math.floor(Math.random() * 10);
  if (rand == 0) {
    res.setHeader("X-Powered-By", "Hopes and dreams");
  } else if (rand == 1) {
    res.setHeader("X-Powered-By", "Pure faith");
  } else if (rand == 2) {
    res.setHeader("X-Powered-By", "Caffeine");
  } else {
    res.removeHeader("X-Powered-By");
  }

  next();
});

app.get("/", (_req, res) => {
  res.sendFile(__dirname + "/html/home.html");
});

app.listen(3000, () => {
  console.log("Running on port 3000! :3");
});
