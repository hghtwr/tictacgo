package pkg

import (
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/dynamodb"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func DynamoDb(ctx *pulumi.Context) (db *dynamodb.Table, err error) {
	db, err = dynamodb.NewTable(ctx, "tictacgo", &dynamodb.TableArgs{
		Attributes: dynamodb.TableAttributeArray{
			&dynamodb.TableAttributeArgs{
				Name: pulumi.String("GameId"),
				Type: pulumi.String("S"),
			},
		},
		BillingMode:  pulumi.String("PROVISIONED"),
		HashKey:      pulumi.String("GameId"),
		ReadCapacity: pulumi.Int(20),
		Tags: pulumi.StringMap{
			"Environment": pulumi.String("development"),
			"Name":        pulumi.String("tictacgo"),
		},
		WriteCapacity: pulumi.Int(20),
	})
	if err != nil {
		return db, err
	} else {
		return db, err
	}

}
