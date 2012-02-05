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

cl_context_properties PlatformToContextParameter(cl_platform_id platform) { return (cl_context_properties)platform; }

*/
import "C"

import (
	"runtime"
	"unsafe"
)

type ContextParameter C.cl_context_properties

const (
	CONTEXT_PLATFORM ContextParameter = C.CL_CONTEXT_PLATFORM
)

type ContextProperty C.cl_context_info

const (
	CONTEXT_REFERENCE_COUNT ContextProperty = C.CL_CONTEXT_REFERENCE_COUNT
	CONTEXT_NUM_DEVICES     ContextProperty = C.CL_CONTEXT_NUM_DEVICES
	CONTEXT_DEVICES         ContextProperty = C.CL_CONTEXT_DEVICES
	CONTEXT_PROPERTIES      ContextProperty = C.CL_CONTEXT_PROPERTIES
	//CONTEXT_D3D10_PREFER_SHARED_RESOURCES_KHR ContextProperty = C.CL_CONTEXT_D3D10_PREFER_SHARED_RESOURCES_KHR
)

func ContextProperties() []ContextProperty {
	return []ContextProperty{
		CONTEXT_REFERENCE_COUNT,
		CONTEXT_NUM_DEVICES,
		CONTEXT_DEVICES,
		CONTEXT_PROPERTIES}
}

type Context struct {
	id C.cl_context
}

func createParameters(params map[ContextParameter]interface{}) ([]C.cl_context_properties, error) {
	c_params := make([]C.cl_context_properties, (len(params)<<1)+1)
	i := 0
	for param, value := range params {
		c_params[i] = C.cl_context_properties(param)

		switch param {
		case CONTEXT_PLATFORM:
			if v, ok := value.(Platform); ok {
				c_params[i+1] = C.PlatformToContextParameter(v.id)
			} else {
				return nil, Cl_error(C.CL_INVALID_VALUE)
			}

		default:
			return nil, Cl_error(C.CL_INVALID_VALUE)
		}
		i += 2
	}
	c_params[i] = 0

	return c_params, nil
}

func NewContextOfType(params map[ContextParameter]interface{}, t DeviceType) (*Context, error) {
	var c_params []C.cl_context_properties
	var propErr error
	if c_params, propErr = createParameters(params); propErr != nil {
		return nil, propErr
	}

	var c_context C.cl_context
	var err C.cl_int
	if c_context = C.clCreateContextFromType(&c_params[0], C.cl_device_type(t), nil, nil, &err); err != C.CL_SUCCESS {
		return nil, Cl_error(err)
	}
	context := &Context{id: c_context}
	runtime.SetFinalizer(context, (*Context).release)

	return context, nil
}

func NewContextOfDevices(params map[ContextParameter]interface{}, devices []Device) (*Context, error) {
	var c_params []C.cl_context_properties
	var propErr error
	if c_params, propErr = createParameters(params); propErr != nil {
		return nil, propErr
	}

	c_devices := make([]C.cl_device_id, len(devices))
	for i, device := range devices {
		c_devices[i] = device.id
	}

	var c_context C.cl_context
	var err C.cl_int
	if c_context = C.clCreateContext(&c_params[0], C.cl_uint(len(c_devices)), &c_devices[0], nil, nil, &err); err != C.CL_SUCCESS {
		return nil, Cl_error(err)
	}
	context := &Context{id: c_context}
	runtime.SetFinalizer(context, (*Context).release)

	return context, nil
}

func (context *Context) Properties() (map[ContextProperty]interface{}, error) {
	props := make(map[ContextProperty]interface{})
	for _, prop := range ContextProperties() {
		if data, err := context.Property(prop); err == nil {
			props[prop] = data
		} else if err != Cl_error(C.CL_INVALID_VALUE) {
			return nil, err
		}
	}
	if len(props) == 0 {
		return nil, Cl_error(C.CL_INVALID_VALUE)
	}
	return props, nil
}

func (context *Context) Property(prop ContextProperty) (interface{}, error) {
	var data interface{}
	var length C.size_t
	var ret C.cl_int

	switch prop {
	/*case CONTEXT_D3D10_PREFER_SHARED_RESOURCES_KHR:
	var val C.cl_bool
	ret = C.clGetContextInfo(context.id, C.cl_context_info(prop), C.size_t(unsafe.Sizeof(val)), unsafe.Pointer(&val), &length)
	data = val == C.CL_TRUE
	*/

	case CONTEXT_REFERENCE_COUNT,
		CONTEXT_NUM_DEVICES:
		var val C.cl_uint
		ret = C.clGetContextInfo(context.id, C.cl_context_info(prop), C.size_t(unsafe.Sizeof(val)), unsafe.Pointer(&val), &length)
		data = val

	case CONTEXT_DEVICES:
		if data, err := context.Property(CONTEXT_NUM_DEVICES); err != nil {
			return nil, err
		} else {
			num_devs := data.(C.cl_uint)
			c_devs := make([]C.cl_device_id, num_devs)
			if ret = C.clGetContextInfo(context.id, C.cl_context_info(prop), C.size_t(num_devs*C.cl_uint(unsafe.Sizeof(c_devs[0]))), unsafe.Pointer(&c_devs[0]), &length); ret != C.CL_SUCCESS {
				return nil, Cl_error(ret)
			}
			devs := make([]Device, length/C.size_t(unsafe.Sizeof(c_devs[0])))
			for i, val := range c_devs {
				devs[i].id = val
			}
			return devs, nil
		}

	default:
		return nil, Cl_error(C.CL_INVALID_VALUE)
	}

	if ret != C.CL_SUCCESS {
		return nil, Cl_error(ret)
	}
	return data, nil
}

func (c *Context) release() error {
	if c.id != nil {
		if err := C.clReleaseContext(c.id); err != C.CL_SUCCESS {
			return Cl_error(err)
		}
		c.id = nil
	}
	return nil
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

func (c *Context) NewBuffer(flags BufferFlags, size uint32) (*Buffer, error) {
	var c_buffer C.cl_mem
	var err C.cl_int

	if c_buffer = C.clCreateBuffer(c.id, C.cl_mem_flags(flags), C.size_t(size), nil, &err); err != C.CL_SUCCESS {
		return nil, Cl_error(err)
	}

	buffer := &Buffer{id: c_buffer}
	runtime.SetFinalizer(buffer, (*Buffer).release)

	return buffer, nil
}
