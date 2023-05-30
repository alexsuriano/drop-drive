package main

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/alexsuriano/drop-drive/internal/bucket"
	"github.com/alexsuriano/drop-drive/internal/queue"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
)

func main() {
	//create queue config
	queueConfig := queue.RabbitMQConfig{
		URL:       os.Getenv("RABBITMQ_URL"),
		TopicName: os.Getenv("RABBITMQ_TOPIC_NAME"),
		Timeout:   time.Second * 30,
	}

	//create new queue connection
	queueConnection, err := queue.New(queue.RabbitMQ, queueConfig)
	if err != nil {
		panic(err)
	}

	//create channel to consume messages
	queueChannel := make(chan queue.QueueDTO)
	queueConnection.Consume(queueChannel)

	//create bucket config
	bucketConfig := bucket.AwsConfig{
		Config: &aws.Config{
			Region:      aws.String(os.Getenv("AWS_REGION")),
			Credentials: credentials.NewStaticCredentials(os.Getenv("AWS_KEYS"), os.Getenv("AWS_SCRET"), ""),
		},
		BucketDownload: "drop-drive-raw",
		BucketUpload:   "drop-drive-gzip",
	}

	//create new bucket connection
	bucketConnection, err := bucket.New(bucket.AwsProvider, bucketConfig)
	if err != nil {
		panic(err)
	}

	//looping
	for message := range queueChannel {
		source := fmt.Sprintf("%s/%s", message.Path, message.Filename)
		destiny := fmt.Sprintf("%d_%s", message.ID, message.Filename)

		file, err := bucketConnection.Download(source, destiny)
		if err != nil {
			log.Printf("ERROR: %v", err)
			continue
		}

		body, err := io.ReadAll(file)
		if err != nil {
			log.Printf("ERROR: %v", err)
			continue
		}

		var buffer bytes.Buffer
		gzipWriter := gzip.NewWriter(&buffer)
		_, err = gzipWriter.Write(body)
		if err != nil {
			log.Printf("ERROR %v", err)
			continue
		}

		if err := gzipWriter.Close(); err != nil {
			log.Printf("ERROR %v", err)
			continue
		}

		gzipReader, err := gzip.NewReader(&buffer)
		if err != nil {
			log.Printf("ERROR %v", err)
			continue
		}

		if err = bucketConnection.Upload(gzipReader, source); err != nil {
			log.Printf("ERROR %v", err)
			continue
		}

		if err = os.Remove(destiny); err != nil {
			log.Printf("ERROR %v", err)
			continue
		}
	}
}
