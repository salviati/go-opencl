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

type Program struct {
	id C.cl_program
}

func (c *Context) NewProgramFromSource(prog string) (*Program, error) {
	var c_program C.cl_program
	var err C.cl_int

	cs := C.CString(prog)
	defer C.free(unsafe.Pointer(cs))

	if c_program = C.clCreateProgramWithSource(c.id, 1, &cs, (*C.size_t)(nil), &err); err != C.CL_SUCCESS {
		return nil, Cl_error(err)
	} else if err = C.clBuildProgram(c_program, 0, nil, nil, nil, nil); err != C.CL_SUCCESS {
		C.clReleaseProgram(c_program)
		return nil, Cl_error(err)
	}

	program := &Program{id: c_program}
	runtime.SetFinalizer(program, (*Program).release)

	return program, nil
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
