package main

import (
	"infra/stacks"
	"os"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/joho/godotenv"

	// "github.com/aws/aws-cdk-go/awscdk/v2/awssqs"

	"github.com/aws/jsii-runtime-go"
)

func main() {
	defer jsii.Close()
	godotenv.Load(".env.dev")

	app := awscdk.NewApp(nil)

	godotenv.Load(".env.dev")
	env := os.Getenv("ENV")

	stacks.CognitoStack(app, *jsii.Sprintf("CognitoStack-%s", env), &stacks.CognitoStackProps{
		FACEBOOK_CLIENT_SECRET: os.Getenv("FACEBOOK_CLIENT_SECRET"),
		GOOGLE_CLIENT_SECRET:   os.Getenv("GOOGLE_CLIENT_SECRET"),
	})

	stacks.ApiGw(app, *jsii.Sprintf("APIGW-Stack-%s", env), &stacks.ApiGwStackProps{})

	app.Synth(nil)
}

// env determines the AWS environment (account+region) in which our stack is to
// be deployed. For more information see: https://docs.aws.amazon.com/cdk/latest/guide/environments.html
func env() *awscdk.Environment {
	// If unspecified, this stack will be "environment-agnostic".
	// Account/Region-dependent features and context lookups will not work, but a
	// single synthesized template can be deployed anywhere.
	//---------------------------------------------------------------------------
	return nil

	// Uncomment if you know exactly what account and region you want to deploy
	// the stack to. This is the recommendation for production stacks.
	//---------------------------------------------------------------------------
	// return &awscdk.Environment{
	//  Account: jsii.String("123456789012"),
	//  Region:  jsii.String("us-east-1"),
	// }

	// Uncomment to specialize this stack for the AWS Account and Region that are
	// implied by the current CLI configuration. This is recommended for dev
	// stacks.
	//---------------------------------------------------------------------------
	// return &awscdk.Environment{
	//  Account: jsii.String(os.Getenv("CDK_DEFAULT_ACCOUNT")),
	//  Region:  jsii.String(os.Getenv("CDK_DEFAULT_REGION")),
	// }
}
