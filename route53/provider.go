package route53

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/route53"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

type MySQLConfiguration struct {
	Route53 *route53.Route53
}

func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"access_key": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},

			"secret_key": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},

			"region": {
				Type:     schema.TypeString,
				Required: true,
				DefaultFunc: schema.MultiEnvDefaultFunc([]string{
					"AWS_REGION",
					"AWS_DEFAULT_REGION",
				}, nil),
			},
		},

		DataSourcesMap: map[string]*schema.Resource{
			"route53_record": dataSourceRoute53Record(),
		},

		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	access_key := d.Get("access_key").(string)
	secret_key := d.Get("secret_key").(string)
	region := d.Get("region").(string)
	config := aws.NewConfig().WithRegion(region)

	if access_key != "" && secret_key != "" {
		creds := credentials.NewStaticCredentials(access_key, secret_key, "")
		config = config.WithCredentials(creds)
	} else {
		config = config.WithCredentialsChainVerboseErrors(true)
	}

	sess, err := session.NewSession()

	if err != nil {
		return nil, err
	}

	svc := route53.New(sess, config)

	return svc, nil
}
