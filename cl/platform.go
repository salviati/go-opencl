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
	id      C.cl_platform_id
	Devices []Device
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
