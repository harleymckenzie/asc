package types

type GetVPCsInput struct{}

type GetNACLsInput struct {
	NetworkAclIds []string
}

type GetNatGatewaysInput struct {
	NatGatewayIds []string
}

type GetPrefixListsInput struct {
	PrefixListIds []string
}

type GetManagedPrefixListsInput struct {
	PrefixListIds []string
}

type GetRouteTablesInput struct {
	RouteTableIds []string
}

type GetSubnetsInput struct {
	SubnetIds []string
	VPCIds    []string
}

type GetIGWsInput struct {
	IGWIds []string
}
