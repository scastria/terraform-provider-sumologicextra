package sumologicextra

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/scastria/terraform-provider-sumologicextra/sumologicextra/client"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"num_retries": {
				Type:        schema.TypeInt,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("SUMOLOGIC_NUM_RETRIES", 3),
			},
			"retry_delay": {
				Type:        schema.TypeInt,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("SUMOLOGIC_RETRY_DELAY", 30),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"sumologicextra_collector": resourceCollector(),
		},
		DataSourcesMap:       map[string]*schema.Resource{},
		ConfigureContextFunc: providerConfigure,
	}
}

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	numRetries := d.Get("num_retries").(int)
	retryDelay := d.Get("retry_delay").(int)

	var diags diag.Diagnostics
	c, err := client.NewClient(numRetries, retryDelay)
	if err != nil {
		return nil, diag.FromErr(err)
	}
	return c, diags
}
