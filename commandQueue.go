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

type CommandQueueParameter C.cl_command_queue_properties

const (
	CL_QUEUE_OUT_OF_ORDER_EXEC_MODE_ENABLE CommandQueueParameter = C.CL_QUEUE_OUT_OF_ORDER_EXEC_MODE_ENABLE
	CL_QUEUE_PROFILING_ENABLE              CommandQueueParameter = C.CL_QUEUE_PROFILING_ENABLE
	CL_QUEUE_NIL                           CommandQueueParameter = 0
)

type CommandQueue struct {
	id C.cl_command_queue
}

func (q *CommandQueue) release() error {
	if q.id != nil {
		if err := C.clReleaseCommandQueue(q.id); err != C.CL_SUCCESS {
			return Cl_error(err)
		}
		q.id = nil
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
