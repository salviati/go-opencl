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

import (
	"fmt"
	"testing"
)

func getPlatform(t *testing.T) Platform {
	if len(Platforms) == 0 {
		t.Fatal("No platforms found")
	}
	return Platforms[0]
}

func getCPUDevice(p Platform, t *testing.T) []Device {
	fmt.Println("Devices:", len(p.Devices))
	if len(p.Devices) == 0 {
		t.Fatal("No devices found")
	}
	return p.Devices[:1]
}

func getContext(p Platform, d []Device, t *testing.T) *Context {
	if context, err := NewContextOfDevices(map[ContextParameter]interface{}{CONTEXT_PLATFORM: p}, d); err != nil {
		t.Fatal("Error creating context:", err)
	} else {
		return context
	}
	return nil
}

func getProgram(c *Context, source string, t *testing.T) *Program {
	if program, err := c.NewProgramFromSource(source); err != nil {
		t.Fatal("Error creating program:", err)
	} else {
		return program
	}
	return nil
}

func getKernel(p *Program, name string, t *testing.T) *Kernel {
	if kernel, err := p.NewKernelNamed(name); err != nil {
		t.Fatal("Error creating kernel:", err)
	} else {
		return kernel
	}
	return nil
}

func getQueue(c *Context, d Device, t *testing.T) *CommandQueue {
	if queue, err := c.NewCommandQueue(d, QUEUE_NIL); err != nil {
		t.Fatal("Error creating command queue:", err)
	} else {
		return queue
	}
	return nil
}

func Test_OpenCl(t *testing.T) {
	var square_source = `
__kernel void hello(__global uchar *input, __global uchar *output)
{
   size_t id = get_global_id(0);
   output[id] = input[id] * input[id];
}`

	inData := []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	var err error

	platform := getPlatform(t)
	devices := getCPUDevice(platform, t)
	context := getContext(platform, devices, t)
	queue := getQueue(context, devices[0], t)

	program := getProgram(context, square_source, t)
	kernel := getKernel(program, "hello", t)

	var inBuf, outBuf *Buffer
	if inBuf, err = context.NewBuffer(MEM_READ_ONLY, 100); err != nil {
		t.Fatal("Error creating in buffer:", err)
	} else if outBuf, err = context.NewBuffer(MEM_WRITE_ONLY, 100); err != nil {
		t.Fatal("Error creating out buffer:", err)
	}

	if err = queue.EnqueueWriteBuffer(inBuf, inData, 0); err != nil {
		t.Fatal("Error enquing data:", err)
	}

	if err = kernel.SetArg(0, inBuf); err != nil {
		t.Fatal("Error setting kernel arg 0:", err)
	} else if err = kernel.SetArg(1, outBuf); err != nil {
		t.Fatal("Error setting kernel arg 1:", err)
	} else if err = queue.EnqueueKernel(kernel, 0, uint(len(inData)), uint(len(inData))); err != nil {
		t.Fatal("Error enquing kernel:", err)
	}

	var data []byte
	if data, err = queue.EnqueueReadBuffer(outBuf, 0, uint32(len(inData))); err != nil {
		t.Fatal("Error reading data:", err)
	}

	for i, v := range data {
		fmt.Println("Data[", i, "] = ", v)
		if v != inData[i]*inData[i] {
			t.Fatal("Incorrect results")
		}
	}
}
