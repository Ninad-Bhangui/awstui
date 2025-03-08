package services

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
)

// ListEC2Instances returns a list of EC2 instances
func ListEC2Instances(ctx context.Context, cfg config.Config) ([]EC2Instance, error) {
	// Convert config.Config to aws.Config
	awsCfg := GetAWSConfig(cfg).(aws.Config)
	client := ec2.NewFromConfig(awsCfg)

	resp, err := client.DescribeInstances(ctx, &ec2.DescribeInstancesInput{})
	if err != nil {
		return nil, err
	}

	var instances []EC2Instance
	for _, reservation := range resp.Reservations {
		for _, instance := range reservation.Instances {
			// Get instance name from tags
			var name string
			for _, tag := range instance.Tags {
				if *tag.Key == "Name" {
					name = *tag.Value
					break
				}
			}

			// Create instance info
			inst := EC2Instance{
				ID:        *instance.InstanceId,
				Name:      name,
				Type:      string(instance.InstanceType),
				State:     string(instance.State.Name),
				PrivateIP: stringOrEmpty(instance.PrivateIpAddress),
				PublicIP:  stringOrEmpty(instance.PublicIpAddress),
			}
			instances = append(instances, inst)
		}
	}

	return instances, nil
}

func stringOrEmpty(ptr *string) string {
	if ptr == nil {
		return "-"
	}
	return *ptr
}

// GetInstanceDetail returns detailed information about an EC2 instance
func GetInstanceDetail(cfg aws.Config, instanceID string) (string, error) {
	client := ec2.NewFromConfig(cfg)

	input := &ec2.DescribeInstancesInput{
		InstanceIds: []string{instanceID},
	}

	result, err := client.DescribeInstances(context.TODO(), input)
	if err != nil {
		return "", fmt.Errorf("failed to get instance details: %w", err)
	}

	if len(result.Reservations) == 0 || len(result.Reservations[0].Instances) == 0 {
		return "", fmt.Errorf("instance not found: %s", instanceID)
	}

	// Marshal the first instance to JSON
	jsonBytes, err := json.MarshalIndent(result.Reservations[0].Instances[0], "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal instance details: %w", err)
	}

	return string(jsonBytes), nil
}
