package provider

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// AWSProvider fetches live resource attributes from AWS.
type AWSProvider struct {
	ec2Client *ec2.Client
	s3Client  *s3.Client
}

// NewAWSProvider creates an AWSProvider using the default AWS config chain.
func NewAWSProvider(ctx context.Context) (*AWSProvider, error) {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return nil, fmt.Errorf("loading aws config: %w", err)
	}
	return &AWSProvider{
		ec2Client: ec2.NewFromConfig(cfg),
		s3Client:  s3.NewFromConfig(cfg),
	}, nil
}

// GetAttributes returns the live attributes for the given resource type and ID.
func (p *AWSProvider) GetAttributes(ctx context.Context, resourceType, resourceID string) (map[string]string, error) {
	switch resourceType {
	case "aws_instance":
		return p.getEC2InstanceAttributes(ctx, resourceID)
	case "aws_s3_bucket":
		return p.getS3BucketAttributes(ctx, resourceID)
	default:
		return nil, fmt.Errorf("unsupported resource type: %s", resourceType)
	}
}

func (p *AWSProvider) getEC2InstanceAttributes(ctx context.Context, instanceID string) (map[string]string, error) {
	out, err := p.ec2Client.DescribeInstances(ctx, &ec2.DescribeInstancesInput{
		InstanceIds: []string{instanceID},
	})
	if err != nil {
		return nil, fmt.Errorf("describe instance %s: %w", instanceID, err)
	}
	if len(out.Reservations) == 0 || len(out.Reservations[0].Instances) == 0 {
		return nil, fmt.Errorf("instance %s not found", instanceID)
	}
	inst := out.Reservations[0].Instances[0]
	attrs := map[string]string{
		"instance_type": string(inst.InstanceType),
		"instance_state": string(inst.State.Name),
	}
	if inst.PublicIpAddress != nil {
		attrs["public_ip"] = *inst.PublicIpAddress
	}
	if inst.PrivateIpAddress != nil {
		attrs["private_ip"] = *inst.PrivateIpAddress
	}
	return attrs, nil
}

func (p *AWSProvider) getS3BucketAttributes(ctx context.Context, bucketName string) (map[string]string, error) {
	_, err := p.s3Client.HeadBucket(ctx, &s3.HeadBucketInput{
		Bucket: &bucketName,
	})
	if err != nil {
		return nil, fmt.Errorf("head bucket %s: %w", bucketName, err)
	}
	return map[string]string{
		"bucket": bucketName,
		"exists":  "true",
	}, nil
}
