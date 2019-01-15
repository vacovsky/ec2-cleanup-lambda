package main

import (
	"context"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

// HandleRequest is the entrypoint for AWS Lambda
func HandleRequest(ctx context.Context) (bool, error) {
	log.Println("Beginning execution of Ephemeral Instance Lifecycle Lambda function")
	result, err := reapInstances()
	log.Println("Ephemeral Instance Lifecycle Lambda function has completed:", err)
	return result, err
}

func main() {
	lambda.Start(HandleRequest)
}

func reapInstances() (bool, error) {
	// Initialize a session in us-west-2 that the SDK will use to load
	// credentials from the shared credentials file ~/.aws/credentials.
	session, err := session.NewSession(&aws.Config{
		Region: aws.String("us-west-2")},
	)

	// Create service client
	svc := ec2.New(session)

	var result *ec2.DescribeInstancesOutput

	result, err = svc.DescribeInstances(&ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			&ec2.Filter{
				Name: aws.String("tag:Ephemeral"),
				Values: []*string{
					aws.String("True"),
				},
			},
			// &ec2.Filter{
			// 	Name: aws.String("tag:DateCreated"),
			// 	Values: []*string{
			// 		aws.String(time.Now().Unix() - int64(86400)),
			// 	},
			// },
		},
	})

	for _, instances := range result.Reservations {
		for _, instance := range instances.Instances {
			stateCode := strconv.Itoa(int(aws.Int64Value(instance.State.Code)))
			if stateCode == "16" { // 16 represents the "running" instance state
				for _, tag := range instance.Tags {
					if aws.StringValue(tag.Key) == "DateCreated" { // could use the built-in metadata, but
						// the tag lets us adjust the value to keep it around longer if we want
						strV, err := strconv.Atoi(aws.StringValue(tag.Value))
						if err != nil {
							log.Println(err)
						}
						if int64(strV)+int64(86400) < time.Now().Unix() {
							shutdownInstanceByID(aws.StringValue(instance.InstanceId))
							log.Println("Shutting down short-lived instance:", aws.StringValue(instance.InstanceId))
						}
						if err != nil {
							log.Println("Unable to parse DateTime tag value for instance: ", instance.InstanceId)
						}
					}
				}
			}
		}
	}

	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	return true, err
}

func shutdownInstanceByID(instanceID string) (*ec2.TerminateInstancesOutput, error) {
	session, err := session.NewSession(&aws.Config{
		Region: aws.String("us-west-2")},
	)

	// Create service client
	svc := ec2.New(session)

	req := &ec2.TerminateInstancesInput{
		InstanceIds: []*string{
			aws.String(instanceID),
		},
	}
	output, err := svc.TerminateInstances(req)

	return output, err
}
