package pkg

import (
	"github.com/pulumi/pulumi-aws-native/sdk/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/apigatewayv2"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/iam"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/lambda"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
	"strings"
)

func CreateIam(ctx *pulumi.Context, routeKey string) (iamRole *iam.Role, err error) {
	iamRole, err = iam.NewRole(ctx, "tictacgo-"+strings.Trim(routeKey, "$"), &iam.RoleArgs{
		AssumeRolePolicy: pulumi.String(`{
				"Version": "2012-10-17",
				"Statement": [{
					"Sid": "",
					"Effect": "Allow",
					"Principal": {
						"Service": "lambda.amazonaws.com"
					},
					"Action": "sts:AssumeRole"
				}]
			}`),
	})
	return iamRole, err
}

func CreateFunction(ctx *pulumi.Context, role *iam.Role, routeKey string) (function *lambda.Function, err error) {
	//role, err := OnConnectIam(ctx)

	args := &lambda.FunctionArgs{
		Handler: pulumi.String(strings.Trim(routeKey, "$")),
		Runtime: pulumi.String("go1.x"),
		Role:    role.Arn,
		Code:    pulumi.NewFileArchive("./handler/" + strings.Trim(routeKey, "$") + ".zip"),

		//Code:    pulumi.NewFileArchive("./handler/handler.zip"),
	}

	// Create the lambda using the args.
	function, err = lambda.NewFunction(
		ctx,
		"tictacgo-"+strings.Trim(routeKey, "$"),
		args)

	return function, err
}
func CreateLambdaPermission(ctx *pulumi.Context, function *lambda.Function, route *apigatewayv2.Route, gw *apigatewayv2.Api, routeKey string) (permission *lambda.Permission, err error) {
	accountId, err := aws.GetAccountId(ctx)
	if err != nil {
		return permission, err
	}
	arn := gw.ID().ApplyT(func(id pulumi.ID) string {
		return string(id)
	}).(pulumi.StringOutput)
	target := pulumi.Sprintf("arn:aws:execute-api:eu-central-1:%s:%s/*/%s", accountId.AccountId, arn, routeKey)

	permission, err = lambda.NewPermission(ctx, "tictacgo-"+strings.Trim(routeKey, "$")+"-apigateway", &lambda.PermissionArgs{
		Action:    pulumi.String("lambda:InvokeFunction"),
		Function:  function.Name,
		Principal: pulumi.String("apigateway.amazonaws.com"),
		SourceArn: target,
	})

	return permission, err
}
