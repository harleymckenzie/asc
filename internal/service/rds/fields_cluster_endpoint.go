package rds

import (
	"sort"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/rds/types"
)

// ClusterEndpointRow is a synthetic list row that surfaces an Aurora cluster's
// writer or reader connection endpoint alongside instance rows in `rds ls -e`.
type ClusterEndpointRow struct {
	ClusterIdentifier string
	Role              string // "Writer" or "Reader"
	Engine            string
	Endpoint          string
}

// getClusterEndpointRowFieldValue returns the value of a list field for a synthetic
// cluster endpoint row. Columns that do not apply to an endpoint return an empty string.
func getClusterEndpointRowFieldValue(fieldName string, row ClusterEndpointRow) (string, error) {
	switch fieldName {
	case "Cluster Identifier":
		return row.ClusterIdentifier, nil
	case "Identifier":
		return "▸ " + row.Role + " endpoint", nil
	case "Role":
		return row.Role, nil
	case "Engine":
		return row.Engine, nil
	case "Endpoint":
		return row.Endpoint, nil
	default:
		return "", nil
	}
}

// BuildEndpointListData returns the rows for `rds ls -e`: instances grouped by cluster
// (ascending by cluster identifier), with each Aurora cluster's writer and reader
// endpoints injected as synthetic rows at the top of its group. Standalone instances
// follow. The second return value is false when no endpoint rows were injected, in which
// case the caller should fall back to the normal (sorted) instance listing.
func BuildEndpointListData(instances []types.DBInstance, clusters []types.DBCluster) ([]any, bool) {
	clusterByID := make(map[string]types.DBCluster, len(clusters))
	for _, c := range clusters {
		clusterByID[aws.ToString(c.DBClusterIdentifier)] = c
	}

	grouped := make(map[string][]types.DBInstance)
	var order []string
	var standalone []types.DBInstance
	for _, inst := range instances {
		cid := aws.ToString(inst.DBClusterIdentifier)
		if cid == "" {
			standalone = append(standalone, inst)
			continue
		}
		if _, seen := grouped[cid]; !seen {
			order = append(order, cid)
		}
		grouped[cid] = append(grouped[cid], inst)
	}
	sort.Strings(order)

	var data []any
	injected := false
	for _, cid := range order {
		if c, ok := clusterByID[cid]; ok {
			engine := aws.ToString(c.Engine)
			if ep := aws.ToString(c.Endpoint); ep != "" {
				data = append(data, ClusterEndpointRow{ClusterIdentifier: cid, Role: "Writer", Engine: engine, Endpoint: ep})
				injected = true
			}
			if ep := aws.ToString(c.ReaderEndpoint); ep != "" {
				data = append(data, ClusterEndpointRow{ClusterIdentifier: cid, Role: "Reader", Engine: engine, Endpoint: ep})
				injected = true
			}
		}
		for _, inst := range grouped[cid] {
			data = append(data, inst)
		}
	}
	for _, inst := range standalone {
		data = append(data, inst)
	}
	return data, injected
}
