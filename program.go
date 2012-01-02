package opencl

/*
#cgo CFLAGS: -I CL
#cgo LDFLAGS: -lOpenCL

#include "CL/opencl.h"

*/
import "C"

import (
	"runtime"
	"unsafe"
)

type Program struct {
	id C.cl_program
}

func (p *Program) NewKernelNamed(name string) (*Kernel, error) {
	var c_kernel C.cl_kernel
	var err C.cl_int

	cs := C.CString(name)
	defer C.free(unsafe.Pointer(cs))

	if c_kernel = C.clCreateKernel(p.id, cs, &err); err != C.CL_SUCCESS {
		return nil, Cl_error(err)
	}

	kernel := &Kernel{id: c_kernel}
	runtime.SetFinalizer(kernel, (*Kernel).release)

	return kernel, nil
}

func (p *Program) release() error {
	if p.id != nil {
		if err := C.clReleaseProgram(p.id); err != C.CL_SUCCESS {
			return Cl_error(err)
		}
		p.id = nil
	}
	return nil
}
