package opc

import (
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/hashicorp/go-oracle-terraform/compute"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

func orchestrationInstanceSchema() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Required: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"id": {
					Type:     schema.TypeString,
					Computed: true,
				},

				"name": {
					Type:     schema.TypeString,
					Required: true,
					ForceNew: true,
				},

				"shape": {
					Type:     schema.TypeString,
					Required: true,
					ForceNew: true,
				},

				/////////////////////////
				// Optional Attributes //
				/////////////////////////
				"persistent": {
					Type:     schema.TypeBool,
					Optional: true,
					Default:  false,
				},

				"instance_attributes": {
					Type:     schema.TypeString,
					Computed: true,
				},

				"boot_order": {
					Type:     schema.TypeList,
					Optional: true,
					ForceNew: true,
					Elem:     &schema.Schema{Type: schema.TypeInt},
				},

				"hostname": {
					Type:     schema.TypeString,
					Optional: true,
					Computed: true,
					ForceNew: true,
				},

				"image_list": {
					Type:     schema.TypeString,
					Optional: true,
					ForceNew: true,
				},

				"label": {
					Type:     schema.TypeString,
					Optional: true,
					Computed: true,
					ForceNew: true,
				},

				"networking_info": {
					Type:     schema.TypeList,
					Optional: true,
					Computed: true,
					ForceNew: true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"dns": {
								// Required for Shared Network Interface, will default if unspecified, however
								// Optional for IP Network Interface
								Type:     schema.TypeList,
								Optional: true,
								Computed: true,
								ForceNew: true,
								Elem:     &schema.Schema{Type: schema.TypeString},
							},

							"index": {
								Type:     schema.TypeInt,
								ForceNew: true,
								Required: true,
							},

							"ip_address": {
								// Optional, IP Network only
								Type:     schema.TypeString,
								ForceNew: true,
								Optional: true,
							},

							"ip_network": {
								// Required for an IP Network Interface
								Type:     schema.TypeString,
								ForceNew: true,
								Optional: true,
							},

							"is_default_gateway": {
								// Optional, IP Network only
								Type:     schema.TypeBool,
								ForceNew: true,
								Optional: true,
							},

							"mac_address": {
								// Optional, IP Network Only
								Type:     schema.TypeString,
								ForceNew: true,
								Computed: true,
								Optional: true,
							},

							"name_servers": {
								// Optional, IP Network + Shared Network
								Type:     schema.TypeList,
								Optional: true,
								ForceNew: true,
								Elem:     &schema.Schema{Type: schema.TypeString},
							},

							"nat": {
								// Optional for IP Network
								// Required for Shared Network
								Type:     schema.TypeList,
								Optional: true,
								ForceNew: true,
								Elem:     &schema.Schema{Type: schema.TypeString},
							},

							"search_domains": {
								// Optional, IP Network + Shared Network
								Type:     schema.TypeList,
								Optional: true,
								ForceNew: true,
								Elem:     &schema.Schema{Type: schema.TypeString},
							},

							"sec_lists": {
								// Required, Shared Network only. Will default if unspecified however
								Type:     schema.TypeList,
								Optional: true,
								Computed: true,
								ForceNew: true,
								Elem:     &schema.Schema{Type: schema.TypeString},
							},

							"shared_network": {
								Type:     schema.TypeBool,
								Optional: true,
								ForceNew: true,
								Default:  false,
							},

							"vnic": {
								// Optional, IP Network only.
								Type:     schema.TypeString,
								ForceNew: true,
								Optional: true,
							},

							"vnic_sets": {
								// Optional, IP Network only.
								Type:     schema.TypeList,
								Optional: true,
								ForceNew: true,
								Elem:     &schema.Schema{Type: schema.TypeString},
							},
						},
					},
				},

				"reverse_dns": {
					Type:     schema.TypeBool,
					Optional: true,
					Default:  true,
					ForceNew: true,
				},

				"ssh_keys": {
					Type:     schema.TypeList,
					Optional: true,
					ForceNew: true,
					Elem:     &schema.Schema{Type: schema.TypeString},
				},

				"storage": {
					Type:     schema.TypeList,
					Optional: true,
					ForceNew: true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"index": {
								Type:         schema.TypeInt,
								Required:     true,
								ForceNew:     true,
								ValidateFunc: validation.IntBetween(1, 10),
							},
							"volume": {
								Type:     schema.TypeString,
								Required: true,
								ForceNew: true,
							},
							"name": {
								Type:     schema.TypeString,
								Computed: true,
							},
						},
					},
				},

				"tags": tagsForceNewSchema(),

				/////////////////////////
				// Computed Attributes //
				/////////////////////////
				"attributes": {
					Type:     schema.TypeString,
					Computed: true,
				},

				"availability_domain": {
					Type:     schema.TypeString,
					Computed: true,
				},

				"domain": {
					Type:     schema.TypeString,
					Computed: true,
				},

				"entry": {
					Type:     schema.TypeInt,
					Computed: true,
				},

				"fingerprint": {
					Type:     schema.TypeString,
					Computed: true,
				},

				"fqdn": {
					Type:     schema.TypeString,
					Computed: true,
				},

				"image_format": {
					Type:     schema.TypeString,
					Computed: true,
				},

				"ip_address": {
					Type:     schema.TypeString,
					Computed: true,
				},

				"placement_requirements": {
					Type:     schema.TypeList,
					Computed: true,
					Elem:     &schema.Schema{Type: schema.TypeString},
				},

				"platform": {
					Type:     schema.TypeString,
					Computed: true,
				},

				"priority": {
					Type:     schema.TypeString,
					Computed: true,
				},

				"quota_reservation": {
					Type:     schema.TypeString,
					Computed: true,
				},

				"relationships": {
					Type:     schema.TypeList,
					Computed: true,
					Elem:     &schema.Schema{Type: schema.TypeString},
				},

				"resolvers": {
					Type:     schema.TypeList,
					Computed: true,
					Elem:     &schema.Schema{Type: schema.TypeString},
				},

				"site": {
					Type:     schema.TypeString,
					Computed: true,
				},

				"start_time": {
					Type:     schema.TypeString,
					Computed: true,
				},

				"state": {
					Type:     schema.TypeString,
					Computed: true,
				},

				"vcable": {
					Type:     schema.TypeString,
					Computed: true,
				},

				"virtio": {
					Type:     schema.TypeBool,
					Computed: true,
				},

				"vnc_address": {
					Type:     schema.TypeString,
					Computed: true,
				},
			},
		},
	}
}

