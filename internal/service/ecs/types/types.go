package types

type ListClustersInput struct {
}

type DescribeClustersInput struct {
	ClusterARNs []string
}

type ListServicesInput struct {
	Cluster string
}

type DescribeServicesInput struct {
	Cluster  string
	Services []string
}

type ListTasksInput struct {
	Cluster     string
	ServiceName string
}

type DescribeTasksInput struct {
	Cluster string
	Tasks   []string
}

type ListTaskDefinitionFamiliesInput struct {
}

type ListTaskDefinitionRevisionsInput struct {
	FamilyName string
}

type DescribeTaskDefinitionInput struct {
	TaskDefinition string
}
