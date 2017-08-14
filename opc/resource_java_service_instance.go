package opc

import (
	"fmt"
	"log"
	"strings"

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
				Default:  java.ServiceInstanceEditionEE,
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
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"uri": {
							Type:     schema.TypeString,
							Computed: true,
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
							Type:     schema.TypeInt,
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
			"public_key_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"backup_destination": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "BOTH",
				ValidateFunc: validation.StringInSlice([]string{
					string(java.ServiceInstanceBackupDestinationBoth),
					string(java.ServiceInstanceBackupDestinationNone),
				}, false),
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"enable_admin_console": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Default:  false,
			},
			"ip_network": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"provision_otd": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"public_network": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"sample_app_deployment_requested": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Default:  false,
			},
			"app_db": {
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"username": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"password": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"name": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"pdb_name": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
					},
				},
			},
			"backup_volume_size": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"cluster_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"connect_string": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"content_port": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
				Default:  8001,
			},
			"deployment_channel_port": {
				Type:     schema.TypeInt,
				Optional: true,
				ForceNew: true,
			},
			"domain": {
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"mode": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
							Default:  "PRODUCTION",
							ValidateFunc: validation.StringInSlice([]string{
								string(java.ServiceInstanceDomainModeDev),
								string(java.ServiceInstanceDomainModePro),
							}, false),
						},
						"name": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
							Computed: true,
						},
						"partition_count": {
							Type:         schema.TypeInt,
							Optional:     true,
							ForceNew:     true,
							ValidateFunc: validation.IntBetween(0, 4),
						},
						"volume_size": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
					},
				},
			},
			"high_availability": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Default:  false,
			},
			"ip_reservations": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"listener": {
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"port": {
							Type:     schema.TypeInt,
							Optional: true,
							ForceNew: true,
							Default:  8080,
						},
						"enabled": {
							Type:     schema.TypeBool,
							Optional: true,
							ForceNew: true,
							Default:  true,
						},
						"type": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
							Default:  "http",
						},
					},
				},
			},
			"load_balancing_policy": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(java.ServiceInstanceLoadBalancingPolicyLCC),
					string(java.ServiceInstanceLoadBalancingPolicyLRT),
					string(java.ServiceInstanceLoadBalancingPolicyRR),
				}, false),
			},
			"managed_servers": {
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"count": {
							Type:         schema.TypeInt,
							Optional:     true,
							ForceNew:     true,
							ValidateFunc: validation.IntBetween(1, 8),
							Default:      1,
						},
						"initial_heap_size": {
							Type:     schema.TypeInt,
							Optional: true,
							ForceNew: true,
						},
						"max_heap_size": {
							Type:     schema.TypeInt,
							Optional: true,
							ForceNew: true,
						},
						"jvm_args": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"initial_permanent_generation": {
							Type:     schema.TypeInt,
							Optional: true,
							ForceNew: true,
						},
						"max_permanent_generation": {
							Type:     schema.TypeInt,
							Optional: true,
							ForceNew: true,
						},
						"overwrite_jvm_args": {
							Type:     schema.TypeBool,
							Optional: true,
							ForceNew: true,
							Default:  false,
						},
					},
				},
			},
			"mw_volume_size": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"node_manager": {
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
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
						"port": {
							Type:     schema.TypeInt,
							Optional: true,
							ForceNew: true,
							Default:  5556,
						},
					},
				},
			},
			"pdb_service_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"privileged_ports": {
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"content_port": {
							Type:     schema.TypeInt,
							Optional: true,
							ForceNew: true,
							Default:  80,
						},
						"listener_port": {
							Type:     schema.TypeInt,
							Optional: true,
							ForceNew: true,
							Default:  80,
						},
						"secured_content_port": {
							Type:     schema.TypeInt,
							Optional: true,
							ForceNew: true,
							Default:  443,
						},
						"secured_listener_port": {
							Type:     schema.TypeInt,
							Optional: true,
							ForceNew: true,
							Default:  443,
						},
					},
				},
			},
			"scaling_units": {
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"unit": {
							Type:     schema.TypeSet,
							Required: true,
							ForceNew: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"heap_size": {
										Type:         schema.TypeInt,
										Required:     true,
										ForceNew:     true,
										ValidateFunc: validation.IntBetween(1, 16),
									},
									"jvm_count": {
										Type:         schema.TypeInt,
										Required:     true,
										ForceNew:     true,
										ValidateFunc: validation.IntBetween(1, 8),
									},
									"shape": {
										Type:     schema.TypeString,
										Required: true,
										ForceNew: true,
									},
									"vm_count": {
										Type:         schema.TypeInt,
										Required:     true,
										ForceNew:     true,
										ValidateFunc: validation.IntBetween(1, 3),
									},
								},
							},
						},
					},
				},
			},
			"secured_ports": {
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"admin_port": {
							Type:     schema.TypeInt,
							Optional: true,
							ForceNew: true,
							Default:  7002,
						},
						"content_port": {
							Type:     schema.TypeInt,
							Optional: true,
							ForceNew: true,
							Default:  8002,
						},
						"listener_port": {
							Type:     schema.TypeInt,
							Optional: true,
							ForceNew: true,
							Default:  8081,
						},
					},
				},
			},
			"upper_stack_product_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(java.ServiceInstanceUpperStackProductNameODI),
					string(java.ServiceInstanceUpperStackProductNameWCP),
				}, false),
			},
			// Values below are passed back from the GET API call
			"auto_update": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"compliance_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"compliance_status_description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_by": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"creation_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"db_associations": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"db_apex_url": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"db_app": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"db_connect_string": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"db_em_url": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"db_infra": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"db_monitor_url": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"db_service_level": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"db_service_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"db_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"pdb_service_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"db_info": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"fmw_control_url": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"identity_domain": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"is_app_2_cloud": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"life_cycle_control_job_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"memory_size": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"ocpu_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"options": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cluster": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"heap_increments": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"heap_size": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"jvm_count": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"max_heap": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"max_primary": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"max_scaling_unit": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"primary_increments": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"scaling_unit_count": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"scaling_unit_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"shape": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"total_heap": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"total_primary": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"vm_count": {
										Type:     schema.TypeInt,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
			"otd_admin_url": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"otd_provisioned": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"otd_shape": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"otd_storage_size": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"psm_plugin_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"sample_app_url": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"secure_content_url": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"service_component": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"version": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"uri": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"storage_size": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"wls_admin_url": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"wls_deployment_channel_port": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"wls_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"oracle_version": {
				Type:     schema.TypeString,
				Computed: true,
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
		ServiceName:                  d.Get("name").(string),
		Level:                        java.ServiceInstanceLevel(d.Get("level").(string)),
		SubscriptionType:             java.ServiceInstanceSubscriptionType(d.Get("subscription_type").(string)),
		BackupDestination:            java.ServiceInstanceBackupDestination(d.Get("backup_destination").(string)),
		EnableAdminConsole:           d.Get("enable_admin_console").(bool),
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
	d.Set("auto_update", result.AutoUpdate)
	d.Set("cluster_name", result.ClusterName)
	d.Set("compliance_status", result.ComplianceStatus)
	d.Set("compliance_status_description", result.ComplianceStatusDescription)
	d.Set("created_by", result.CreatedBy)
	d.Set("creation_time", result.CreationTime)
	d.Set("db_info", result.DBInfo)
	d.Set("description", result.Description)
	d.Set("edition", result.Edition)
	d.Set("fmw_control_url", result.FMWControlURL)
	d.Set("identity_domain", result.IdentityDomain)
	d.Set("ip_network", result.IPNetwork)
	d.Set("is_app_2_cloud", result.IsApp2Cloud)
	d.Set("level", result.Level)
	d.Set("life_cycle_control_job_id", result.LifecycleControlJobID)
	d.Set("memory_size", result.MemorySize)
	d.Set("ocpu_count", result.OCPUCount)
	d.Set("otd_admin_url", result.OTDAdminURL)
	d.Set("otd_provisioned", result.OTDProvisioned)
	d.Set("otd_shape", result.OTDShape)
	d.Set("otd_storage_size", result.OTDStorageSize)
	d.Set("psm_plugin_version", result.PSMPluginVersion)
	d.Set("region", result.Region)
	d.Set("sample_app_url", result.SampleAppURL)
	d.Set("secure_content_url", result.SecureContentURL)
	d.Set("uri", result.ServiceURI)
	d.Set("shape", result.Shape)
	d.Set("storage_size", result.StorageSize)
	// TODO Changed subscriptionType to subscription_type in the golang sdk
	// d.Set("subscription_type", result.SubscriptionType)
	d.Set("upper_stack_product_name", result.UpperStackProductName)
	// The version you recieve is different than the version you pass in
	// so a different variable name is needed
	d.Set("oracle_version", result.Version)
	d.Set("wls_admin_url", result.WLSAdminURL)
	d.Set("wls_deployment_channel_port", result.WLSDeploymentChannelPort)
	d.Set("wls_version", result.WLSVersion)

	if err := readDBAssociations(d, result.DBAssociations); err != nil {
		return err
	}
	if err := readDatabase(d, result.DBServiceName, result.DBServiceURI); err != nil {
		return err
	}
	if err := readDomain(d, result.DomainName, result.DomainMode); err != nil {
		return err
	}
	if err := readIPReservations(d, result.IPReservations); err != nil {
		return err
	}
	if err := readOptions(d, result.Options); err != nil {
		return err
	}
	if err := readServiceComponents(d, result.ServiceComponents); err != nil {
		return err
	}

	return nil
}

func resourceOPCJavaServiceInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Resource state: %#v", d.State())
	client := meta.(*OPCClient).javaClient.ServiceInstanceClient()
	name := d.Id()

	log.Printf("[DEBUG] Deleting JavaServiceInstance: %v", name)

	// Need to get the dba username and password to delete the service instance
	dbaInfo := d.Get("database").(*schema.Set)
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
		Type:        java.ServiceInstanceType(d.Get("type").(string)),
		Shape:       java.ServiceInstanceShape(d.Get("shape").(string)),
		Version:     java.ServiceInstanceVersion(d.Get("version").(string)),
		ContentPort: d.Get("content_port").(int),
		HAEnabled:   d.Get("high_availability").(bool),
	}

	if val, ok := d.GetOk("public_key"); ok {
		parameter.VMsPublicKey = val.(string)
	}
	if val, ok := d.GetOk("public_key_name"); ok {
		parameter.VMsPublicKeyName = val.(string)
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
		parameter.Edition = java.ServiceInstanceEdition(val.(string))
	}
	if val, ok := d.GetOk("load_balancing_policy"); ok {
		parameter.LoadBalancingPolicy = java.ServiceInstanceLoadBalancingPolicy(val.(string))
	}
	if val, ok := d.GetOk("mw_volume_size"); ok {
		parameter.MWVolumeSize = val.(string)
	}
	if val, ok := d.GetOk("pdb_service_name"); ok {
		parameter.PDBServiceName = val.(string)
	}
	if val, ok := d.GetOk("upper_stack_product_name"); ok {
		parameter.UpperStackProductName = java.ServiceInstanceUpperStackProductName(val.(string))
	}

	if ipReservations := getStringList(d, "ip_reservations"); len(ipReservations) > 0 {
		parameter.IPReservations = strings.Join(ipReservations, ",")
	}

	expandDB(d, &parameter)
	expandAdmin(d, &parameter)
	expandAppDBs(d, &parameter)
	expandDomain(d, &parameter)
	expandListener(d, &parameter)
	expandManagedServers(d, &parameter)
	expandNodeManager(d, &parameter)
	expandPrivilegedPorts(d, &parameter)
	expandScalingUnit(d, &parameter)
	expandSecuredPorts(d, &parameter)
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
	dbaInfo := d.Get("database").(*schema.Set)
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
	appDBs := make([]java.AppDB, appDBInfo.Len())
	for i, val := range appDBInfo.List() {
		attrs := val.(map[string]interface{})
		appDB := java.AppDB{
			DBAName:       attrs["username"].(string),
			DBAPassword:   attrs["password"].(string),
			DBServiceName: attrs["name"].(string),
		}
		appDBs[i] = appDB
	}
	parameter.AppDBs = appDBs
}

