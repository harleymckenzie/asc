package profile

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Profile represents an AWS CLI profile or sso-session parsed from config/credentials files.
type Profile struct {
	Name                  string
	Type                  string // "profile" or "sso-session"
	Region                string
	Output                string
	SSOSession            string
	SSOStartURL           string
	SSORegion             string
	SSOAccountID          string
	SSORoleName           string
	SSORegistrationScopes string
	SourceProfile         string
	RoleARN               string
}

// ListProfilesOptions controls which entries are included in the result.
type ListProfilesOptions struct {
	IncludeSSOSessions bool
}

// ListProfiles reads and merges profiles from ~/.aws/config and ~/.aws/credentials.
func ListProfiles(opts ListProfilesOptions) ([]Profile, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("get home directory: %w", err)
	}

	configPath := filepath.Join(homeDir, ".aws", "config")
	credentialsPath := filepath.Join(homeDir, ".aws", "credentials")

	profiles := make(map[string]*Profile)

	if err := parseConfigFile(configPath, profiles, true); err != nil {
		return nil, fmt.Errorf("parse config file: %w", err)
	}

	if err := parseConfigFile(credentialsPath, profiles, false); err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("parse credentials file: %w", err)
	}

	result := make([]Profile, 0, len(profiles))
	for _, p := range profiles {
		if p.Type == "sso-session" && !opts.IncludeSSOSessions {
			continue
		}
		result = append(result, *p)
	}
	return result, nil
}

// parseConfigFile parses an INI-style AWS config or credentials file.
// isConfig indicates whether the file uses [profile name] prefix (config) or plain [name] (credentials).
func parseConfigFile(path string, profiles map[string]*Profile, isConfig bool) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	var currentKey string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if line == "" || strings.HasPrefix(line, "#") || strings.HasPrefix(line, ";") {
			continue
		}

		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			section := line[1 : len(line)-1]
			// Skip non-profile/non-sso-session sections
			if strings.HasPrefix(section, "granted_registry") {
				currentKey = ""
				continue
			}
			entryType := "profile"
			if strings.HasPrefix(section, "sso-session ") {
				section = strings.TrimPrefix(section, "sso-session ")
				entryType = "sso-session"
			} else if isConfig {
				section = strings.TrimPrefix(section, "profile ")
			}
			currentKey = entryType + ":" + section
			if _, ok := profiles[currentKey]; !ok {
				profiles[currentKey] = &Profile{Name: section, Type: entryType}
			}
			continue
		}

		if currentKey == "" {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		p := profiles[currentKey]
		switch key {
		case "region":
			p.Region = value
		case "output":
			p.Output = value
		case "sso_session":
			p.SSOSession = value
		case "sso_start_url":
			p.SSOStartURL = value
		case "sso_region":
			p.SSORegion = value
		case "sso_account_id":
			p.SSOAccountID = value
		case "sso_role_name":
			p.SSORoleName = value
		case "sso_registration_scopes":
			p.SSORegistrationScopes = value
		case "source_profile":
			p.SourceProfile = value
		case "role_arn":
			p.RoleARN = value
		}
	}
	return scanner.Err()
}
