package vpc

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/harleymckenzie/asc/internal/shared/format"
)

// FieldValueGetter is a function that returns the value of a field for a given instance.
type RouteTableFieldValueGetter func(instance any) (string, error)
type RouteFieldValueGetter func(instance any) (string, error)

// Route Table field getters
var routeTableFieldValueGetters = map[string]RouteTableFieldValueGetter{
	"Route Table ID":    getRouteTableID,
	"VPC ID":            getRouteTableVPCID,
	"Association Count": getRouteTableAssociationCount,
	"Route Count":       getRouteTableRouteCount,
	"Main":              getRouteTableIsMain,
	"Owner":             getRouteTableOwner,
}

// Route field getters
var routeFieldValueGetters = map[string]RouteFieldValueGetter{
	"Destination": getRouteDestination,
	"Target":      getRouteTarget,
	"Status":      getRouteStatus,
	"Propagated":  getRoutePropagated,
	"Origin":      getRouteOrigin,
}

// GetRouteTableFieldValue returns the value of a field for the given Route Table instance.
func GetRouteTableFieldValue(fieldName string, instance any) (string, error) {
	switch v := instance.(type) {
	case types.RouteTable:
		return getRouteTableFieldValue(fieldName, v)
	case types.Route:
		return getRouteFieldValue(fieldName, v)
	default:
		return "", fmt.Errorf("unsupported instance type: %T", instance)
	}
}

// getRouteTableFieldValue returns the value of a field for a Route Table
func getRouteTableFieldValue(fieldName string, routeTable types.RouteTable) (string, error) {
	if getter, exists := routeTableFieldValueGetters[fieldName]; exists {
		return getter(routeTable)
	}
	return "", fmt.Errorf("field %s not found in routeTableFieldValueGetters", fieldName)
}

// getRouteFieldValue returns the value of a field for a Route
func getRouteFieldValue(fieldName string, route types.Route) (string, error) {
	if getter, exists := routeFieldValueGetters[fieldName]; exists {
		return getter(route)
	}
	return "", fmt.Errorf("field %s not found in routeFieldValueGetters", fieldName)
}

// GetRouteTableTagValue returns the value of a tag for the given instance.
func GetRouteTableTagValue(tagKey string, instance any) (string, error) {
	switch v := instance.(type) {
	case types.RouteTable:
		for _, tag := range v.Tags {
			if aws.ToString(tag.Key) == tagKey {
				return aws.ToString(tag.Value), nil
			}
		}
	default:
		return "", fmt.Errorf("unsupported instance type for tags: %T", instance)
	}
	return "", nil
}

// -----------------------------------------------------------------------------
// Route Table field getters
// -----------------------------------------------------------------------------

func getRouteTableID(instance any) (string, error) {
	return format.StringOrEmpty(instance.(types.RouteTable).RouteTableId), nil
}

func getRouteTableVPCID(instance any) (string, error) {
	return format.StringOrEmpty(instance.(types.RouteTable).VpcId), nil
}

func getRouteTableAssociationCount(instance any) (string, error) {
	rt := instance.(types.RouteTable)
	return fmt.Sprintf("%d", len(rt.Associations)), nil
}

func getRouteTableRouteCount(instance any) (string, error) {
	rt := instance.(types.RouteTable)
	return fmt.Sprintf("%d", len(rt.Routes)), nil
}

func getRouteTableIsMain(instance any) (string, error) {
	rt := instance.(types.RouteTable)
	if len(rt.Associations) == 0 {
		return "No", nil
	}
	return format.BoolToLabel(rt.Associations[0].Main, "Yes", "No"), nil
}

func getRouteTableOwner(instance any) (string, error) {
	return format.StringOrEmpty(instance.(types.RouteTable).OwnerId), nil
}

// -----------------------------------------------------------------------------
// Route field getters
// -----------------------------------------------------------------------------

func getRouteDestination(instance any) (string, error) {
	route := instance.(types.Route)
	if route.DestinationCidrBlock != nil {
		return *route.DestinationCidrBlock, nil
	}
	if route.DestinationIpv6CidrBlock != nil {
		return *route.DestinationIpv6CidrBlock, nil
	}
	if route.DestinationPrefixListId != nil {
		return *route.DestinationPrefixListId, nil
	}
	return "", nil
}

func getRouteTarget(instance any) (string, error) {
	route := instance.(types.Route)
	if route.NatGatewayId != nil {
		return format.StringOrEmpty(route.NatGatewayId), nil
	}
	if route.NetworkInterfaceId != nil {
		return format.StringOrEmpty(route.NetworkInterfaceId), nil
	}
	if route.GatewayId != nil {
		return format.StringOrEmpty(route.GatewayId), nil
	}
	if route.InstanceId != nil {
		return format.StringOrEmpty(route.InstanceId), nil
	}
	if route.VpcPeeringConnectionId != nil {
		return format.StringOrEmpty(route.VpcPeeringConnectionId), nil
	}
	if route.TransitGatewayId != nil {
		return format.StringOrEmpty(route.TransitGatewayId), nil
	}
	if route.EgressOnlyInternetGatewayId != nil {
		return format.StringOrEmpty(route.EgressOnlyInternetGatewayId), nil
	}
	if route.CarrierGatewayId != nil {
		return format.StringOrEmpty(route.CarrierGatewayId), nil
	}
	if route.LocalGatewayId != nil {
		return format.StringOrEmpty(route.LocalGatewayId), nil
	}
	return "local", nil
}

func getRouteStatus(instance any) (string, error) {
	route := instance.(types.Route)
	return format.Status(string(route.State)), nil
}

func getRoutePropagated(instance any) (string, error) {
	route := instance.(types.Route)
	if route.Origin == types.RouteOriginCreateRoute {
		return "No", nil
	}
	return "Yes", nil
}

func getRouteOrigin(instance any) (string, error) {
	route := instance.(types.Route)
	return string(route.Origin), nil
}
