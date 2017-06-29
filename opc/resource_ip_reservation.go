package opc

import (
	"fmt"
	"log"

	"github.com/hashicorp/go-oracle-terraform/client"
	"github.com/hashicorp/go-oracle-terraform/compute"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

func resourceOPCIPReservation() *schema.Resource {
	return &schema.Resource{
		Create: resourceOPCIPReservationCreate,
		Read:   resourceOPCIPReservationRead,
		Delete: resourceOPCIPReservationDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"permanent": {
				Type:     schema.TypeBool,
				Required: true,
				ForceNew: true,
			},
			"parent_pool": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  string(compute.PublicReservationPool),
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(compute.PublicReservationPool),
				}, true),
			},
			"tags": tagsForceNewSchema(),
			"ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceOPCIPReservationCreate(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Resource state: %#v", d.State())

	reservation := compute.CreateIPReservationInput{
		Name:       d.Get("name").(string),
		ParentPool: compute.IPReservationPool(d.Get("parent_pool").(string)),
		Permanent:  d.Get("permanent").(bool),
	}

	tags := getStringList(d, "tags")
	if len(tags) != 0 {
		reservation.Tags = tags
	}

	log.Printf("[DEBUG] Creating ip reservation from parent_pool %s with tags=%s",
		reservation.ParentPool, reservation.Tags)

	client := meta.(*OPCClient).computeClient.IPReservations()
	info, err := client.CreateIPReservation(&reservation)
	if err != nil {
		return fmt.Errorf("Error creating ip reservation from parent_pool %s with tags=%s: %s",
			reservation.ParentPool, reservation.Tags, err)
	}

	d.SetId(info.Name)
	return resourceOPCIPReservationRead(d, meta)
}

func resourceOPCIPReservationRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Resource state: %#v", d.State())
	computeClient := meta.(*OPCClient).computeClient.IPReservations()

	log.Printf("[DEBUG] Reading state of ip reservation %s", d.Id())
	input := compute.GetIPReservationInput{
		Name: d.Id(),
	}

	result, err := computeClient.GetIPReservation(&input)
	if err != nil {
		// IP Reservation does not exist
		if client.WasNotFoundError(err) {
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error reading ip reservation %s: %s", d.Id(), err)
	}

	if result == nil {
		d.SetId("")
		return nil
	}

	log.Printf("[DEBUG] Read state of ip reservation %s: %#v", d.Id(), result)
	d.Set("name", result.Name)
	d.Set("parent_pool", result.ParentPool)
	d.Set("permanent", result.Permanent)

	if err := setStringList(d, "tags", result.Tags); err != nil {
		return err
	}

	d.Set("ip", result.IP)
	return nil
}

func resourceOPCIPReservationDelete(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[DEBUG] Resource state: %#v", d.State())
	client := meta.(*OPCClient).computeClient.IPReservations()

	log.Printf("[DEBUG] Deleting ip reservation %s", d.Id())

	input := compute.DeleteIPReservationInput{
		Name: d.Id(),
	}
	if err := client.DeleteIPReservation(&input); err != nil {
		return fmt.Errorf("Error deleting ip reservation %s", d.Id())
	}
	return nil
}
