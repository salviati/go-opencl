package opencl

import (
	"testing"
)

func Test_clGetPlatformInfo(t *testing.T) {
	platforms, err := Platforms()
	if err != nil {
		t.Fatal("Error getting platforms:", err)
	}

	for _, platform := range platforms {
		if properties, err := platform.Properties(); err != nil {
			t.Fatal("Error retrieving platform properties:", err)
		} else {
			for _, property := range PlatformProperties() {
				if value, err := platform.Property(property); err != nil {
					t.Fatal("Error retrieving platform property:", err)
				} else if value != properties[property] {
					t.Fatal("Platform info disagrees -", property, "-", value, "!=", properties[property])
				}
			}
		}
	}
}
