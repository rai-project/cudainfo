// +build cgo,linux,amd64

package cudainfo

// This file provides CGO flags to find CUDA libraries and headers.

//#cgo LDFLAGS: -lcudart_static -lnvidia-ml -ldl -Wl,--unresolved-symbols=ignore-in-object-files
//#cgo CFLAGS: -I/usr/local/cuda/include -I /usr/include/nvidia/gdk
//#cgo LDFLAGS: -L/usr/local/cuda/lib64 -L /usr/src/gdk/nvml/lib/ -L /usr/lib/nvidia-367 -L /usr/lib/nvidia-377 -L /usr/lib/nvidia-378
import "C"
