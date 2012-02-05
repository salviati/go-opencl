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
#cgo CFLAGS: -I .
#cgo LDFLAGS: -lOpenCL

#include "CL/opencl.h"

*/
import "C"

import ()

type BufferFlags C.cl_mem_flags

const (
	MEM_READ_WRITE BufferFlags = C.CL_MEM_READ_WRITE
	MEM_WRITE_ONLY BufferFlags = C.CL_MEM_WRITE_ONLY
	MEM_READ_ONLY  BufferFlags = C.CL_MEM_READ_ONLY
	//MEM_USE_HOST_PTR   BufferFlags = C.CL_MEM_USE_HOST_PTR
	MEM_ALLOC_HOST_PTR BufferFlags = C.CL_MEM_ALLOC_HOST_PTR
	//MEM_COPY_HOST_PTR  BufferFlags = C.CL_MEM_COPY_HOST_PTR
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