func expandDomain(d *schema.ResourceData, parameter *java.Parameter) {
	domainInfo := d.Get("domain").(*schema.Set)
	for _, i := range domainInfo.List() {
		attrs := i.(map[string]interface{})

		parameter.DomainMode = java.ServiceInstanceDomainMode(attrs["mode"].(string))
		if val, ok := attrs["name"].(string); ok && val != "" {
			parameter.DomainName = val
		}
		if val, ok := attrs["partition_count"].(int); ok {
			parameter.DomainPartitionCount = val
		}
		if val, ok := attrs["volume_size"].(string); ok && val != "" {
			parameter.DomainVolumeSize = val
		}
	}
}

func expandListener(d *schema.ResourceData, parameter *java.Parameter) {
	listenerInfo := d.Get("listener").(*schema.Set)
	for _, i := range listenerInfo.List() {
		attrs := i.(map[string]interface{})
		parameter.ListenerPort = attrs["port"].(int)
		parameter.ListenerPortEnabled = attrs["enabled"].(bool)
		parameter.ListenerType = attrs["type"].(string)
	}
}

func expandManagedServers(d *schema.ResourceData, parameter *java.Parameter) {
	msInfo := d.Get("managed_servers").(*schema.Set)
	for _, i := range msInfo.List() {
		attrs := i.(map[string]interface{})
		if val, ok := attrs["count"].(int); ok {
			parameter.ManagedServerCount = val
		}
		if val, ok := attrs["initial_heap_size"].(int); ok {
			parameter.MSInitialHeapMB = val
		}
		if val, ok := attrs["max_heap_size"].(int); ok {
			parameter.MSMaxHeapMB = val
		}
		if val, ok := attrs["jvm_args"].(string); ok && val != "" {
			parameter.MSJvmArgs = val
		}
		if val, ok := attrs["initial_permanent_generation"].(int); ok {
			parameter.MSPermMB = val
		}
		if val, ok := attrs["max_permanent_generation"].(int); ok {
			parameter.MSMaxPermMB = val
		}
		if val, ok := attrs["overwrite_jvm_args"].(bool); ok {
			parameter.OverwriteMsJVMArgs = val
		}
	}
}

