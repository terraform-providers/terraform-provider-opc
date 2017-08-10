package opc

import (
	"fmt"
	"log"

	opcClient "github.com/hashicorp/go-oracle-terraform/client"
	"github.com/hashicorp/go-oracle-terraform/java"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

func resourceOPCJavaServiceInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceOPCJavaServiceInstanceCreate,
		Read:   resourceOPCJavaServiceInstanceRead,
		Delete: resourceOPCJavaServiceInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"edition": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(java.ServiceInstanceEditionSE),
					string(java.ServiceInstanceEditionEE),
					string(java.ServiceInstanceEditionSuite),
				}, false),
			},
			"level": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "PAAS",
				ValidateFunc: validation.StringInSlice([]string{
					string(java.ServiceInstanceLevelPAAS),
					string(java.ServiceInstanceLevelBasic),
				}, false),
			},
			"cloud_storage": {
				Type:     schema.TypeSet,
				Required: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"container": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"create_if_missing": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
							ForceNew: true,
						},
						"username": {
							Type:      schema.TypeString,
							Optional:  true,
							ForceNew:  true,
							Computed:  true,
							Sensitive: true,
						},
						"password": {
							Type:      schema.TypeString,
							Optional:  true,
							ForceNew:  true,
							Computed:  true,
							Sensitive: true,
						},
					},
				},
			},
			"subscription_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(java.ServiceInstanceSubscriptionTypeHourly),
					string(java.ServiceInstanceSubscriptionTypeMonthly),
				}, false),
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(java.ServiceInstanceTypeWebLogic),
					string(java.ServiceInstanceTypeDataGrid),
					string(java.ServiceInstanceTypeOTD),
				}, false),
			},
			"database": {
				Type:     schema.TypeSet,
				Required: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"username": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"password": {
							Type:      schema.TypeString,
							Required:  true,
							ForceNew:  true,
							Sensitive: true,
						},
						"name": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"network": {
							Type: schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
					},
				},
			},
			"shape": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"version": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"admin": {
				Type:     schema.TypeSet,
				Required: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"username": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"password": {
							Type:      schema.TypeString,
							Required:  true,
							ForceNew:  true,
							Sensitive: true,
						},
						"port": {
							Type: schema.TypeInt,
							Optional: true,
							ForceNew: true,
						},
					},
				},
			},
			"public_key": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
		},
		"backup_destination": {
			Type: schema.TypeString,
			Optional: true,
			ForceNew: true,
			Default: "BOTH"
			ValidateFunc: validation.StringInSlice([]string{
				string(java.ServiceInstanceBackupDestinationBoth),
				string(java.ServiceInstanceBackupDestinationNone),
			}, false),
		},
		"description": {
			Type: schema.TypeString,
			Optional: true,
			ForceNew: true,
		},
		"enable_admin_console": {
			Type: schema.TypeBool,
			Optional: true,
			ForceNew: true,
			Default: false,
		},
		"ip_network": {
			Type: schema.TypeString,
			Optional: true,
			ForceNew: true,
		},
		"provision_otd": {
			Type: schema.TypeBool,
			Optional: true,
			ForceNew: true,
		},
		"public_network": {
			Type: schema.TypeString,
			Optional: true,
			ForceNew: true,
		},
		"region": {
			Type: schema.TypeString,
			Optional: true,
			ForceNew: true,
		},
		"sample_app_deployment_requested": {
			Type: schema.TypeBool,
			Optional: true,
			ForceNew: true,
			Default: false,
		},
		"app_db": {
			Type:     schema.TypeSet,
			Optional: true,
			ForceNew: true,
			Elem: &schema.Resource{
				Schema: map[string]*schema.Schema{
					"username": {
						Type: schema.TypeString,
						Required: true,
						ForceNew: true,
					},
					"password": {
						Type: schema.TypeString,
						Required: true,
						ForceNew: true,
					},
					"name": {
						Type: schema.TypeString,
						Required: true,
						ForceNew: true,
					},
					"pdb_name": {
						Type: schema.TypeString,
						Optional: true,
						ForceNew: true,
					},
				},
			},
			"backup_volume_size": {
				Type: schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"cluster_name": {
				Type: schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"connect_string": {
				Type: schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"content_port": {
				Type: schema.TypeInt,
				Optional: true,
				ForceNew: true,
				Default: 8001,
			},
			"deployment_channel_port": {
				Type: schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"domain": {
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"mode": {
							Type: schema.TypeString,
							Optional: true,
							ForceNew: true,
							Default: "PRODUCTION"
							ValidateFunc: validation.StringInSlice([]string{
								string(java.ServiceInstanceDomainModeDev),
								string(java.ServiceInstanceDomainModePro),
							}, false),
						},
						"name": {
							Type: schema.TypeString,
							Optional: true,
							ForceNew: true,
							Computed: true,
						},
						"partition_count": {
							Type: schema.TypeInt,
							Optional: true,
							ForceNew: true,
							ValidateFunc: validation.InBetween(0, 4)
						},
						"volume_size": {
							Type: schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
					},
				},
			},
		},
	}
}

func resourceOPCJavaServiceInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Resource state: %#v", d.State())

	log.Print("[DEBUG] Creating JavaServiceInstance")

	client := meta.(*OPCClient).javaClient.ServiceInstanceClient()
	if client == nil {
		return fmt.Errorf("Java Client is not initialized. Make sure to use `java_endpoint` variable or `OPC_JAVA_ENDPOINT` env variable")
	}
	input := java.CreateServiceInstanceInput{
		ServiceName:      d.Get("name").(string),
		Level:            java.ServiceInstanceLevel(d.Get("level").(string)),
		SubscriptionType: java.ServiceInstanceSubscriptionType(d.Get("subscription_type").(string)),
		BackupDestination: java.ServiceInstanceBackupDestination(d.Get("backup_destination").(string)),
		EnableAdminConsole: d.Get("enable_admin_console").(bool),
		SampleAppDeploymentRequested: d.Get("sample_app_deployment_requested").(bool),
	}

	if val, ok := d.GetOk("description"); ok {
		input.Description = val.(string)
	}
	if val, ok := d.GetOk("ip_network"); ok {
		input.IPNetwork = val.(string)
	}
	if val, ok := d.GetOk("provision_otd"); ok {
		input.ProvisionOTD = val.(bool)
	}
	if val, ok := d.GetOk("public_network"); ok {
		input.PublicNetwork = val.(string)
	}
	if val, ok := d.GetOk("region"); ok {
		input.Region = val.(string)
	}
	expandJavaCloudStorage(d, &input)

	input.Parameters = expandJavaParameter(client, d)

	info, err := client.CreateServiceInstance(&input)
	if err != nil {
		return fmt.Errorf("Error creating JavaServiceInstance: %s", err)
	}

	d.SetId(info.ServiceName)
	return resourceOPCJavaServiceInstanceRead(d, meta)
}

func resourceOPCJavaServiceInstanceRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Resource state: %#v", d.State())
	client := meta.(*OPCClient).javaClient.ServiceInstanceClient()

	log.Printf("[DEBUG] Reading state of ip reservation %s", d.Id())
	getInput := java.GetServiceInstanceInput{
		Name: d.Id(),
	}

	result, err := client.GetServiceInstance(&getInput)
	if err != nil {
		// Java Service Instance does not exist
		if opcClient.WasNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error reading JavaServiceInstance %s: %s", d.Id(), err)
	}

	if result == nil {
		d.SetId("")
		return nil
	}

	log.Printf("[DEBUG] Read state of JavaServiceInstance %s: %#v", d.Id(), result)
	d.Set("name", result.ServiceName)

	return nil
}

func resourceOPCJavaServiceInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Resource state: %#v", d.State())
	client := meta.(*OPCClient).javaClient.ServiceInstanceClient()
	name := d.Id()

	log.Printf("[DEBUG] Deleting JavaServiceInstance: %v", name)

	// Need to get the dba username and password to delete the service instance
	dbaInfo := d.Get("dba").(*schema.Set)
	var username, password string
	for _, i := range dbaInfo.List() {
		attrs := i.(map[string]interface{})
		username = attrs["username"].(string)
		password = attrs["password"].(string)
	}

	input := java.DeleteServiceInstanceInput{
		Name:        name,
		DBAUsername: username,
		DBAPassword: password,
	}
	if err := client.DeleteServiceInstance(&input); err != nil {
		return fmt.Errorf("Error deleting JavaServiceInstance")
	}
	return nil
}

