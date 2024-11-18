import express from "express";
import { lambdaToExpress } from "../utils/adapter";
import dotenv from "dotenv";
import path from "path";
import { handler } from "./handler";

dotenv.config({ path: path.resolve(__dirname, "./.env") });

const app = express();
const PORT = process.env.PORT || 3000;

// Use EJS as the templating engine
app.set("view engine", "ejs");
app.set("views", __dirname + "/views");

// Route all requests to the Lambda handler through the adapter
app.get("*", (req, res) => lambdaToExpress(req, res, handler));

app.listen(PORT, () => {
  console.log(`Server running at http://localhost:${PORT}`);
});
