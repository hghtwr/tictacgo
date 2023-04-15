package pkg

import (
	"github.com/pulumi/pulumi-aws-native/sdk/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/apigatewayv2"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/iam"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/lambda"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func OnConnectIam(ctx *pulumi.Context) (onConnectIamRole *iam.Role, err error) {
	onConnectIamRole, err = iam.NewRole(ctx, "tictacgo-onconnect", &iam.RoleArgs{
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
	return onConnectIamRole, err
}

func OnConnect(ctx *pulumi.Context, role *iam.Role) (function *lambda.Function, err error) {
	//role, err := OnConnectIam(ctx)

	args := &lambda.FunctionArgs{
		Handler: pulumi.String("handler"),
		Runtime: pulumi.String("go1.x"),
		Role:    role.Arn,
		Code:    pulumi.NewFileArchive("./handler/handler.zip"),
	}

	// Create the lambda using the args.
	function, err = lambda.NewFunction(
		ctx,
		"tictacgo-onconnect",
		args)

	return function, err
}

func OnConnectPermission(ctx *pulumi.Context, function *lambda.Function, route *apigatewayv2.Route, gw *apigatewayv2.Api) (permission *lambda.Permission, err error) {
	accountId, err := aws.GetAccountId(ctx)
	if err != nil {
		return permission, err
	}
	arn := gw.ID().ApplyT(func(id pulumi.ID) string {
		return string(id)
	}).(pulumi.StringOutput)
	target := pulumi.Sprintf("arn:aws:execute-api:eu-central-1:%s:%s/*/$connect", accountId.AccountId, arn)

	permission, err = lambda.NewPermission(ctx, "tictacgo-connect-apigateway", &lambda.PermissionArgs{
		Action:    pulumi.String("lambda:InvokeFunction"),
		Function:  function.Name,
		Principal: pulumi.String("apigateway.amazonaws.com"),
		SourceArn: target,
	})

	return permission, err
}
