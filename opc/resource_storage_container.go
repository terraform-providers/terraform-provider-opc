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
        Elem:     &schema.Schema{Type: schema.TypeString},
      },
      "write_acls": {
        Type:     schema.TypeList,
        Optional: true,
        Elem:     &schema.Schema{Type: schema.TypeString},
      },
      "allowed_origins": {
        Type:     schema.TypeList,
        Optional: true,
        Elem:     &schema.Schema{Type: schema.TypeString},
      },
      "primary_key": {
        Type: schema.TypeString,
        Optional: true,
      },
      "secondary_key": {
        Type: schema.TypeString,
        Optional: true,
      },
      "max_age": {
        Type: schema.TypeInt,
        Optional: true,
      },
		},
	}
}

func resourceOPCStorageContainerCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*OPCClient).storageClient

	input := storage.CreateContainerInput{
		Name: d.Get("name").(string)
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
  if primaryKey, ok := d.GetOk("primary_key"); ok {
    input.PrimaryKey = primaryKey.(string)
  }
  if secondaryKey, ok := d.GetOk("secondary_key"); ok {
    input.SecondaryKey = secondaryKey.(string)
  }
  if maxAge, ok := d.GetOk("max_age"); ok {
    input.MaxAge = maxAge.(int)
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

	name := d.Id()
	input := compute.GetStorageContainerInput{
		Name: name,
	}

	result, err := storageClient.GetStorageContainer(&input)
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
  if err := setStringList(d, "read_acls", result.ReadACLs); err != nil {
    return err
  }
  if err := setStringList(d, "write_acls", result.WriteACLs); err != nil {
    return err
  }
  if err := setStringList(d, "allowed_origins", result.AllowedOrigins); err != nil {
    return err
  }

	return nil
}

func resourceOPCStorageContainerDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*OPCClient).storageClient

	name := d.Id()
	input := compute.DeleteStorageContainerInput{
		Name: name,
	}
	if err := client.DeleteStorageContainer(&input); err != nil {
		return fmt.Errorf("Error deleting Storage Container '%s': %s", name, err)
	}

	return nil
}
