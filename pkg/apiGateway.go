package pkg

import (
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/apigatewayv2"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/lambda"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func ApiGateway(ctx *pulumi.Context) (gw *apigatewayv2.Api, err error) {
	gw, err = apigatewayv2.NewApi(ctx, "tictacgo", &apigatewayv2.ApiArgs{
		ProtocolType:             pulumi.String("WEBSOCKET"),
		RouteSelectionExpression: pulumi.String("$request.body.action"),
		Tags:                     pulumi.StringMap{"app": pulumi.String("tictacgo")},
	})

	return gw, err

}

func ApiGatewayStage(ctx *pulumi.Context, gw *apigatewayv2.Api, stageName string) (stage *apigatewayv2.Stage, err error) {
	stage, err = apigatewayv2.NewStage(ctx, "tictacgo-"+stageName, &apigatewayv2.StageArgs{
		ApiId:      gw.ID(),
		AutoDeploy: pulumi.Bool(true),
	})
	return stage, err
}

func CreateIntegration(ctx *pulumi.Context, gw *apigatewayv2.Api, function *lambda.Function, routeKey string) (onConnectIntegration *apigatewayv2.Integration, err error) {

	onConnectIntegration, err = apigatewayv2.NewIntegration(ctx, "tictacgo-"+routeKey+"-integration", &apigatewayv2.IntegrationArgs{
		ApiId:                   gw.ID(),
		IntegrationType:         pulumi.String("AWS_PROXY"),
		ConnectionType:          pulumi.String("INTERNET"),
		ContentHandlingStrategy: pulumi.String("CONVERT_TO_TEXT"),
		Description:             pulumi.String("Lambda example"),
		IntegrationMethod:       pulumi.String("POST"),
		IntegrationUri:          function.InvokeArn,
		PassthroughBehavior:     pulumi.String("WHEN_NO_MATCH"),
	})
	return onConnectIntegration, err
}

func CreateRoute(ctx *pulumi.Context, gw *apigatewayv2.Api, integration *apigatewayv2.Integration, routeKey string) (route *apigatewayv2.Route, err error) {

	intId := integration.ID().ApplyT(func(id pulumi.ID) string {
		return string(id)
	}).(pulumi.StringOutput)

	target := pulumi.Sprintf("integrations/%s", intId)
	route, err = apigatewayv2.NewRoute(ctx, "tictacgo-"+routeKey+"-route", &apigatewayv2.RouteArgs{
		ApiId:                            gw.ID(),
		RouteKey:                         pulumi.String(routeKey),
		Target:                           target,
		RouteResponseSelectionExpression: pulumi.String("$default"),
	})
	return route, err
}

func CreateResponse(ctx *pulumi.Context, gw *apigatewayv2.Api, route *apigatewayv2.Route, routeKey string) (response *apigatewayv2.RouteResponse, err error) {

	integrationId := route.ID().ApplyT(func(id pulumi.ID) string {
		return string(id)
	}).(pulumi.StringOutput)

	response, err = apigatewayv2.NewRouteResponse(ctx, "tictacgo-"+routeKey+"-response", &apigatewayv2.RouteResponseArgs{
		ApiId:            gw.ID(),
		RouteId:          integrationId,
		RouteResponseKey: pulumi.String("$default"),
	})
	return response, err
}