// We can create multiple instances with an orchestration so we pass a prefix in to obtain
// the CreateInput for each instance.
func expandCreateInstanceInput(prefix string, d *schema.ResourceData) (*compute.CreateInstanceInput, error) {
	// Get Required Attributes
	input := &compute.CreateInstanceInput{
		Name:  d.Get(fmt.Sprintf("%s.name", prefix)).(string),
		Shape: d.Get(fmt.Sprintf("%s.shape", prefix)).(string),
	}

	bootOrder := getIntList(d, fmt.Sprintf("%s.boot_order", prefix))
	imageList := d.Get(fmt.Sprintf("%s.image_list", prefix))
	if len(bootOrder) == 0 && imageList.(string) == "" {
		return nil, fmt.Errorf("One of boot_order or image_list must be set for instance to be created")
	}

	if len(bootOrder) > 0 {
		input.BootOrder = bootOrder
	}

	if imageList != nil {
		input.ImageList = imageList.(string)
	}

	if v, ok := d.GetOk(fmt.Sprintf("%s.hostname", prefix)); ok {
		input.Hostname = v.(string)
	}

	if v, ok := d.GetOk(fmt.Sprintf("%s.label", prefix)); ok {
		input.Label = v.(string)
	}

	interfaces, err := expandNetworkInterfacesFromConfig(d, prefix)
	if err != nil {
		return nil, err
	}
	if interfaces != nil {
		input.Networking = interfaces
	}

	if v, ok := d.GetOk(fmt.Sprintf("%s.reverse_dns", prefix)); ok {
		input.ReverseDNS = v.(bool)
	}

	if sshKeys := getStringList(d, fmt.Sprintf("%s.ssh_keys", prefix)); len(sshKeys) > 0 {
		input.SSHKeys = sshKeys
	}

	storage := expandStorageAttachments(d, prefix)
	if len(storage) > 0 {
		input.Storage = storage
	}

	if tags := getStringList(d, fmt.Sprintf("%s.tags", prefix)); len(tags) > 0 {
		input.Tags = tags
	}

	return input, nil
}

