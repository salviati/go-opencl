package cl

/*
#cgo CFLAGS: -I CL
#cgo LDFLAGS: -lOpenCL

#include "CL/opencl.h"

*/
import "C"

import (
	"unsafe"
)

type Kernel struct {
	id C.cl_kernel
}

func (k *Kernel) SetArg(index uint, arg interface{}) error {
	var ret C.cl_int

	switch t := arg.(type) {
	case *Buffer:
		ret = C.clSetKernelArg(k.id, C.cl_uint(index), C.size_t(unsafe.Sizeof(t.id)), unsafe.Pointer(&t.id))

	default:
		return Cl_error(C.CL_INVALID_VALUE)
	}

	if ret != C.CL_SUCCESS {
		return Cl_error(ret)
	}
	return nil
}

func (k *Kernel) release() error {
	if k.id != nil {
		if err := C.clReleaseKernel(k.id); err != C.CL_SUCCESS {
			return Cl_error(err)
		}
		k.id = nil
	}
	return nil
}
