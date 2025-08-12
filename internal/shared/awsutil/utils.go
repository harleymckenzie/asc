package awsutil

import (
	"fmt"
	"reflect"

	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/harleymckenzie/asc/internal/shared/tableformat"
)

// NormalizeTags turns the provided tags into a map
func NormalizeTags(tags any) ([]tableformat.Tag, error) {
	var result []tableformat.Tag

	// EC2
	if t, ok := tags.([]types.Tag); ok {
		for _, tag := range t {
			if tag.Key != nil && tag.Value != nil {
				result = append(result, tableformat.Tag{
					Name:  *tag.Key,
					Value: *tag.Value,
				})
			}
		}
	} else {
		return nil, fmt.Errorf("provided tag type %s is currently not supported", reflect.TypeOf(tags))
	}
	return result, nil
}
