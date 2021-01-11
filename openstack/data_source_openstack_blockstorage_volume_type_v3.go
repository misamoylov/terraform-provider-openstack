package openstack

import (
	"fmt"
	"log"

	"github.com/gophercloud/gophercloud/openstack/blockstorage/v3/volumetypes"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func dataSourceBlockStorageVolumeTypeV3() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceBlockStorageVolumeTypeV3Read,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"extra_specs": {
				Type:     schema.TypeMap,
				Optional: true,
				Computed: true,
			},

			"qos_specs_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"os_volume_type_access": {
				Type:     schema.TypeBool,
				Computed: true,
			},

			"is_public": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func dataSourceBlockStorageVolumeTypeV3Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	client, err := config.BlockStorageV3Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating OpenStack block storage client: %s", err)
	}

	allPages, err := volumetypes.List(client, nil).AllPages()
	if err != nil {
		return fmt.Errorf("Unable to query openstack_blockstorage_volume_type_v3: %s", err)
	}

	allVolumetypes, err := volumetypes.ExtractVolumeTypes(allPages)
	fmt.Printf("Volumetypes: %+v\n", allVolumetypes)
	if err != nil {
		return fmt.Errorf("Unable to retrieve openstack_blockstorage_volume_type_v3: %s", err)
	}

	for v := range allVolumetypes {
		if allVolumetypes[v].Name == d.Get("name").(string) {
			return dataSourceBlockStorageVolumeTypeV3Attributes(d, allVolumetypes[v])
		}
	}
	return nil
}

func dataSourceBlockStorageVolumeTypeV3Attributes(d *schema.ResourceData, volumetype volumetypes.VolumeType) error {
	d.SetId(volumetype.ID)
	d.Set("name", volumetype.Name)
	d.Set("description", volumetype.Description)
	d.Set("os_volume_type_access", volumetype.PublicAccess)
	d.Set("is_public", volumetype.IsPublic)
	d.Set("qos_specs_id", volumetype.QosSpecID)

	if err := d.Set("extra_specs", volumetype.ExtraSpecs); err != nil {
		log.Printf("[DEBUG] Unable to set extra_specs for openstack_blockstorage_volume_type_v3 %s: %s", volumetype.ID, err)
	}

	return nil
}
