package efs

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/efs"
	"github.com/aws/aws-sdk-go-v2/service/efs/types"

	"github.com/harleymckenzie/asc/internal/shared/awsutil"
)

type EFSClientAPI interface {
	DescribeFileSystems(context.Context, *efs.DescribeFileSystemsInput, ...func(*efs.Options)) (*efs.DescribeFileSystemsOutput, error)
}

type EFSService struct {
	Client EFSClientAPI
}

func NewEFSService(ctx context.Context, profile string, region string) (*EFSService, error) {
	cfg, err := awsutil.LoadDefaultConfig(ctx, profile, region)
	if err != nil {
		return nil, err
	}

	client := efs.NewFromConfig(cfg.Config)
	return &EFSService{Client: client}, nil
}

func (svc *EFSService) GetFileSystems(ctx context.Context) ([]types.FileSystemDescription, error) {
	output, err := svc.Client.DescribeFileSystems(ctx, &efs.DescribeFileSystemsInput{})
	if err != nil {
		return nil, err
	}

	var fileSystems []types.FileSystemDescription
	fileSystems = append(fileSystems, output.FileSystems...)
	return fileSystems, nil
}

func (svc *EFSService) GetFileSystem(ctx context.Context, identifier string) (types.FileSystemDescription, error) {
	if strings.HasPrefix(identifier, "fs-") {
		return svc.getFileSystemByID(ctx, identifier)
	}
	return svc.getFileSystemByName(ctx, identifier)
}

func (svc *EFSService) getFileSystemByID(ctx context.Context, fileSystemID string) (types.FileSystemDescription, error) {
	output, err := svc.Client.DescribeFileSystems(ctx, &efs.DescribeFileSystemsInput{
		FileSystemId: aws.String(fileSystemID),
	})
	if err != nil {
		return types.FileSystemDescription{}, err
	}

	if len(output.FileSystems) == 0 {
		return types.FileSystemDescription{}, fmt.Errorf("file system %s not found", fileSystemID)
	}
	return output.FileSystems[0], nil
}

func (svc *EFSService) getFileSystemByName(ctx context.Context, name string) (types.FileSystemDescription, error) {
	fileSystems, err := svc.GetFileSystems(ctx)
	if err != nil {
		return types.FileSystemDescription{}, err
	}

	var matches []types.FileSystemDescription
	for _, fs := range fileSystems {
		if aws.ToString(fs.Name) == name {
			matches = append(matches, fs)
		}
	}

	switch len(matches) {
	case 0:
		return types.FileSystemDescription{}, fmt.Errorf("no file system found with name %q", name)
	case 1:
		return matches[0], nil
	default:
		return types.FileSystemDescription{}, fmt.Errorf("multiple file systems found with name %q, use a file system ID (fs-...) instead", name)
	}
}
