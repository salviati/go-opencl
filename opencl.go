package opencl

/*
#cgo LDFLAGS: -lOpenCL

#include "CL/cl.h"

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
	C.CL_SUCCESS:             "Success",
	C.CL_DEVICE_NOT_FOUND:    "Device not Found",
	C.CL_INVALID_DEVICE:      "Invalid Device",
	C.CL_INVALID_DEVICE_TYPE: "Invalid Device Type",
	C.CL_INVALID_PLATFORM:    "Invalid Platform",
	C.CL_INVALID_VALUE:       "Invalid Value",
	C.CL_OUT_OF_HOST_MEMORY:  "Out of Host Memory",
	C.CL_OUT_OF_RESOURCES:    "Out of Resources",
}
