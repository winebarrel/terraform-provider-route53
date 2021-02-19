package route53

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/route53"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceRoute53Record() *schema.Resource {
	return &schema.Resource{
		Read: dataRoute53RecordRead,
		Schema: map[string]*schema.Schema{
			"hosted_zone_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"ttl": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"records": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataRoute53RecordRead(d *schema.ResourceData, meta interface{}) error {
	svc := meta.(*route53.Route53)
	hostedZoneId := d.Get("hosted_zone_id").(string)
	recordName := d.Get("name").(string)
	recordType := d.Get("type").(string)

	input := &route53.ListResourceRecordSetsInput{
		HostedZoneId:    aws.String(hostedZoneId),
		StartRecordName: aws.String(recordName),
		StartRecordType: aws.String(recordType),
	}

	output, err := svc.ListResourceRecordSets(input)

	if err != nil {
		return err
	}

	if len(output.ResourceRecordSets) == 0 {
		return fmt.Errorf(
			"Error listing Route53 resource record sets: Record not found (hosted_zone_id=%s, name=%s, type=%s)",
			hostedZoneId, recordName, recordType,
		)
	}

	rrSet := output.ResourceRecordSets[0]
	records := []string{}

	for _, v := range rrSet.ResourceRecords {
		records = append(records, *v.Value)
	}

	d.SetId(resource.UniqueId())
	d.Set("name", *rrSet.Name)
	d.Set("type", *rrSet.Type)
	d.Set("ttl", *rrSet.Type)
	d.Set("records", records)

	return nil
}
