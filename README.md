# tictacgo

This is a serverless implementation of TicTacGo. It's techstack is: 
- AWS 
- Pulumi for IaC
- AWS Lambda: Serverless functions in Go (Lambda)
- AWS API Gateway: Exposing the websocket Api
- DynamoDB for storing the data of the set
- SAM for local testing of lambda functions

## Overall architecture

```mermaid
sequenceDiagram
    Client->>ApiGateway: $connect to websocket
    ApiGateway->>Client: Successfully connected
    Client->>ApiGateway: setup -- playerName, gameId
    ApiGateway->>Lambda:  setup -- Triage Game
    Lambda->>DynamoDB: get game
    DynamoDB->>Lambda: GameItem  
    alt game already full
      Lambda->>ApiGateway: Game already occupied
      ApiGateway->>Client: Game already occupied
    else game does not exist yet
      Lambda->>DynamoDB: Add player + game
      DynamoDB->>Lambda: Playeradded
      Lambda->>ApiGateway: Waiting for second player
      ApiGateway->>Client: Waiting for second player
    else player 1 is waiting
      Lambda->>DynamoDB: Add second player
      DynamoDB->>Lambda: Player added
      Lambda->>ApiGateway: Lets go and play 
      ApiGateway->>Client: Lets go and play
    end    
```

## Install 

1. Create the Lambda Go binaries
```sh 
make build
```
2. Install the whole thing to your AWS Account, just run the pulumi IaC stack: 
```sh
pulumi up 
```

## Development of lambda functions

To develop the lambda functions locally, we rely on [AWS SAM](https://docs.aws.amazon.com/serverless-application-model/latest/developerguide/what-is-sam.html)
Imagine, developing the setup lambda function: 
```sh
cd handler
sam build
sam local invoke setup -e events/setup-event.json
```
For this you have to have the ```template.yaml``` set up fitting the lambda functions you have created until now.








