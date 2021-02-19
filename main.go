package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/plugin"
	"github.com/winebarrel/terraform-provider-route53/route53"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: route53.Provider})
}
