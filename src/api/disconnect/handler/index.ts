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
    // Remove the connection ID from DynamoDB
    await dynamoDB
      .delete({
        TableName: CONNECTIONS_TABLE,
        Key: {
          connectionId: connectionId,
        },
      })
      .promise();

    console.log(`Client disconnected: ${connectionId}`);
    return { statusCode: 200, body: "Disconnected" };
  } catch (err) {
    console.error(`Failed to disconnect: ${(err as Error).message}`);
    return { statusCode: 500, body: "Failed to disconnect" };
  }
};
