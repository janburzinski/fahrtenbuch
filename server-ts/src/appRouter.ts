import { Router, Request, Response } from "express";
import { authMiddleware } from "./middleware/authMiddleware";

// setup the app router
// func return the router to use in the index.ts file
export const appRouter = (): Router => {
  const router = Router();

  router.get("/test", authMiddleware, (req: Request, res: Response) => {
    res.send({ message: "miau" });
  });

  return router;
};
