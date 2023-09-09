package gce

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	compute "cloud.google.com/go/compute/apiv1"
	"cloud.google.com/go/compute/apiv1/computepb"
	"golang.org/x/oauth2"
	auth "golang.org/x/oauth2/google"
	compute_engine "google.golang.org/api/compute/v1"
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

// Start will start an instance.
func (i *Instances) Start(inst string) error {
	ctx := context.Background()
	instancesClient, err := compute.NewInstancesRESTClient(ctx)
	if err != nil {
		return fmt.Errorf("NewInstancesRESTClient: %w", err)
	}
	defer instancesClient.Close()

	req := &computepb.StartInstanceRequest{
		Project:  i.Project,
		Zone:     i.Zone,
		Instance: inst,
	}

	op, err := instancesClient.Start(ctx, req)
	if err != nil {
		return fmt.Errorf("start instance: %w", err)
	}

	if err = op.Wait(ctx); err != nil {
		return fmt.Errorf("wait operation: %w", err)
	}

	return nil
}

// Status returns the status of the instance.
func (i *Instances) Status(inst string) (string, error) {
	ctx := context.Background()
	computeService, err := compute_engine.NewService(ctx)
	if err != nil {
		return "", fmt.Errorf("get status instance %q: %w", inst, err)
	}

	resp, err := computeService.Instances.Get(i.Project, i.Zone, inst).Context(ctx).Do()
	if err != nil {
		return "", fmt.Errorf("get status instance %q: %w", inst, err)
	}

	return resp.Status, nil
}

// Stop will stop the instance.
func (i *Instances) Stop(inst string) error {
	ctx := context.Background()
	instancesClient, err := compute.NewInstancesRESTClient(ctx)
	if err != nil {
		return fmt.Errorf("NewInstancesRESTClient: %w", err)
	}
	defer instancesClient.Close()

	req := &computepb.StopInstanceRequest{
		Project:  i.Project,
		Zone:     i.Zone,
		Instance: inst,
	}

	op, err := instancesClient.Stop(ctx, req)
	if err != nil {
		return fmt.Errorf("stop instance: %w", err)
	}

	if err = op.Wait(ctx); err != nil {
		return fmt.Errorf("wait operation: %w", err)
	}

	return nil
}

// GetInstanceExternalIP returns the external IP address of the instance in specific project and zone.
func GetInstanceExternalIP(inst string, project string, zone string) (string, error) {
	url := fmt.Sprintf("https://compute.googleapis.com/compute/v1/projects/%s/zones/%s/instances/%s", project, zone, inst)
	token := getToken()

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("new request: %w", err)
	}
	req.Header.Add("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("client response: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read response body: %w", err)
	}

	return getNATIP(body), nil
}

// getToken returns an API token.
func getToken() string {
	var token *oauth2.Token
	ctx := context.Background()
	scopes := []string{
		"https://www.googleapis.com/auth/cloud-platform",
	}
	credentials, err := auth.FindDefaultCredentials(ctx, scopes...)
	if err == nil {
		token, err = credentials.TokenSource.Token()
		if err != nil {
			fmt.Print(err)
		}
	}

	return token.AccessToken
}

type InstanceDetails struct {
	NIC []Interface `json:"networkInterfaces"`
}

type Interface struct {
	AccessConfig []Config `json:"accessConfigs"`
}

type Config struct {
	NATIP string `json:"natIP"`
}

// getNATIP returns the NAT IP address.
func getNATIP(b []byte) string {
	var data InstanceDetails
	if err := json.Unmarshal(b, &data); err != nil {
		fmt.Println("unmarshal instance details:", err)
		return ""
	}
	return data.NIC[0].AccessConfig[0].NATIP
}
