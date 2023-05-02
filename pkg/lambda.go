package pkg

import (
	"fmt"
	"strings"

	"github.com/pulumi/pulumi-aws-native/sdk/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/apigatewayv2"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/iam"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/lambda"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func CreateIam(ctx *pulumi.Context, routeKey string, gw *apigatewayv2.Api) (iamRole *iam.Role, err error) {
	accountId, err := aws.GetAccountId(ctx)
	if err != nil {
		return iamRole, err
	}
	arn := gw.ID().ApplyT(func(id pulumi.ID) string {
		return string(id)
	}).(pulumi.StringOutput)

	//target := pulumi.Sprintf("arn:aws:execute-api:eu-central-1:%s:%s/*/%s", accountId.AccountId, arn, routeKey)

	iamRole, err = iam.NewRole(ctx, "tictacgo-"+strings.Trim(routeKey, "$"), &iam.RoleArgs{
		InlinePolicies: iam.RoleInlinePolicyArray{
			&iam.RoleInlinePolicyArgs{
				Name: pulumi.String(fmt.Sprintf("manage-api-%s", strings.Trim(routeKey, "$"))),
				Policy: pulumi.Sprintf(`{
				"Version": "2012-10-17",
				"Statement": [
					{
						"Effect": "Allow",
						"Action": "execute-api:ManageConnections",
						"Resource": "arn:aws:execute-api:eu-central-1:%s:%s/tictacgo-development-*/POST/@connections/*"
					}
				]
			}`, accountId.AccountId, arn),
			},
			&iam.RoleInlinePolicyArgs{
				Name: pulumi.String(fmt.Sprintf("logs-%s", strings.Trim(routeKey, "$"))),
				Policy: pulumi.String(fmt.Sprintf(`{
					"Version": "2012-10-17",
					"Statement": [
						{
							"Effect": "Allow",
							"Action": "logs:CreateLogGroup",
							"Resource": "arn:aws:logs:eu-central-1:%s:*"
						},
						{
							"Effect": "Allow",
							"Action": [
								"logs:CreateLogStream",
								"logs:PutLogEvents"
							],
							"Resource": [
								"arn:aws:logs:eu-central-1:572618378599:log-group:/aws/lambda/tictacgo-%s-*:*"
							]
						}
					]
				}`, accountId.AccountId, strings.Trim(routeKey, "$"))),
			},
			&iam.RoleInlinePolicyArgs{
				Name: pulumi.String(fmt.Sprintf("dynamodb-%s", strings.Trim(routeKey, "$"))),
				Policy: pulumi.String(fmt.Sprintf(`{
					"Version": "2012-10-17",
					"Statement": [
						{
							"Sid": "VisualEditor0",
							"Effect": "Allow",
							"Action": [
								"dynamodb:PutItem",
								"dynamodb:DeleteItem",
								"dynamodb:GetItem",
								"dynamodb:Scan",
								"dynamodb:Query",
								"dynamodb:UpdateItem"
							],
							"Resource": "arn:aws:dynamodb:eu-central-1:%s:table/*"
						}
					]
				}`, accountId.AccountId)),
			},
		},
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

func CreateFunction(ctx *pulumi.Context, role *iam.Role, routeKey string, dbId pulumi.StringOutput, gwStage *apigatewayv2.Stage) (function *lambda.Function, err error) {
	//role, err := OnConnectIam(ctx)

	args := &lambda.FunctionArgs{
		Handler: pulumi.String(strings.Trim(routeKey, "$")),
		Runtime: pulumi.String("go1.x"),
		Environment: &lambda.FunctionEnvironmentArgs{
			Variables: pulumi.StringMap{
				"TABLE_NAME": pulumi.Sprintf("%s", dbId),
				"APIGATEWAY_ADDRESS": gwStage.InvokeUrl.ApplyT(func(url string) string {
					return strings.TrimPrefix(url, "wss://")
				}).(pulumi.StringOutput),
			},
		},
		Role: role.Arn,
		Code: pulumi.NewFileArchive("./handler/" + strings.Trim(routeKey, "$") + ".zip"),
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
