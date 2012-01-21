package main

import (
	"fmt"
	"github.com/tones111/go-opencl/cl"
	"syscall"
)

func main() {
	platforms, err := cl.Platforms()
	if err != nil {
		fmt.Println("Error probing platforms:", err)
		syscall.Exit(1)
	}

	fmt.Println("Number of Platforms:", len(platforms))
	for _, platform := range platforms {
		properties, err := platform.Properties()
		if err != nil {
			fmt.Println("Error reading platform properties:", err, "\n\n")
			continue
		}
		fmt.Println("  Platform Profile:", properties[cl.PLATFORM_PROFILE])
		fmt.Println("  Platform Version:", properties[cl.PLATFORM_VERSION])
		fmt.Println("  Platform Name:", properties[cl.PLATFORM_NAME])
		fmt.Println("  Platform Vendor:", properties[cl.PLATFORM_VENDOR])
		fmt.Println("  Platform Extensions:", properties[cl.PLATFORM_EXTENSIONS], "\n\n")
		fmt.Println("  Platform Name:", properties[cl.PLATFORM_NAME])

		devices, err := platform.Devices(cl.DEVICE_TYPE_ALL)
		if err != nil {
			fmt.Println("Error probing devices:", err, "\n\n")
			continue
		}

		fmt.Println("Number of devices:", len(devices))
		for _, device := range devices {
			properties, err := device.Properties()
			if err != nil {
				fmt.Println("Error reading device properties:", err, "\n\n")
				continue
			}
			fmt.Println("  Device Type:", properties[cl.DEVICE_TYPE])
			//fmt.Println("  Device ID:", "TODO")
			//fmt.Println("  Board name:", "TODO")
			fmt.Println("  Max compute units:", properties[cl.DEVICE_MAX_COMPUTE_UNITS])
			fmt.Println("  Max work items dimensions:", properties[cl.DEVICE_MAX_WORK_ITEM_DIMENSIONS])
			//fmt.Println("    Max work items[]", "TODO")
			fmt.Println("  Max work group size:", properties[cl.DEVICE_MAX_WORK_GROUP_SIZE])
			fmt.Println("  Preferred vector width char:", properties[cl.DEVICE_PREFERRED_VECTOR_WIDTH_CHAR])
			fmt.Println("  Preferred vector width short:", properties[cl.DEVICE_PREFERRED_VECTOR_WIDTH_SHORT])
			fmt.Println("  Preferred vector width int:", properties[cl.DEVICE_PREFERRED_VECTOR_WIDTH_INT])
			fmt.Println("  Preferred vector width long:", properties[cl.DEVICE_PREFERRED_VECTOR_WIDTH_LONG])
			fmt.Println("  Preferred vector width float:", properties[cl.DEVICE_PREFERRED_VECTOR_WIDTH_FLOAT])
			fmt.Println("  Preferred vector width double:", properties[cl.DEVICE_PREFERRED_VECTOR_WIDTH_DOUBLE])
			fmt.Println("  Native vector width char:", properties[cl.DEVICE_NATIVE_VECTOR_WIDTH_CHAR])
			fmt.Println("  Native vector width short:", properties[cl.DEVICE_NATIVE_VECTOR_WIDTH_SHORT])
			fmt.Println("  Native vector width int:", properties[cl.DEVICE_NATIVE_VECTOR_WIDTH_INT])
			fmt.Println("  Native vector width long:", properties[cl.DEVICE_NATIVE_VECTOR_WIDTH_LONG])
			fmt.Println("  Native vector width float:", properties[cl.DEVICE_NATIVE_VECTOR_WIDTH_FLOAT])
			fmt.Println("  Native vector width double:", properties[cl.DEVICE_NATIVE_VECTOR_WIDTH_DOUBLE])
			fmt.Printf("  Max clock frequency: %dMhz\n", properties[cl.DEVICE_MAX_CLOCK_FREQUENCY])
			fmt.Println("  Address bits:", properties[cl.DEVICE_ADDRESS_BITS])
			fmt.Println("  Max memory allocation:", properties[cl.DEVICE_MAX_MEM_ALLOC_SIZE])
			fmt.Println("  Image support:", properties[cl.DEVICE_IMAGE_SUPPORT])
			fmt.Println("  Max number of images read arguments:", properties[cl.DEVICE_MAX_READ_IMAGE_ARGS])
			fmt.Println("  Max number of images write arguments:", properties[cl.DEVICE_MAX_WRITE_IMAGE_ARGS])
			fmt.Println("  Max image 2D width:", properties[cl.DEVICE_IMAGE2D_MAX_WIDTH])
			fmt.Println("  Max image 2D height:", properties[cl.DEVICE_IMAGE2D_MAX_HEIGHT])
			fmt.Println("  Max image 3D width:", properties[cl.DEVICE_IMAGE3D_MAX_WIDTH])
			fmt.Println("  Max image 3D height:", properties[cl.DEVICE_IMAGE3D_MAX_HEIGHT])
			fmt.Println("  Max image 3D depth:", properties[cl.DEVICE_IMAGE3D_MAX_DEPTH])
			fmt.Println("  Max samplers within kernel:", properties[cl.DEVICE_MAX_SAMPLERS])
			fmt.Println("  Max size of kernel argument:", properties[cl.DEVICE_MAX_PARAMETER_SIZE])
			fmt.Println("  Alignment (bits) of base address:", properties[cl.DEVICE_MEM_BASE_ADDR_ALIGN])
			fmt.Println("  Minimum alignment (bytes) for any datatype:", properties[cl.DEVICE_MIN_DATA_TYPE_ALIGN_SIZE])

			/*fmt.Println("  Single precision floating point capability")
			fmt.Println("    Denorms:", "TODO")
			fmt.Println("    Quiet NaNs:", "TODO")
			fmt.Println("    Round to nearest even:", "TODO")
			fmt.Println("    Round to zero:", "TODO")
			fmt.Println("    Round to +ve and infinity:", "TODO")
			fmt.Println("    IEEE754-2008 fused multiply-add:", "TODO")
			*/

			//fmt.Println("  Cache type:", "TODO" /*properties[cl.DEVICE_GLOBAL_MEM_CACHE_TYPE]*/ )
			//fmt.Println("  Cache line size:", "TODO" /*properties[cl.DEVICE_GLOBAL_MEM_CACHELINE_SIZE]*/ )
			fmt.Println("  Cache size:", properties[cl.DEVICE_GLOBAL_MEM_CACHE_SIZE])
			fmt.Println("  Global memory size:", properties[cl.DEVICE_GLOBAL_MEM_SIZE])
			fmt.Println("  Constant buffer size:", properties[cl.DEVICE_MAX_CONSTANT_BUFFER_SIZE])
			fmt.Println("  Max number of constant args:", properties[cl.DEVICE_MAX_CONSTANT_ARGS])
			//fmt.Println("  Local memory type:", "TODO" /*properties[cl.DEVICE_LOCAL_MEM_TYPE]*/ )
			fmt.Println("  Local memory size:", properties[cl.DEVICE_LOCAL_MEM_SIZE])
			//fmt.Println("  Kernel Preferred work group size multiple:", "TODO")
			fmt.Println("  Error correction support:", properties[cl.DEVICE_ERROR_CORRECTION_SUPPORT])
			fmt.Println("  Unified memory for Host and Device:", properties[cl.DEVICE_HOST_UNIFIED_MEMORY])
			fmt.Println("  Profiling timer resolution:", properties[cl.DEVICE_PROFILING_TIMER_RESOLUTION])
			fmt.Println("  Little endian:", properties[cl.DEVICE_ENDIAN_LITTLE])
			fmt.Println("  Available:", properties[cl.DEVICE_AVAILABLE])
			fmt.Println("  Compiler available:", properties[cl.DEVICE_COMPILER_AVAILABLE])

			/*fmt.Println("  Execution capabilities:")
			fmt.Println("    Execute OpenCL kernels:", "TODO")
			fmt.Println("    Execute native function:", "TODO")
			*/

			/*fmt.Println("  Queue properties:")
			fmt.Println("    Out-of-Order:", "TODO")
			fmt.Println("    Profiling:", "TODO")
			*/

			//fmt.Println("  Platform ID:", "TODO" /* properties[cl.DEVICE_PLATFORM]*/ )
			fmt.Println("  Name:", properties[cl.DEVICE_NAME])
			fmt.Println("  Vendor:", properties[cl.DEVICE_VENDOR])
			fmt.Println("  Device OpenCL C version:", properties[cl.DEVICE_OPENCL_C_VERSION])
			fmt.Println("  Driver version:", properties[cl.DRIVER_VERSION])
			fmt.Println("  Profile:", properties[cl.DEVICE_PROFILE])
			fmt.Println("  Version:", properties[cl.DEVICE_VERSION])
			fmt.Println("  Extensions:", properties[cl.DEVICE_EXTENSIONS])
		}
	}
}
