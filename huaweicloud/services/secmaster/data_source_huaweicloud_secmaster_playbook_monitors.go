// Generated by PMS #345
package secmaster

import (
	"context"
	"strings"

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

func DataSourceSecmasterPlaybookMonitors() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSecmasterPlaybookMonitorsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the region in which to query the resource. If omitted, the provider-level region will be used.`,
			},
			"workspace_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the workspace ID to which the playbook belongs.`,
			},
			"playbook_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the playbook ID.`,
			},
			"start_time": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the start time.`,
			},
			"end_time": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the end time.`,
			},
			"version_query_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the playbook version type.`,
			},
			"data": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The playbook running monitor details.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"total_instance_run_num": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The total running times.`,
						},
						"average_run_time": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: `The average duration.`,
						},
						"min_run_time_instance": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: `The workflow with the shortest running duration.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"playbook_instance_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The playbook instance ID.`,
									},
									"playbook_instance_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The playbook instance name.`,
									},
									"playbook_instance_run_time": {
										Type:        schema.TypeFloat,
										Computed:    true,
										Description: `The playbook instance running time.`,
									},
								},
							},
						},
						"success_instance_num": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The number of successful instances.`,
						},
						"terminate_instance_num": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The number of terminated instances.`,
						},
						"running_instance_num": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The number of running instances.`,
						},
						"schedule_instance_run_num": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The number of scheduled trigger executions.`,
						},
						"event_instance_run_num": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The time-triggered executions.`,
						},
						"max_run_time_instance": {
							Type:        schema.TypeList,
							Computed:    true,
							Description: `The workflow with the longest running duration.`,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"playbook_instance_id": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The playbook instance ID.`,
									},
									"playbook_instance_name": {
										Type:        schema.TypeString,
										Computed:    true,
										Description: `The playbook instance name.`,
									},
									"playbook_instance_run_time": {
										Type:        schema.TypeFloat,
										Computed:    true,
										Description: `The playbook instance running time.`,
									},
								},
							},
						},
						"total_instance_num": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The total number of playbook instances.`,
						},
						"fail_instance_num": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The number of failed instances.`,
						},
					},
				},
			},
		},
	}
}

type PlaybookMonitorsDSWrapper struct {
	*schemas.ResourceDataWrapper
	Config *config.Config
}

func newPlaybookMonitorsDSWrapper(d *schema.ResourceData, meta interface{}) *PlaybookMonitorsDSWrapper {
	return &PlaybookMonitorsDSWrapper{
		ResourceDataWrapper: schemas.NewSchemaWrapper(d),
		Config:              meta.(*config.Config),
	}
}

func dataSourceSecmasterPlaybookMonitorsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	wrapper := newPlaybookMonitorsDSWrapper(d, meta)
	shoPlaMonRst, err := wrapper.ShowPlaybookMonitors()
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(id)

	err = wrapper.showPlaybookMonitorsToSchema(shoPlaMonRst)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

// @API SecMaster GET /v1/{project_id}/workspaces/{workspace_id}/soc/playbooks/{playbook_id}/monitor
func (w *PlaybookMonitorsDSWrapper) ShowPlaybookMonitors() (*gjson.Result, error) {
	client, err := w.NewClient(w.Config, "secmaster")
	if err != nil {
		return nil, err
	}

	uri := "/v1/{project_id}/workspaces/{workspace_id}/soc/playbooks/{playbook_id}/monitor"
	uri = strings.ReplaceAll(uri, "{workspace_id}", w.Get("workspace_id").(string))
	uri = strings.ReplaceAll(uri, "{playbook_id}", w.Get("playbook_id").(string))
	params := map[string]any{
		"start_time":         w.Get("start_time"),
		"version_query_type": w.Get("version_query_type"),
		"end_time":           w.Get("end_time"),
	}
	params = utils.RemoveNil(params)
	return httphelper.New(client).
		Method("GET").
		URI(uri).
		Query(params).
		Request().
		Result()
}

func (w *PlaybookMonitorsDSWrapper) showPlaybookMonitorsToSchema(body *gjson.Result) error {
	d := w.ResourceData
	mErr := multierror.Append(nil,
		d.Set("region", w.Config.GetRegion(w.ResourceData)),
		d.Set("data", schemas.ObjectToList(body.Get("data"),
			func(data gjson.Result) any {
				return map[string]any{
					"total_instance_run_num": data.Get("total_instance_run_num").Value(),
					"average_run_time":       data.Get("average_run_time").Value(),
					"min_run_time_instance": schemas.SliceToList(data.Get("min_run_time_instance"),
						func(minRunTimeInstance gjson.Result) any {
							return map[string]any{
								"playbook_instance_id":       minRunTimeInstance.Get("playbook_instance_id").Value(),
								"playbook_instance_name":     minRunTimeInstance.Get("playbook_instance_name").Value(),
								"playbook_instance_run_time": minRunTimeInstance.Get("playbook_instance_run_time").Value(),
							}
						},
					),
					"success_instance_num":      data.Get("success_instance_num").Value(),
					"terminate_instance_num":    data.Get("terminate_instance_num").Value(),
					"running_instance_num":      data.Get("running_instance_num").Value(),
					"schedule_instance_run_num": data.Get("schedule_instance_run_num").Value(),
					"event_instance_run_num":    data.Get("event_instance_run_num").Value(),
					"max_run_time_instance": schemas.SliceToList(data.Get("max_run_time_instance"),
						func(maxRunTimeInstance gjson.Result) any {
							return map[string]any{
								"playbook_instance_id":       maxRunTimeInstance.Get("playbook_instance_id").Value(),
								"playbook_instance_name":     maxRunTimeInstance.Get("playbook_instance_name").Value(),
								"playbook_instance_run_time": maxRunTimeInstance.Get("playbook_instance_run_time").Value(),
							}
						},
					),
					"total_instance_num": data.Get("total_instance_num").Value(),
					"fail_instance_num":  data.Get("fail_instance_num").Value(),
				}
			},
		)),
	)
	return mErr.ErrorOrNil()
}