func expandNodeManager(d *schema.ResourceData, parameter *java.Parameter) {
	nmInfo := d.Get("node_manager").(*schema.Set)
	for _, i := range nmInfo.List() {
		attrs := i.(map[string]interface{})
		parameter.NodeManagerPort = attrs["port"].(int)
		if val, ok := attrs["password"].(string); ok && val != "" {
			parameter.NodeManagerPassword = val
		}
		if val, ok := attrs["username"].(string); ok && val != "" {
			parameter.NodeManagerUsername = val
		}
	}
}

func expandPrivilegedPorts(d *schema.ResourceData, parameter *java.Parameter) {
	portInfo := d.Get("privileged_ports").(*schema.Set)
	for _, i := range portInfo.List() {
		attrs := i.(map[string]interface{})
		parameter.PrivilegedContentPort = attrs["content_port"].(int)
		parameter.PrivilegedListenerPort = attrs["listener_port"].(int)
		parameter.PrivilegedSecuredContentPort = attrs["secured_content_port"].(int)
		parameter.PrivilegedSecuredListenerPort = attrs["secured_listener_port"].(int)
	}
}

func expandScalingUnit(d *schema.ResourceData, parameter *java.Parameter) {
	scalingUnitInfo := d.Get("scaling_units").(*schema.Set)
	for _, i := range scalingUnitInfo.List() {
		attrs := i.(map[string]interface{})
		if val, ok := attrs["name"].(string); ok && val != "" {
			parameter.ScalingUnitName = java.ServiceInstanceScalingUnitName(val)
		}
		expandUnits(d, parameter)
	}
}

func expandUnits(d *schema.ResourceData, parameter *java.Parameter) {
	unitInfo := d.Get("scaling_units.0.unit").(*schema.Set)
	parameter.ScalingUnitCount = unitInfo.Len()
	units := make([]java.ScalingUnit, unitInfo.Len())
	for i, val := range unitInfo.List() {
		attrs := val.(map[string]interface{})
		unit := java.ScalingUnit{
			HeapSize: attrs["heap_size"].(int),
			JVMCount: attrs["jvm_count"].(int),
			Shape:    java.ServiceInstanceShape(attrs["shape"].(string)),
			VMCount:  attrs["vm_count"].(int),
		}
		units[i] = unit
	}

	// TODO Fix java service instance `ScalingUnit` to []ScalingUnit
	// TODO Change units[0] to units
	parameter.ScalingUnit = units[0]
}

func expandSecuredPorts(d *schema.ResourceData, parameter *java.Parameter) {
	portInfo := d.Get("secured_ports").(*schema.Set)
	for _, i := range portInfo.List() {
		attrs := i.(map[string]interface{})
		parameter.SecuredAdminPort = attrs["admin_port"].(int)
		parameter.SecuredContentPort = attrs["content_port"].(int)
		parameter.SecuredListenerPort = attrs["listener_port"].(int)
	}
}

