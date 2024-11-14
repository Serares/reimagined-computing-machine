import express from "express";
import { lambdaToExpress } from "./adapter";
import dotenv from "dotenv";
import path from "path";

dotenv.config({ path: path.resolve(__dirname, "./.env") });

const app = express();
const PORT = process.env.PORT || 3000;

// Use EJS as the templating engine
app.set("view engine", "ejs");
app.set("views", __dirname + "/views");

// Route all requests to the Lambda handler through the adapter
app.get("*", lambdaToExpress);

app.listen(PORT, () => {
  console.log(`Server running at http://localhost:${PORT}`);
});
