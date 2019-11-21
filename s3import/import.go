package s3import

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"

	"github.com/hazelcast/hazelcast-cloud-go-demo/common"
	"github.com/hazelcast/hazelcast-go-client/core"
)

// ImportAirports copies the airports data from file into a Hazelcast map
func ImportAirports(s3Client *s3.S3, bucket string, key string, airportsMap core.Map) error {
	objOutput, err := s3Client.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return err
	}
	defer objOutput.Body.Close()

	dat, err := ioutil.ReadAll(objOutput.Body)
	if err != nil {
		return fmt.Errorf("Failed to read S3 object body: %w", err)
	}

	var airports []common.Airport
	err = json.Unmarshal(dat, &airports)
	if err != nil {
		return fmt.Errorf("Failed to unmarshal the airports array: %w", err)
	}

	for _, a := range airports {
		airportsMap.Put(a.Code, &a)
	}

	size, err := airportsMap.Size()
	if err != nil {
		return fmt.Errorf("Failed to get a map size: %w", err)
	}

	fmt.Println("Successfully imported the data about ", size, " airports")
	return nil
}
