package opc

import (
	"fmt"
	"log"
	"strings"
	"time"

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

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(120 * time.Minute),
			Delete: schema.DefaultTimeout(120 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
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
			"weblogic": {
				Type:     schema.TypeSet,
				Required: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
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
										Computed: true,
										ForceNew: true,
									},
								},
							},
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
										Computed: true,
										ForceNew: true,
									},
								},
							},
						},
						"backup_volume_size": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
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
						"deployment_channel_port": {
							Type:     schema.TypeInt,
							Optional: true,
							ForceNew: true,
							Default:  9001,
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
						"edition": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(java.ServiceInstanceEditionSE),
								string(java.ServiceInstanceEditionEE),
								string(java.ServiceInstanceEditionSuite),
							}, false),
						},
						"ip_reservations": {
							Type:     schema.TypeList,
							Optional: true,
							ForceNew: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"managed_servers": {
							Type:     schema.TypeSet,
							Optional: true,
							ForceNew: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"server_count": {
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
							Computed: true,
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
									"secured_content_port": {
										Type:     schema.TypeInt,
										Optional: true,
										ForceNew: true,
										Default:  443,
									},
								},
							},
						},
						"secured_admin_port": {
							Type:     schema.TypeInt,
							Optional: true,
							ForceNew: true,
							Default:  7002,
						},
						"shape": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
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
						"version": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"public_key": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
					},
				},
			},
			"otd": {
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
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
										Computed: true,
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
							Default:  java.ServiceInstanceLoadBalancingPolicyLCC,
							ValidateFunc: validation.StringInSlice([]string{
								string(java.ServiceInstanceLoadBalancingPolicyLCC),
								string(java.ServiceInstanceLoadBalancingPolicyLRT),
								string(java.ServiceInstanceLoadBalancingPolicyRR),
							}, false),
						},
						"privileged_ports": {
							Type:     schema.TypeSet,
							Optional: true,
							ForceNew: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"listener_port": {
										Type:     schema.TypeInt,
										Optional: true,
										ForceNew: true,
										Default:  80,
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
						"secured_listener_port": {
							Type:     schema.TypeInt,
							Optional: true,
							ForceNew: true,
							Default:  8081,
						},
						"shape": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"public_key": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
					},
				},
			},
			"datagrid": {
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cluster_name": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
							Computed: true,
						},
						"scaling_unit_count": {
							Type:     schema.TypeInt,
							Required: true,
							ForceNew: true,
							Computed: true,
						},
						"scaling_unit": {
							Type:     schema.TypeSet,
							Optional: true,
							ForceNew: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Required: true,
										ForceNew: true,
									},
									"unit": {
										Type:     schema.TypeSet,
										Optional: true,
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
					},
				},
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
			"oracle_middleware_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceOPCJavaServiceInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Resource state: %#v", d.State())

	log.Print("[DEBUG] Creating JavaServiceInstance")

	jClient, err := getJavaClient(meta)
	if err != nil {
		return err
	}
	client := jClient.ServiceInstanceClient()

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
	if val, ok := d.GetOk("public_network"); ok {
		input.PublicNetwork = val.(string)
	}
	if val, ok := d.GetOk("region"); ok {
		input.Region = val.(string)
	}
	expandJavaCloudStorage(d, &input)

	expandJavaParameter(d, &input)

	info, err := client.CreateServiceInstance(&input)
	if err != nil {
		return fmt.Errorf("Error creating JavaServiceInstance: %s", err)
	}

	d.SetId(info.ServiceName)
	return resourceOPCJavaServiceInstanceRead(d, meta)
}

func resourceOPCJavaServiceInstanceRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Resource state: %#v", d.State())
	jClient, err := getJavaClient(meta)
	if err != nil {
		return err
	}
	client := jClient.ServiceInstanceClient()

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
	d.Set("subscription_type", result.SubscriptionType)
	d.Set("upper_stack_product_name", result.UpperStackProductName)
	// The version you recieve is different than the version you pass in
	// so a different variable name is needed
	d.Set("oracle_middleware_version", result.Version)
	d.Set("wls_admin_url", result.WLSAdminURL)
	d.Set("wls_deployment_channel_port", result.WLSDeploymentChannelPort)
	d.Set("wls_version", result.WLSVersion)

	if err := readDBAssociations(d, result.DBAssociations); err != nil {
		return err
	}
	if err := readServiceComponents(d, result.ServiceComponents); err != nil {
		return err
	}

	return nil
}

func resourceOPCJavaServiceInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Resource state: %#v", d.State())
	jClient, err := getJavaClient(meta)
	if err != nil {
		return err
	}
	client := jClient.ServiceInstanceClient()
	name := d.Id()

	log.Printf("[DEBUG] Deleting JavaServiceInstance: %v", name)

	// Need to get the dba username and password to delete the service instance
	webLogicInfo := d.Get("weblogic").(*schema.Set).List()
	webLogicConfig := webLogicInfo[0].(map[string]interface{})
	dbaInfo := webLogicConfig["database"].(*schema.Set)

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

func expandJavaParameter(d *schema.ResourceData, input *java.CreateServiceInstanceInput) {
	parameters := []java.Parameter{expandWebLogicParameter(d)}

	if _, ok := d.GetOk("otd"); ok {
		input.ProvisionOTD = true
		parameters = append(parameters, expandOTDParameter(d))
	}

	if _, ok := d.GetOk("datagrid"); ok {
		parameters = append(parameters, expandDatagridParameter(d))
	}

	input.Parameters = parameters

}

func expandWebLogicParameter(d *schema.ResourceData) java.Parameter {
	webLogicParam := java.Parameter{}
	webLogicInfo := d.Get("weblogic").(*schema.Set).List()
	webLogicConfig := webLogicInfo[0].(map[string]interface{})

	webLogicParam.Type = java.ServiceInstanceTypeWebLogic
	webLogicParam.Edition = java.ServiceInstanceEdition(webLogicConfig["edition"].(string))
	webLogicParam.Shape = java.ServiceInstanceShape(webLogicConfig["shape"].(string))
	webLogicParam.Version = java.ServiceInstanceVersion(webLogicConfig["version"].(string))
	webLogicParam.VMsPublicKey = webLogicConfig["public_key"].(string)

	expandAdmin(&webLogicParam, webLogicConfig)
	expandDB(&webLogicParam, webLogicConfig)
	expandDomain(&webLogicParam, webLogicConfig)
	expandManagedServers(&webLogicParam, webLogicConfig)
	expandNodeManager(&webLogicParam, webLogicConfig)
	expandPrivilegedPorts(&webLogicParam, webLogicConfig)

	if v := webLogicConfig["app_db"]; v != nil {
		expandAppDBs(&webLogicParam, webLogicConfig)
	}

	if v := webLogicConfig["backup_volume_size"]; v != nil {
		webLogicParam.BackupVolumeSize = v.(string)
	}
	if v := webLogicConfig["cluster_name"]; v != nil {
		webLogicParam.ClusterName = v.(string)
	}
	if v := webLogicConfig["connect_string"]; v != nil {
		webLogicParam.ConnectString = v.(string)
	}
	if v := webLogicConfig["content_port"]; v != nil {
		webLogicParam.ContentPort = v.(int)
	}
	if v := webLogicConfig["deployment_channel_port"]; v != nil {
		webLogicParam.DeploymentChannelPort = v.(int)
	}
	if ipReservations := getStringList(d, "weblogic.0.ip_reservations"); len(ipReservations) > 0 {
		webLogicParam.IPReservations = strings.Join(ipReservations, ",")
	}
	if v := webLogicConfig["mw_volume_size"]; v != nil {
		webLogicParam.MWVolumeSize = v.(string)
	}
	if v := webLogicConfig["pdb_service_name"]; v != nil {
		webLogicParam.PDBServiceName = v.(string)
	}
	if v := webLogicConfig["secured_admin_port"]; v != nil {
		webLogicParam.SecuredAdminPort = v.(int)
	}
	if v := webLogicConfig["upper_stack_product_name"]; v != nil {
		webLogicParam.UpperStackProductName = java.ServiceInstanceUpperStackProductName(v.(string))
	}

	return webLogicParam
}

func expandOTDParameter(d *schema.ResourceData) java.Parameter {
	otdParam := java.Parameter{}
	otdInfo := d.Get("otd").(*schema.Set).List()
	otdConfig := otdInfo[0].(map[string]interface{})

	otdParam.Type = java.ServiceInstanceTypeOTD
	otdParam.HAEnabled = otdConfig["high_availability"].(bool)
	otdParam.LoadBalancingPolicy = java.ServiceInstanceLoadBalancingPolicy(otdConfig["load_balancing_policy"].(string))
	otdParam.SecuredListenerPort = otdConfig["secured_listener_port"].(int)
	otdParam.Shape = java.ServiceInstanceShape(otdConfig["shape"].(string))
	otdParam.VMsPublicKey = otdConfig["public_key"].(string)

	expandListener(&otdParam, otdConfig)

	return otdParam
}

