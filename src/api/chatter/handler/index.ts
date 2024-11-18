import { APIGatewayEvent, APIGatewayProxyResult } from "aws-lambda";
import AWS from "aws-sdk";

const apiGateway = new AWS.ApiGatewayManagementApi({
  endpoint: process.env.WEBSOCKET_API_ENDPOINT, // Set in environment variables
});

export const handler = async (
  event: APIGatewayEvent
): Promise<APIGatewayProxyResult> => {
  const connectionId = event.requestContext.connectionId;
  const body = JSON.parse(event.body || "{}");
  const { message } = body;

  try {
    // Echo the message back to the sender
    await apiGateway
      .postToConnection({
        ConnectionId: connectionId as string,
        Data: `You said: ${message}`,
      })
      .promise();

    console.log(`Message sent to ${connectionId}: ${message}`);
    return { statusCode: 200, body: "Message sent" };
  } catch (err) {
    console.error(`Failed to send message: ${(err as Error).message}`);
    return { statusCode: 500, body: "Failed to send message" };
  }
};
