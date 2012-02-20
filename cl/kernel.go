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
	"unsafe"
)

type Kernel struct {
	id C.cl_kernel
}

func (k *Kernel) SetArg(index uint, arg interface{}) error {
	var ret C.cl_int

	switch t := arg.(type) {
	case *Buffer: // can't use *Image here, compiler bug?
		ret = C.clSetKernelArg(k.id, C.cl_uint(index), C.size_t(unsafe.Sizeof(t.id)), unsafe.Pointer(&t.id))
	case *Image:
		ret = C.clSetKernelArg(k.id, C.cl_uint(index), C.size_t(unsafe.Sizeof(t.id)), unsafe.Pointer(&t.id))
	case float32:
		ret = C.clSetKernelArg(k.id, C.cl_uint(index), C.size_t(4), unsafe.Pointer(&t))
	case []float32:
		ret = C.clSetKernelArg(k.id, C.cl_uint(index), C.size_t(len(t)*4), unsafe.Pointer(&t[0]))
	case *float32:
		var p *C.float
		ret = C.clSetKernelArg(k.id, C.cl_uint(index), C.size_t(unsafe.Sizeof(p)), unsafe.Pointer(&t))

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
