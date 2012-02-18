package main

import (
	"flag"
	cl "github.com/salviati/go-opencl/cl"
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
	c  *cl.Context
	cq *cl.CommandQueue
	p  *cl.Program
	k  *cl.Kernel
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

	p, err = c.NewProgramFromSource(mustReadFile("rotate.cl"))
	if err != nil {
		return err
	}

	k, err = p.NewKernelNamed("rotateImage")
	if err != nil {
		return err
	}

	return nil
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

	imageIn, err := c.NewImage2D(cl.MEM_READ_ONLY|cl.MEM_USE_HOST_PTR, cl.RGBA, cl.UNSIGNED_INT8,
		uint32(image.W), uint32(image.H), uint32(image.Pitch), image.Pixels)
	check(err)

	imageOut, err := c.NewImage2D(cl.MEM_WRITE_ONLY, cl.RGBA, cl.UNSIGNED_INT8,
		uint32(image.W), uint32(image.H), 0, nil)
	check(err)

	elemSize, err := imageOut.Info(cl.IMAGE_ELEMENT_SIZE)

	/*
		https://github.com/banthar/Go-SDL doesn't implement CreateRGBSurfaceFrom, since pixels may be garbage collected in general.

		func CreateRGBSurfaceFrom(pixels *byte, width int, height int, depth int, pitch int, Rmask uint32, Gmask uint32, Bmask uint32, Amask uint32) *Surface {
			p := C.SDL_CreateRGBSurfaceFrom(unsafe.Pointer(pixels), C.int(width), C.int(height), C.int(depth), C.int(pitch),
				C.Uint32(Rmask), C.Uint32(Gmask), C.Uint32(Bmask), C.Uint32(Amask))
			return (*Surface)(cast(p))
		}
	*/
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

		err = k.SetArg(0, imageIn)
		check(err)

		err = k.SetArg(1, imageOut)
		check(err)

		err = k.SetArg(2, angle)
		check(err)

		empty := make([]cl.Size, 0)
		gsize := []cl.Size{cl.Size(image.W), cl.Size(image.H)}
		err = cq.EnqueueKernel(k, empty, gsize, empty)

		pixels, err := cq.EnqueueReadImage(imageOut, [3]cl.Size{0, 0, 0}, [3]cl.Size{cl.Size(image.W), cl.Size(image.H), 1}, 0, 0)
		check(err)

		rotated := sdl.CreateRGBSurfaceFrom(&pixels[0], int(image.W), int(image.H), int(elemSize*8), int(elemSize)*int(image.W), image.Format.Rmask, image.Format.Gmask, image.Format.Bmask, image.Format.Amask)

		screen.FillRect(nil, 0)
		screen.Blit(&sdl.Rect{0, 0, 0, 0}, rotated, nil)
		screen.Flip()
		sdl.Delay(25)
	}
}
