include $(GOROOT)/src/Make.inc

TARG=opencl
CGOFILES=\
   buffer.go\
   commandQueue.go\
   context.go\
   device.go\
   opencl.go\
   platform.go\
   program.go\

include $(GOROOT)/src/Make.pkg