func expandJavaParameter(client *java.ServiceInstanceClient, d *schema.ResourceData) []java.Parameter {
	parameter := java.Parameter{
		Type:    java.ServiceInstanceType(d.Get("type").(string)),
		Shape:   java.ServiceInstanceShape(d.Get("shape").(string)),
		Version: java.ServiceInstanceVersion(d.Get("version").(string)),
		ContentPort: d.Get("content_port").(int)
	}

	if val, ok := d.GetOk("public_key"); ok {
		parameter.VMsPublicKey = val.(string)
	}
	if val, ok := d.GetOk("backup_volume_size"); ok {
		parameter.BackupVolumeSize = val.(string)
	}
	if val, ok := d.GetOk("cluster_name"); ok {
		parameter.ClusterName = val.(string)
	}
	if val, ok := d.GetOk("connect_string"); ok {
		parameter.ConnectString = val.(string)
	}
	if val, ok := d.GetOk("deployment_channel_port"); ok {
		parameter.DeploymentChannelPort = val.(int)
	}
	if val, ok := d.GetOk("edition"); ok {
		parameter.Edition = val.(string)
	}

	expandDB(d, &parameter)
	expandAdmin(d, &parameter)
	expandAppDBs(d, &parameter)
	expandDomain(d, &parameter)
	return []java.Parameter{parameter}
}

func expandJavaCloudStorage(d *schema.ResourceData, input *java.CreateServiceInstanceInput) {
	cloudStorageInfo := d.Get("cloud_storage").(*schema.Set)
	for _, i := range cloudStorageInfo.List() {
		attrs := i.(map[string]interface{})
		input.CloudStorageContainer = attrs["container"].(string)
		input.CreateStorageContainerIfMissing = attrs["create_if_missing"].(bool)
		if val, ok := attrs["username"].(string); ok && val != "" {
			input.CloudStorageUsername = val
		}
		if val, ok := attrs["password"].(string); ok && val != "" {
			input.CloudStoragePassword = val
		}
	}
}

func expandDB(d *schema.ResourceData, parameter *java.Parameter) {
	dbaInfo := d.Get("db").(*schema.Set)
	for _, i := range dbaInfo.List() {
		attrs := i.(map[string]interface{})
		parameter.DBServiceName = attrs["name"].(string)
		parameter.DBAName = attrs["username"].(string)
		parameter.DBAPassword = attrs["password"].(string)
		if val, ok := attrs["network"].(string); ok && val != "" {
			parameter.DBNetwork = val
		}
	}
}

func expandAdmin(d *schema.ResourceData, parameter *java.Parameter) {
	adminInfo := d.Get("admin").(*schema.Set)
	for _, i := range adminInfo.List() {
		attrs := i.(map[string]interface{})
		parameter.AdminUsername = attrs["username"].(string)
		parameter.AdminPassword = attrs["password"].(string)
		if val, ok := attrs["port"].(int); ok && val != 0 {
			parameter.AdminPort = val
		}
	}
}

func expandAppDBs(d *schema.ResourceData, parameter *java.Parameter) {
	appDBInfo := d.Get("app_db").(*schema.Set)
	appDBs := make([]java.AppDB, len(appDBInfo))
	for i, val := range appDBInfo.List() {
		attrs := val.(map[string]interface{})
		appDB := java.AppDB {
			DBAName: attrs["username"].(string)
			DBAPassword: attrs["password"].(string)
			DBAServiceName: attrs["name"].(string)
		}
		appDBs[i] = appDB
	}
}

func expandDomain(d *schema.ResourceData, parameter *java.Parameter) {
	domainInfo := d.Get("domain").(*schema.Set)
	for _, i := range domainInfo.List() {
		attrs := i.(map[string]interface{})

		parameter.DomainMode = java.ServiceInstanceDomainMode(attrs["mode"].(string))
		if val, ok := attrs["name"].(string); ok && val != "" {
			parameter.DomainName = val
		}
		if val, ok := attrs["partition_count"].(int); ok && val != "" {
			parameter.DomainPartitionCount = val
		}
		if val, ok := attrs["volume_size"].(string); ok && val != "" {
			parameter.DomainVolumeSize = val
		}
	}
}
