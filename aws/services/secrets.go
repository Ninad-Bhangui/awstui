package services

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

// ListSecrets returns a list of secrets
func ListSecrets(ctx context.Context, cfg config.Config) ([]Secret, error) {
	// Convert config.Config to aws.Config
	awsCfg := GetAWSConfig(cfg).(aws.Config)
	client := secretsmanager.NewFromConfig(awsCfg)

	resp, err := client.ListSecrets(ctx, &secretsmanager.ListSecretsInput{})
	if err != nil {
		return nil, err
	}

	var secrets []Secret
	for _, s := range resp.SecretList {
		// Calculate days until rotation
		var daysUntilRotation int64 = -1
		if s.NextRotationDate != nil {
			daysUntilRotation = int64(time.Until(*s.NextRotationDate).Hours() / 24)
		}

		secrets = append(secrets, Secret{
			Name:              *s.Name,
			LastModified:      *s.LastChangedDate,
			DaysUntilRotation: daysUntilRotation,
		})
	}

	return secrets, nil
}

// GetSecretDetail returns detailed information about a secret (excluding the secret value)
func GetSecretDetail(cfg aws.Config, secretID string) (string, error) {
	client := secretsmanager.NewFromConfig(cfg)

	// Get secret metadata
	input := &secretsmanager.DescribeSecretInput{
		SecretId: aws.String(secretID),
	}

	result, err := client.DescribeSecret(context.TODO(), input)
	if err != nil {
		return "", fmt.Errorf("failed to get secret details: %w", err)
	}

	// Get secret policy
	policyInput := &secretsmanager.GetResourcePolicyInput{
		SecretId: aws.String(secretID),
	}

	policy, err := client.GetResourcePolicy(context.TODO(), policyInput)
	// Ignore error as policy might not exist

	// Combine details (excluding secret value)
	details := struct {
		Metadata *secretsmanager.DescribeSecretOutput
		Policy   *secretsmanager.GetResourcePolicyOutput
	}{
		Metadata: result,
		Policy:   policy,
	}

	// Marshal to JSON
	jsonBytes, err := json.MarshalIndent(details, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal secret details: %w", err)
	}

	return string(jsonBytes), nil
}

// GetSecretValue retrieves the actual secret value
// Note: This is separated from GetSecretDetail for security reasons
func GetSecretValue(cfg aws.Config, secretID string) (string, error) {
	client := secretsmanager.NewFromConfig(cfg)

	input := &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(secretID),
	}

	result, err := client.GetSecretValue(context.TODO(), input)
	if err != nil {
		return "", fmt.Errorf("failed to get secret value: %w", err)
	}

	return *result.SecretString, nil
}
