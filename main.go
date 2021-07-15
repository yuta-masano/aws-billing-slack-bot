package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/costexplorer"
)

const (
	Region string = "ap-northeast-1"
)

type MyResponse struct {
	Message []byte `json:"message"`
}

func BillingCheck(ctx context.Context) (*MyResponse, error) {
	// Must be in YYYY-MM-DD Format
	start := "2019-06-01"
	end := "2019-07-01"
	granularity := "MONTHLY"
	metrics := []string{"AmortizedCost"}

	// Initialize a session in us-east-1 that the SDK will use to load credentials
	cfg := new(aws.Config)
	cfg.Region = aws.String(Region)

	sess, err := session.NewSession(cfg)
	if err != nil {
		return nil, fmt.Errorf("could not create new session: %w", err)
	}

	// Create Cost Explorer Service Client
	svc := costexplorer.New(sess)
	input := new(costexplorer.GetCostAndUsageInput)
	input.TimePeriod = &costexplorer.DateInterval{
		Start: aws.String(start),
		End:   aws.String(end),
	}
	input.Granularity = aws.String(granularity)
	input.Metrics = aws.StringSlice(metrics)

	result, err := svc.GetCostAndUsage(input)
	if err != nil {
		return nil, fmt.Errorf("could not get cost info: %w", err)
	}

	json, err := json.Marshal(result)
	if err != nil {
		return nil, fmt.Errorf("could not get marshal: %w", err)
	}

	return &MyResponse{Message: json}, nil
}

func main() {
	lambda.Start(BillingCheck)
}
