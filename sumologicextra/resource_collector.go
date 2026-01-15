package sumologicextra

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/scastria/terraform-provider-sumologicextra/sumologicextra/client"
)

func resourceCollector() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCollectorCreate,
		ReadContext:   resourceCollectorRead,
		UpdateContext: resourceCollectorUpdate,
		DeleteContext: resourceCollectorDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"time_zone": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"collector_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceCollectorCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	return diags
}

func resourceCollectorRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*client.Client)
	if d.Id() == "" {
		return diags
	}
	requestPath := fmt.Sprintf("collectors/%s", d.Id())
	body, err := c.HttpRequest(ctx, http.MethodGet, requestPath, nil, nil, &bytes.Buffer{})
	if err != nil {
		d.SetId("")
		var re *client.RequestError
		ok := errors.As(err, &re)
		if ok && re.StatusCode == http.StatusNotFound {
			return diags
		}
		return diag.FromErr(err)
	}
	var resp struct {
		Collector struct {
			Name          string `json:"name"`
			TimeZone      string `json:"timeZone,omitempty"`
			CollectorType string `json:"collectorType,omitempty"`
		} `json:"collector"`
	}
	err = json.NewDecoder(body).Decode(&resp)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}

	_ = d.Set("name", resp.Collector.Name)
	_ = d.Set("time_zone", resp.Collector.TimeZone)
	_ = d.Set("collector_type", resp.Collector.CollectorType)
	return diags
}

func resourceCollectorUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	return diags
}

func resourceCollectorDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	return diags
}
