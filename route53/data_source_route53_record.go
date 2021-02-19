package route53

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/route53"
	"github.com/aws/aws-sdk-go-v2/service/route53/types"
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
				Type:     schema.TypeInt,
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
	svc := meta.(*route53.Client)
	hostedZoneId := d.Get("hosted_zone_id").(string)
	recordName := d.Get("name").(string)
	recordType := d.Get("type").(string)

	input := &route53.ListResourceRecordSetsInput{
		HostedZoneId:    aws.String(hostedZoneId),
		StartRecordName: aws.String(recordName),
		StartRecordType: types.RRType(recordType),
		MaxItems:        aws.Int32(1),
	}

	output, err := svc.ListResourceRecordSets(context.TODO(), input)

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
		records = append(records, aws.ToString(v.Value))
	}

	d.SetId(resource.UniqueId())
	d.Set("name", aws.ToString(rrSet.Name))
	d.Set("type", string(rrSet.Type))
	d.Set("ttl", *rrSet.TTL)
	d.Set("records", records)

	return nil
}
