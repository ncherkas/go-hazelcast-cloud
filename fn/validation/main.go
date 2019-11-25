package main

import (
	"context"
	"encoding/gob"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/hazelcast/hazelcast-cloud-go-demo/common"
	"github.com/hazelcast/hazelcast-cloud-go-demo/validation"
	hazelcast "github.com/hazelcast/hazelcast-go-client"
)

type RequestType string

const (
	KEEP_ALIVE RequestType = "KEEP_ALIVE"
	VALIDATE   RequestType = "VALIDATE"
)

// Input comes with a type
type Input struct {
	Type RequestType `json:"type"`
	*validation.Request
}

// Output with a message
type Output struct {
	Message string `json:"result"`
}

var hzClient hazelcast.Client

func init() {
	clusterName := os.Getenv("CLUSTER_NAME")
	password := os.Getenv("CLUSTER_PASSWORD")
	token := os.Getenv("DISCOVERY_TOKEN")

	var err error
	hzClient, err = common.NewHazelcastClient(clusterName, password, token)
	if err != nil {
		panic(err)
	}

	gob.Register(&common.User{})
	gob.Register(&common.Airport{})
}

// HandleRequest handles API GW events
func HandleRequest(ctx context.Context, input Input) (*Output, error) {
	switch input.Type {
	case KEEP_ALIVE:
		return &Output{"Keep Alive request, nothing to validate"}, nil
	case VALIDATE:
		usersMap, err := common.GetMap(hzClient, "users")
		if err != nil {
			return nil, err
		}

		airportsMap, err := common.GetMap(hzClient, "airports")
		if err != nil {
			return nil, err
		}

		resp, err := validation.Apply(input.Request, usersMap, airportsMap)
		if err != nil {
			return nil, err
		}

		return &Output{resp.Message}, nil
	default:
		return nil, fmt.Errorf("Unsupported request type: %v", input.Type)
	}
}

func main() {
	lambda.Start(HandleRequest)
}
