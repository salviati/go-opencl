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

type PlatformProperty C.cl_platform_info

const (
	PLATFORM_PROFILE    PlatformProperty = C.CL_PLATFORM_PROFILE
	PLATFORM_VERSION    PlatformProperty = C.CL_PLATFORM_VERSION
	PLATFORM_NAME       PlatformProperty = C.CL_PLATFORM_NAME
	PLATFORM_VENDOR     PlatformProperty = C.CL_PLATFORM_VENDOR
	PLATFORM_EXTENSIONS PlatformProperty = C.CL_PLATFORM_EXTENSIONS
)

func PlatformProperties() []PlatformProperty {
	return []PlatformProperty{
		PLATFORM_PROFILE,
		PLATFORM_VERSION,
		PLATFORM_NAME,
		PLATFORM_VENDOR,
		PLATFORM_EXTENSIONS}
}

type Platform struct {
	id C.cl_platform_id
}

func Platforms() ([]Platform, error) {
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

	platforms := make([]Platform, num_platforms)
	for i, v := range c_platforms {
		platforms[i].id = v
	}
	return platforms, nil
}

func (platform *Platform) Properties() (map[PlatformProperty]string, error) {
	infos := make(map[PlatformProperty]string)
	for _, param := range PlatformProperties() {
		var err error
		if infos[param], err = platform.Property(param); err != nil {
			return nil, err
		}
	}
	return infos, nil
}

func (platform *Platform) Property(prop PlatformProperty) (string, error) {
	const bufsize = 1024
	var buf [bufsize]C.char
	var length C.size_t
	if ret := C.clGetPlatformInfo(platform.id, C.cl_platform_info(prop),
		bufsize, unsafe.Pointer(&buf[0]), &length); ret != C.CL_SUCCESS {
		return "", Cl_error(ret)
	}
	return C.GoStringN(&buf[0], C.int(length)), nil
}

func (platform *Platform) Devices(t DeviceType) ([]Device, error) {
	var num_devices C.cl_uint
	if ret := C.clGetDeviceIDs(platform.id, C.cl_device_type(t),
		0, (*C.cl_device_id)(nil), &num_devices); ret != C.CL_SUCCESS {
		return nil, Cl_error(ret)
	} else if num_devices == 0 {
		return nil, Cl_error(C.CL_INVALID_VALUE)
	}

	c_devices := make([]C.cl_device_id, num_devices)
	if ret := C.clGetDeviceIDs(platform.id, C.cl_device_type(t),
		num_devices, &c_devices[0], &num_devices); ret != C.CL_SUCCESS {
		return nil, Cl_error(ret)
	} else if num_devices == 0 {
		return nil, Cl_error(C.CL_INVALID_VALUE)
	}

	devices := make([]Device, num_devices)
	for i, v := range c_devices {
		devices[i].id = v
	}
	return devices, nil
}
