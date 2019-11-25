package main

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"

	"github.com/hazelcast/hazelcast-cloud-go-demo/common"
	"github.com/hazelcast/hazelcast-cloud-go-demo/s3import"
)

// HandleRequest handles S3 events
func HandleRequest(ctx context.Context, s3Event events.S3Event) error {
	clusterName := os.Getenv("CLUSTER_NAME")
	password := os.Getenv("CLUSTER_PASSWORD")
	token := os.Getenv("DISCOVERY_TOKEN")

	hzClient, err := common.NewHazelcastClient(clusterName, password, token)
	if err != nil {
		return err
	}

	airportsMap, err := common.GetMap(hzClient, "airports")
	if err != nil {
		return err
	}

	session, err := session.NewSession()
	if err != nil {
		return err
	}
	s3 := s3.New(session)
	bucket := s3Event.Records[0].S3.Bucket.Name
	key := s3Event.Records[0].S3.Object.Key

	fmt.Printf("Handling upload into S3 bucket '%v'\n", bucket)
	if err = s3import.ImportAirports(s3, bucket, key, airportsMap); err != nil {
		return err
	}

	return nil
}

func main() {
	lambda.Start(HandleRequest)
}
