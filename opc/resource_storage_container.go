package opc

import (
	"fmt"

	"github.com/hashicorp/go-oracle-terraform/client"
	"github.com/hashicorp/go-oracle-terraform/storage"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceOPCStorageContainer() *schema.Resource {
	return &schema.Resource{
		Create: resourceOPCStorageContainerCreate,
		Read:   resourceOPCStorageContainerRead,
		Delete: resourceOPCStorageContainerDelete,
		Update: resourceOPCStorageContainerUpdate,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"read_acls": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"write_acls": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"allowed_origins": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"exposed_headers": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"primary_key": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"secondary_key": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"max_age": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"quota_bytes": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"quota_count": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"metadata": {
				Type:     schema.TypeMap,
				Optional: true,
			},
			// "georeplication_policy": {
			// 	Type:     schema.TypeList,
			// 	Optional: true,
			// 	Computed: true,
			// 	Elem:     &schema.Schema{Type: schema.TypeString},
			// },
		},
	}
}

func resourceOPCStorageContainerCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*OPCClient).storageClient
	if client == nil {
		return fmt.Errorf("Storage Client is not initialized. Make sure to use `storage_endpoint` variable or `OPC_STORAGE_ENDPOINT env variable`")
	}

	input := storage.CreateContainerInput{
		Name: d.Get("name").(string),
	}
	if readAcls := getStringList(d, "read_acls"); len(readAcls) > 0 {
		input.ReadACLs = readAcls
	}
	if writeAcls := getStringList(d, "write_acls"); len(writeAcls) > 0 {
		input.WriteACLs = writeAcls
	}
	if allowedOrigins := getStringList(d, "allowed_origins"); len(allowedOrigins) > 0 {
		input.AllowedOrigins = allowedOrigins
	}
	if exposedHeaders := getStringList(d, "exposed_headers"); len(exposedHeaders) > 0 {
		input.ExposedHeaders = exposedHeaders
	}
	if primaryKey, ok := d.GetOk("primary_key"); ok {
		input.PrimaryKey = primaryKey.(string)
	}
	if secondaryKey, ok := d.GetOk("secondary_key"); ok {
		input.SecondaryKey = secondaryKey.(string)
	}
	if maxAge, ok := d.GetOk("max_age"); ok {
		input.MaxAge = maxAge.(int)
	}
	if quotaBytes, ok := d.GetOk("quota_bytes"); ok {
		input.QuotaBytes = quotaBytes.(int)
	}
	if quotaCount, ok := d.GetOk("quota_count"); ok {
		input.QuotaCount = quotaCount.(int)
	}

	if v, ok := d.GetOk("metadata"); ok {
		metadata := make(map[string]string)
		for name, value := range v.(map[string]interface{}) {
			metadata[name] = value.(string)
		}
		input.CustomMetadata = metadata
	}

	info, err := client.CreateContainer(&input)
	if err != nil {
		return fmt.Errorf("Error creating Storage Container: %s", err)
	}

	d.SetId(info.Name)

	return resourceOPCStorageContainerRead(d, meta)
}

func resourceOPCStorageContainerRead(d *schema.ResourceData, meta interface{}) error {
	storageClient := meta.(*OPCClient).storageClient
	if storageClient == nil {
		return fmt.Errorf("Storage Client is not initialized. Make sure to use `storage_endpoint` variable or `OPC_STORAGE_ENDPOINT env variable`")
	}

	name := d.Id()
	input := storage.GetContainerInput{
		Name: name,
	}

	result, err := storageClient.GetContainer(&input)
	if err != nil {
		// Storage Container does not exist
		if client.WasNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error reading Storage Container '%s': %s", name, err)
	}

	if result == nil {
		d.SetId("")
		return nil
	}

	d.Set("name", result.Name)
	d.Set("primary_key", result.PrimaryKey)
	d.Set("secondary_key", result.SecondaryKey)
	d.Set("max_age", result.MaxAge)
	d.Set("quota_bytes", result.QuotaBytes)
	d.Set("quota_count", result.QuotaCount)
	d.Set("metadata", result.CustomMetadata)

	if err := setStringList(d, "read_acls", result.ReadACLs); err != nil {
		return err
	}
	if err := setStringList(d, "write_acls", result.WriteACLs); err != nil {
		return err
	}
	if err := setStringList(d, "allowed_origins", result.AllowedOrigins); err != nil {
		return err
	}
	if err := setStringList(d, "exposed_headers", result.ExposedHeaders); err != nil {
		return err
	}

	return nil
}

func resourceOPCStorageContainerDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*OPCClient).storageClient
	if client == nil {
		return fmt.Errorf("Storage Client is not initialized. Make sure to use `storage_endpoint` variable or `OPC_STORAGE_ENDPOINT env variable`")
	}

	name := d.Id()
	input := storage.DeleteContainerInput{
		Name: name,
	}
	if err := client.DeleteContainer(&input); err != nil {
		return fmt.Errorf("Error deleting Storage Container '%s': %s", name, err)
	}

	return nil
}

func resourceOPCStorageContainerUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*OPCClient).storageClient
	if client == nil {
		return fmt.Errorf("Storage Client is not initialized. Make sure to use `storage_endpoint` variable or `OPC_STORAGE_ENDPOINT env variable`")
	}

	input := storage.UpdateContainerInput{
		Name: d.Get("name").(string),
	}
	if readAcls := getStringList(d, "read_acls"); len(readAcls) > 0 {
		input.ReadACLs = readAcls
	}
	if writeAcls := getStringList(d, "write_acls"); len(writeAcls) > 0 {
		input.WriteACLs = writeAcls
	}
	if allowedOrigins := getStringList(d, "allowed_origins"); len(allowedOrigins) > 0 {
		input.AllowedOrigins = allowedOrigins
	}
	if exposedHeaders := getStringList(d, "exposed_headers"); len(exposedHeaders) > 0 {
		input.ExposedHeaders = exposedHeaders
	}
	if primaryKey, ok := d.GetOk("primary_key"); ok {
		input.PrimaryKey = primaryKey.(string)
	}
	if secondaryKey, ok := d.GetOk("secondary_key"); ok {
		input.SecondaryKey = secondaryKey.(string)
	}
	if maxAge, ok := d.GetOk("max_age"); ok {
		input.MaxAge = maxAge.(int)
	}
	if quotaBytes, ok := d.GetOk("quota_bytes"); ok {
		input.QuotaBytes = quotaBytes.(int)
	}
	if quotaCount, ok := d.GetOk("quota_count"); ok {
		input.QuotaCount = quotaCount.(int)
	}

	if v, ok := d.GetOk("metadata"); ok {
		metadata := make(map[string]string)
		for name, value := range v.(map[string]interface{}) {
			metadata[name] = value.(string)
		}
		input.CustomMetadata = metadata
	}

	info, err := client.UpdateContainer(&input)
	if err != nil {
		return fmt.Errorf("Error updating Storage Container: %s", err)
	}

	d.SetId(info.Name)

	return resourceOPCStorageContainerRead(d, meta)
}
