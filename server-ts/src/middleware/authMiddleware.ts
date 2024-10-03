import supabase from "src/utils/supachaseClient";
import { Request, Response, NextFunction } from "express";

export const authMiddleware = async (
  req: Request,
  res: Response,
  next: NextFunction
) => {
  // skip auth middleware if not in prod
  if (process.env.NODE_ENV !== "prod") {
    next();
    return;
  }

  //proceed if in prod
  const token = req.headers.authorization?.split(" ")[1];

  if (!token) {
    res.status(401).json({ error: "Unauthorized" });
    return;
  }

  try {
    const {
      data: { user },
      error,
    } = await supabase.auth.getUser(token);
    if (error) {
      res.status(401).json({ error: "invalid auth token" });
      return;
    }

    if (!user) {
      res.status(401).json({ error: "user not found" });
      return;
    }

    (req as any).user = user;
    next();
  } catch (error) {
    console.error("error in auth middleware: ", error);
    res.status(401).json({ error: "internal server error" });
  }
};