// Populates and validates shared network and ip network interfaces to return the of map
// objects needed to create/update an instance's networking_info
func expandNetworkInterfacesFromConfig(d *schema.ResourceData, prefix string) (map[string]compute.NetworkingInfo, error) {
	interfaces := make(map[string]compute.NetworkingInfo)

	if v, ok := d.GetOk(fmt.Sprintf("%s.networking_info", prefix)); ok {
		vL := v.([]interface{})
		for _, v := range vL {
			ni := v.(map[string]interface{})
			index, ok := ni["index"].(int)
			if !ok {
				return nil, fmt.Errorf("Index not specified for network interface: %v", ni)
			}

			deviceIndex := fmt.Sprintf("eth%d", index)

			// Verify that the network interface doesn't already exist
			if _, ok := interfaces[deviceIndex]; ok {
				return nil, fmt.Errorf("Duplicate Network interface at eth%d already specified", index)
			}

			// Determine if we're creating a shared network interface or an IP Network interface
			info := compute.NetworkingInfo{}
			var err error
			if ni["shared_network"].(bool) {
				// Populate shared network parameters
				info, err = readSharedNetworkFromConfig(ni)
				// Set 'model' since we're configuring a shared network interface
				info.Model = compute.NICDefaultModel
			} else {
				// Populate IP Network Parameters
				info, err = readIPNetworkFromConfig(ni)
			}
			if err != nil {
				return nil, err
			}
			// And you may find yourself in a beautiful house, with a beautiful wife
			// And you may ask yourself, well, how did I get here?
			interfaces[deviceIndex] = info
		}
	}

	return interfaces, nil
}

func expandStorageAttachments(d *schema.ResourceData, prefix string) []compute.StorageAttachmentInput {
	storageAttachments := []compute.StorageAttachmentInput{}
	storage := d.Get(fmt.Sprintf("%s.storage", prefix)).([]interface{})
	for _, i := range storage {
		attrs := i.(map[string]interface{})
		storageAttachments = append(storageAttachments, compute.StorageAttachmentInput{
			Index:  attrs["index"].(int),
			Volume: attrs["volume"].(string),
		})
	}
	return storageAttachments
}

func flattenOrchestratedInstances(d *schema.ResourceData, meta interface{}, objects []compute.Object) (interface{}, error) {
	instanceClient := meta.(*OPCClient).computeClient.Instances()

	result := make([]interface{}, len(objects))
	for i := range objects {
		v := make(map[string]interface{})
		getIdInput := &compute.GetInstanceIdInput{
			// Oracle's api returns an unordered list so we'll find out instances through the config file name
			Name: d.Get(fmt.Sprintf("instance.%d.name", i)).(string),
		}
		instance, err := instanceClient.GetInstanceFromName(getIdInput)
		if err != nil {
			return nil, err
		}

		v["name"] = instance.Name
		v["persistent"] = objects[i].Persistent
		v["shape"] = instance.Shape
		v["id"] = instance.ID

		instanceAttributes, err := flattenInstanceAttributes(instance.Attributes)
		if err != nil {
			return nil, err
		}
		v["instance_attributes"] = instanceAttributes

		sort.Ints(instance.BootOrder)
		v["boot_order"] = instance.BootOrder

		split_hostname := strings.Split(instance.Hostname, ".")
		if len(split_hostname) == 0 {
			return nil, fmt.Errorf("Unable to parse hostname: %s", instance.Hostname)
		}
		v["hostname"] = split_hostname[0]
		v["fqdn"] = instance.Hostname

		v["image_list"] = instance.ImageList
		v["label"] = instance.Label

		networkInterfaces, err := flattenNetworkInterfaces(instance.Networking)
		if err != nil {
			return nil, err
		}
		if len(networkInterfaces) > 0 {
			v["networking_info"] = networkInterfaces
		}

		sort.Strings(instance.SSHKeys)
		v["ssh_keys"] = instance.SSHKeys

		v["reverse_dns"] = instance.ReverseDNS

		v["storage"] = flattenInstanceStorageAttachments(instance.Storage)

		sort.Strings(instance.Tags)
		v["tags"] = instance.Tags

		v["availability_domain"] = instance.AvailabilityDomain
		v["domain"] = instance.Domain
		v["entry"] = instance.Entry
		v["fingerprint"] = instance.Fingerprint
		v["image_format"] = instance.ImageFormat
		v["ip_address"] = instance.IPAddress

		sort.Strings(instance.PlacementRequirements)
		v["placement_requirements"] = instance.PlacementRequirements

		v["platform"] = instance.Platform
		v["priority"] = instance.Priority
		v["quota_reservation"] = instance.QuotaReservation

		sort.Strings(instance.Relationships)
		v["relationships"] = instance.Relationships

		sort.Strings(instance.Resolvers)
		v["resolvers"] = instance.Resolvers

		v["site"] = instance.Site
		v["start_time"] = instance.StartTime
		v["state"] = instance.State

		v["vcable"] = instance.VCableID
		v["virtio"] = instance.Virtio
		v["vnc_address"] = instance.VNC

		result[i] = v
	}

	return result, nil
}

