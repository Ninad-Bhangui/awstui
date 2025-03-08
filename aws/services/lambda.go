package services

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
)

// ListLambdaFunctions returns a list of Lambda functions
func ListLambdaFunctions(ctx context.Context, cfg config.Config) ([]LambdaFunction, error) {
	// Convert config.Config to aws.Config
	awsCfg := GetAWSConfig(cfg).(aws.Config)
	client := lambda.NewFromConfig(awsCfg)

	resp, err := client.ListFunctions(ctx, &lambda.ListFunctionsInput{})
	if err != nil {
		return nil, err
	}

	var functions []LambdaFunction
	for _, fn := range resp.Functions {
		lastMod, _ := time.Parse(time.RFC3339, *fn.LastModified)
		functions = append(functions, LambdaFunction{
			Name:         *fn.FunctionName,
			Runtime:      string(fn.Runtime),
			MemorySize:   int64(*fn.MemorySize),
			LastModified: lastMod,
		})
	}

	return functions, nil
}

// GetFunctionDetail returns detailed information about a Lambda function
func GetFunctionDetail(cfg aws.Config, functionName string) (string, error) {
	client := lambda.NewFromConfig(cfg)

	// Get function configuration
	input := &lambda.GetFunctionInput{
		FunctionName: aws.String(functionName),
	}

	result, err := client.GetFunction(context.TODO(), input)
	if err != nil {
		return "", fmt.Errorf("failed to get function details: %w", err)
	}

	// Get function policy
	policyInput := &lambda.GetPolicyInput{
		FunctionName: aws.String(functionName),
	}

	policy, err := client.GetPolicy(context.TODO(), policyInput)
	// Ignore error as policy might not exist

	// Get function concurrency
	concurrencyInput := &lambda.GetFunctionConcurrencyInput{
		FunctionName: aws.String(functionName),
	}

	concurrency, err := client.GetFunctionConcurrency(context.TODO(), concurrencyInput)
	// Ignore error as concurrency might not be set

	// Combine all details
	details := struct {
		Configuration *lambda.GetFunctionOutput
		Policy        *lambda.GetPolicyOutput
		Concurrency   *lambda.GetFunctionConcurrencyOutput
	}{
		Configuration: result,
		Policy:        policy,
		Concurrency:   concurrency,
	}

	// Marshal to JSON
	jsonBytes, err := json.MarshalIndent(details, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal function details: %w", err)
	}

	return string(jsonBytes), nil
}
