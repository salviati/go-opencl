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
	C.CL_SUCCESS:          "Success",
	C.CL_INVALID_PLATFORM: "Invalid Platform",
	C.CL_INVALID_VALUE:    "Invalid Value",
}
