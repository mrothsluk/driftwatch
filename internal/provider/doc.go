// Package provider defines the Provider interface and built-in implementations
// for fetching live cloud resource attributes used during drift detection.
//
// Supported providers:
//   - AWSProvider: queries EC2, S3 and other AWS services via the AWS SDK v2.
//   - MockProvider: in-memory stub for use in unit tests.
//
// Usage:
//
//	prov, err := provider.NewAWSProvider(ctx)
//	if err != nil { ... }
//	attrs, err := prov.GetAttributes(ctx, "aws_instance", "i-0abc123")
package provider
