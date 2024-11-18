package stacks

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/aws-cdk-go/awscdk/v2/awss3assets"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

// lambdaRole := utils.CreateLambdaBasicRole(stack, "lambdaRole", env)
// // Define the Lambda function
// awslambda.NewFunction(stack, aws.String("a1-pages"), &awslambda.FunctionProps{
// 	Runtime: awslambda.Runtime_NODEJS_LATEST(),
// 	Handler: aws.String("pages/handler/index.handler"),            // Adjust to your compiled JS handler path
// 	Code:    awslambda.Code_FromAsset(aws.String("../dist"), nil), // Point to dist folder
// 	Role:    lambdaRole,
// })

type FeLambdaStackProps struct {
	awscdk.StackProps
	Env string
}

// ❗This is just a template
func FELambdaStack(scope constructs.Construct, id string, props *FeLambdaStackProps) awscdk.Stack {
	stack := awscdk.NewStack(scope, &id, &props.StackProps)

	// Import Cognito Hosted UI URL
	cognitoHostedUIUrl := awscdk.Fn_ImportValue(jsii.String("CognitoHostedUIUrl"))

	// Import CloudFront URL
	cloudFrontUrl := awscdk.Fn_ImportValue(jsii.String("CloudFrontUrl"))

	// Lambda function code
	function := awslambda.NewFunction(stack, jsii.String("RedirectLambda"), &awslambda.FunctionProps{
		Runtime: awslambda.Runtime_NODEJS_14_X(),
		Handler: jsii.String("pages/handler/index.handler"),
		Code:    awslambda.AssetCode_FromAsset(jsii.String("../dist"), &awss3assets.AssetOptions{}),
		Environment: &map[string]*string{
			"COGNITO_HOSTED_UI_URL": cognitoHostedUIUrl,
			"CLOUDFRONT_URL":        cloudFrontUrl,
		},
	})

	lambdaURL := function.AddFunctionUrl(&awslambda.FunctionUrlOptions{
		AuthType: awslambda.FunctionUrlAuthType_NONE,
	})
	awscdk.NewCfnOutput(stack, jsii.String("FrontEndLambda"), &awscdk.CfnOutputProps{
		ExportName: jsii.Sprintf("FrontEndLambda-%s", props.Env),
		Value:      lambdaURL.Url(),
	})

	return stack
}
