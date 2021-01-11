package openstack

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	"github.com/gophercloud/gophercloud/openstack/blockstorage/v3/volumetypes"
)

func TestAccBlockStorageV3VolumeTypeDataSource_basic(t *testing.T) {
	resourceName := "data.openstack_blockstorage_volume_type_v3.volume_type_1"
	volumetypeName := acctest.RandomWithPrefix("tf-acc-volume-type")

	var volumetypeID string
	if os.Getenv("TF_ACC") != "" {
		var err error
		volumetypeID, err = testAccBlockStorageV3CreateVolumeType(volumetypeName)
		if err != nil {
			t.Fatal(err)
		}
		defer testAccBlockStorageV3DeleteVolumeType(t, volumetypeID)
	}

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccBlockStorageV3VolumeTypeDataSourceBasic(volumetypeName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckBlockStorageV3VolumeTypeDataSourceID(resourceName),
					resource.TestCheckResourceAttr(resourceName, "name", volumetypeName),
					resource.TestCheckResourceAttr(resourceName, "description", "1"),
				),
			},
		},
	})
}

func testAccBlockStorageV3CreateVolumeType(volumetypeName string) (string, error) {
	config, err := testAccAuthFromEnv()
	if err != nil {
		return "", err
	}

	bsClient, err := config.BlockStorageV3Client(osRegionName)
	if err != nil {
		return "", err
	}

	var isPublic = true

	volCreateOpts := volumetypes.CreateOpts{
		Description: "1",
		IsPublic:    &isPublic,
		Name:        volumetypeName,
		ExtraSpecs:  map[string]string{"volume_backend_name": "lvmdriver-1"},
	}

	volumetype, err := volumetypes.Create(bsClient, volCreateOpts).Extract()
	if err != nil {
		return "", err
	}

	return volumetype.ID, nil
}

func testAccBlockStorageV3DeleteVolumeType(t *testing.T, volumetypeID string) {
	config, err := testAccAuthFromEnv()
	if err != nil {
		t.Fatal(err)
	}

	bsClient, err := config.BlockStorageV3Client(osRegionName)
	if err != nil {
		t.Fatal(err)
	}

	err = volumetypes.Delete(bsClient, volumetypeID).ExtractErr()
	if err != nil {
		t.Fatal(err)
	}
}

func testAccCheckBlockStorageV3VolumeTypeDataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find volumetype data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Volume Type data source ID not set")
		}

		return nil
	}
}

func testAccBlockStorageV3VolumeTypeDataSourceBasic(name string) string {
	return fmt.Sprintf(`
    data "openstack_blockstorage_volume_type_v3" "volume_type_1" {
      name = "%s"
    }
  `, name)
}
