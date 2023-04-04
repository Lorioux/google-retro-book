// ----------------------------------------------------------------------------
//
//     ***     AUTO GENERATED CODE    ***    Type: MMv1     ***
//
// ----------------------------------------------------------------------------
//
//     This file is automatically generated by Magic Modules and manual
//     changes will be clobbered when the file is regenerated.
//
//     Please read more about how to change this file in
//     .github/CONTRIBUTING.md.
//
// ----------------------------------------------------------------------------

package google

import (
	"fmt"
	"log"
	"reflect"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceMonitoringGenericService() *schema.Resource {
	return &schema.Resource{
		Create: resourceMonitoringGenericServiceCreate,
		Read:   resourceMonitoringGenericServiceRead,
		Update: resourceMonitoringGenericServiceUpdate,
		Delete: resourceMonitoringGenericServiceDelete,

		Importer: &schema.ResourceImporter{
			State: resourceMonitoringGenericServiceImport,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Update: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"service_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				Description: `An optional service ID to use. If not given, the server will generate a
service ID.`,
			},
			"basic_service": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Description: `A well-known service type, defined by its service type and service labels.
Valid values are described at
https://cloud.google.com/stackdriver/docs/solutions/slo-monitoring/api/api-structures#basic-svc-w-basic-sli`,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"service_labels": {
							Type:     schema.TypeMap,
							Optional: true,
							ForceNew: true,
							Description: `Labels that specify the resource that emits the monitoring data 
which is used for SLO reporting of this 'Service'.`,
							Elem: &schema.Schema{Type: schema.TypeString},
						},
						"service_type": {
							Type:     schema.TypeString,
							Optional: true,
							Description: `The type of service that this basic service defines, e.g. 
APP_ENGINE service type`,
						},
					},
				},
			},
			"display_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Name used for UI elements listing this Service.`,
			},
			"user_labels": {
				Type:     schema.TypeMap,
				Optional: true,
				Description: `Labels which have been used to annotate the service. Label keys must start
with a letter. Label keys and values may contain lowercase letters,
numbers, underscores, and dashes. Label keys and values have a maximum
length of 63 characters, and must be less than 128 bytes in size. Up to 64
label entries may be stored. For labels which do not have a semantic value,
the empty string may be supplied for the label value.`,
				Elem: &schema.Schema{Type: schema.TypeString},
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
				Description: `The full resource name for this service. The syntax is:
projects/[PROJECT_ID]/services/[SERVICE_ID].`,
			},
			"telemetry": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `Configuration for how to query telemetry on a Service.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"resource_name": {
							Type:     schema.TypeString,
							Optional: true,
							Description: `The full name of the resource that defines this service.
Formatted as described in
https://cloud.google.com/apis/design/resource_names.`,
						},
					},
				},
			},
			"project": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
		},
		UseJSONNumber: true,
	}
}

func resourceMonitoringGenericServiceCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	userAgent, err := generateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	obj := make(map[string]interface{})
	displayNameProp, err := expandMonitoringGenericServiceDisplayName(d.Get("display_name"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("display_name"); !isEmptyValue(reflect.ValueOf(displayNameProp)) && (ok || !reflect.DeepEqual(v, displayNameProp)) {
		obj["displayName"] = displayNameProp
	}
	userLabelsProp, err := expandMonitoringGenericServiceUserLabels(d.Get("user_labels"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("user_labels"); ok || !reflect.DeepEqual(v, userLabelsProp) {
		obj["userLabels"] = userLabelsProp
	}
	basicServiceProp, err := expandMonitoringGenericServiceBasicService(d.Get("basic_service"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("basic_service"); !isEmptyValue(reflect.ValueOf(basicServiceProp)) && (ok || !reflect.DeepEqual(v, basicServiceProp)) {
		obj["basicService"] = basicServiceProp
	}

	url, err := replaceVars(d, config, "{{MonitoringBasePath}}v3/projects/{{project}}/services?serviceId={{service_id}}")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Creating new GenericService: %#v", obj)
	billingProject := ""

	project, err := getProject(d, config)
	if err != nil {
		return fmt.Errorf("Error fetching project for GenericService: %s", err)
	}
	billingProject = project

	// err == nil indicates that the billing_project value was found
	if bp, err := getBillingProject(d, config); err == nil {
		billingProject = bp
	}

	res, err := SendRequestWithTimeout(config, "POST", billingProject, url, userAgent, obj, d.Timeout(schema.TimeoutCreate), isMonitoringConcurrentEditError)
	if err != nil {
		return fmt.Errorf("Error creating GenericService: %s", err)
	}
	if err := d.Set("name", flattenMonitoringGenericServiceName(res["name"], d, config)); err != nil {
		return fmt.Errorf(`Error setting computed identity field "name": %s`, err)
	}

	// Store the ID now
	id, err := replaceVars(d, config, "projects/{{project}}/services/{{service_id}}")
	if err != nil {
		return fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	log.Printf("[DEBUG] Finished creating GenericService %q: %#v", d.Id(), res)

	return resourceMonitoringGenericServiceRead(d, meta)
}

func resourceMonitoringGenericServiceRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	userAgent, err := generateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	url, err := replaceVars(d, config, "{{MonitoringBasePath}}v3/projects/{{project}}/services/{{service_id}}")
	if err != nil {
		return err
	}

	billingProject := ""

	project, err := getProject(d, config)
	if err != nil {
		return fmt.Errorf("Error fetching project for GenericService: %s", err)
	}
	billingProject = project

	// err == nil indicates that the billing_project value was found
	if bp, err := getBillingProject(d, config); err == nil {
		billingProject = bp
	}

	res, err := SendRequest(config, "GET", billingProject, url, userAgent, nil, isMonitoringConcurrentEditError)
	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("MonitoringGenericService %q", d.Id()))
	}

	if err := d.Set("project", project); err != nil {
		return fmt.Errorf("Error reading GenericService: %s", err)
	}

	if err := d.Set("name", flattenMonitoringGenericServiceName(res["name"], d, config)); err != nil {
		return fmt.Errorf("Error reading GenericService: %s", err)
	}
	if err := d.Set("display_name", flattenMonitoringGenericServiceDisplayName(res["displayName"], d, config)); err != nil {
		return fmt.Errorf("Error reading GenericService: %s", err)
	}
	if err := d.Set("user_labels", flattenMonitoringGenericServiceUserLabels(res["userLabels"], d, config)); err != nil {
		return fmt.Errorf("Error reading GenericService: %s", err)
	}
	if err := d.Set("telemetry", flattenMonitoringGenericServiceTelemetry(res["telemetry"], d, config)); err != nil {
		return fmt.Errorf("Error reading GenericService: %s", err)
	}
	if err := d.Set("basic_service", flattenMonitoringGenericServiceBasicService(res["basicService"], d, config)); err != nil {
		return fmt.Errorf("Error reading GenericService: %s", err)
	}

	return nil
}

func resourceMonitoringGenericServiceUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	userAgent, err := generateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	billingProject := ""

	project, err := getProject(d, config)
	if err != nil {
		return fmt.Errorf("Error fetching project for GenericService: %s", err)
	}
	billingProject = project

	obj := make(map[string]interface{})
	displayNameProp, err := expandMonitoringGenericServiceDisplayName(d.Get("display_name"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("display_name"); !isEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, displayNameProp)) {
		obj["displayName"] = displayNameProp
	}
	userLabelsProp, err := expandMonitoringGenericServiceUserLabels(d.Get("user_labels"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("user_labels"); ok || !reflect.DeepEqual(v, userLabelsProp) {
		obj["userLabels"] = userLabelsProp
	}

	url, err := replaceVars(d, config, "{{MonitoringBasePath}}v3/projects/{{project}}/services/{{service_id}}")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Updating GenericService %q: %#v", d.Id(), obj)
	updateMask := []string{}

	if d.HasChange("display_name") {
		updateMask = append(updateMask, "displayName")
	}

	if d.HasChange("user_labels") {
		updateMask = append(updateMask, "userLabels")
	}
	// updateMask is a URL parameter but not present in the schema, so replaceVars
	// won't set it
	url, err = addQueryParams(url, map[string]string{"updateMask": strings.Join(updateMask, ",")})
	if err != nil {
		return err
	}

	// err == nil indicates that the billing_project value was found
	if bp, err := getBillingProject(d, config); err == nil {
		billingProject = bp
	}

	res, err := SendRequestWithTimeout(config, "PATCH", billingProject, url, userAgent, obj, d.Timeout(schema.TimeoutUpdate), isMonitoringConcurrentEditError)

	if err != nil {
		return fmt.Errorf("Error updating GenericService %q: %s", d.Id(), err)
	} else {
		log.Printf("[DEBUG] Finished updating GenericService %q: %#v", d.Id(), res)
	}

	return resourceMonitoringGenericServiceRead(d, meta)
}

func resourceMonitoringGenericServiceDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	userAgent, err := generateUserAgentString(d, config.UserAgent)
	if err != nil {
		return err
	}

	billingProject := ""

	project, err := getProject(d, config)
	if err != nil {
		return fmt.Errorf("Error fetching project for GenericService: %s", err)
	}
	billingProject = project

	url, err := replaceVars(d, config, "{{MonitoringBasePath}}v3/projects/{{project}}/services/{{service_id}}")
	if err != nil {
		return err
	}

	var obj map[string]interface{}
	log.Printf("[DEBUG] Deleting GenericService %q", d.Id())

	// err == nil indicates that the billing_project value was found
	if bp, err := getBillingProject(d, config); err == nil {
		billingProject = bp
	}

	res, err := SendRequestWithTimeout(config, "DELETE", billingProject, url, userAgent, obj, d.Timeout(schema.TimeoutDelete), isMonitoringConcurrentEditError)
	if err != nil {
		return handleNotFoundError(err, d, "GenericService")
	}

	log.Printf("[DEBUG] Finished deleting GenericService %q: %#v", d.Id(), res)
	return nil
}

func resourceMonitoringGenericServiceImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	config := meta.(*Config)
	if err := parseImportId([]string{
		"projects/(?P<project>[^/]+)/services/(?P<service_id>[^/]+)",
		"(?P<project>[^/]+)/(?P<service_id>[^/]+)",
		"(?P<service_id>[^/]+)",
	}, d, config); err != nil {
		return nil, err
	}

	// Replace import id for the resource id
	id, err := replaceVars(d, config, "projects/{{project}}/services/{{service_id}}")
	if err != nil {
		return nil, fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	return []*schema.ResourceData{d}, nil
}

func flattenMonitoringGenericServiceName(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	return v
}

func flattenMonitoringGenericServiceDisplayName(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	return v
}

func flattenMonitoringGenericServiceUserLabels(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	return v
}

func flattenMonitoringGenericServiceTelemetry(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	if v == nil {
		return nil
	}
	original := v.(map[string]interface{})
	if len(original) == 0 {
		return nil
	}
	transformed := make(map[string]interface{})
	transformed["resource_name"] =
		flattenMonitoringGenericServiceTelemetryResourceName(original["resourceName"], d, config)
	return []interface{}{transformed}
}
func flattenMonitoringGenericServiceTelemetryResourceName(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	return v
}

func flattenMonitoringGenericServiceBasicService(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	if v == nil {
		return nil
	}
	original := v.(map[string]interface{})
	if len(original) == 0 {
		return nil
	}
	transformed := make(map[string]interface{})
	transformed["service_type"] =
		flattenMonitoringGenericServiceBasicServiceServiceType(original["serviceType"], d, config)
	transformed["service_labels"] =
		flattenMonitoringGenericServiceBasicServiceServiceLabels(original["serviceLabels"], d, config)
	return []interface{}{transformed}
}
func flattenMonitoringGenericServiceBasicServiceServiceType(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	return v
}

func flattenMonitoringGenericServiceBasicServiceServiceLabels(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	return v
}

func expandMonitoringGenericServiceDisplayName(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandMonitoringGenericServiceUserLabels(v interface{}, d TerraformResourceData, config *Config) (map[string]string, error) {
	if v == nil {
		return map[string]string{}, nil
	}
	m := make(map[string]string)
	for k, val := range v.(map[string]interface{}) {
		m[k] = val.(string)
	}
	return m, nil
}

func expandMonitoringGenericServiceBasicService(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := make(map[string]interface{})

	transformedServiceType, err := expandMonitoringGenericServiceBasicServiceServiceType(original["service_type"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedServiceType); val.IsValid() && !isEmptyValue(val) {
		transformed["serviceType"] = transformedServiceType
	}

	transformedServiceLabels, err := expandMonitoringGenericServiceBasicServiceServiceLabels(original["service_labels"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedServiceLabels); val.IsValid() && !isEmptyValue(val) {
		transformed["serviceLabels"] = transformedServiceLabels
	}

	return transformed, nil
}

func expandMonitoringGenericServiceBasicServiceServiceType(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandMonitoringGenericServiceBasicServiceServiceLabels(v interface{}, d TerraformResourceData, config *Config) (map[string]string, error) {
	if v == nil {
		return map[string]string{}, nil
	}
	m := make(map[string]string)
	for k, val := range v.(map[string]interface{}) {
		m[k] = val.(string)
	}
	return m, nil
}
