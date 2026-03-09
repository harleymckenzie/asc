package profile

import "fmt"

// GetFieldValue returns the value of a field for the given profile.
func GetFieldValue(fieldName string, instance any) (string, error) {
	p, ok := instance.(Profile)
	if !ok {
		return "", fmt.Errorf("unsupported type: %T", instance)
	}

	switch fieldName {
	case "Name":
		return p.Name, nil
	case "Type":
		return p.Type, nil
	case "Region":
		if p.Region != "" {
			return p.Region, nil
		}
		return p.SSORegion, nil
	case "Output":
		return p.Output, nil
	case "SSO Session":
		return p.SSOSession, nil
	case "SSO Start URL":
		return p.SSOStartURL, nil
	case "SSO Account ID":
		return p.SSOAccountID, nil
	case "SSO Role Name":
		return p.SSORoleName, nil
	case "SSO Registration Scopes":
		return p.SSORegistrationScopes, nil
	case "Source Profile":
		return p.SourceProfile, nil
	case "Role ARN":
		return p.RoleARN, nil
	default:
		return "", fmt.Errorf("unknown field: %s", fieldName)
	}
}

// GetTagValue is a no-op for profiles (profiles don't have tags).
func GetTagValue(tagKey string, instance any) (string, error) {
	return "", nil
}
