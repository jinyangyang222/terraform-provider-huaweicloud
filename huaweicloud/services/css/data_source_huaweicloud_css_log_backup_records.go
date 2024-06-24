// Generated by PMS #202
package css

import (
	"context"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tidwall/gjson"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/filters"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/httphelper"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/schemas"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func DataSourceCssLogBackupRecords() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCssLogBackupRecordsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the region in which to query the resource. If omitted, the provider-level region will be used.`,
			},
			"cluster_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the ID of the CSS cluster.`,
			},
			"job_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the ID of the log backup job.`,
			},
			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the type of the log backup job.`,
			},
			"status": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the status of the log backup job.`,
			},
			"records": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The list of the CSS cluster log backup records.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the log backup job.`,
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The type of the log backup job.`,
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The status of the log backup job.`,
						},
						"cluster_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the CSS cluster.`,
						},
						"log_path": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The storage path of backed up logs in the OBS bucket.`,
						},
						"create_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The creation time.`,
						},
						"finished_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The end time.`,
						},
						"failed_msg": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The error information.`,
						},
					},
				},
			},
		},
	}
}

type LogBackupRecordsDSWrapper struct {
	*schemas.ResourceDataWrapper
	Config *config.Config
}

func newLogBackupRecordsDSWrapper(d *schema.ResourceData, meta interface{}) *LogBackupRecordsDSWrapper {
	return &LogBackupRecordsDSWrapper{
		ResourceDataWrapper: schemas.NewSchemaWrapper(d),
		Config:              meta.(*config.Config),
	}
}

func dataSourceCssLogBackupRecordsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	wrapper := newLogBackupRecordsDSWrapper(d, meta)
	listLogsJobRst, err := wrapper.ListLogsJob()
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(id)

	err = wrapper.listLogsJobToSchema(listLogsJobRst)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

// @API CSS GET /v1.0/{project_id}/clusters/{cluster_id}/logs/records
func (w *LogBackupRecordsDSWrapper) ListLogsJob() (*gjson.Result, error) {
	client, err := w.NewClient(w.Config, "css")
	if err != nil {
		return nil, err
	}

	uri := "/v1.0/{project_id}/clusters/{cluster_id}/logs/records"
	uri = strings.ReplaceAll(uri, "{cluster_id}", w.Get("cluster_id").(string))
	return httphelper.New(client).
		Method("GET").
		URI(uri).
		OffsetStart(1).
		OffsetPager("clusterLogRecord", "start", "limit", 10).
		Filter(
			filters.New().From("clusterLogRecord").
				Where("jobId", "=", w.Get("job_id")).
				Where("jobTypes", "=", w.Get("type")).
				Where("status", "=", w.Get("status")),
		).
		Request().
		Result()
}

func (w *LogBackupRecordsDSWrapper) listLogsJobToSchema(body *gjson.Result) error {
	d := w.ResourceData
	mErr := multierror.Append(nil,
		d.Set("region", w.Config.GetRegion(w.ResourceData)),
		d.Set("records", schemas.SliceToList(body.Get("clusterLogRecord"),
			func(records gjson.Result) any {
				return map[string]any{
					"id":          records.Get("jobId").Value(),
					"type":        records.Get("jobTypes").Value(),
					"status":      records.Get("status").Value(),
					"cluster_id":  records.Get("clusterId").Value(),
					"log_path":    records.Get("logPath").Value(),
					"create_at":   w.setCluLogRecCreAt(records),
					"finished_at": w.setCluLogRecFinAt(records),
					"failed_msg":  records.Get("failedMsg").Value(),
				}
			},
		)),
	)
	return mErr.ErrorOrNil()
}

func (*LogBackupRecordsDSWrapper) setCluLogRecCreAt(data gjson.Result) string {
	return utils.FormatTimeStampRFC3339(data.Get("createAt").Int()/1000, false)
}

func (*LogBackupRecordsDSWrapper) setCluLogRecFinAt(data gjson.Result) string {
	finishedAt := data.Get("finishedAt").Value()
	if finishedAt != nil {
		return utils.FormatTimeStampRFC3339(int64(finishedAt.(float64))/1000, false)
	}

	return ""
}