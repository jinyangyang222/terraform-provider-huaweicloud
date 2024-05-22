// Generated by PMS #4
package dds

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tidwall/gjson"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/filters"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/httphelper"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/schemas"
)

func DataSourceDdsParameterTemplates() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDdsParameterTemplatesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the region in which to query the resource. If omitted, the provider-level region will be used.`,
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the parameter template name.`,
			},
			"node_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The the node type of parameter template.`,
			},
			"datastore_version": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the database (DB Engine) version.`,
			},
			"configurations": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `DDS parameter templates list.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The parameter template ID.`,
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The parameter template name.`,
						},
						"node_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The the node type of parameter template.`,
						},
						"datastore_version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Database (DB Engine) version.`,
						},
						"datastore_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Database (DB Engine) type.`,
						},
						"user_defined": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: `Whether the parameter template is a custom template.`,
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The parameter template description.`,
						},
						"updated_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The creation time of the parameter template.`,
						},
						"created_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The update time of the parameter template.`,
						},
					},
				},
			},
		},
	}
}

type ParameterTemplatesDSWrapper struct {
	*schemas.ResourceDataWrapper
	Config *config.Config
}

func newParameterTemplatesDSWrapper(d *schema.ResourceData, meta interface{}) *ParameterTemplatesDSWrapper {
	return &ParameterTemplatesDSWrapper{
		ResourceDataWrapper: schemas.NewSchemaWrapper(d),
		Config:              meta.(*config.Config),
	}
}

func dataSourceDdsParameterTemplatesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	wrapper := newParameterTemplatesDSWrapper(d, meta)
	lisConRst, err := wrapper.ListConfigurations()
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(id)

	err = wrapper.listConfigurationsToSchema(lisConRst)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

// @API DDS GET /v3/{project_id}/configurations
func (w *ParameterTemplatesDSWrapper) ListConfigurations() (*gjson.Result, error) {
	client, err := w.NewClient(w.Config, "dds")
	if err != nil {
		return nil, err
	}

	uri := "/v3/{project_id}/configurations"
	return httphelper.New(client).
		Method("GET").
		URI(uri).
		OffsetPager("configurations", "offset", "limit", 0).
		Filter(
			filters.New().From("configurations").
				Where("name", "=", w.Get("name")).
				Where("node_type", "=", w.Get("node_type")).
				Where("datastore_version", "=", w.Get("datastore_version")),
		).
		Request().
		Result()
}

func (w *ParameterTemplatesDSWrapper) listConfigurationsToSchema(body *gjson.Result) error {
	d := w.ResourceData
	mErr := multierror.Append(nil,
		d.Set("region", w.Config.GetRegion(w.ResourceData)),
		d.Set("configurations", schemas.SliceToList(body.Get("configurations"),
			func(configurations gjson.Result) any {
				return map[string]any{
					"id":                configurations.Get("id").Value(),
					"name":              configurations.Get("name").Value(),
					"node_type":         configurations.Get("node_type").Value(),
					"datastore_version": configurations.Get("datastore_version").Value(),
					"datastore_name":    configurations.Get("datastore_name").Value(),
					"user_defined":      configurations.Get("user_defined").Value(),
					"description":       configurations.Get("description").Value(),
					"updated_at":        configurations.Get("updated").Value(),
					"created_at":        configurations.Get("created").Value(),
				}
			},
		)),
	)
	return mErr.ErrorOrNil()
}