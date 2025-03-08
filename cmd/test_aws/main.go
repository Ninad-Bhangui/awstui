package main

import (
	"context"
	"fmt"
	"os"

	"github.com/Ninad-Bhangui/awstui/aws/services"
	"github.com/aws/aws-sdk-go-v2/config"
)

func main() {
	// Load AWS configuration
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		fmt.Printf("Failed to load AWS config: %v\n", err)
		os.Exit(1)
	}

	// Test EC2 functions
	fmt.Println("Testing EC2 functions...")
	instances, err := services.ListEC2Instances(cfg)
	if err != nil {
		fmt.Printf("Failed to list EC2 instances: %v\n", err)
	} else {
		fmt.Printf("Found %d EC2 instances\n", len(instances))
		for _, instance := range instances {
			fmt.Printf("- %s (%s): %s\n", instance.Name, instance.ID, instance.State)

			// Get details for first instance
			if details, err := services.GetInstanceDetail(cfg, instance.ID); err == nil {
				fmt.Printf("  Details: %s\n", details)
				break // Just show first one
			}
		}
	}

	// Test ECR functions
	fmt.Println("\nTesting ECR functions...")
	repos, err := services.ListECRRepos(cfg)
	if err != nil {
		fmt.Printf("Failed to list ECR repositories: %v\n", err)
	} else {
		fmt.Printf("Found %d ECR repositories\n", len(repos))
		for _, repo := range repos {
			fmt.Printf("- %s: %d images\n", repo.Name, repo.ImageCount)

			// Get details for first repo
			if details, err := services.GetRepoDetail(cfg, repo.Name); err == nil {
				fmt.Printf("  Details: %s\n", details)
				break // Just show first one
			}
		}
	}

	// Test Lambda functions
	fmt.Println("\nTesting Lambda functions...")
	functions, err := services.ListLambdaFunctions(cfg)
	if err != nil {
		fmt.Printf("Failed to list Lambda functions: %v\n", err)
	} else {
		fmt.Printf("Found %d Lambda functions\n", len(functions))
		for _, fn := range functions {
			fmt.Printf("- %s (%s): %s\n", fn.Name, fn.Runtime, fn.State)

			// Get details for first function
			if details, err := services.GetFunctionDetail(cfg, fn.Name); err == nil {
				fmt.Printf("  Details: %s\n", details)
				break // Just show first one
			}
		}
	}

	// Test Secrets Manager functions
	fmt.Println("\nTesting Secrets Manager functions...")
	secrets, err := services.ListSecrets(cfg)
	if err != nil {
		fmt.Printf("Failed to list secrets: %v\n", err)
	} else {
		fmt.Printf("Found %d secrets\n", len(secrets))
		for _, secret := range secrets {
			fmt.Printf("- %s: %s\n", secret.Name, secret.Description)

			// Get details for first secret
			if details, err := services.GetSecretDetail(cfg, secret.ARN); err == nil {
				fmt.Printf("  Details: %s\n", details)
				// Note: Not showing secret value for security
				break // Just show first one
			}
		}
	}
}
