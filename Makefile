include $(GOROOT)/src/Make.inc

TARG=opencl
CGOFILES=\
   commandQueue.go\
   context.go\
   device.go\
   opencl.go\
   platform.go\

include $(GOROOT)/src/Make.pkg
