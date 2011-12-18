package opencl

/*
#cgo LDFLAGS: -lOpenCL

#include "CL/cl.h"

*/
import "C"

import (
	"runtime"
	"unsafe"
)

type Kernel struct {
	id C.cl_kernel
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

func (cq *CommandQueue) EnqueueKernel(k *Kernel, offset uint, gsize uint, lsize uint) error {

	c_offset := C.size_t(offset)
	c_gsize := C.size_t(gsize)
	c_lsize := C.size_t(lsize)
	if ret := C.clEnqueueNDRangeKernel(cq.id, k.id, 1, &c_offset, &c_gsize, &c_lsize, 0, nil, nil); ret != C.CL_SUCCESS {
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
