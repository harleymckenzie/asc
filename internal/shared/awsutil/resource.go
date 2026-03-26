package awsutil

import "strings"

// ResourceInfo describes an AWS resource type identified by its ID prefix.
type ResourceInfo struct {
	Service      string
	ResourceType string
}

var prefixMap = []struct {
	prefix string
	info   ResourceInfo
}{
	{"i-", ResourceInfo{Service: "ec2", ResourceType: "instance"}},
	{"vol-", ResourceInfo{Service: "ec2", ResourceType: "volume"}},
	{"sg-", ResourceInfo{Service: "ec2", ResourceType: "security-group"}},
	{"snap-", ResourceInfo{Service: "ec2", ResourceType: "snapshot"}},
	{"ami-", ResourceInfo{Service: "ec2", ResourceType: "image"}},
	{"eni-", ResourceInfo{Service: "ec2", ResourceType: "network-interface"}},
	{"vpc-", ResourceInfo{Service: "vpc", ResourceType: "vpc"}},
	{"subnet-", ResourceInfo{Service: "vpc", ResourceType: "subnet"}},
	{"igw-", ResourceInfo{Service: "vpc", ResourceType: "internet-gateway"}},
	{"nat-", ResourceInfo{Service: "vpc", ResourceType: "nat-gateway"}},
	{"rtb-", ResourceInfo{Service: "vpc", ResourceType: "route-table"}},
	{"acl-", ResourceInfo{Service: "vpc", ResourceType: "network-acl"}},
	{"pl-", ResourceInfo{Service: "vpc", ResourceType: "prefix-list"}},
}

// IdentifyResource returns the service and resource type for a given AWS
// resource ID based on its prefix. Returns nil for unrecognised IDs.
func IdentifyResource(id string) *ResourceInfo {
	for _, entry := range prefixMap {
		if strings.HasPrefix(id, entry.prefix) {
			return &entry.info
		}
	}
	return nil
}
