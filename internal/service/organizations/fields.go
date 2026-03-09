package organizations

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	ascTypes "github.com/harleymckenzie/asc/internal/service/organizations/types"
	"github.com/harleymckenzie/asc/internal/shared/format"
)

// accountFieldValueGetters maps field names to functions that extract the field value from an AccountWithOU.
var accountFieldValueGetters = map[string]func(ascTypes.AccountWithOU) (string, error){
	"OU": func(a ascTypes.AccountWithOU) (string, error) {
		return a.OUName, nil
	},
	"OU Path": func(a ascTypes.AccountWithOU) (string, error) {
		return a.OUPath, nil
	},
	"ID": func(a ascTypes.AccountWithOU) (string, error) {
		return aws.ToString(a.Id), nil
	},
	"Name": func(a ascTypes.AccountWithOU) (string, error) {
		return aws.ToString(a.Name), nil
	},
	"Email": func(a ascTypes.AccountWithOU) (string, error) {
		return aws.ToString(a.Email), nil
	},
	"Status": func(a ascTypes.AccountWithOU) (string, error) {
		status := string(a.Status)
		return format.Status(strings.ToLower(status)), nil
	},
	"Joined Method": func(a ascTypes.AccountWithOU) (string, error) {
		return string(a.JoinedMethod), nil
	},
	"Joined": func(a ascTypes.AccountWithOU) (string, error) {
		if a.JoinedTimestamp != nil {
			return a.JoinedTimestamp.Format("2006-01-02 15:04:05"), nil
		}
		return "", nil
	},
}

// GetFieldValue returns the value of a field for a given AccountWithOU.
func GetFieldValue(fieldName string, instance any) (string, error) {
	account, ok := instance.(ascTypes.AccountWithOU)
	if !ok {
		return "", fmt.Errorf("expected types.AccountWithOU, got %T", instance)
	}
	if getter, exists := accountFieldValueGetters[fieldName]; exists {
		return getter(account)
	}
	return "", fmt.Errorf("unknown field: %s", fieldName)
}
