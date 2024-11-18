import { APIGatewayEvent, APIGatewayProxyResult } from "aws-lambda";
import AWS from "aws-sdk";

const dynamoDB = new AWS.DynamoDB.DocumentClient();
const CONNECTIONS_TABLE =
  process.env.CONNECTIONS_TABLE || "WebSocketConnections";

export const handler = async (
  event: APIGatewayEvent
): Promise<APIGatewayProxyResult> => {
  const connectionId = event.requestContext.connectionId;

  try {
    await dynamoDB
      .put({
        TableName: CONNECTIONS_TABLE,
        Item: {
          connectionId: connectionId,
          timestamp: Date.now(),
        },
      })
      .promise();

    console.log(`Client connected: ${connectionId}`);
    return { statusCode: 200, body: "Connected" };
  } catch (err) {
    console.error(`Failed to connect: ${(err as Error).message}`);
    return { statusCode: 500, body: "Failed to connect" };
  }
};
