package pkg

import (
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
