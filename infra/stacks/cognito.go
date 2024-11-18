package stacks

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awscognito"
	"github.com/aws/aws-cdk-go/awscdk/v2/awssecretsmanager"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

const (
	USER_TIER_BASIC = "basic"
	USER_TIER_PRO   = "pro"
)

type CognitoStackProps struct {
	FACEBOOK_CLIENT_SECRET string
	GOOGLE_CLIENT_SECRET   string
}

func CognitoStack(scope constructs.Construct, id string, props *CognitoStackProps) awscdk.Stack {
	stack := awscdk.NewStack(scope, &id, nil)

	// Create the Cognito User Pool
	userPool := awscognito.NewUserPool(stack, jsii.String("UserPool"), &awscognito.UserPoolProps{
		UserPoolName: jsii.String("CustomOIDCUserPool"),
		StandardAttributes: &awscognito.StandardAttributes{
			Email: &awscognito.StandardAttribute{Required: jsii.Bool(true), Mutable: jsii.Bool(true)},
		},
		SignInAliases: &awscognito.SignInAliases{
			Email: jsii.Bool(true),
		},
		CustomAttributes: &map[string]awscognito.ICustomAttribute{
			"user_tier": awscognito.NewStringAttribute(&awscognito.StringAttributeProps{
				MinLen: jsii.Number(0),
				MaxLen: jsii.Number(100),
			}),
		},
		RemovalPolicy: awscdk.RemovalPolicy_DESTROY,
	})

	// Create a User Pool Client
	userPoolClient := awscognito.NewUserPoolClient(stack, jsii.String("UserPoolClient"), &awscognito.UserPoolClientProps{
		UserPool:       userPool,
		GenerateSecret: jsii.Bool(true),
		AuthFlows: &awscognito.AuthFlow{
			UserPassword: jsii.Bool(true),
		},
	})

	// Add Google Identity Provider
	googleSecret := awssecretsmanager.NewSecret(stack, jsii.String("GoogleSecret"), &awssecretsmanager.SecretProps{
		SecretName: jsii.String("google-oauth"),
		GenerateSecretString: &awssecretsmanager.SecretStringGenerator{
			SecretStringTemplate: jsii.Sprintf(`{"clientId": "%s"}`, props.GOOGLE_CLIENT_SECRET),
			GenerateStringKey:    jsii.String("clientSecret"),
		},
	})

	awscognito.NewUserPoolIdentityProviderGoogle(stack, jsii.String("GoogleProvider"), &awscognito.UserPoolIdentityProviderGoogleProps{
		ClientId:          googleSecret.SecretValueFromJson(jsii.String("clientId")).UnsafeUnwrap(),
		ClientSecretValue: googleSecret.SecretValueFromJson(jsii.String("clientSecret")),
		UserPool:          userPool,
		AttributeMapping: &awscognito.AttributeMapping{
			Email:    awscognito.ProviderAttribute_GOOGLE_EMAIL(),
			Fullname: awscognito.ProviderAttribute_GOOGLE_NAMES(),
		},
	})

	// Add Facebook Identity Provider
	facebookSecret := awssecretsmanager.NewSecret(stack, jsii.String("FacebookSecret"), &awssecretsmanager.SecretProps{
		SecretName: jsii.String("facebook-oauth"),
		GenerateSecretString: &awssecretsmanager.SecretStringGenerator{
			SecretStringTemplate: jsii.Sprintf(`{"clientId": "%s"}`, props.FACEBOOK_CLIENT_SECRET),
			GenerateStringKey:    jsii.String("clientSecret"),
		},
	})

	awscognito.NewUserPoolIdentityProviderFacebook(stack, jsii.String("FacebookProvider"), &awscognito.UserPoolIdentityProviderFacebookProps{
		ClientId:     facebookSecret.SecretValueFromJson(jsii.String("clientId")).UnsafeUnwrap(),
		ClientSecret: facebookSecret.SecretValueFromJson(jsii.String("clientSecret")).UnsafeUnwrap(),
		UserPool:     userPool,
		AttributeMapping: &awscognito.AttributeMapping{
			Email: awscognito.ProviderAttribute_FACEBOOK_EMAIL()},
	})

	domain := userPool.AddDomain(jsii.String("UserPoolDomain"), &awscognito.UserPoolDomainOptions{
		CognitoDomain: &awscognito.CognitoDomainOptions{
			DomainPrefix: jsii.String("chat-app"), // Customize as needed
		},
	})

	awscdk.NewCfnOutput(stack, jsii.String("CognitoHostedUIUrl"), &awscdk.CfnOutputProps{
		Value:      domain.DomainName(),
		ExportName: jsii.String("CognitoHostedUIUrl"),
	})

	// Export User Pool ID and Client ID for use in other stacks
	awscdk.NewCfnOutput(stack, jsii.String("UserPoolId"), &awscdk.CfnOutputProps{
		Value: userPool.UserPoolId(),
	})
	awscdk.NewCfnOutput(stack, jsii.String("UserPoolClientId"), &awscdk.CfnOutputProps{
		Value:      userPoolClient.UserPoolClientId(),
		ExportName: jsii.String("UserPoolClientId"),
	})

	return stack
}
