package ecs

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecs/types"
	"github.com/harleymckenzie/asc/internal/shared/format"
)

// Task Definition field getters (for show)
var taskDefinitionFieldValueGetters = map[string]FieldValueGetter{
	"Family":               getTaskDefFamily,
	"Revision":             getTaskDefRevision,
	"ARN":                  getTaskDefARN,
	"Status":               getTaskDefStatus,
	"Network Mode":         getTaskDefNetworkMode,
	"Requires Compatibilities": getTaskDefRequiresCompatibilities,
	"vCPU":                 getTaskDefCPU,
	"Memory":               getTaskDefMemory,
	"Task Role ARN":        getTaskDefTaskRoleARN,
	"Execution Role ARN":   getTaskDefExecutionRoleARN,
	"Containers":           getTaskDefContainers,
	"Registered At":        getTaskDefRegisteredAt,
}

func getTaskDefinitionFieldValue(fieldName string, td types.TaskDefinition) (string, error) {
	if getter, exists := taskDefinitionFieldValueGetters[fieldName]; exists {
		return getter(td)
	}
	return "", fmt.Errorf("field %s not found in taskDefinitionFieldValueGetters", fieldName)
}

func getTaskDefFamily(instance any) (string, error) {
	return aws.ToString(instance.(types.TaskDefinition).Family), nil
}

func getTaskDefRevision(instance any) (string, error) {
	return fmt.Sprintf("%d", instance.(types.TaskDefinition).Revision), nil
}

func getTaskDefARN(instance any) (string, error) {
	return aws.ToString(instance.(types.TaskDefinition).TaskDefinitionArn), nil
}

func getTaskDefStatus(instance any) (string, error) {
	return format.Status(string(instance.(types.TaskDefinition).Status)), nil
}

func getTaskDefNetworkMode(instance any) (string, error) {
	return string(instance.(types.TaskDefinition).NetworkMode), nil
}

func getTaskDefRequiresCompatibilities(instance any) (string, error) {
	td := instance.(types.TaskDefinition)
	if len(td.RequiresCompatibilities) == 0 {
		return "-", nil
	}
	var compat []string
	for _, c := range td.RequiresCompatibilities {
		compat = append(compat, string(c))
	}
	return strings.Join(compat, ", "), nil
}

func getTaskDefCPU(instance any) (string, error) {
	return aws.ToString(instance.(types.TaskDefinition).Cpu), nil
}

func getTaskDefMemory(instance any) (string, error) {
	return aws.ToString(instance.(types.TaskDefinition).Memory), nil
}

func getTaskDefTaskRoleARN(instance any) (string, error) {
	return format.StringOrDefault(instance.(types.TaskDefinition).TaskRoleArn, "-"), nil
}

func getTaskDefExecutionRoleARN(instance any) (string, error) {
	return format.StringOrDefault(instance.(types.TaskDefinition).ExecutionRoleArn, "-"), nil
}

func getTaskDefContainers(instance any) (string, error) {
	td := instance.(types.TaskDefinition)
	if len(td.ContainerDefinitions) == 0 {
		return "-", nil
	}
	var names []string
	for _, c := range td.ContainerDefinitions {
		names = append(names, aws.ToString(c.Name))
	}
	return strings.Join(names, ", "), nil
}

func getTaskDefRegisteredAt(instance any) (string, error) {
	return format.TimeToStringOrEmpty(instance.(types.TaskDefinition).RegisteredAt), nil
}

// Task Definition Family field getters (for ls - families)
var taskDefinitionFamilyFieldValueGetters = map[string]FieldValueGetter{
	"Family": getTaskDefFamilyName,
}

func getTaskDefinitionFamilyFieldValue(fieldName string, family TaskDefinitionFamily) (string, error) {
	if getter, exists := taskDefinitionFamilyFieldValueGetters[fieldName]; exists {
		return getter(family)
	}
	return "", fmt.Errorf("field %s not found in taskDefinitionFamilyFieldValueGetters", fieldName)
}

func getTaskDefFamilyName(instance any) (string, error) {
	return instance.(TaskDefinitionFamily).Name, nil
}

// Task Definition Revision field getters (for ls <family>)
var taskDefinitionRevisionFieldValueGetters = map[string]FieldValueGetter{
	"ARN":      getTaskDefRevisionARN,
	"Family":   getTaskDefRevisionFamily,
	"Revision": getTaskDefRevisionRevision,
}

func getTaskDefinitionRevisionFieldValue(fieldName string, rev TaskDefinitionRevision) (string, error) {
	if getter, exists := taskDefinitionRevisionFieldValueGetters[fieldName]; exists {
		return getter(rev)
	}
	return "", fmt.Errorf("field %s not found in taskDefinitionRevisionFieldValueGetters", fieldName)
}

func getTaskDefRevisionARN(instance any) (string, error) {
	return instance.(TaskDefinitionRevision).ARN, nil
}

func getTaskDefRevisionFamily(instance any) (string, error) {
	return instance.(TaskDefinitionRevision).Family, nil
}

func getTaskDefRevisionRevision(instance any) (string, error) {
	return instance.(TaskDefinitionRevision).Revision, nil
}
