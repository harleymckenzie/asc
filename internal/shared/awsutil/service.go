package awsutil

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsmiddleware "github.com/aws/aws-sdk-go-v2/aws/middleware"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/smithy-go/middleware"
)

type BaseService struct {
	Config aws.Config
}

var Version = "dev"

func LoadDefaultConfig(ctx context.Context, profile string, region string) (*BaseService, error) {
    opts := []func(*config.LoadOptions) error{
        config.WithAPIOptions([]func(*middleware.Stack) error{
            awsmiddleware.AddUserAgentKeyValue("asc", Version),
        }),
    }

    if profile != "" {
        opts = append(opts, config.WithSharedConfigProfile(profile))
    }

    if region != "" {
        opts = append(opts, config.WithRegion(region))
    }

    cfg, err := config.LoadDefaultConfig(ctx, opts...)
    if err != nil {
		return nil, err
	}

	return &BaseService{Config: cfg}, nil
}
