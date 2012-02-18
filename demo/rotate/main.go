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
	c        *cl.Context
	cq       *cl.CommandQueue
	p        *cl.Program
	kernels  map[string]*cl.Kernel
)

var kernelNames = []string{
	"image_shrink", "image_enlarge", "image_rotate", 
	"image_flip_h", "image_flip_v", "image_flip_hv",
}

func init() {
	kernels = make(map[string]*cl.Kernel)
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

	addKernel := func(name string) (err error) {
		kernels[name], err = p.NewKernelNamed(name)
		return
	}

	for _, kernelName := range kernelNames {
		if err := addKernel(kernelName); err != nil {
			return err
		}
	}

	return nil
}

func imageCall(kernelName string, src, dst *cl.Image, dstw, dsth uint32, va ...interface{}) ([]byte, error) {
	k, ok := kernels[kernelName]
	if !ok {
		return nil, errors.New("unknown kernel " + kernelName)
	}

	err := k.SetArg(0, src)
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
	gsize := []cl.Size{cl.Size(dstw), cl.Size(dsth)}
	err = cq.EnqueueKernel(k, empty, gsize, empty)

	pixels, err := cq.EnqueueReadImage(dst, [3]cl.Size{0, 0, 0}, [3]cl.Size{cl.Size(dstw), cl.Size(dsth), 1}, 0, 0)
	if err != nil {
		return nil, err
	}

	return pixels, nil
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

	dw := uint32(image0.W)
	dh := uint32(image0.H)

	screen := sdl.SetVideoMode(int(dw), int(dh), 32, sdl.DOUBLEBUF)

	if screen == nil {
		panic(sdl.GetError())
	}

	image := sdl.DisplayFormat(image0)
	format := image.Format

	err := initAndPrepCL()
	check(err)

	order := cl.RGBA
	elemSize := 4

	src, err := c.NewImage2D(cl.MEM_READ_ONLY|cl.MEM_USE_HOST_PTR, order, cl.UNSIGNED_INT8,
		uint32(image.W), uint32(image.H), uint32(image.Pitch), image.Pixels)
	check(err)

	dst, err := c.NewImage2D(cl.MEM_WRITE_ONLY, order, cl.UNSIGNED_INT8,
		dw, dh, 0, nil)
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

		pixels, err := imageCall("image_rotate", src, dst, dw, dh, angle)
		check(err)

		news := sdl.CreateRGBSurfaceFrom(&pixels[0],
			int(dw), int(dh), int(elemSize*8), int(elemSize)*int(dw),
			format.Rmask, format.Gmask, format.Bmask, format.Amask,
		)

		if news == nil {
			log.Fatal(sdl.GetError())
		}

		screen.FillRect(nil, 0)
		screen.Blit(&sdl.Rect{0, 0, 0, 0}, news, nil)
		screen.Flip()
		sdl.Delay(25)
	}
}
