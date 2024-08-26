// Generated by PMS #304
package ces

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tidwall/gjson"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/httphelper"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/schemas"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func DataSourceCesDashboards() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCesDashboardsRead,

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
				Description: `Specifies the dashboard name.`,
			},
			"dashboard_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the dashboard ID.`,
			},
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the enterprise project ID.`,
			},
			"is_favorite": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: `Specifies whether a dashboard in an enterprise project is added to favorites.`,
			},
			"dashboards": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The dashboard list.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"dashboard_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The dashboard ID.`,
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the dashboard.`,
						},
						"enterprise_project_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The enterprise project ID.`,
						},
						"creator_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The creator of the dashboard.`,
						},
						"created_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The creation time of the dashboard.`,
						},
						"row_widget_num": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The monitoring view display mode.`,
						},
						"is_favorite": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: `Whether a dashboard is added to favorites.`,
						},
					},
				},
			},
		},
	}
}

type DashboardsDSWrapper struct {
	*schemas.ResourceDataWrapper
	Config *config.Config
}

func newDashboardsDSWrapper(d *schema.ResourceData, meta interface{}) *DashboardsDSWrapper {
	return &DashboardsDSWrapper{
		ResourceDataWrapper: schemas.NewSchemaWrapper(d),
		Config:              meta.(*config.Config),
	}
}

func dataSourceCesDashboardsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	wrapper := newDashboardsDSWrapper(d, meta)
	lisDasInfRst, err := wrapper.ListDashboardInfos()
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(id)

	err = wrapper.listDashboardInfosToSchema(lisDasInfRst)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

// @API CES GET /v2/{project_id}/dashboards
func (w *DashboardsDSWrapper) ListDashboardInfos() (*gjson.Result, error) {
	client, err := w.NewClient(w.Config, "ces")
	if err != nil {
		return nil, err
	}

	uri := "/v2/{project_id}/dashboards"
	params := map[string]any{
		"enterprise_id":  w.Get("enterprise_project_id"),
		"is_favorite":    w.Get("is_favorite"),
		"dashboard_name": w.Get("name"),
		"dashboard_id":   w.Get("dashboard_id"),
	}
	params = utils.RemoveNil(params)
	return httphelper.New(client).
		Method("GET").
		URI(uri).
		Query(params).
		Request().
		Result()
}

func (w *DashboardsDSWrapper) listDashboardInfosToSchema(body *gjson.Result) error {
	d := w.ResourceData
	mErr := multierror.Append(nil,
		d.Set("region", w.Config.GetRegion(w.ResourceData)),
		d.Set("dashboards", schemas.SliceToList(body.Get("dashboards"),
			func(dashboards gjson.Result) any {
				return map[string]any{
					"dashboard_id":          dashboards.Get("dashboard_id").Value(),
					"name":                  dashboards.Get("dashboard_name").Value(),
					"enterprise_project_id": dashboards.Get("enterprise_id").Value(),
					"creator_name":          dashboards.Get("creator_name").Value(),
					"created_at":            w.setDashboardsCreateTime(dashboards),
					"row_widget_num":        dashboards.Get("row_widget_num").Value(),
					"is_favorite":           dashboards.Get("is_favorite").Value(),
				}
			},
		)),
	)
	return mErr.ErrorOrNil()
}

func (*DashboardsDSWrapper) setDashboardsCreateTime(data gjson.Result) string {
	return utils.FormatTimeStampRFC3339(int64(data.Get("create_time").Float()/1000), false)
}