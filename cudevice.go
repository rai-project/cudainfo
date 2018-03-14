// +build cgo,linux,amd64 cgo,linux,ppc64le

package cudainfo

// #include <stdlib.h>
// #include <cuda_runtime_api.h>
import "C"
import (
	"errors"
	"fmt"
	"unsafe"
)

type cuDeviceHandle C.int

func cudaErr(ret C.cudaError_t) error {
	if ret == C.cudaSuccess {
		return nil
	}
	err := C.GoString(C.cudaGetErrorString(ret))
	return errors.New(err)
}

// GetCUDAVersion ...
func GetCUDAVersion() (string, error) {
	var driver C.int
	err := cudaErr(C.cudaDriverGetVersion(&driver))
	d := fmt.Sprintf("%d.%d", int(driver)/1000, int(driver)%100/10)
	return d, err
}

// NewCUDADeviceByIdx ...
func NewCUDADeviceByIdx(id int) (*cuDevice, error) {
	var prop C.struct_cudaDeviceProp

	dev := C.int(id)

	if err := cudaErr(C.cudaGetDeviceProperties(&prop, dev)); err != nil {
		return nil, err
	}
	arch := fmt.Sprintf("%d.%d", prop.major, prop.minor)
	cores, ok := archToCoresPerSM[arch]
	if !ok {
		return nil, fmt.Errorf("unsupported CUDA arch: %s", arch)
	}
	hyperq, ok := archToHyperQ[arch]
	if !ok {
		return nil, fmt.Errorf("unsupported CUDA arch: %s", arch)
	}
	// Destroy the active CUDA context
	cudaErr(C.cudaDeviceReset())
	return &cuDevice{
		handle: cuDeviceHandle(dev),
		Family: archToFamily[arch[:1]],
		Arch:   arch,
		Cores:  cores * uint(prop.multiProcessorCount),
		HyperQ: hyperq,
		Memory: cuMemoryInfo{
			ECC:       bool(prop.ECCEnabled != 0),
			Global:    uint(prop.totalGlobalMem / (1024 * 1024)),
			Shared:    uint(prop.sharedMemPerMultiprocessor / 1024),
			Constant:  uint(prop.totalConstMem / 1024),
			L2Cache:   uint(prop.l2CacheSize / 1024),
			Bandwidth: 2 * uint((prop.memoryClockRate/1000)*(prop.memoryBusWidth/8)) / 1000,
		},
	}, nil
}

// NewCUDADevice ...
func NewCUDADevice(busID string) (*cuDevice, error) {
	var dev C.int
	id := C.CString(busID)
	if err := cudaErr(C.cudaDeviceGetByPCIBusId(&dev, id)); err != nil {
		return nil, err
	}
	C.free(unsafe.Pointer(id))
	return NewCUDADeviceByIdx(int(dev))
}

// CanAccessPeer ...
func CanAccessPeer(dev1, dev2 *cuDevice) (bool, error) {
	var ok C.int
	err := cudaErr(C.cudaDeviceCanAccessPeer(&ok, C.int(dev1.handle), C.int(dev2.handle)))
	return (ok != 0), err
}

// Returns the number of devices with compute capability greater than or equal to 1.0 that are available for execution.
func DeviceGetCount() (int, error) {
	var count C.int
	err := cudaErr(C.cudaGetDeviceCount(&count))
	if err != nil {
		return 0, err
	}
	return int(count), nil
}
