package opencl

/*
#cgo LDFLAGS: -lOpenCL

#include "CL/cl.h"

*/
import "C"

import (
	"unsafe"
)

type Cl_platform_info C.cl_platform_info

const (
	CL_PLATFORM_PROFILE    Cl_platform_info = C.CL_PLATFORM_PROFILE
	CL_PLATFORM_VERSION    Cl_platform_info = C.CL_PLATFORM_VERSION
	CL_PLATFORM_NAME       Cl_platform_info = C.CL_PLATFORM_NAME
	CL_PLATFORM_VENDOR     Cl_platform_info = C.CL_PLATFORM_VENDOR
	CL_PLATFORM_EXTENSIONS Cl_platform_info = C.CL_PLATFORM_EXTENSIONS
)

func PlatformInfoParams() []Cl_platform_info {
	return []Cl_platform_info{
		CL_PLATFORM_PROFILE,
		CL_PLATFORM_VERSION,
		CL_PLATFORM_NAME,
		CL_PLATFORM_VENDOR,
		CL_PLATFORM_EXTENSIONS}
}

type Cl_platform_id struct {
	id C.cl_platform_id
}

func ClGetPlatformIDs() ([]Cl_platform_id, error) {
	var num_platforms C.cl_uint
	if ret := C.clGetPlatformIDs(0, (*C.cl_platform_id)(nil), &num_platforms); ret != C.CL_SUCCESS {
		return nil, Cl_error(ret)
	} else if num_platforms == 0 {
		return nil, Cl_error(C.CL_INVALID_VALUE)
	}

	c_platforms := make([]C.cl_platform_id, num_platforms)
	if ret := C.clGetPlatformIDs(num_platforms, &c_platforms[0], &num_platforms); ret != C.CL_SUCCESS {
		return nil, Cl_error(ret)
	} else if num_platforms == 0 {
		return nil, Cl_error(C.CL_INVALID_VALUE)
	}

	platforms := make([]Cl_platform_id, num_platforms)
	for i, v := range c_platforms {
		platforms[i].id = v
	}
	return platforms, nil
}

func GetPlatformInfo(platform Cl_platform_id) (map[Cl_platform_info]string, error) {
	params := make(map[Cl_platform_info]string)
	for _, param := range PlatformInfoParams() {
		var err error
		if params[param], err = ClGetPlatformInfo(platform, param); err != nil {
			return nil, err
		}
	}
	return params, nil
}

func ClGetPlatformInfo(platform Cl_platform_id, info Cl_platform_info) (string, error) {
	const bufsize = 1024
	var buf [bufsize]C.char
	var length C.size_t
	if ret := C.clGetPlatformInfo(platform.id, C.cl_platform_info(info),
		bufsize, unsafe.Pointer(&buf[0]), &length); ret != C.CL_SUCCESS {
		return "", Cl_error(ret)
	}
	return C.GoStringN(&buf[0], C.int(length)), nil
}
