package services

import (
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
)

// EC2Instance represents simplified EC2 instance information
type EC2Instance struct {
	ID        string
	Name      string
	Type      string
	State     string
	PrivateIP string
	PublicIP  string
}

// ECRRepository represents simplified ECR repository information
type ECRRepository struct {
	Name       string
	URI        string
	ImageCount int64
	CreatedAt  time.Time
}

// LambdaFunction represents simplified Lambda function information
type LambdaFunction struct {
	Name         string
	Runtime      string
	MemorySize   int64
	LastModified time.Time
}

// Secret represents simplified Secrets Manager secret information
type Secret struct {
	Name              string
	LastModified      time.Time
	DaysUntilRotation int64
}

// InstanceInfo represents simplified EC2 instance information
type InstanceInfo struct {
	ID         string
	Name       string
	Type       string
	State      string
	PublicIP   string
	PrivateIP  string
	LaunchTime time.Time
	Tags       map[string]string
}

// RepoInfo represents simplified ECR repository information
type RepoInfo struct {
	Name       string
	URI        string
	ImageCount int64
	CreatedAt  time.Time
	LastPushAt *time.Time
	Tags       map[string]string
}

// FunctionInfo represents simplified Lambda function information
type FunctionInfo struct {
	Name        string
	Runtime     string
	Memory      int64
	Timeout     int64
	LastMod     time.Time
	State       string
	Description string
}

// SecretInfo represents simplified Secrets Manager secret information
type SecretInfo struct {
	Name        string
	ARN         string
	Description string
	LastChanged time.Time
	Tags        map[string]string
}

// GetAWSConfig returns the AWS config as aws.Config
func GetAWSConfig(cfg config.Config) interface{} {
	return cfg
}
