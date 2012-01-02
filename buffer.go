package opencl

/*
#cgo CFLAGS: -I .
#cgo LDFLAGS: -lOpenCL

#include "CL/opencl.h"

*/
import "C"

import ()

type BufferFlags C.cl_mem_flags

const (
	CL_MEM_READ_WRITE     BufferFlags = C.CL_MEM_READ_WRITE
	CL_MEM_WRITE_ONLY     BufferFlags = C.CL_MEM_WRITE_ONLY
	CL_MEM_READ_ONLY      BufferFlags = C.CL_MEM_READ_ONLY
	cl_MEM_ALLOC_HOST_PTR BufferFlags = C.CL_MEM_ALLOC_HOST_PTR
)

type Buffer struct {
	id C.cl_mem
}

func (b *Buffer) release() error {
	if b.id != nil {
		if err := C.clReleaseMemObject(b.id); err != C.CL_SUCCESS {
			return Cl_error(err)
		}
		b.id = nil
	}
	return nil
}
