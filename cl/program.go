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
