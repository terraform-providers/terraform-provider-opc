package opc

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	"github.com/hashicorp/go-oracle-terraform/storage"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/mitchellh/go-homedir"
)

func resourceOPCStorageObject() *schema.Resource {
	return &schema.Resource{
		Create: resourceOPCStorageObjectCreate,
		Read:   resourceOPCStorageObjectRead,
		Delete: resourceOPCStorageObjectDelete,
		// TODO: Add Import

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of the storage object",
			},
			"container": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Name of the storage container",
			},
			"content": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				Description:   "Raw content in string-form of the data",
				ConflictsWith: []string{"copy_from", "file"},
			},
			"file": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				Description:   "File path for the content to use for data",
				ConflictsWith: []string{"copy_from", "content"},
			},
			"content_disposition": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Overrides the behavior of the browser",
			},
			"content_encoding": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Set the content-encoding metadata",
			},
			"content_type": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Set the MIME type for the object",
			},
			"copy_from": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"content", "file"},
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					value := v.(string)
					if !strings.Contains(value, "/") {
						errors = append(errors, fmt.Errorf(
							"%q does not contain both a container and object name",
							k))
					}
					return
				},
			},
			"delete_at": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "Specify the number of seconds after which the system deletes the object",
			},
			"etag": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "MD5 checksum value of the request body. Unquoted. Strongly Recommended",
			},
			"transfer_encoding": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "Sets the transfer encoding. Can only be 'chunked' or Nil, requires Content-Length to be 0 if set",
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					value := v.(string)
					if value != "chunked" {
						errors = append(errors, fmt.Errorf(
							"%q must be either 'chunked' or nil", k))
					}
					return
				},
			},

			// Computed Attributes
			"accept_ranges": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Type of ranges that the object accepts",
			},
			"content_length": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Length of the object in bytes",
			},
			"last_modified": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Date and Time that the object was created/modified in ISO 8601",
			},
			"object_manifest": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The dynamic large-object manifest object",
			},
			"timestamp": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Date and Time in UNIX EPOCH when the account, container, or object was initially created at the current version",
			},
			"transaction_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Transaction ID of the request. Used for bug reports",
			},
		},
	}
}

func resourceOPCStorageObjectCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*OPCClient).storageClient.Objects()
	if client == nil {
		return fmt.Errorf("Storage client is not initialized. Make sure to use `storage_endpoint` variable or the `OPC_STORAGE_ENDPOINT` environment variable")
	}

	// Populate required attr
	input := &storage.CreateObjectInput{
		Name:      d.Get("name").(string),
		Container: d.Get("container").(string),
	}

	// Check for `content` or `file`.
	if v, ok := d.GetOk("content"); ok {
		// Read content as io.ReadSeeker
		content := v.(string)
		input.Body = bytes.NewReader([]byte(content))
	} else if v, ok := d.GetOk("file"); ok {
		// Read raw file
		source := v.(string)
		path, err := homedir.Expand(source)
		if err != nil {
			return fmt.Errorf("Error expanding homedir in file (%s): %s", source, err)
		}
		file, err := os.Open(path)
		if err != nil {
			return fmt.Errorf("Error opening Storage Object file (%s): %s", source, err)
		}
		input.Body = file
	} else if v, ok := d.GetOk("copy_from"); ok {
		input.CopyFrom = v.(string)
	} else {
		// One of the three attributes are required
		return fmt.Errorf("Must specify %q, %q, or %q field", "file", "copy_from", "content")
	}

	if v, ok := d.GetOk("content_disposition"); ok {
		input.ContentDisposition = v.(string)
	}

	if v, ok := d.GetOk("content_encoding"); ok {
		input.ContentEncoding = v.(string)
	}

	if v, ok := d.GetOk("content_type"); ok {
		input.ContentType = v.(string)
	}

	if v, ok := d.GetOk("delete_at"); ok {
		input.DeleteAt = v.(int)
	}

	if v, ok := d.GetOk("etag"); ok {
		input.ETag = v.(string)
	}

	if v, ok := d.GetOk("transfer_encoding"); ok {
		input.TransferEncoding = v.(string)
	}

	result, err := client.CreateObject(input)
	if err != nil {
		return fmt.Errorf("Error creating Object: %s", err)
	}

	d.SetId(result.ID)
	return resourceOPCStorageObjectRead(d, meta)
}

func resourceOPCStorageObjectRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*OPCClient).storageClient.Objects()
	if client == nil {
		return fmt.Errorf("Storage client is not initialized. Make sure to use `storage_endpoint` variable or the `OPC_STORAGE_ENDPOINT` environment variable")
	}

	input := &storage.GetObjectInput{
		ID: d.Id(),
	}

	result, err := client.GetObject(input)
	if err != nil {
		return fmt.Errorf("Error reading Storage Container Object (%s): %s", d.Id(), err)
	}

	if result == nil {
		d.SetId("")
		return nil
	}

	d.Set("name", result.Name)
	d.Set("container", result.Container)
	d.Set("content_disposition", result.ContentDisposition)
	d.Set("content_encoding", result.ContentEncoding)
	d.Set("content_length", result.ContentLength)
	d.Set("content_type", result.ContentType)
	d.Set("date", result.Date)
	d.Set("etag", result.Etag)
	d.Set("last_modified", result.LastModified)
	d.Set("delete_at", result.DeleteAt)
	d.Set("object_manifest", result.ObjectManifest)
	d.Set("timestamp", result.Timestamp)
	d.Set("transaction_id", result.TransactionID)

	return nil
}

func resourceOPCStorageObjectDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*OPCClient).storageClient.Objects()
	if client == nil {
		return fmt.Errorf("Storage client is not initialized. Make sure to use `storage_endpoint` variable or the `OPC_STORAGE_ENDPOINT` environment variable")
	}

	input := &storage.DeleteObjectInput{
		ID: d.Id(),
	}
	if err := client.DeleteObject(input); err != nil {
		return fmt.Errorf("Error deleting Storage Container Object (%s): %s", d.Id(), err)
	}

	return nil
}
