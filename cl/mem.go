package cl

/*
#cgo CFLAGS: -I CL
#cgo LDFLAGS: -lOpenCL

#include "CL/opencl.h"
*/
import "C"

type MemFlags C.cl_mem_flags

const (
	MEM_READ_WRITE     MemFlags = C.CL_MEM_READ_WRITE
	MEM_WRITE_ONLY     MemFlags = C.CL_MEM_WRITE_ONLY
	MEM_READ_ONLY      MemFlags = C.CL_MEM_READ_ONLY
	MEM_USE_HOST_PTR   MemFlags = C.CL_MEM_USE_HOST_PTR
	MEM_ALLOC_HOST_PTR MemFlags = C.CL_MEM_ALLOC_HOST_PTR
	MEM_COPY_HOST_PTR  MemFlags = C.CL_MEM_COPY_HOST_PTR
)

const (
	MEM_OBJECT_BUFFER  = C.CL_MEM_OBJECT_BUFFER
	MEM_OBJECT_IMAGE2D = C.CL_MEM_OBJECT_IMAGE2D
	MEM_OBJECT_IMAGE3D = C.CL_MEM_OBJECT_IMAGE3D
)

func releaseMemObject(p C.cl_mem) error {
	if p != nil {
		if err := C.clReleaseMemObject(p); err != C.CL_SUCCESS {
			return Cl_error(err)
		}
	}
	return nil
}
