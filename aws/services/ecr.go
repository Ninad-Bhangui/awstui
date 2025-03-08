package services

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ecr"
)

// ListECRRepositories returns a list of ECR repositories
func ListECRRepositories(ctx context.Context, cfg config.Config) ([]ECRRepository, error) {
	// Convert config.Config to aws.Config
	awsCfg := GetAWSConfig(cfg).(aws.Config)
	client := ecr.NewFromConfig(awsCfg)

	// Get repositories
	repoResp, err := client.DescribeRepositories(ctx, &ecr.DescribeRepositoriesInput{})
	if err != nil {
		return nil, err
	}

	var repos []ECRRepository
	for _, repo := range repoResp.Repositories {
		// Get image count
		imgResp, err := client.DescribeImages(ctx, &ecr.DescribeImagesInput{
			RepositoryName: repo.RepositoryName,
		})
		if err != nil {
			// Skip repositories with errors
			continue
		}

		repos = append(repos, ECRRepository{
			Name:       *repo.RepositoryName,
			URI:        *repo.RepositoryUri,
			ImageCount: int64(len(imgResp.ImageDetails)),
			CreatedAt:  *repo.CreatedAt,
		})
	}

	return repos, nil
}

// GetRepoDetail returns detailed information about an ECR repository
func GetRepoDetail(cfg aws.Config, repoName string) (string, error) {
	client := ecr.NewFromConfig(cfg)

	input := &ecr.DescribeRepositoriesInput{
		RepositoryNames: []string{repoName},
	}

	result, err := client.DescribeRepositories(context.TODO(), input)
	if err != nil {
		return "", fmt.Errorf("failed to get repository details: %w", err)
	}

	if len(result.Repositories) == 0 {
		return "", fmt.Errorf("repository not found: %s", repoName)
	}

	// Get image details
	imagesInput := &ecr.DescribeImagesInput{
		RepositoryName: aws.String(repoName),
		MaxResults:     aws.Int32(100),
	}

	imagesResult, err := client.DescribeImages(context.TODO(), imagesInput)
	if err != nil {
		return "", fmt.Errorf("failed to get image details: %w", err)
	}

	// Combine repository and image details
	details := struct {
		Repository *ecr.DescribeRepositoriesOutput
		Images     *ecr.DescribeImagesOutput
	}{
		Repository: result,
		Images:     imagesResult,
	}

	// Marshal to JSON
	jsonBytes, err := json.MarshalIndent(details, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal repository details: %w", err)
	}

	return string(jsonBytes), nil
}
