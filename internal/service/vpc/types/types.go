package types

type GetVPCsInput struct{}

type GetNACLsInput struct {
	NACLIDs []string
}

type GetNatGatewaysInput struct {
	NatGatewayIDs []string
}

type GetPrefixListsInput struct {
	PrefixListIDs []string
}

type GetRouteTablesInput struct {
	RouteTableIDs []string
}

type GetSubnetsInput struct {
	SubnetIDs []string
}

type GetIGWsInput struct {
	IGWIDs []string // List of Internet Gateway IDs to filter
}
