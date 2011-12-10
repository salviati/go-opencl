package opencl

import (
	"testing"
)

func Test_clGetPlatformInfo(t *testing.T) {
	platforms, err := ClGetPlatformIDs()
	if err != nil {
		t.Fatal(err)
	}

	for _, platform := range platforms {
		if params, err := GetPlatformInfo(platform); err != nil {
			t.Fatal(err)
		} else {
			for _, param := range PlatformInfoParams() {
				if info, err := ClGetPlatformInfo(platform, param); err != nil {
					t.Fatal(err)
				} else if info != params[param] {
					t.Fatal("Platform info disagrees -", param, "-", info, "!=", params[param])
				}
			}
		}
	}
}