func readDBAssociations(d *schema.ResourceData, dbAssociations []java.DBAssociation) error {
	result := make([]map[string]interface{}, 0)

	if dbAssociations == nil || len(dbAssociations) == 0 {
		return d.Set("db_associations", nil)
	}

	for _, dbAssociation := range dbAssociations {
		res := make(map[string]interface{})
		res["db_apex_url"] = dbAssociation.DBApexURL
		res["db_app"] = dbAssociation.DBApp
		res["db_connect_string"] = dbAssociation.DBConnectString
		res["db_em_url"] = dbAssociation.DBEmURL
		res["db_infra"] = dbAssociation.DBInfra
		res["db_monitor_url"] = dbAssociation.DBMonitorURL
		res["db_service_level"] = dbAssociation.DBServiceName
		res["db_version"] = dbAssociation.DBVersion
		res["pdb_service_name"] = dbAssociation.PDBServiceName
		result = append(result, res)
	}
	return d.Set("db_associations", result)
}

func readDatabase(d *schema.ResourceData, name string, uri string) error {
	result := make([]map[string]interface{}, 0)

	db := make(map[string]interface{})
	db["name"] = name
	db["uri"] = uri

	// TODO ask about whether or not to seperate the db uri and name into their own
	// seperate attributes or lumping them in and reading the username and password
	dbaInfo := d.Get("database").(*schema.Set)
	for _, i := range dbaInfo.List() {
		attrs := i.(map[string]interface{})
		db["username"] = attrs["username"].(string)
		db["password"] = attrs["password"].(string)
		if val, ok := attrs["network"].(string); ok && val != "" {
			db["network"] = val
		}
	}

	result = append(result, db)
	return d.Set("database", result)
}

func readDomain(d *schema.ResourceData, name string, mode string) error {
	result := make([]map[string]interface{}, 0)

	domain := make(map[string]interface{})
	domain["name"] = name
	domain["mode"] = mode

	// TODO: Check this question from readDatabase method too
	domainInfo := d.Get("domain").(*schema.Set)
	for _, i := range domainInfo.List() {
		attrs := i.(map[string]interface{})
		if val, ok := attrs["partition_count"].(int); ok {
			domain["partition_count"] = val
		}
		if val, ok := attrs["volume_size"].(string); ok && val != "" {
			domain["volume_size"] = val
		}
	}
	result = append(result, domain)
	return d.Set("domain", result)
}

func readIPReservations(d *schema.ResourceData, ipReservations []java.IPReservation) error {
	result := make([]string, 0)

	for _, res := range ipReservations {
		result = append(result, res.Name)
	}
	return setStringList(d, "ip_reservations", result)
}

func readOptions(d *schema.ResourceData, options []java.Option) error {
	result := make([]map[string]interface{}, 0)

	if options == nil || len(options) == 0 {
		return d.Set("options", nil)
	}
	for _, option := range options {
		res := make(map[string]interface{})
		res["type"] = option.Type
		res["clusters"] = readClusters(option.Clusters)
		result = append(result, res)
	}
	return d.Set("options", result)
}

func readClusters(clusters []java.Cluster) []map[string]interface{} {
	result := make([]map[string]interface{}, 0)

	// Shouldn't have to check if clusters is nil or len == 0 because we'll return an empty map
	// anyway
	for _, cluster := range clusters {
		res := make(map[string]interface{})
		res["name"] = cluster.ClusterName
		res["heap_increments"] = cluster.HeapIncrements
		res["heap_size"] = cluster.HeapSize
		res["jvm_count"] = cluster.JVMCount
		res["max_heap"] = cluster.MaxHeap
		res["max_primary"] = cluster.MaxPrimary
		res["max_scaling_unit"] = cluster.MaxScalingUnit
		res["primary_increments"] = cluster.PrimaryIncrements
		res["scaling_unit_count"] = cluster.ScalingUnitCount
		res["scaling_unit_name"] = cluster.ScalingUnitName
		res["shape"] = cluster.Shape
		res["total_heap"] = cluster.TotalHeap
		res["total_primary"] = cluster.TotalPrimary
		res["vm_count"] = cluster.VMCount

		result = append(result, res)
	}
	return result
}

func readServiceComponents(d *schema.ResourceData, serviceComponents []java.ServiceComponent) error {
	result := make([]map[string]interface{}, 0)

	if serviceComponents == nil || len(serviceComponents) == 0 {
		return d.Set("service_component", nil)
	}

	for _, serviceComponent := range serviceComponents {
		res := make(map[string]interface{})
		res["type"] = serviceComponent.Type
		res["version"] = serviceComponent.Version
		result = append(result, res)
	}
	return d.Set("service_component", result)
}
