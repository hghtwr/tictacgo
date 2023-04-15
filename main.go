package main

import (
	"github.com/hghtwr/tictacgo/pkg"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		_, err := pkg.DynamoDb(ctx)
		if err != nil {
			return err
		}

		onConnectIamRole, err := pkg.OnConnectIam(ctx)

		if err != nil {
			return err
		}

		onConnectLambdaFunction, err := pkg.OnConnect(ctx, onConnectIamRole)
		if err != nil {
			return err
		}
		gw, err := pkg.ApiGateway(ctx)
		if err != nil {
			return err
		}
		onConnectIntegration, err := pkg.OnConnectIntegration(ctx, gw, onConnectLambdaFunction)
		if err != nil {
			return err
		}

		onConnectRoute, err := pkg.OnConnectRoute(ctx, gw, onConnectIntegration)
		if err != nil {
			return err
		}

		_, err = pkg.OnConnectResponse(ctx, gw, onConnectIntegration)
		if err != nil {
			return err
		}
		_, err = pkg.OnConnectPermission(ctx, onConnectLambdaFunction, onConnectRoute, gw)
		if err != nil {
			return err
		}
		_, err = pkg.ApiGatewayStage(ctx, gw)
		if err != nil {
			return err
		}
		return nil
	})
}
