/*
 * Copyright © 2012 Paul Sbarra
 *
 * This file is part of go-opencl.
 *
 * go-opencl is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Lesser General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * go-opencl is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Lesser General Public License for more details.
 *
 * You should have received a copy of the GNU Lesser General Public License
 * along with go-opencl.  If not, see <http://www.gnu.org/licenses/>.
 */

package cl

/*
#cgo CFLAGS: -I CL
#cgo LDFLAGS: -lOpenCL

#include "CL/opencl.h"

*/
import "C"

import (
	"github.com/tones111/raw"
	"unsafe"
)

type CommandQueueParameter C.cl_command_queue_properties

const (
	QUEUE_OUT_OF_ORDER_EXEC_MODE_ENABLE CommandQueueParameter = C.CL_QUEUE_OUT_OF_ORDER_EXEC_MODE_ENABLE
	QUEUE_PROFILING_ENABLE              CommandQueueParameter = C.CL_QUEUE_PROFILING_ENABLE
	QUEUE_NIL                           CommandQueueParameter = 0
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

func (cq *CommandQueue) EnqueueReadBuffer(buf *Buffer, offset uint32, data interface{}) error {
	bytes := raw.ByteSlice(data)
	if ret := C.clEnqueueReadBuffer(cq.id, buf.id, C.CL_TRUE, C.size_t(offset), C.size_t(len(bytes)), raw.DataAddress(bytes), 0, nil, nil); ret != C.CL_SUCCESS {
		return Cl_error(ret)
	}
	raw.ByteCopy(data, bytes)
	return nil
}

func (cq *CommandQueue) EnqueueWriteBuffer(buf *Buffer, data interface{}, offset uint32) error {
	bytes := raw.ByteSlice(data)
	if len(bytes) == 0 {
		return Cl_error(C.CL_INVALID_VALUE)
	}

	if ret := C.clEnqueueWriteBuffer(cq.id, buf.id, C.CL_TRUE, C.size_t(offset), C.size_t(len(bytes)), unsafe.Pointer(&bytes[0]), 0, nil, nil); ret != C.CL_SUCCESS {
		return Cl_error(ret)
	}
	return nil
}
