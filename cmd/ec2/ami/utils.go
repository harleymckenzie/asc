package ami

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/ec2/types"

	"github.com/harleymckenzie/asc/internal/service/ec2"
	ascTypes "github.com/harleymckenzie/asc/internal/service/ec2/types"
)

// getImages retrieves EC2 images (AMIs) based on the provided input parameters
func getImages(ctx context.Context, svc *ec2.EC2Service, input *ascTypes.GetImagesInput) ([]types.Image, error) {
	return svc.GetImages(ctx, input)
}
