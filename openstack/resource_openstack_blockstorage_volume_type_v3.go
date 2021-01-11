package openstack

import (
	"fmt"
	"log"

	"github.com/gophercloud/gophercloud/openstack/blockstorage/v3/volumetypes"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceBlockStorageVolumeTypeV3() *schema.Resource {
	return &schema.Resource{
		Create: resourceBlockStorageVolumeTypeV3Create,
		Read:   resourceBlockStorageVolumeTypeV3Read,
		Update: resourceBlockStorageVolumeTypeV3Update,
		Delete: resourceBlockStorageVolumeTypeV3Delete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"qos_specs_id": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},

			"extra_specs": {
				Type:     schema.TypeMap,
				Optional: true,
				Computed: true,
			},

			"os_volume_type_access": {
				Type:     schema.TypeBool,
				Computed: true,
				Optional: true,
			},

			"is_public": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func resourceBlockStorageVolumeTypeV3Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	blockStorageClient, err := config.BlockStorageV3Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenStack block storage client: %s", err)
	}
	isPublic := d.Get("is_public").(bool)
	extraspecs := d.Get("extra_specs").(map[string]interface{})
	createOpts := &volumetypes.CreateOpts{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		IsPublic:    &isPublic,
		ExtraSpecs:  expandToMapStringString(extraspecs),
	}

	log.Printf("[DEBUG] openstack_blockstorage_volume_type_v3 create options: %#v", createOpts)

	v, err := volumetypes.Create(blockStorageClient, createOpts).Extract()
	if err != nil {
		return fmt.Errorf("Error creating openstack_blockstorage_volume_type_v3: %s", err)
	}
	d.SetId(v.ID)

	return resourceBlockStorageVolumeTypeV3Read(d, meta)
}

func resourceBlockStorageVolumeTypeV3Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	blockStorageClient, err := config.BlockStorageV3Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenStack block storage client: %s", err)
	}

	v, err := volumetypes.Get(blockStorageClient, d.Id()).Extract()
	if err != nil {
		return CheckDeleted(d, err, "Error retrieving openstack_blockstorage_volume_type_v3")
	}

	log.Printf("[DEBUG] Retrieved openstack_blockstorage_volume_type_v3 %s: %#v", d.Id(), v)

	d.Set("name", v.Name)
	d.Set("description", v.Description)
	d.Set("name", v.Name)
	d.Set("os_volume_type_access", v.PublicAccess)
	d.Set("qos_specs_id", v.QosSpecID)
	d.Set("is_public", v.IsPublic)
	d.Set("extra_specs", v.ExtraSpecs)
	d.Set("region", GetRegion(d, config))

	return nil
}

func resourceBlockStorageVolumeTypeV3Update(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	blockStorageClient, err := config.BlockStorageV3Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenStack block storage client: %s", err)
	}

	name := d.Get("name").(string)
	description := d.Get("description").(string)
	isPublic := d.Get("is_public").(bool)
	updateOpts := volumetypes.UpdateOpts{
		Name:        &name,
		Description: &description,
		IsPublic:    &isPublic,
	}

	_, err = volumetypes.Update(blockStorageClient, d.Id(), updateOpts).Extract()
	if err != nil {
		return fmt.Errorf("Error updating openstack_blockstorage_volume_type_v3 %s: %s", d.Id(), err)
	}

	return resourceBlockStorageVolumeTypeV3Read(d, meta)
}

func resourceBlockStorageVolumeTypeV3Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	blockStorageClient, err := config.BlockStorageV3Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenStack block storage client: %s", err)
	}

	_, err = volumetypes.Get(blockStorageClient, d.Id()).Extract()
	if err != nil {
		return CheckDeleted(d, err, "Error retrieving openstack_blockstorage_volume_type_v3")
	}

	err = volumetypes.Delete(blockStorageClient, d.Id()).ExtractErr()
	if err != nil {
		return CheckDeleted(d, err, "Error deleting openstack_compute_volume_type_v3")
	}
	return nil
}
