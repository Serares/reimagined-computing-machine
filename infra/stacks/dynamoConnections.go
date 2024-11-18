package stacks

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsdynamodb"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type DynamoDbStackProps struct {
	awscdk.StackProps
	Env string
}

func NewDynamoDbStack(scope constructs.Construct, id string, props *DynamoDbStackProps) awscdk.Stack {
	stack := awscdk.NewStack(scope, &id, &props.StackProps)

	// Define the DynamoDB table
	connectionsTable := awsdynamodb.NewTable(stack, jsii.String("ConnectionsTable"), &awsdynamodb.TableProps{
		TableName:           jsii.String("WebSocketConnections"),
		PartitionKey:        &awsdynamodb.Attribute{Name: jsii.String("connectionId"), Type: awsdynamodb.AttributeType_STRING},
		BillingMode:         awsdynamodb.BillingMode_PAY_PER_REQUEST,
		RemovalPolicy:       awscdk.RemovalPolicy_DESTROY, // Use RETAIN for production
		TimeToLiveAttribute: jsii.String("ttl"),           // Optional TTL attribute
	})

	// Output the table name for other stacks
	awscdk.NewCfnOutput(stack, jsii.Sprintf("DynamoConnectionsTable-%s", props.Env), &awscdk.CfnOutputProps{
		Value:      connectionsTable.TableName(),
		ExportName: jsii.Sprintf("DynamoConnectionsTable-%s", props.Env),
	})

	return stack
}
