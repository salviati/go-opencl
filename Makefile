include $(GOROOT)/src/Make.inc

TARG=opencl
CGOFILES=\
   opencl.go\
   platform.go\

include $(GOROOT)/src/Make.pkg
