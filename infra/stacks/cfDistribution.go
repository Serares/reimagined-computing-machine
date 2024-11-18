package stacks

import (
	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/constructs-go/constructs/v10"
)

type CloudFrontStackProps struct {
	awscdk.StackProps
}

/*
*
❗This is just a template
*
*/
func CloudFrontStack(scope constructs.Construct, id string, props *CloudFrontStackProps) {
	// stack := awscdk.NewStack(scope, &id, &props.StackProps)

	// // Create an S3 bucket for hosting
	// bucket := awss3.NewBucket(stack, jsii.String("FrontendBucket"), &awss3.BucketProps{
	// 	WebsiteIndexDocument: jsii.String("index.html"),
	// 	PublicReadAccess:     jsii.Bool(true),
	// })

	// // Create CloudFront distribution
	// distribution := awscloudfront.NewDistribution(stack, jsii.String("CloudFrontDistribution"), &awscloudfront.DistributionProps{
	// 	DefaultBehavior: &awscloudfront.BehaviorOptions{
	// 		Origin: awscloudfrontorigins.NewS3Origin(bucket, nil),
	// 	},
	// })

	// // Export CloudFront URL
	// awscdk.NewCfnOutput(stack, jsii.String("CloudFrontUrl"), &awscdk.CfnOutputProps{
	// 	Value:      distribution.DomainName(),
	// 	ExportName: jsii.String("CloudFrontUrl"),
	// })

	// return stack
}
