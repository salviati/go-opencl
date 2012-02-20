package main

// #cgo pkg-config: sdl
// #include <SDL/SDL.h>
import "C"

import (
	"errors"
	"github.com/banthar/Go-SDL/sdl"
	"unsafe"
)

func SDL_CreateRGBSurfaceFrom(pixels *byte, width int, height int, depth int, pitch int, Rmask uint32, Gmask uint32, Bmask uint32, Amask uint32) (*sdl.Surface, error) {
	p := C.SDL_CreateRGBSurfaceFrom(unsafe.Pointer(pixels), C.int(width), C.int(height), C.int(depth), C.int(pitch),
		C.Uint32(Rmask), C.Uint32(Gmask), C.Uint32(Bmask), C.Uint32(Amask))
	if p == nil {
		return nil, errors.New("SDL_CreateRGBSurfaceFrom:" + sdl.GetError())
	}
	s := (*sdl.Surface)(unsafe.Pointer(p))
	s.Pixels = pixels // prevent gc
	return s, nil
}
