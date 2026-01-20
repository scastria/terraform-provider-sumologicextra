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
	d.Set("use_existing", d.Get("use_existing").(bool))
}

func resourceCollectorCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*client.Client)
	newCollector := client.CollectorResponse{}
	fillCollector(&newCollector, d)
	var body *bytes.Buffer = nil
	var err error
	if newCollector.Collector.UseExisting {
		const pageSize = 100
		offset := 0
		var foundID int64 = 0
		for {
			q := url.Values{}
			q.Set("limit", strconv.Itoa(pageSize))
			q.Set("offset", strconv.Itoa(offset))
			listBody, err := c.HttpRequest(ctx, http.MethodGet, client.CollectorPath, q, nil, &bytes.Buffer{})
			if err != nil {
				return diag.FromErr(err)
			}
			listRetVal := &client.CollectorsResponse{}
			err = json.NewDecoder(listBody).Decode(listRetVal)
			if err != nil {
				d.SetId("")
				return diag.FromErr(err)
			}
			if len(listRetVal.Collectors) == 0 {
				break
			}
			for _, col := range listRetVal.Collectors {
				if col.Name == newCollector.Collector.Name {
					foundID = col.ID
					break
				}
			}
			if foundID != 0 {
				break
			}
			offset += pageSize
		}
		if foundID != 0 {
			requestPath := fmt.Sprintf(client.CollectorPathGet, strconv.FormatInt(foundID, 10))
			body, err = c.HttpRequest(ctx, http.MethodGet, requestPath, nil, nil, &bytes.Buffer{})
			if err != nil {
				re := err.(*client.RequestError)
				if re.StatusCode != http.StatusNotFound {
					return diag.FromErr(err)
				}
				body = nil
			}
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
		body, err = c.HttpRequest(ctx, http.MethodPost, client.CollectorPath, nil, requestHeaders, &buf)
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
	d.SetId(strconv.FormatInt(retVal.Collector.ID, 10))
	return diags
}

func resourceCollectorRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	c := m.(*client.Client)

	if d.Id() == "" {
		return diags
	}
	requestPath := fmt.Sprintf(client.CollectorPathGet, d.Id())
	body, err := c.HttpRequest(ctx, http.MethodGet, requestPath, nil, nil, &bytes.Buffer{})
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
	if d.Id() == "" {
		return diags
	}
	if !d.HasChange("name") {
		return diags
	}
	requestPath := fmt.Sprintf(client.CollectorPathGet, d.Id())
	etag, err := c.GetEtag(requestPath)
	if err != nil {
		return diag.FromErr(err)
	}
	if etag == "" {
		return diag.Errorf("SumoLogic API did not return ETag for %s", requestPath)
	}
	upd := client.CollectorResponse{}
	fillCollector(&upd, d) // sets Name + Hosted/Ephemeral etc
	upd.Collector.ID, _ = strconv.ParseInt(d.Id(), 10, 64)
	buf := bytes.Buffer{}
	err = json.NewEncoder(&buf).Encode(upd)
	if err != nil {
		return diag.FromErr(err)
	}
	requestHeaders := http.Header{
		headers.ContentType: []string{client.ApplicationJson},
		headers.IfMatch:     []string{etag},
	}

	body, err := c.HttpRequest(ctx, http.MethodPut, requestPath, nil, requestHeaders, &buf)
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
	if d.Id() == "" {
		return diags
	}
	requestPath := fmt.Sprintf(client.CollectorPathGet, d.Id())
	_, err := c.HttpRequest(ctx, http.MethodDelete, requestPath, nil, nil, &bytes.Buffer{})
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId("")
	return diags
}
