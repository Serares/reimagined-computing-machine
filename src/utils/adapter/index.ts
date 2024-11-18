import { Request, Response } from "express";

export type IHandlerResponse = {
  statusCode: number;
  body: string | Record<string, string>;
  headers?: Record<string, string> | undefined;
};

export type IHandler = (event: any) => Promise<IHandlerResponse>;

// Adapter to convert Lambda response to Express response
export const lambdaToExpress = async (
  req: Request,
  res: Response,
  handler: IHandler
) => {
  // Create a Lambda-like event object from the Express request
  const event = {
    path: req.path,
    httpMethod: req.method,
    headers: req.headers,
    queryStringParameters: req.query,
    body: req.body ? JSON.stringify(req.body) : null,
  };

  // Call the Lambda handler and get the result
  const result = await handler(event);

  // Set status code and headers from the Lambda response
  res.status(result.statusCode);
  if (result.headers) {
    for (const [key, value] of Object.entries(result.headers)) {
      res.setHeader(key, value as string);
    }
  }

  // Send the response body
  res.send(result.body);
};
