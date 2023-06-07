// Package record implements functions to manage DNS records.
package record

import (
	"fmt"
	"net"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/route53"
)

// Record is a struct that holds required data and methods to create a DNS record.
type Record struct {
	Domain string
	IP     []string
	Type   string
	Zone   string
}

// NewRecord returns a Record struct.
func NewRecord(ip []string, rt string, zone string) *Record {
	return &Record{
		IP:     ip,
		Type:   rt,
		Zone:   zone,
		Domain: "marincroitoru.com",
	}
}

// Route53 creates a DNS record in Route 53 DNS service.
func (r *Record) Route53() error {
	if err := awsCredentials(); err != nil {
		return fmt.Errorf("aws credentials: %w", err)
	}

	sess, err := session.NewSession()
	if err != nil {
		return fmt.Errorf("new session: %w", err)
	}

	svc := route53.New(sess)

	hostedZoneNameInput := &route53.ListHostedZonesByNameInput{
		DNSName: &r.Domain,
	}

	list, err := svc.ListHostedZonesByName(hostedZoneNameInput)
	if err != nil {
		return fmt.Errorf("hosted zone by name: %w", err)
	}

	hostedZoneID := getHostedZoneID(r.Domain+".", list.HostedZones)

	hostedZoneInput := &route53.GetHostedZoneInput{
		Id: hostedZoneID,
	}

	hostedZone, err := svc.GetHostedZone(hostedZoneInput)
	if err != nil {
		return fmt.Errorf("hosted zone: %w", err)
	}

	zoneID := hostedZone.HostedZone.Id
	if err != nil {
		return fmt.Errorf("zone id: %w", err)
	}

	var ttl int64 = 10

	var ipList []*route53.ResourceRecord

	for _, v := range r.IP {
		ipList = append(ipList, &route53.ResourceRecord{Value: aws.String(v)})
	}

	params := &route53.ChangeResourceRecordSetsInput{
		ChangeBatch: &route53.ChangeBatch{
			Changes: []*route53.Change{
				{
					Action: aws.String("UPSERT"),
					ResourceRecordSet: &route53.ResourceRecordSet{
						Name:            aws.String(r.Zone),
						Type:            aws.String(r.Type),
						ResourceRecords: ipList,
						TTL:             aws.Int64(ttl),
					},
				},
			},
			Comment: aws.String("Update record to reflect new IP address for a system"),
		},
		HostedZoneId: aws.String(*zoneID),
	}

	resp, err := svc.ChangeResourceRecordSets(params)
	if err != nil {
		return fmt.Errorf("resp: %w", err)
	}

	changeID := &route53.GetChangeInput{
		Id: resp.ChangeInfo.Id,
	}

	return checkChange(svc, changeID)
}

// awsCredentials ensures that AWS credentials are present.
func awsCredentials() error {
	home, _ := os.UserHomeDir()
	path := home + "/.aws/credentials"
	if _, err := os.Stat(path); err != nil {
		return fmt.Errorf("env file stat: %w", err)
	}

	return nil
}

// getHostedZoneID returns pointer string of the Hosted Zone ID
func getHostedZoneID(d string, l []*route53.HostedZone) *string {
	for _, hz := range l {
		if *hz.Name == d {
			return hz.Id
		}
	}
	return nil
}

// checkChange checks status of change in record.
func checkChange(svc *route53.Route53, id *route53.GetChangeInput) error {
	for i := 0; i < 31; i++ {
		if i == 30 {
			return fmt.Errorf("check change status: %s", "1 minute timeout")
		}

		change, err := svc.GetChange(id)
		if err != nil {
			return fmt.Errorf("change status: %w", err)
		}

		if *change.ChangeInfo.Status == "INSYNC" {
			return nil
		}

		time.Sleep(2 * time.Second)
	}

	return nil
}

// CheckRecordIP will check is the record IP addresses are up-to-date.
func CheckRecordIP(zone string, recordIP []string) bool {
	IPs, err := net.LookupIP(zone)
	if err != nil {
		fmt.Println("lookup ip:", err)
	}

	if len(IPs) != len(recordIP) {
		return false
	}

	count := 0
	for _, ip := range IPs {
		for _, rec := range recordIP {
			if ip.String() == rec {
				count++
			}
		}
	}

	return count == len(IPs)
}
