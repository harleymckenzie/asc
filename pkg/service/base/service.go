package base

import (
	"context"

	"github.com/harleymckenzie/asc/pkg/shared/tableformat"
)

// AWSService provides common functionality for AWS services
type AWSService struct {
	TableRenderer tableformat.TableRenderer
	Config        tableformat.TableConfig
}

// NewAWSService creates a new base AWS service
func NewAWSService(config tableformat.TableConfig) *AWSService {
	return &AWSService{
		TableRenderer: tableformat.NewBaseTableRenderer(),
		Config:        config,
	}
}

// RenderResourceTable renders a table of AWS resources
func (s *AWSService) RenderResourceTable(
	ctx context.Context,
	data interface{},
	options tableformat.TableOptions,
) error {
	return s.TableRenderer.RenderTable(data, s.Config, options)
}
