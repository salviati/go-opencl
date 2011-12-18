package opencl

/*
#cgo LDFLAGS: -lOpenCL

#include "CL/cl.h"

*/
import "C"

import (
	"runtime"
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

func (c *Context) NewCommandQueue(device Device, param CommandQueueParameter) (*CommandQueue, error) {
	var c_queue C.cl_command_queue
	var err C.cl_int
	if c_queue = C.clCreateCommandQueue(c.id, device.id, C.cl_command_queue_properties(param), &err); err != C.CL_SUCCESS {
		return nil, Cl_error(err)
	}
	queue := &CommandQueue{id: c_queue}
	runtime.SetFinalizer(queue, (*CommandQueue).release)

	return queue, nil
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
