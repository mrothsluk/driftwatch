// Package config handles loading and validating the driftwatch YAML
// configuration file. It defines which Terraform state file to inspect,
// which cloud provider to query, the desired output format, and optional
// resource filters.
//
// Example configuration file (driftwatch.yaml):
//
//	statefile: terraform.tfstate
//	provider: aws
//	output: text
//	filters:
//	  resource_types:
//	    - aws_instance
//	  exclude_ids:
//	    - i-0abc123def456
package config
