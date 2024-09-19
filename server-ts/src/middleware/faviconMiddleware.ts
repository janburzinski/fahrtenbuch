import { Request, Response, NextFunction } from "express";

export const faviconMiddleware = (
  req: Request,
  res: Response,
  next: NextFunction
) => {
  if (req.url === "/favicon.ico") {
    res.status(204).end();
  } else {
    next();
  }
};
