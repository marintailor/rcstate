package gce

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	compute "cloud.google.com/go/compute/apiv1"
	"cloud.google.com/go/compute/apiv1/computepb"
	"google.golang.org/api/iterator"
)

// Instance stores details of an instance.
type Instance struct {
	Name        string `json:"name"`
	Status      string `json:"status"`
	Internal    string `json:"internal"`
	External    string `json:"external"`
	Type        string `json:"type"`
	Preemptible bool   `json:"preemptible"`
}

// Instances stores list of instances in specific project and zone.
type Instances struct {
	List    []Instance
	Project string
	Zone    string
}

// NewInstances returns an Instances struct with provided project and zone.
func NewInstances(project string, zone string) *Instances {
	return &Instances{
		Project: project,
		Zone:    zone,
	}
}

// GetInstancesList returns a JSON formatted string with instances.
func (i *Instances) GetInstancesList() (string, error) {
	ctx := context.Background()
	instancesClient, err := compute.NewInstancesRESTClient(ctx)

	if err != nil {
		return "", fmt.Errorf("NewInstancesRESTClient: %w", err)
	}
	defer instancesClient.Close()

	req := &computepb.ListInstancesRequest{
		Project: i.Project,
		Zone:    i.Zone,
	}

	it := instancesClient.List(ctx, req)
	for {
		inst, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return "", fmt.Errorf("iterate instances: %w", err)
		}

		d := getInstanceDetails(inst)
		i.addInstance(d)
	}

	j, err := json.Marshal(i.List)
	if err != nil {
		return "", fmt.Errorf("marshal instances list: %w", err)
	}

	return string(j), nil
}

// GetList returns a slice of Instance.
func (i *Instances) GetList() ([]Instance, error) {
	ctx := context.Background()
	instancesClient, err := compute.NewInstancesRESTClient(ctx)

	if err != nil {
		return []Instance{}, fmt.Errorf("NewInstancesRESTClient: %w", err)
	}
	defer instancesClient.Close()

	req := &computepb.ListInstancesRequest{
		Project: i.Project,
		Zone:    i.Zone,
	}

	it := instancesClient.List(ctx, req)
	for {
		inst, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return []Instance{}, fmt.Errorf("iterate instances: %w", err)
		}

		d := getInstanceDetails(inst)
		i.addInstance(d)
	}

	return i.List, nil
}

// getInstanceDetails returns a Instance struct with instance's details.
func getInstanceDetails(inst *computepb.Instance) Instance {
	network := inst.GetNetworkInterfaces()
	schedule := inst.GetScheduling()

	externalIP := ""
	if inst.GetStatus() == "RUNNING" {
		externalIP = *network[0].AccessConfigs[0].NatIP
	}

	vmType := strings.Split(inst.GetMachineType(), "/")[len(strings.Split(inst.GetMachineType(), "/"))-1]

	return Instance{
		Name:        inst.GetName(),
		Status:      inst.GetStatus(),
		Internal:    *network[0].NetworkIP,
		External:    externalIP,
		Type:        vmType,
		Preemptible: *schedule.Preemptible,
	}
}

// addInstance will add an instance to the instances list.
func (i *Instances) addInstance(inst Instance) {
	i.List = append(i.List, inst)
}