package opencl

import (
	"testing"
)

var source = `
__kernel void hello(__global float *input, __global float *output)
{
   size_t id = get_global_id(0);
   output[id] = input[id] * input[id];
}`

func Test_Program(t *testing.T) {
	platforms, err := Platforms()
	if err != nil {
		t.Fatal("Error getting platform:", err)
	} else if len(platforms) == 0 {
		t.Fatal("No platforms found:", err)
	}

	var context *Context
	if context, err = NewContextOfType(
		map[ContextParameter]interface{}{
			CL_CONTEXT_PLATFORM: platforms[0]},
		CL_DEVICE_TYPE_ALL); err != nil {
		t.Fatal("Error creating context of type:", err)
	}

	if _, err = context.NewProgramFromSource(source); err != nil {
		t.Fatal("Error creating program:", err)
	}
}
