package stacks

import (
	"infra/utils"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsapigatewayv2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type ApiGwStackProps struct {
	awscdk.StackProps
	Env string
}

func ApiGw(scope constructs.Construct, id string, props *ApiGwStackProps) {
	stack := awscdk.NewStack(scope, &id, &props.StackProps)
	// Define Lambda functions
	// cognitoPoolId := awscdk.Fn_ImportValue(jsii.Sprintf("UserPoolId"))
	userPoolClientId := awscdk.Fn_ImportValue(jsii.Sprintf("UserPoolClientId"))

	dynamoTableName := awscdk.Fn_ImportValue(jsii.Sprintf("DynamoConnectionsTable-%s", props.Env))

	lambdaRole := utils.CreateLambdaBasicRole(stack, "lambdaRole", props.Env)

	connectHandler := awslambda.NewFunction(stack, jsii.String("ConnectHandler"), &awslambda.FunctionProps{
		Runtime: awslambda.Runtime_NODEJS_14_X(),
		Handler: jsii.String("api/connect/handler/index.handler"),
		Code:    awslambda.Code_FromAsset(jsii.String("../dist"), nil),
		Role:    lambdaRole,
		Environment: &map[string]*string{
			"CONNECTIONS_TABLE": dynamoTableName,
		},
	})

	disconnectHandler := awslambda.NewFunction(stack, jsii.String("DisconnectHandler"), &awslambda.FunctionProps{
		Runtime: awslambda.Runtime_NODEJS_14_X(),
		Handler: jsii.String("api/diconnect/handler/index.handler"),
		Code:    awslambda.Code_FromAsset(jsii.String("../dist"), nil),
		Role:    lambdaRole,
		Environment: &map[string]*string{
			"CONNECTIONS_TABLE": dynamoTableName,
		},
	})

	// Create the WebSocket API
	webSocketApi := awsapigatewayv2.NewCfnApi(stack, jsii.String("WebSocketApi"), &awsapigatewayv2.CfnApiProps{
		Name:                     jsii.String("MyWebSocketApi"),
		ProtocolType:             jsii.String("WEBSOCKET"),
		RouteSelectionExpression: jsii.String("$request.body.action"),
	})

	// Capture the WebSocket API endpoint
	webSocketApiEndpoint := awscdk.Fn_Sub(
		jsii.String("wss://${ApiId}.execute-api.${Region}.amazonaws.com/$default"),
		&map[string]*string{
			"ApiId":  webSocketApi.Ref(),
			"Region": stack.Region(),
		},
	)

	messageHandler := awslambda.NewFunction(stack, jsii.String("MessageHandler"), &awslambda.FunctionProps{
		Runtime: awslambda.Runtime_NODEJS_14_X(),
		Handler: jsii.String("api/chatter/handler/index.handler"),
		Code:    awslambda.Code_FromAsset(jsii.String("../dist"), nil),
		Role:    lambdaRole,
		Environment: &map[string]*string{
			"WEBSOCKET_API_ENDPOINT": webSocketApiEndpoint,
		},
	})

	// Create the Cognito Authorizer
	authorizer := awsapigatewayv2.NewCfnAuthorizer(stack, jsii.String("CognitoAuthorizer"), &awsapigatewayv2.CfnAuthorizerProps{
		ApiId:          webSocketApi.Node().Id(),
		AuthorizerType: jsii.String("JWT"),
		IdentitySource: &[]*string{
			jsii.String("$request.header.Authorization"),
		},
		Name: jsii.String("CognitoAuthorizer"),
		JwtConfiguration: &awsapigatewayv2.CfnAuthorizer_JWTConfigurationProperty{
			Audience: &[]*string{jsii.String("YOUR_COGNITO_APP_CLIENT_ID")},
			Issuer:   jsii.Sprintf("https://cognito-idp.YOUR_REGION.amazonaws.com/%s", *userPoolClientId),
		},
	})

	connectIntegration := awsapigatewayv2.NewCfnIntegration(stack, jsii.String("ConnectIntegration"), &awsapigatewayv2.CfnIntegrationProps{
		ApiId:             webSocketApi.Ref(),
		IntegrationType:   jsii.String("AWS_PROXY"),
		IntegrationUri:    connectHandler.FunctionArn(),
		IntegrationMethod: jsii.String("POST"),
	})

	awsapigatewayv2.NewCfnRoute(stack, jsii.String("ConnectRoute"), &awsapigatewayv2.CfnRouteProps{
		ApiId:        webSocketApi.Ref(),
		RouteKey:     jsii.String("$connect"),
		Target:       jsii.String("integrations/" + *connectIntegration.Ref()),
		AuthorizerId: authorizer.ApiId(),
	})

	// Create the $disconnect route and attach the Lambda integration
	disconnectIntegration := awsapigatewayv2.NewCfnIntegration(stack, jsii.String("DisconnectIntegration"), &awsapigatewayv2.CfnIntegrationProps{
		ApiId:             webSocketApi.Ref(),
		IntegrationType:   jsii.String("AWS_PROXY"),
		IntegrationUri:    disconnectHandler.FunctionArn(),
		IntegrationMethod: jsii.String("POST"),
	})

	awsapigatewayv2.NewCfnRoute(stack, jsii.String("DisconnectRoute"), &awsapigatewayv2.CfnRouteProps{
		ApiId:    webSocketApi.Ref(),
		RouteKey: jsii.String("$disconnect"),
		Target:   jsii.String("integrations/" + *disconnectIntegration.Ref()),
	})

	defaultIntegration := awsapigatewayv2.NewCfnIntegration(stack, jsii.String("DefaultIntegration"), &awsapigatewayv2.CfnIntegrationProps{
		ApiId:             webSocketApi.Ref(),
		IntegrationType:   jsii.String("AWS_PROXY"),
		IntegrationUri:    messageHandler.FunctionArn(),
		IntegrationMethod: jsii.String("POST"),
	})

	awsapigatewayv2.NewCfnRoute(stack, jsii.String("DefaultRoute"), &awsapigatewayv2.CfnRouteProps{
		ApiId:    webSocketApi.Ref(),
		RouteKey: jsii.String("$default"),
		Target:   jsii.String("integrations/" + *defaultIntegration.Ref()),
	})

	// Create a deployment
	deployment := awsapigatewayv2.NewCfnDeployment(stack, jsii.String("WebSocketDeployment"), &awsapigatewayv2.CfnDeploymentProps{
		ApiId: webSocketApi.Ref(),
	})

	// Create a stage
	awsapigatewayv2.NewCfnStage(stack, jsii.String("WebSocketStage"), &awsapigatewayv2.CfnStageProps{
		ApiId:        webSocketApi.Ref(),
		StageName:    jsii.String("$default"),
		DeploymentId: deployment.Ref(),
		AutoDeploy:   jsii.Bool(true),
	})

}
