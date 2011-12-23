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

func (c *Context) NewBuffer(flags BufferFlags, size uint32) (*Buffer, error) {
	var c_buffer C.cl_mem
	var err C.cl_int

	if c_buffer = C.clCreateBuffer(c.id, C.cl_mem_flags(flags|cl_MEM_ALLOC_HOST_PTR), C.size_t(size), nil, &err); err != C.CL_SUCCESS {
		return nil, Cl_error(err)
	}

	buffer := &Buffer{id: c_buffer}
	runtime.SetFinalizer(buffer, (*Buffer).release)

	return buffer, nil
}

func (cq *CommandQueue) EnqueueReadBuffer(buf *Buffer, offset uint32, size uint32) ([]byte, error) {
	c_bytes := make([]byte, size)
	if ret := C.clEnqueueReadBuffer(cq.id, buf.id, C.CL_TRUE, C.size_t(offset), C.size_t(size), unsafe.Pointer(&c_bytes[0]), 0, nil, nil); ret != C.CL_SUCCESS {
		return nil, Cl_error(ret)
	}

	// Copy the buffer in case the garbage collector moves the slice in memory
	bytes := make([]byte, size)
	for i, v := range c_bytes {
		bytes[i] = v
	}
	return bytes, nil
}

func (cq *CommandQueue) EnqueueWriteBuffer(buf *Buffer, data []byte, offset uint32) error {

	if ret := C.clEnqueueWriteBuffer(cq.id, buf.id, C.CL_TRUE, C.size_t(offset), C.size_t(len(data)), unsafe.Pointer(&data[0]), 0, nil, nil); ret != C.CL_SUCCESS {
		return Cl_error(ret)
	}
	return nil
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
