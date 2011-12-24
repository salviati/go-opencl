package opencl

/*
#cgo CFLAGS: -I CL
#cgo LDFLAGS: -lOpenCL

#include "CL/opencl.h"

*/
import "C"

import (
	"fmt"
)

type Cl_error C.cl_int

func (e Cl_error) Error() string {
	mesg := errMesg[e]
	if mesg == "" {
		return fmt.Sprintf("error %d", int(e))
	}
	return mesg
}

var errMesg = map[Cl_error]string{
	C.CL_SUCCESS: "Success",

	C.CL_BUILD_PROGRAM_FAILURE:  "Build Program Failure",
	C.CL_COMPILER_NOT_AVAILABLE: "Compiler not Available",
	C.CL_DEVICE_NOT_AVAILABLE:   "Device not Available",
	C.CL_DEVICE_NOT_FOUND:       "Device not Found",
	C.CL_INVALID_BINARY:         "Invalid Binary",
	C.CL_INVALID_BUFFER_SIZE:    "Invalid Buffer Size",
	C.CL_INVALID_BUILD_OPTIONS:  "Invalid Build Options",
	C.CL_INVALID_COMMAND_QUEUE:  "Invalid Command Queue",
	C.CL_INVALID_CONTEXT:        "Invalid Context",
	//C.CL_INVALID_D3D10_DEVICE_KHR:            "Invalid D3D10 Device (KHR)",
	C.CL_INVALID_DEVICE:          "Invalid Device",
	C.CL_INVALID_DEVICE_TYPE:     "Invalid Device Type",
	C.CL_INVALID_EVENT_WAIT_LIST: "Invalid Event Wait List",
	//C.CL_INVALID_GL_SHAREGROUP_REFERENCE_KHR: "Invalid GL Sharegroup Reference (KHR)",
	C.CL_INVALID_GLOBAL_OFFSET:         "Invalid Global Offset",
	C.CL_INVALID_GLOBAL_WORK_SIZE:      "Invalid Global Work Size",
	C.CL_INVALID_HOST_PTR:              "Invalid Host Pointer",
	C.CL_INVALID_IMAGE_SIZE:            "Invalid Image Size",
	C.CL_INVALID_KERNEL:                "Invalid Kernel",
	C.CL_INVALID_KERNEL_ARGS:           "Invalid Kernel Arguments",
	C.CL_INVALID_KERNEL_DEFINITION:     "Invalid Kernel Definition",
	C.CL_INVALID_KERNEL_NAME:           "Invalid Kernel Name",
	C.CL_INVALID_MEM_OBJECT:            "Invalid Memory Object",
	C.CL_INVALID_OPERATION:             "Invalid Operation",
	C.CL_INVALID_PLATFORM:              "Invalid Platform",
	C.CL_INVALID_PROPERTY:              "Invalid Property",
	C.CL_INVALID_PROGRAM:               "Invalid Program",
	C.CL_INVALID_PROGRAM_EXECUTABLE:    "Invalid Program Executable",
	C.CL_INVALID_VALUE:                 "Invalid Value",
	C.CL_INVALID_WORK_DIMENSION:        "Invalid Work Dimension",
	C.CL_INVALID_WORK_GROUP_SIZE:       "Invalid Work Group Size",
	C.CL_INVALID_WORK_ITEM_SIZE:        "Invalid Work Item Size",
	C.CL_MEM_OBJECT_ALLOCATION_FAILURE: "Memory Object Allocation Failure",
	C.CL_MISALIGNED_SUB_BUFFER_OFFSET:  "Misaligned Sub-Buffer Offset",
	C.CL_OUT_OF_HOST_MEMORY:            "Out of Host Memory",
	C.CL_OUT_OF_RESOURCES:              "Out of Resources",
}
