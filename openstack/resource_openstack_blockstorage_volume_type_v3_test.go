package openstack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	"github.com/gophercloud/gophercloud/openstack/blockstorage/v3/volumetypes"
)

func TestAccBlockStorageV3VolumeType_basic(t *testing.T) {
	var volumetype volumetypes.VolumeType
	var volumetypeName = acctest.RandomWithPrefix("tf-acc-volume-type")

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckBlockStorageV3VolumeTypeDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccBlockStorageV3VolumeTypeBasic(volumetypeName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBlockStorageV3VolumeTypeExists("openstack_blockstorage_volume_type_v3.volume_type_1", &volumetype),
					testAccCheckBlockStorageV3VolumeTypeExtraSpecs(&volumetype, "foo", "bar"),
					resource.TestCheckResourceAttr(
						"openstack_blockstorage_volume_type_v3.volume_type_1", "name", volumetypeName),
					resource.TestCheckResourceAttr(
						"openstack_blockstorage_volume_type_v3.volume_type_1", "description", "volumetype description"),
				),
			},
			{
				Config: testAccBlockStorageV3VolumeTypeUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBlockStorageV3VolumeTypeExists("openstack_blockstorage_volume_type_v3.volume_type_1", &volumetype),
					testAccCheckBlockStorageV3VolumeTypeExtraSpecs(&volumetype, "foo", "bar2"),
					resource.TestCheckResourceAttr(
						"openstack_blockstorage_volume_type_v3.volume_type_1", "name", "volume_type_1_updated"),
					resource.TestCheckResourceAttr(
						"openstack_blockstorage_volume_type_v3.volume_type_1", "description", "update test volume type"),
				),
			},
		},
	})
}

func testAccCheckBlockStorageV3VolumeTypeDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	blockStorageClient, err := config.BlockStorageV3Client(osRegionName)
	if err != nil {
		return fmt.Errorf("Error creating OpenStack block storage client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "openstack_blockstorage_volume_type_v3" {
			continue
		}

		_, err := volumetypes.Get(blockStorageClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("Volumetype still exists")
		}
	}

	return nil
}

func testAccCheckBlockStorageV3VolumeTypeExists(n string, volumetype *volumetypes.VolumeType) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		blockStorageClient, err := config.BlockStorageV3Client(osRegionName)
		if err != nil {
			return fmt.Errorf("Error creating OpenStack block storage client: %s", err)
		}

		found, err := volumetypes.Get(blockStorageClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("Volumetype not found")
		}

		*volumetype = *found

		return nil
	}
}

func testAccCheckBlockStorageV3VolumeTypeExtraSpecs(
	volumetypes *volumetypes.VolumeType, k string, v string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if volumetypes.ExtraSpecs == nil {
			return fmt.Errorf("No Extra Specs")
		}

		for key, value := range volumetypes.ExtraSpecs {
			if k != key {
				continue
			}

			if v == value {
				return nil
			}

			return fmt.Errorf("Bad value for %s: %s", k, value)
		}

		return fmt.Errorf("Extra Specs not found: %s", k)
	}
}

func testAccBlockStorageV3VolumeTypeBasic(volumeTypeName string) string {
	return fmt.Sprintf(`
		resource "openstack_blockstorage_volume_type_v3" "volume_type_1" {
		  name = "%s"
			is_public = true
		  description = "volumetype description"
		  extra_specs = {
		    foo = "bar"
		  }
		}
		`, volumeTypeName)
}

const testAccBlockStorageV3VolumeTypeUpdate = `
resource "openstack_blockstorage_volume_type_v3" "volume_type_1" {
  name = "volume_type_1_updated"
  description = "update test volume type"
  extra_specs = {
    foo = "bar2"
  }
}
`
