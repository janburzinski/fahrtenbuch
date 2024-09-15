import express, { Request, Response } from "express";
import dotenv from "dotenv";
import path from "path";
import cors from "cors";
import { appRouter } from "./appRouter";

const main = () => {
  const app = express();

  // Load environment variables from .env file
  dotenv.config({ path: path.resolve(__dirname, "..", ".env") });

  // Verify that required environment variables are set
  const requiredEnvVars = [
    "SUPABASE_ANON_PUBLIC_KEY",
    "SUPABASE_SECRET_KEY",
    "SUPABASE_JWT_SECRET",
    "SUPABASE_PROJECT_URL",
    "SUPABASE_PROJECT_ID",
    "ALLOWED_ORIGIN",
  ];

  requiredEnvVars.forEach((varName) => {
    if (!process.env[varName]) {
      throw new Error(`Missing required environment variable: ${varName}`);
    }
  });

  const isProd = process.env.NODE_ENV === "production";

  const corsOptions = {
    origin: isProd ? process.env.ALLOWED_ORIGIN : "*",
    methods: ["GET", "POST", "PUT", "DELETE"],
    allowedHeaders: ["Content-Type", "Authorization"],
  };

  app.use(cors(corsOptions));
  app.use(express.json());

  // init app router
  app.use("/api/v1", appRouter);

  const port = process.env.SERVER_PORT || 8080;
  app.listen(port, () => {
    console.log(`Server is running on http://localhost:${port}`);
  });
};

main();
