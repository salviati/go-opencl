package opencl

import (
	"testing"
)

func Test_clGetDeviceInfo(t *testing.T) {
	devices, err := Devices(CL_DEVICE_TYPE_ALL)
	if err != nil {
		t.Fatal(err)
	}

	for _, device := range devices {
		if properties, err := device.Properties(); err != nil {
			t.Fatal("error mapping device properties:", err)
		} else {
			for _, property := range DeviceProperties() {
				if value, _ := device.Property(property); value != properties[property] {
					t.Fatal("Device property disagrees -", property, "-", value, "!=", properties[property])
				}
			}
		}
	}
}
