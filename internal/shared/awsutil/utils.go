package awsutil

import (
	"fmt"
	"reflect"

	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/harleymckenzie/asc/internal/shared/tablewriter"
)

type Tag struct {
	Name  string
	Value string
}

// NormalizeTags turns the provided tags into a map
func NormalizeTags(tags any) ([]Tag, error) {
	var result []Tag

	// EC2
	if t, ok := tags.([]types.Tag); ok {
		for _, tag := range t {
			if tag.Key != nil && tag.Value != nil {
				result = append(result, Tag{
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

func PopulateTagFields(tags []types.Tag) ([]tablewriter.Field, error) {
	normalizedTags, err := NormalizeTags(tags)
	if err != nil {
		return nil, err
	}

	var fields []tablewriter.Field
	for _, tag := range normalizedTags {
		fields = append(fields, tablewriter.Field{
			Category: "Tag",
			Name:     tag.Name,
			Value:    tag.Value,
		})
	}
	return fields, nil
}
