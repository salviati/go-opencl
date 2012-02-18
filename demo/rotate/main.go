/*
	  NOTE:

		https://github.com/banthar/Go-SDL doesn't implement CreateRGBSurfaceFrom, since pixels may be garbage collected in general.

		func CreateRGBSurfaceFrom(pixels *byte, width int, height int, depth int, pitch int, Rmask uint32, Gmask uint32, Bmask uint32, Amask uint32) *Surface {
			p := C.SDL_CreateRGBSurfaceFrom(unsafe.Pointer(pixels), C.int(width), C.int(height), C.int(depth), C.int(pitch),
				C.Uint32(Rmask), C.Uint32(Gmask), C.Uint32(Bmask), C.Uint32(Amask))
			return (*Surface)(cast(p))
		}
*/
package main

import (
	"errors"
	"flag"
	"github.com/salviati/go-opencl/cl"
	"io/ioutil"
	"log"
	"os"
	"sdl"
)

var file = flag.String("t", "lenna.png", "Test file")

var (
	platform cl.Platform
)

func init() {
	flag.Parse()

	platforms := cl.GetPlatforms()
	platform = platforms[0]

	if ImageSupport() == false {
		log.Fatal("Your device doesn't support images through OpenCL")
	}
}

func ImageSupport() bool {
	return platform.Devices[0].Property(cl.DEVICE_IMAGE_SUPPORT).(bool)
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
func mustReadFile(path string) string {
	fi, err := os.Open(path)
	check(err)
	src, err := ioutil.ReadAll(fi)
	check(err)
	return string(src)
}

var (
	c                             *cl.Context
	cq                            *cl.CommandQueue
	p                             *cl.Program
	k_shrink, k_enlarge, k_rotate *cl.Kernel
)

// init OpenCL & load the program
func initAndPrepCL() error {
	var err error

	params := make(map[cl.ContextParameter]interface{})
	c, err = cl.NewContextOfDevices(params, platform.Devices[0:1])
	if err != nil {
		return err
	}

	cq, err = c.NewCommandQueue(platform.Devices[0], cl.QUEUE_PROFILING_ENABLE|cl.QUEUE_OUT_OF_ORDER_EXEC_MODE_ENABLE)
	if err != nil {
		return err
	}

	p, err = c.NewProgramFromSource(mustReadFile("image.cl"))
	if err != nil {
		return err
	}

	k_shrink, err = p.NewKernelNamed("image_shrink")
	if err != nil {
		return err
	}

	k_enlarge, err = p.NewKernelNamed("image_enlarge")
	if err != nil {
		return err
	}

	k_rotate, err = p.NewKernelNamed("image_rotate")
	if err != nil {
		return err
	}

	return nil
}

func imageCall(k *cl.Kernel, s *sdl.Surface, dw, dh uint32, va ...interface{}) (*sdl.Surface, error) {
	order := cl.RGBA
	elemSize := 4

	src, err := c.NewImage2D(cl.MEM_READ_ONLY|cl.MEM_USE_HOST_PTR, order, cl.UNSIGNED_INT8,
		uint32(s.W), uint32(s.H), uint32(s.Pitch), s.Pixels)
	if err != nil {
		return nil, err
	}

	dst, err := c.NewImage2D(cl.MEM_WRITE_ONLY, order, cl.UNSIGNED_INT8,
		dw, dh, 0, nil)
	if err != nil {
		return nil, err
	}

	err = k.SetArg(0, src)
	if err != nil {
		return nil, err
	}

	err = k.SetArg(1, dst)
	if err != nil {
		return nil, err
	}

	for i, v := range va {
		err = k.SetArg(uint(i+2), v)
		if err != nil {
			return nil, err
		}
	}

	empty := make([]cl.Size, 0)
	gsize := []cl.Size{cl.Size(dw), cl.Size(dh)}
	err = cq.EnqueueKernel(k, empty, gsize, empty)

	pixels, err := cq.EnqueueReadImage(dst, [3]cl.Size{0, 0, 0}, [3]cl.Size{cl.Size(dw), cl.Size(dh), 1}, 0, 0)
	if err != nil {
		return nil, err
	}

	d := sdl.CreateRGBSurfaceFrom(&pixels[0],
		int(dw), int(dh), int(elemSize*8), int(elemSize)*int(dw),
		s.Format.Rmask, s.Format.Gmask, s.Format.Bmask, s.Format.Amask,
	)

	if d == nil {
		return nil, errors.New(sdl.GetError())
	}

	return d, nil
}

func shrink(s *sdl.Surface, factorx, factory float32) (*sdl.Surface, error) {
	dw := uint32(float32(s.W) / factorx)
	dh := uint32(float32(s.H) / factory)
	return imageCall(k_shrink, s, dw, dh, factorx, factory)
}

func enlarge(s *sdl.Surface, factorx, factory float32) (*sdl.Surface, error) {
	dw := uint32(float32(s.W) * factorx)
	dh := uint32(float32(s.H) * factory)
	return imageCall(k_enlarge, s, dw, dh, factorx, factory)
}

func rotate(s *sdl.Surface, angle float32) (*sdl.Surface, error) {
	dw := uint32(s.W)
	dh := uint32(s.H)
	return imageCall(k_rotate, s, dw, dh, angle)
}

func main() {
	if sdl.Init(sdl.INIT_EVERYTHING) != 0 {
		panic(sdl.GetError())
	}
	defer sdl.Quit()

	image0 := sdl.Load(*file)
	if image0 == nil {
		panic(sdl.GetError())
	}

	screen := sdl.SetVideoMode(int(image0.W), int(image0.H), 32, sdl.RESIZABLE|sdl.DOUBLEBUF)

	if screen == nil {
		panic(sdl.GetError())
	}

	image := sdl.DisplayFormat(image0)

	err := initAndPrepCL()
	check(err)

	e := new(sdl.Event)
	angle := float32(0)
	for running := true; running; angle += 0.001 {
		e.Poll()

		switch e.Type {
		case sdl.QUIT:
			running = false
			break
		case sdl.KEYDOWN:
			if e.Keyboard().Keysym.Sym == sdl.K_ESCAPE {
				running = false
				break
			}
		}

		news, err := rotate(image, angle)
		check(err)

		screen.FillRect(nil, 0)
		screen.Blit(&sdl.Rect{0, 0, 0, 0}, news, nil)
		screen.Flip()
		sdl.Delay(25)
	}
}
