package sumologicextra

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/go-http-utils/headers"
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
			"use_existing": {
				Type:             schema.TypeBool,
				Optional:         true,
				Default:          false,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool { return d.Id() != "" },
			},
		},
	}
}

func fillCollector(c *client.CollectorResponse, d *schema.ResourceData) {
	c.Collector.Name = d.Get("name").(string)
	c.Collector.UseExisting = d.Get("use_existing").(bool)
	c.Collector.CollectorType = "Hosted"
	c.Collector.Ephemeral = false
}

func fillResourceDataFromCollector(c *client.CollectorResponse, d *schema.ResourceData) {
	d.Set("name", c.Collector.Name)
	d.Set("use_existing", c.Collector.UseExisting)
}

func resourceCollectorCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*client.Client)
	newCollector := client.CollectorResponse{}
	fillCollector(&newCollector, d)
	var body *bytes.Buffer = nil
	var err error
	if newCollector.Collector.UseExisting {
		requestPath := fmt.Sprintf(client.CollectorPathNameGet, url.PathEscape(newCollector.Collector.Name))
		body, _, err = c.HttpRequest(ctx, http.MethodGet, requestPath, nil, nil, &bytes.Buffer{})
		if err != nil {
			re := err.(*client.RequestError)
			if re.StatusCode != http.StatusNotFound {
				return diag.FromErr(err)
			}
			body = nil
		}
	}
	if body == nil {
		buf := bytes.Buffer{}
		err := json.NewEncoder(&buf).Encode(newCollector)
		if err != nil {
			d.SetId("")
			return diag.FromErr(err)
		}
		requestHeaders := http.Header{
			headers.ContentType: []string{client.ApplicationJson},
		}
		body, _, err = c.HttpRequest(ctx, http.MethodPost, client.CollectorPath, nil, requestHeaders, &buf)
		if err != nil {
			d.SetId("")
			return diag.FromErr(err)
		}
	}
	retVal := &client.CollectorResponse{}
	err = json.NewDecoder(body).Decode(retVal)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	fillResourceDataFromCollector(retVal, d)
	d.SetId(strconv.Itoa(retVal.Collector.ID))
	return diags
}

func resourceCollectorRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*client.Client)
	requestPath := fmt.Sprintf(client.CollectorPathGet, d.Id())
	body, _, err := c.HttpRequest(ctx, http.MethodGet, requestPath, nil, nil, &bytes.Buffer{})
	if err != nil {
		d.SetId("")
		re, ok := err.(*client.RequestError)
		if ok && re.StatusCode == http.StatusNotFound {
			return diags
		}
		return diag.FromErr(err)
	}
	retVal := &client.CollectorResponse{}
	err = json.NewDecoder(body).Decode(retVal)
	if err != nil {
		d.SetId("")
		return diag.FromErr(err)
	}
	fillResourceDataFromCollector(retVal, d)
	return diags
}

func resourceCollectorUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*client.Client)
	if !d.HasChange("name") {
		return diags
	}
	requestPath := fmt.Sprintf(client.CollectorPathGet, d.Id())
	_, respHeaders, err := c.HttpRequest(ctx, http.MethodGet, requestPath, nil, nil, &bytes.Buffer{})
	if err != nil {
		return diag.FromErr(err)
	}
	etag := respHeaders.Get(headers.ETag)
	if etag == "" {
		return diag.Errorf("SumoLogic API did not return ETag for %s", requestPath)
	}
	upd := client.CollectorResponse{}
	fillCollector(&upd, d)
	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}
	upd.Collector.ID = id
	buf := bytes.Buffer{}
	err = json.NewEncoder(&buf).Encode(upd)
	if err != nil {
		return diag.FromErr(err)
	}
	requestHeaders := http.Header{
		headers.ContentType: []string{client.ApplicationJson},
		headers.IfMatch:     []string{etag},
	}
	body, _, err := c.HttpRequest(ctx, http.MethodPut, requestPath, nil, requestHeaders, &buf)
	if err != nil {
		return diag.FromErr(err)
	}
	retVal := &client.CollectorResponse{}
	err = json.NewDecoder(body).Decode(retVal)
	if err != nil {
		return diag.FromErr(err)
	}
	fillResourceDataFromCollector(retVal, d)
	return diags
}

func resourceCollectorDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*client.Client)
	requestPath := fmt.Sprintf(client.CollectorPathGet, d.Id())
	_, _, err := c.HttpRequest(ctx, http.MethodDelete, requestPath, nil, nil, &bytes.Buffer{})
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	return diags
}
