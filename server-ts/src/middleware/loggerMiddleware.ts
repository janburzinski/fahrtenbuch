import { Request, Response, NextFunction } from "express";

const formatDate = (date: Date): string => {
  return new Intl.DateTimeFormat("en-US", {
    timeZone: "Europe/Berlin",
    year: "numeric",
    month: "2-digit",
    day: "2-digit",
    hour: "2-digit",
    minute: "2-digit",
    second: "2-digit",
    hour12: false,
  }).format(date);
};

export const loggerMiddleware = (
  req: Request,
  res: Response,
  next: NextFunction
) => {
  const start = Date.now();

  res.on("finish", () => {
    const duration = Date.now() - start;
    const { method, url } = req;
    const { statusCode } = res;
    const timestamp = formatDate(new Date());
    console.log(
      `${timestamp} - ${method} | ${url} (${statusCode}) | ${duration}ms`
    );
  });

  next();
};
