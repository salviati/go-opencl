package opencl

import (
	"testing"
)

func Test_Context(t *testing.T) {
	platforms, err := Platforms()
	if err != nil {
		t.Fatal("Error getting platform:", err)
	} else if len(platforms) == 0 {
		t.Fatal("No platforms found:", err)
	}

	devices, err := Devices(CL_DEVICE_TYPE_ALL)
	if err != nil {
		t.Fatal(err)
	}

	var params = map[ContextParameter]interface{}{
		CL_CONTEXT_PLATFORM: platforms[0]}

	if context, err := NewContextOfDevices(params, devices); err != nil {
		t.Fatal("Error creating context of devices:", err)
	} else if _, err := context.Properties(); err != nil {
		t.Fatal("Error getting properties for context (1):", err)
	}

	if _, err := NewContextOfType(params, CL_DEVICE_TYPE_ALL); err != nil {
		t.Fatal("Error creating context of type:", err)
	}
}