// Flattens attributes from the returned instance object, and sets the computed attributes string
// as JSON
func flattenInstanceAttributes(attributes map[string]interface{}) (string, error) {
	// Shouldn't ever get nil attributes on an instance, but protect against the case either way
	if attributes == nil {
		return "", nil
	}

	b, err := json.Marshal(attributes)
	if err != nil {
		return "", fmt.Errorf("Error marshalling returned attributes: %s", err)
	}
	return string(b), nil
}

// Flattens network interfaces
func flattenNetworkInterfaces(ifaces map[string]compute.NetworkingInfo) ([]map[string]interface{}, error) {
	result := make([]map[string]interface{}, 0)

	// Nil check for import case
	if ifaces == nil {
		return result, nil
	}

	for index, iface := range ifaces {
		res := make(map[string]interface{})
		// The index returned from the SDK holds the full device_index from the instance.
		// For users convenience, we simply allow them to specify the integer equivalent of the device_index
		// so a user could implement several network interfaces via `count`.
		// Convert the full device_index `ethN` to `N` as an integer.
		index := strings.TrimPrefix(index, "eth")
		indexInt, err := strconv.Atoi(index)
		if err != nil {
			return nil, err
		}
		res["index"] = indexInt
		res["is_default_gateway"] = iface.IsDefaultGateway

		// Set the proper attributes for this specific network interface
		if iface.DNS != nil {
			res["dns"] = iface.DNS
		}
		if iface.IPAddress != "" {
			res["ip_address"] = iface.IPAddress
		}
		if iface.IPNetwork != "" {
			res["ip_network"] = iface.IPNetwork
		}
		if iface.MACAddress != "" {
			res["mac_address"] = iface.MACAddress
		}
		if iface.Model != "" {
			// Model can only be set on Shared networks
			res["shared_network"] = true
		}
		if iface.NameServers != nil {
			res["name_servers"] = iface.NameServers
		}
		if iface.Nat != nil {
			res["nat"] = iface.Nat
		}
		if iface.SearchDomains != nil {
			res["search_domains"] = iface.SearchDomains
		}
		if iface.SecLists != nil {
			res["sec_lists"] = iface.SecLists
		}
		if iface.Vnic != "" {
			res["vnic"] = iface.Vnic
			// VNIC can only be set on an IP Network
			res["shared_network"] = false
		}
		if iface.VnicSets != nil {
			res["vnic_sets"] = iface.VnicSets
		}

		result = append(result, res)
	}

	return result, nil
}

// Flattens the returned slice of storage attachments to a map
func flattenInstanceStorageAttachments(attachments []compute.StorageAttachment) []map[string]interface{} {
	result := make([]map[string]interface{}, 0)

	if attachments == nil || len(attachments) == 0 {
		return nil
	}

	for _, attachment := range attachments {
		res := make(map[string]interface{})
		res["index"] = attachment.Index
		res["volume"] = attachment.StorageVolumeName
		res["name"] = attachment.Name
		result = append(result, res)
	}
	return result
}
