package main

import (
	"fmt"

	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/apigatewaymanagementapi"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

// It setups up the user for the game
// Checks if game is still joinable
// Checks if both players are already there
// Otherwise it needs to wait for the user to join
// Wait should be implemented on clientside as otherwise lambda will run forever
// Maybe set a timeout on the wait, like max 1 min or so to reduce calls
// Goes into response.body and will be parsed by the client. Based on internal status the client will now what to do
// 1 ->
// 2 ->
// 3 ->
// 4 ->
type ResponseBody struct {
	InternalStatus int    `json:internalStatus`
	ClientMessage  string `json:clientMessage`
}
type Item struct {
	GameId  string
	Player1 Player
	Player2 Player
}
type Player struct {
	Name         string
	ConnectionId string
}
type RequestBody struct {
	Action string
	Params map[string]string
}

func handler(request events.APIGatewayWebsocketProxyRequest) (events.APIGatewayProxyResponse, error) {

	var requestBody RequestBody
	err := json.Unmarshal([]byte(request.Body), &requestBody)
	if err != nil {
		fmt.Println(err)
	}

	tableName := "tictacgo-e8e3a4a"
	gameId := requestBody.Params["gameId"]
	playerName := requestBody.Params["playerName"]
	svc := dynamodb.New(session.New())

	query := &dynamodb.GetItemInput{
		TableName: &tableName,
		Key: map[string]*dynamodb.AttributeValue{
			"GameId": {
				S: &gameId,
			},
		},
	}

	result, err := svc.GetItem(query)
	// Check if there is a problem with database call
	if err != nil {

		return createResponse(500, 1, "Error during game setup. Couldn't connect to database"+err.Error())

		// Check if game is maybe already full
	} else if len(result.Item) > 2 {
		return createResponse(403, 2, "Can't connect to already occupied game")
	}

	//Now proceed with either creating a game or updating the existing one
	if len(result.Item) == 0 {
		createGame := dynamodb.PutItemInput{
			TableName: &tableName,
			Item: map[string]*dynamodb.AttributeValue{
				"GameId": {
					S: &gameId,
				},
				"Player1": {
					M: map[string]*dynamodb.AttributeValue{
						"ConnectionId": {S: &request.RequestContext.ConnectionID},
						"name":         {S: &playerName},
					},
				},
			},
		}
		_, err := svc.PutItem(&createGame)
		if err != nil {
			return createResponse(500, 3, "Error during game setup. Couldn't create game "+err.Error())
		}
		return createResponse(201, 4, "Successfully connected to the game. Waiting for the second player...")

	}
	if len(result.Item) == 2 {
		item := result.Item
		item["Player2"] = &dynamodb.AttributeValue{
			M: map[string]*dynamodb.AttributeValue{
				"ConnectionId": {S: &request.RequestContext.ConnectionID},
				"name":         {S: &playerName},
			},
		}
		putItem := dynamodb.PutItemInput{
			TableName: &tableName,
			Item:      item,
		}
		_, err := svc.PutItem(&putItem)

		if err != nil {
			return createResponse(500, 5, "Error during game setup. Couldn't update game "+err.Error())
		} else {
			// In this case it's player 2 and this should send an update to player 1 as well that player 2 has joined
			// use APIGateway.Post for this.
			// If the connection doesn't work, we need to drop it as the player might have disconnected
			//TO-DO: Make this endpoint dynamic!
			endpoint := "4cpq656h77.execute-api.eu-central-1.amazonaws.com/tictacgo-development-232d7a6"
			gwApi := apigatewaymanagementapi.New(session.New(&aws.Config{
				Endpoint: &endpoint,
				Region:   aws.String("eu-central-1")},
			))
			values := Item{}
			//TO-DO: Error handling
			_ = dynamodbattribute.UnmarshalMap(item, &values)
			fmt.Println(values)
			body, _ := json.Marshal(ResponseBody{
				InternalStatus: 2,
				ClientMessage:  "Your partner has connected. Let's go :))",
			})
			_, err = gwApi.PostToConnection(&apigatewaymanagementapi.PostToConnectionInput{
				ConnectionId: &values.Player1.ConnectionId,
				Data:         body,
			})
			if err != nil {
				return createResponse(500, 6, "Couldn't connect to the other player. Abort"+err.Error())
			}

			// here make sure that we send the right update to the client, about everything being fine and the game can start

			return createResponse(201, 7, "Successfully connected to the game")

		}

	}
	// Catching all the other cases
	return events.APIGatewayProxyResponse{
		StatusCode: 500,
		Body:       "Error during game setup" + err.Error(),
	}, nil
}

func createResponse(statusCode int, internalStatus int, clientMessage string) (events.APIGatewayProxyResponse, error) {
	body, err := json.Marshal(ResponseBody{
		InternalStatus: internalStatus,
		ClientMessage:  clientMessage,
	})
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: statusCode,
			Body:       err.Error(),
		}, nil
	}
	return events.APIGatewayProxyResponse{
		StatusCode: statusCode,
		Body:       string(body),
	}, nil

}

func main() {
	lambda.Start(handler)
}
