package main

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

// It setups up the user for the game
// Checks if game is still joinable
// Checks if both players are already there
// Otherwise it needs to wait for the user to join
// Wait should be implemented on clientside as otherwise lambda will run forever
// Maybe set a timeout on the wait, like max 1 min or so to reduce calls
func handler(request events.APIGatewayWebsocketProxyRequest) (events.APIGatewayProxyResponse, error) {
	params := request.QueryStringParameters
	tableName := "tictacgo"
	gameId := params["gameId"]
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
	if err == nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "Error during game setup. Couldn't connect to database" + err.Error(),
		}, nil

		// Check if game is maybe already full
	} else if len(result.Item) > 1 {
		return events.APIGatewayProxyResponse{
			StatusCode: 403,
			Body:       "Can't connect to already occupied game",
		}, nil
	}

	//Now proceed with either creating a game or updating the existing one

	if len(result.Item) == 0 {
		createGame := dynamodb.PutItemInput{
			TableName: &tableName,
			Item: map[string]*dynamodb.AttributeValue{
				"GameId": gameId,
				"Player1": map[string]string{
					"ConnectionId": request.RequestContext.ConnectionID,
					"name":         params["playerName"],
				},
			},
		}

		_, err := svc.PutItem(&createGame)
		if err != nil {
			return events.APIGatewayProxyResponse{
				StatusCode: 500,
				Body:       "Error during game setup. Couldn't create game" + err.Error(),
			}, nil
		}
	}

	if result.Item.lenght == 1 {

		//update, err := dynamodb.putItem(dynamodb.PutItemInput{})
	}

}

func main() {
	lambda.Start(handler)
}
