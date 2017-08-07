package opc

import (
	"fmt"
	"log"
	"strconv"
	"time"

	opcClient "github.com/hashicorp/go-oracle-terraform/client"
	"github.com/hashicorp/go-oracle-terraform/database"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

func resourceOPCDatabaseServiceInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceOPCDatabaseServiceInstanceCreate,
		Read:   resourceOPCDatabaseServiceInstanceRead,
		Delete: resourceOPCDatabaseServiceInstanceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Delete: schema.DefaultTimeout(60 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"edition": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(database.ServiceInstanceStandardEdition),
					string(database.ServiceInstanceEnterpriseEdition),
					string(database.ServiceInstanceEnterpriseEditionHighPerformance),
					string(database.ServiceInstanceEnterpriseEditionExtremePerformance),
				}, true),
			},
			"level": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(database.ServiceInstanceLevelPAAS),
					string(database.ServiceInstanceLevelBasic),
				}, true),
			},
			"shape": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"subscription_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(database.ServiceInstanceSubscriptionTypeHourly),
					string(database.ServiceInstanceSubscriptionTypeMonthly),
				}, true),
			},
			"version": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"vm_public_key": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"parameter": {
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"db_demo": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"admin_password": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"backup_destination": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
							ValidateFunc: validation.StringInSlice([]string{
								string(database.ServiceInstanceBackupDestinationBoth),
								string(database.ServiceInstanceBackupDestinationOSS),
								string(database.ServiceInstanceBackupDestinationNone),
							}, true),
						},
						"char_set": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
							Default:  "AL32UTF8",
						},
						"disaster_recovery": {
							Type:     schema.TypeBool,
							Optional: true,
							ForceNew: true,
							Default:  false,
						},
						"failover_database": {
							Type:     schema.TypeBool,
							Optional: true,
							ForceNew: true,
							Default:  false,
						},
						"golden_gate": {
							Type:     schema.TypeBool,
							Optional: true,
							ForceNew: true,
							Default:  false,
						},
						"is_rac": {
							Type:     schema.TypeBool,
							Optional: true,
							ForceNew: true,
							Default:  false,
						},
						"n_char_set": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
							Default:  database.ServiceInstanceNCharSetUTF16,
							ValidateFunc: validation.StringInSlice([]string{
								string(database.ServiceInstanceNCharSetUTF16),
								string(database.ServiceInstanceNCharSetUTF8),
							}, true),
						},
						"pdb_name": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
							Default:  "pdb1",
						},
						"sid": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
							Default:  "ORCL",
						},
						"snapshot_name": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"source_service_name": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"timezone": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
							Default:  "UTC",
						},
						"type": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
							Default:  database.ServiceInstanceTypeDB,
							ValidateFunc: validation.StringInSlice([]string{
								string(database.ServiceInstanceTypeDB),
							}, true),
						},
						"usable_storage": {
							Type:         schema.TypeInt,
							Required:     true,
							ForceNew:     true,
							ValidateFunc: validation.IntBetween(15, 2048),
						},
					},
				},
			},
			"ibkup": {
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cloud_storage_password": {
							Type:      schema.TypeString,
							Optional:  true,
							ForceNew:  true,
							Sensitive: true,
							Computed:  true,
						},
						"cloud_storage_username": {
							Type:      schema.TypeString,
							Optional:  true,
							ForceNew:  true,
							Computed:  true,
							Sensitive: true,
						},
						"database_id": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"decryption_key": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
						"ibkup_wallet_file_content": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
					},
				},
			},
			"cloud_storage": {
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"container": {
							Type:     schema.TypeString,
							Required: true,
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
						"create_if_missing": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
							ForceNew: true,
						},
					},
				},
			},
			"apex_url": {
				Type:     schema.TypeString,
				ForceNew: true,
				Computed: true,
			},
			"backup_supported_version": {
				Type:     schema.TypeString,
				ForceNew: true,
				Computed: true,
			},
			"cloud_storage_container": {
				Type:     schema.TypeString,
				ForceNew: true,
				Computed: true,
			},
			"compute_site_name": {
				Type:     schema.TypeString,
				ForceNew: true,
				Computed: true,
			},
			"connect_descriptor_with_public_ip": {
				Type:     schema.TypeString,
				ForceNew: true,
				Computed: true,
			},
			"created_by": {
				Type:     schema.TypeString,
				ForceNew: true,
				Computed: true,
			},
			"creation_time": {
				Type:     schema.TypeString,
				ForceNew: true,
				Computed: true,
			},
			"current_version": {
				Type:     schema.TypeString,
				ForceNew: true,
				Computed: true,
			},
			"dbaas_monitor_url": {
				Type:     schema.TypeString,
				ForceNew: true,
				Computed: true,
			},
			"em_url": {
				Type:     schema.TypeString,
				ForceNew: true,
				Computed: true,
			},
			"failover_database": {
				Type:     schema.TypeBool,
				ForceNew: true,
				Computed: true,
			},
			"glassfish_url": {
				Type:     schema.TypeString,
				ForceNew: true,
				Computed: true,
			},
			"hdg_prem_ip": {
				Type:     schema.TypeString,
				ForceNew: true,
				Computed: true,
			},
			"hybrid_db": {
				Type:     schema.TypeString,
				ForceNew: true,
				Computed: true,
			},
			"identity_domain": {
				Type:     schema.TypeString,
				ForceNew: true,
				Computed: true,
			},
			"ip_network": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ip_reservations": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"jaas_instances_using_service": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"listener_port": {
				Type:     schema.TypeString,
				ForceNew: true,
				Computed: true,
			},
			"num_ip_reservations": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"num_nodes": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"rac_database": {
				Type:     schema.TypeBool,
				ForceNew: true,
				Computed: true,
			},
			"region": {
				Type:     schema.TypeString,
				ForceNew: true,
				Computed: true,
			},
			"uri": {
				Type:     schema.TypeString,
				ForceNew: true,
				Computed: true,
			},
			"sm_plugin_version": {
				Type:     schema.TypeString,
				ForceNew: true,
				Computed: true,
			},
			"total_shared_storage": {
				Type:     schema.TypeString,
				ForceNew: true,
				Computed: true,
			},
		},
	}
}

func resourceOPCDatabaseServiceInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Resource state: %#v", d.State())

	log.Print("[DEBUG] Creating database service instance")

	client := meta.(*OPCClient).databaseClient.ServiceInstanceClient()
	if client == nil {
		return fmt.Errorf("Database Client is not initialized. Make sure to use `database_endpoint` variable or `OPC_DATABASE_ENDPOINT` env variable")
	}
	input := database.CreateServiceInstanceInput{
		Name:             d.Get("name").(string),
		Edition:          database.ServiceInstanceEdition(d.Get("edition").(string)),
		Level:            database.ServiceInstanceLevel(d.Get("level").(string)),
		Shape:            database.ServiceInstanceShape(d.Get("shape").(string)),
		SubscriptionType: database.ServiceInstanceSubscriptionType(d.Get("subscription_type").(string)),
		Version:          database.ServiceInstanceVersion(d.Get("version").(string)),
		VMPublicKey:      d.Get("vm_public_key").(string),
	}
	if description, ok := d.GetOk("description"); ok {
		input.Description = description.(string)
	}

	// Only the PaaS level can have a parameter.
	if input.Level == database.ServiceInstanceLevelPAAS {
		input.Parameter = expandParameter(client, d)
	}

	info, err := client.CreateServiceInstance(&input)
	if err != nil {
		return fmt.Errorf("Error creating DatabaseServiceInstance: %+v", err)
	}

	d.SetId(info.Name)
	return resourceOPCDatabaseServiceInstanceRead(d, meta)
}

func resourceOPCDatabaseServiceInstanceRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Resource state: %#v", d.State())
	client := meta.(*OPCClient).databaseClient.ServiceInstanceClient()
	if client == nil {
		return fmt.Errorf("Database Client is not initialized. Make sure to use `database_endpoint` variable or `OPC_DATABASE_ENDPOINT` env variable")
	}

	log.Printf("[DEBUG] Reading state of ip reservation %s", d.Id())
	getInput := database.GetServiceInstanceInput{
		Name: d.Id(),
	}

	result, err := client.GetServiceInstance(&getInput)
	if err != nil {
		// DatabaseServiceInstance does not exist
		if opcClient.WasNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error reading database service instance %s: %+v", d.Id(), err)
	}

	if result == nil {
		d.SetId("")
		return nil
	}

	log.Printf("[DEBUG] Read state of database service instance %s: %#v", d.Id(), result)
	d.Set("name", result.Name)
	d.Set("description", result.Description)
	d.Set("apex_url", result.ApexURL)
	d.Set("backup_destination", result.BackupDestination)
	d.Set("backup_supported_version", result.BackupSupportedVersion)
	d.Set("char_set", result.CharSet)
	d.Set("cloud_storage_container", result.CloudStorageContainer)
	d.Set("compute_site_name", result.ComputeSiteName)
	d.Set("connect_descriptor", result.ConnectDescriptor)
	d.Set("connect_descriptor_with_public_ip", result.ConnectorDescriptorWithPublicIP)
	d.Set("created_by", result.CreatedBy)
	d.Set("creation_time", result.CreationTime)
	d.Set("current_version", result.CurrentVersion)
	d.Set("dbaas_monitor_url", result.DBAASMonitorURL)
	d.Set("edition", result.Edition)
	d.Set("em_url", result.EMURL)
	d.Set("failover_database", result.FailoverDatabase)
	d.Set("glassfish_url", result.GlassFishURL)
	d.Set("hdg_prem_ip", result.HDGPremIP)
	d.Set("hybrid_db", result.HybridDG)
	d.Set("identity_domain", result.IdentityDomain)
	d.Set("ip_network", result.IPNetwork)
	d.Set("ip_reservations", result.IPReservations)
	d.Set("jaas_instances_using_service", result.JAASInstancesUsingService)
	d.Set("level", result.Level)
	d.Set("listener_port", result.ListenerPort)
	d.Set("n_char_set", result.NCharSet)
	d.Set("num_ip_reservations", result.NumIPReservations)
	d.Set("num_nodes", result.NumNodes)
	d.Set("pdb_name", result.PDBName)
	d.Set("rac_database", result.PDBName)
	d.Set("region", result.Region)
	d.Set("uri", result.URI)
	d.Set("shape", result.Shape)
	d.Set("sid", result.SID)
	d.Set("sm_plugin_version", result.SMPluginVersion)
	d.Set("subscription_type", result.SubscriptionType)
	d.Set("timezone", result.Timezone)
	d.Set("total_shared_storage", result.TotalSharedStorage)
	d.Set("version", result.Version)
	// TODO Add parameter and vm_public_key to read when they get added to api.
	return nil
}

func resourceOPCDatabaseServiceInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Resource state: %#v", d.State())
	client := meta.(*OPCClient).databaseClient.ServiceInstanceClient()
	if client == nil {
		return fmt.Errorf("Database Client is not initialized. Make sure to use `database_endpoint` variable or `OPC_DATABASE_ENDPOINT` env variable")
	}
	name := d.Id()

	log.Printf("[DEBUG] Deleting DatabaseServiceInstance: %v", name)

	input := database.DeleteServiceInstanceInput{
		Name: name,
	}
	if err := client.DeleteServiceInstance(&input); err != nil {
		return fmt.Errorf("Error deleting DatabaseServiceInstance: %+v", err)
	}
	return nil
}

func expandParameter(client *database.ServiceInstanceClient, d *schema.ResourceData) database.ParameterInput {
	parameterInfo := d.Get("parameter").(*schema.Set)
	var parameter database.ParameterInput
	for _, i := range parameterInfo.List() {
		attrs := i.(map[string]interface{})
		parameter = database.ParameterInput{
			AdminPassword:     attrs["admin_password"].(string),
			BackupDestination: database.ServiceInstanceBackupDestination(attrs["backup_destination"].(string)),
			CharSet:           attrs["char_set"].(string),
			DisasterRecovery:  attrs["disaster_recovery"].(bool),
			FailoverDatabase:  attrs["failover_database"].(bool),
			GoldenGate:        attrs["golden_gate"].(bool),
			IsRAC:             attrs["is_rac"].(bool),
			NCharSet:          database.ServiceInstanceNCharSet(attrs["n_char_set"].(string)),
			PDBName:           attrs["pdb_name"].(string),
			SID:               attrs["sid"].(string),
			Timezone:          attrs["timezone"].(string),
			Type:              database.ServiceInstanceType(attrs["type"].(string)),
			UsableStorage:     strconv.Itoa(attrs["usable_storage"].(int)),
		}

		if val, ok := attrs["snapshot_name"].(string); ok && val != "" {
			parameter.SnapshotName = val
		}
		if val, ok := attrs["source_service_name"].(string); ok && val != "" {
			parameter.SourceServiceName = val
		}
		if val, ok := attrs["db_demo"].(string); ok {
			addParam := database.AdditionalParameters{
				DBDemo: val,
			}
			parameter.AdditionalParameters = addParam
		}
	}
	expandIbkup(d, &parameter)
	expandCloudStorage(d, &parameter)
	return parameter
}

func expandIbkup(d *schema.ResourceData, parameter *database.ParameterInput) {
	ibkupInfo := d.Get("ibkup").(*schema.Set)
	for _, i := range ibkupInfo.List() {
		attrs := i.(map[string]interface{})
		parameter.IBKUP = true
		parameter.IBKUPDatabaseID = attrs["cloud_storage_username"].(string)
		if val, ok := attrs["decryption_key"].(string); ok && val != "" {
			parameter.IBKUPCloudStorageUser = val
		}
		if val, ok := attrs["cloud_storage_password"].(string); ok && val != "" {
			parameter.IBKUPCloudStoragePassword = val
		}
		if val, ok := attrs["decryption_key"].(string); ok && val != "" {
			parameter.IBKUPDecryptionKey = val
		}
		if val, ok := attrs["wallet_file_content"].(string); ok && val != "" {
			parameter.IBKUPWalletFileContent = val
		}
	}
}

func expandCloudStorage(d *schema.ResourceData, parameter *database.ParameterInput) {
	cloudStorageInfo := d.Get("cloud_storage").(*schema.Set)
	for _, i := range cloudStorageInfo.List() {
		attrs := i.(map[string]interface{})
		parameter.CloudStorageContainer = attrs["container"].(string)
		parameter.CreateStorageContainerIfMissing = attrs["create_if_missing"].(bool)
		if val, ok := attrs["username"].(string); ok && val != "" {
			parameter.CloudStorageUsername = val
		}
		if val, ok := attrs["password"].(string); ok && val != "" {
			parameter.CloudStoragePassword = val
		}
	}
}
