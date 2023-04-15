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
		gw, err := pkg.ApiGateway(ctx)
		if err != nil {
			return err
		}

		//############################ OnConnect #################################
		onConnectIamRole, err := pkg.CreateIam(ctx, "$connect")

		if err != nil {
			return err
		}

		onConnectLambdaFunction, err := pkg.CreateFunction(ctx, onConnectIamRole, "$connect")
		if err != nil {
			return err
		}

		onConnectIntegration, err := pkg.CreateIntegration(ctx, gw, onConnectLambdaFunction, "$connect")
		if err != nil {
			return err
		}

		onConnectRoute, err := pkg.CreateRoute(ctx, gw, onConnectIntegration, "$connect")
		if err != nil {
			return err
		}

		_, err = pkg.CreateResponse(ctx, gw, onConnectRoute, "$connect")
		if err != nil {
			return err
		}
		_, err = pkg.CreateLambdaPermission(ctx, onConnectLambdaFunction, onConnectRoute, gw, "$connect")
		if err != nil {
			return err
		}
		//############################ Disconnect #################################
		onDisconnectIamRole, err := pkg.CreateIam(ctx, "$disconnect")

		if err != nil {
			return err
		}

		onDisconnectLambdaFunction, err := pkg.CreateFunction(ctx, onDisconnectIamRole, "$disconnect")
		if err != nil {
			return err
		}

		onDisconnectIntegration, err := pkg.CreateIntegration(ctx, gw, onDisconnectLambdaFunction, "$disconnect")
		if err != nil {
			return err
		}

		onDisconnectRoute, err := pkg.CreateRoute(ctx, gw, onDisconnectIntegration, "$disconnect")
		if err != nil {
			return err
		}

		_, err = pkg.CreateResponse(ctx, gw, onDisconnectRoute, "$disconnect")
		if err != nil {
			return err
		}
		_, err = pkg.CreateLambdaPermission(ctx, onDisconnectLambdaFunction, onDisconnectRoute, gw, "$disconnect")
		if err != nil {
			return err
		}
		//############################ Turn #################################
		onTurnIamRole, err := pkg.CreateIam(ctx, "turn")

		if err != nil {
			return err
		}

		onTurnLambdaFunction, err := pkg.CreateFunction(ctx, onTurnIamRole, "turn")
		if err != nil {
			return err
		}

		onTurnIntegration, err := pkg.CreateIntegration(ctx, gw, onTurnLambdaFunction, "turn")
		if err != nil {
			return err
		}

		onTurnRoute, err := pkg.CreateRoute(ctx, gw, onTurnIntegration, "turn")
		if err != nil {
			return err
		}

		_, err = pkg.CreateResponse(ctx, gw, onTurnRoute, "turn")
		if err != nil {
			return err
		}
		_, err = pkg.CreateLambdaPermission(ctx, onTurnLambdaFunction, onTurnRoute, gw, "turn")
		if err != nil {
			return err
		}

		_, err = pkg.ApiGatewayStage(ctx, gw, "development")
		if err != nil {
			return err
		}
		return nil
	})
}
