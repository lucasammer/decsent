const express = require("express");
const app = express();
const rateLimit = require("express-rate-limit");
const fs = require("fs");
const dataloc = __dirname + "/../findingsTEST.csv";

const limiter = rateLimit({
  windowMs: 15 * 60 * 1000,
  max: 100,
  standardHeaders: "draft-7",
  legacyHeaders: false,
});

let data = fs.readFileSync(dataloc, "utf8");
data = data.toString("utf8");
data = data.split(/\r?\n/);
for (let i = 0; i < data.length; i++) {
  const element = data[i];
  let split = element.split(",");
  data[i] = { address: split[1], description: split[0] };
}

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

app.get("/search/raw", (req, res) => {
  let reqTime = Date.now();
  if (req.query.q == null) {
    res.redirect("/");
    return;
  }
  if (typeof req.query.q != "string") {
    res.sendStatus(400);
    return;
  }
  let found = data.filter((item) => {
    return item.description.includes(req.query.q);
  });
  res.setHeader("Content-Type", "application/json");
  res.json({ results: found, time: Date.now() - reqTime });
});

app.get("/search", (req, res) => {
  if (req.query.q == null) {
    res.redirect("/");
    return;
  }
  if (typeof req.query.q != "string") {
    res.sendStatus(400);
    return;
  }
  res.sendFile(__dirname + "/html/results.html");
});

app.get("/", (_req, res) => {
  res.sendFile(__dirname + "/html/home.html");
});

app.listen(3000, () => {
  console.log("Running on port 3000! :3");
});