func expandDatagridParameter(d *schema.ResourceData) java.Parameter {
	datagridParam := java.Parameter{}
	datagridInfo := d.Get("datagrid").(*schema.Set).List()
	datagridConfig := datagridInfo[0].(map[string]interface{})

	datagridParam.Type = java.ServiceInstanceTypeDataGrid
	if v := datagridConfig["scaling_unit_count"]; v != nil {
		datagridParam.ScalingUnitCount = v.(int)
	}
	if v := datagridConfig["cluster_name"]; v != nil {
		datagridParam.ClusterName = v.(string)
	}

	expandScalingUnit(&datagridParam, datagridConfig)

	return datagridParam
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

func expandDB(parameter *java.Parameter, config map[string]interface{}) {
	dbaInfo := config["database"].(*schema.Set)
	for _, i := range dbaInfo.List() {
		attrs := i.(map[string]interface{})
		parameter.DBServiceName = attrs["name"].(string)
		parameter.DBAName = attrs["username"].(string)
		parameter.DBAPassword = attrs["password"].(string)
		if v := attrs["network"]; v != nil {
			parameter.DBNetwork = v.(string)
		}
	}
}

func expandAdmin(parameter *java.Parameter, config map[string]interface{}) {
	adminInfo := config["admin"].(*schema.Set).List()
	attrs := adminInfo[0].(map[string]interface{})

	parameter.AdminUsername = attrs["username"].(string)
	parameter.AdminPassword = attrs["password"].(string)
	if v := attrs["port"]; v != nil {
		parameter.AdminPort = v.(int)
	}
}

func expandAppDBs(parameter *java.Parameter, config map[string]interface{}) {
	appDBInfo := config["app_db"].(*schema.Set)
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

func expandDomain(parameter *java.Parameter, config map[string]interface{}) {
	domainInfo := config["domain"].(*schema.Set)
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

func expandListener(parameter *java.Parameter, config map[string]interface{}) {
	listenerInfo := config["listener"].(*schema.Set)
	for _, i := range listenerInfo.List() {
		attrs := i.(map[string]interface{})
		parameter.ListenerPort = attrs["port"].(int)
		parameter.ListenerPortEnabled = attrs["enabled"].(bool)
		parameter.ListenerType = attrs["type"].(string)
	}
}

func expandManagedServers(parameter *java.Parameter, config map[string]interface{}) {
	msInfo := config["managed_servers"].(*schema.Set)
	for _, i := range msInfo.List() {
		attrs := i.(map[string]interface{})
		if val, ok := attrs["server_count"].(int); ok {
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

func expandNodeManager(parameter *java.Parameter, config map[string]interface{}) {
	nmInfo := config["node_manager"].(*schema.Set)
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

func expandPrivilegedPorts(parameter *java.Parameter, config map[string]interface{}) {
	portInfo := config["privileged_ports"].(*schema.Set)
	for _, i := range portInfo.List() {
		attrs := i.(map[string]interface{})
		if v := attrs["content_port"]; v != nil {
			parameter.PrivilegedContentPort = v.(int)
		}
		if v := attrs["listener_port"]; v != nil {
			parameter.PrivilegedListenerPort = v.(int)
		}
		if v := attrs["secured_content_port"]; v != nil {
			parameter.PrivilegedSecuredContentPort = v.(int)
		}
		if v := attrs["secured_listener_port"]; v != nil {
			parameter.PrivilegedSecuredListenerPort = v.(int)
		}
	}
}

func expandScalingUnit(parameter *java.Parameter, config map[string]interface{}) {
	scalingUnitInfo := config["scaling_unit"].(*schema.Set)
	for _, i := range scalingUnitInfo.List() {
		attrs := i.(map[string]interface{})
		if val, ok := attrs["name"].(string); ok && val != "" {
			parameter.ScalingUnitName = java.ServiceInstanceScalingUnitName(val)
		}
		if val := attrs["units"]; val != nil {
			expandUnits(parameter, val.(*schema.Set))
		}
	}
}

func expandUnits(parameter *java.Parameter, unitInfo *schema.Set) {
	parameter.ScalingUnitCount = unitInfo.Len()
	units := make([]java.ScalingUnit, unitInfo.Len())
	for i, val := range unitInfo.List() {
		attrs := val.(map[string]interface{})
		unit := java.ScalingUnit{
			HeapSize: attrs["heap_size"].(string),
			JVMCount: attrs["jvm_count"].(int),
			Shape:    java.ServiceInstanceShape(attrs["shape"].(string)),
			VMCount:  attrs["vm_count"].(int),
		}
		units[i] = unit
	}

	parameter.ScalingUnits = units
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

	result = append(result, db)
	return d.Set("database", result)
}

func readDomain(d *schema.ResourceData, name string, mode string) error {
	result := make([]map[string]interface{}, 0)

	domain := make(map[string]interface{})
	domain["name"] = name
	domain["mode"] = mode

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
