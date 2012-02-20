package main

import (
	"errors"
	"flag"
	"github.com/banthar/Go-SDL/sdl"
	"github.com/salviati/go-opencl/cl"
	"io/ioutil"
	"log"
	"os"
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
	"image_recscale", "image_rotate",
	"image_flip_h", "image_flip_v", "image_flip_hv",
	"image_affine", "image_affine2",
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

// Will call kernelName(src, dst, va...) on the device side.
func imageCall(dstW, dstH uint32, kernelName string, src, dst *cl.Image, va ...interface{}) ([]byte, error) {
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
	gsize := []cl.Size{cl.Size(dstW), cl.Size(dstH)}
	err = cq.EnqueueKernel(k, empty, gsize, empty)

	pixels, err := cq.EnqueueReadImage(dst, [3]cl.Size{0, 0, 0}, [3]cl.Size{cl.Size(dstW), cl.Size(dstH), 1}, 0, 0)
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

	factorx := float32(1)
	factory := float32(1)

	dstW := uint32(float32(image0.W) * factorx)
	dstH := uint32(float32(image0.H) * factory)

	screen := sdl.SetVideoMode(int(dstW), int(dstH), 32, sdl.DOUBLEBUF)

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
		dstW, dstH, 0, nil)
	check(err)

	angle := float32(0)
	for running := true; running; angle += 0.001 {
		for e := sdl.PollEvent(); e != nil; e = sdl.PollEvent() {
			switch ev := e.(type) {
			case *sdl.QuitEvent:
				running = false
				break
			case *sdl.KeyboardEvent:
				if ev.Keysym.Sym == sdl.K_ESCAPE {
					running = false
					break
				}
			}
		}

		//pixels, err := imageCall(dstW, dstH, "image_recscale", src, dst, 1/factorx, 1/factory)

		//pixels, err := imageCall(dstW, dstH, "image_rotate", src, dst, angle)
		m := mul(S(1.2, 1.2), R(angle))
		off := []float32{float32(dstW / 2), float32(dstH / 2)}
		pixels, err := imageCall(dstW, dstH, "image_affine2", src, dst,
			[]float32(m[0:2]), []float32(m[2:4]),
			off, off)
		check(err)

		news := sdl.CreateRGBSurfaceFrom(&pixels[0],
			int(dstW), int(dstH), int(elemSize*8), int(elemSize)*int(dstW),
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
