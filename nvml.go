// +build cgo,linux,amd64 cgo,linux,ppc64le

// Copyright (c) 2015-2016, NVIDIA CORPORATION. All rights reserved.

package cudainfo

// #include <nvml.h>
import "C"

import (
	"errors"
	"fmt"
	"time"
)

type nvmlDeviceHandle C.nvmlDevice_t

const (
	szDriver   = C.NVML_SYSTEM_DRIVER_VERSION_BUFFER_SIZE
	szModel    = C.NVML_DEVICE_NAME_BUFFER_SIZE
	szUUID     = C.NVML_DEVICE_UUID_BUFFER_SIZE
	szProcs    = 32
	szProcName = 64
)

func nvmlErr(ret C.nvmlReturn_t) error {
	if ret == C.NVML_SUCCESS {
		return nil
	}
	err := C.GoString(C.nvmlErrorString(ret))
	return errors.New(err)
}

func check(ret C.nvmlReturn_t) {
	if err := nvmlErr(ret); err != nil {
		panic(err)
	}
}

// GetDeviceCount ...
func GetDeviceCount() (uint, error) {
	var n C.uint

	err := nvmlErr(C.nvmlDeviceGetCount(&n))
	return uint(n), err
}

// GetDriverVersion ...
func GetDriverVersion() (string, error) {
	var driver [szDriver]C.char

	err := nvmlErr(C.nvmlSystemGetDriverVersion(&driver[0], szDriver))
	return C.GoString(&driver[0]), err
}

var pcieGenToBandwidth = map[int]uint{
	1: 250, // MB/s
	2: 500,
	3: 985,
	4: 1969,
}

// NewNvmlDevice ...
func NewNvmlDevice(idx uint) (device *NVMLDevice, err error) {
	var (
		dev   C.nvmlDevice_t
		model [szModel]C.char
		uuid  [szUUID]C.char
		pci   C.nvmlPciInfo_t
		minor C.uint
		bar1  C.nvmlBAR1Memory_t
		power C.uint
		clock [2]C.uint
		pciel [2]C.uint
		//cpus  cpuSet
	)

	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()

	check(C.nvmlDeviceGetHandleByIndex(C.uint(idx), &dev))
	check(C.nvmlDeviceGetName(dev, &model[0], szModel))
	check(C.nvmlDeviceGetUUID(dev, &uuid[0], szUUID))
	check(C.nvmlDeviceGetPciInfo(dev, &pci))
	check(C.nvmlDeviceGetMinorNumber(dev, &minor))
	check(C.nvmlDeviceGetBAR1MemoryInfo(dev, &bar1))
	check(C.nvmlDeviceGetPowerManagementLimit(dev, &power))
	check(C.nvmlDeviceGetMaxClockInfo(dev, C.NVML_CLOCK_SM, &clock[0]))
	check(C.nvmlDeviceGetMaxClockInfo(dev, C.NVML_CLOCK_MEM, &clock[1]))
	check(C.nvmlDeviceGetMaxPcieLinkGeneration(dev, &pciel[0]))
	check(C.nvmlDeviceGetMaxPcieLinkWidth(dev, &pciel[1]))
	//check(C.nvmlDeviceGetCpuAffinity(dev, C.uint(len(cpus)), (*C.ulong)(&cpus[0])))
	//node, err := getCPUNode(cpus)
	//if err != nil {
	//	return nil, err
	//}

	device = &NVMLDevice{
		Handle: nvmlDeviceHandle(dev),
		Model:  C.GoString(&model[0]),
		UUID:   C.GoString(&uuid[0]),
		Path:   fmt.Sprintf("/dev/nvidia%d", uint(minor)),
		Power:  uint(power / 1000),
		//	CPUAffinity: node,
		PCI: PCIInfo{
			BusID:     C.GoString(&pci.busId[0]),
			BAR1:      uint64(bar1.bar1Total / (1024 * 1024)),
			Bandwidth: pcieGenToBandwidth[int(pciel[0])] * uint(pciel[1]) / 1000,
		},
		Clocks: ClockInfo{
			Core:   uint(clock[0]),
			Memory: uint(clock[1]),
		},
	}
	return
}

// Status ...
func (d *NVMLDevice) Status() (status *DeviceStatus, err error) {
	var (
		power      C.uint
		temp       C.uint
		usage      C.nvmlUtilization_t
		encoder    [2]C.uint
		decoder    [2]C.uint
		mem        C.nvmlMemory_t
		ecc        [3]C.ulonglong
		clock      [2]C.uint
		bar1       C.nvmlBAR1Memory_t
		throughput [2]C.uint
		procname   [szProcName]C.char
		procs      [szProcs]C.nvmlProcessInfo_t
		nprocs     = C.uint(szProcs)
		timestamp  = time.Now()
	)

	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()

	check(C.nvmlDeviceGetPowerUsage(d.Handle, &power))
	check(C.nvmlDeviceGetTemperature(d.Handle, C.NVML_TEMPERATURE_GPU, &temp))
	check(C.nvmlDeviceGetUtilizationRates(d.Handle, &usage))
	check(C.nvmlDeviceGetEncoderUtilization(d.Handle, &encoder[0], &encoder[1]))
	check(C.nvmlDeviceGetDecoderUtilization(d.Handle, &decoder[0], &decoder[1]))
	check(C.nvmlDeviceGetMemoryInfo(d.Handle, &mem))
	check(C.nvmlDeviceGetClockInfo(d.Handle, C.NVML_CLOCK_SM, &clock[0]))
	check(C.nvmlDeviceGetClockInfo(d.Handle, C.NVML_CLOCK_MEM, &clock[1]))
	check(C.nvmlDeviceGetBAR1MemoryInfo(d.Handle, &bar1))
	check(C.nvmlDeviceGetComputeRunningProcesses(d.Handle, &nprocs, &procs[0]))

	status = &DeviceStatus{
		TimeStamp:   timestamp,
		Power:       uint(power / 1000),
		Temperature: uint(temp),
		Utilization: UtilizationInfo{
			GPU:     uint(usage.gpu),
			Memory:  uint(usage.memory),
			Encoder: uint(encoder[0]),
			Decoder: uint(decoder[0]),
		},
		Memory: nvmlMemoryStatus{
			Used: uint64(mem.used),
			Free: uint64(mem.free),
		},
		Clocks: ClockInfo{
			Core:   uint(clock[0]),
			Memory: uint(clock[1]),
		},
		PCI: PCIStatusInfo{
			BAR1Used: uint64(bar1.bar1Used / (1024 * 1024)),
		},
	}

	r := C.nvmlDeviceGetMemoryErrorCounter(d.Handle, C.NVML_MEMORY_ERROR_TYPE_UNCORRECTED, C.NVML_VOLATILE_ECC,
		C.NVML_MEMORY_LOCATION_L1_CACHE, &ecc[0])
	if r != C.NVML_ERROR_NOT_SUPPORTED { // only supported on Tesla cards
		check(r)
		check(C.nvmlDeviceGetMemoryErrorCounter(d.Handle, C.NVML_MEMORY_ERROR_TYPE_UNCORRECTED, C.NVML_VOLATILE_ECC,
			C.NVML_MEMORY_LOCATION_L2_CACHE, &ecc[1]))
		check(C.nvmlDeviceGetMemoryErrorCounter(d.Handle, C.NVML_MEMORY_ERROR_TYPE_UNCORRECTED, C.NVML_VOLATILE_ECC,
			C.NVML_MEMORY_LOCATION_DEVICE_MEMORY, &ecc[2]))
		status.Memory.ECCErrors = ECCErrorsInfo{uint64(ecc[0]), uint64(ecc[1]), uint64(ecc[2])}
	}

	r = C.nvmlDeviceGetPcieThroughput(d.Handle, C.NVML_PCIE_UTIL_RX_BYTES, &throughput[0])
	if r != C.NVML_ERROR_NOT_SUPPORTED { // only supported on Maxwell or newer
		check(r)
		check(C.nvmlDeviceGetPcieThroughput(d.Handle, C.NVML_PCIE_UTIL_TX_BYTES, &throughput[1]))
		status.PCI.Throughput = PCIThroughputInfo{uint(throughput[0]), uint(throughput[1])}
	}

	status.Processes = make([]ProcessInfo, nprocs)
	for i := range status.Processes {
		status.Processes[i].PID = uint(procs[i].pid)
		check(C.nvmlSystemGetProcessName(procs[i].pid, &procname[0], szProcName))
		status.Processes[i].Name = C.GoString(&procname[0])
	}
	return
}

// GetP2PLink ...
func GetP2PLink(dev1, dev2 *NVMLDevice) (link P2PLinkType, err error) {
	var level C.nvmlGpuTopologyLevel_t

	r := C.nvmlDeviceGetTopologyCommonAncestor(dev1.Handle, dev2.Handle, &level)
	if r == C.NVML_ERROR_FUNCTION_NOT_FOUND {
		return P2PLinkUnknown, nil
	}
	if err = nvmlErr(r); err != nil {
		return
	}
	switch level {
	case C.NVML_TOPOLOGY_INTERNAL:
		link = P2PLinkSameBoard
	case C.NVML_TOPOLOGY_SINGLE:
		link = P2PLinkSingleSwitch
	case C.NVML_TOPOLOGY_MULTIPLE:
		link = P2PLinkMultiSwitch
	case C.NVML_TOPOLOGY_HOSTBRIDGE:
		link = P2PLinkHostBridge
	case C.NVML_TOPOLOGY_CPU:
		link = P2PLinkSameCPU
	case C.NVML_TOPOLOGY_SYSTEM:
		link = P2PLinkCrossCPU
	default:
		err = ErrUnsupportedP2PLink
	}
	return
}

// GetDevicePath ...
func GetDevicePath(idx uint) (path string, err error) {
	var dev C.nvmlDevice_t
	var minor C.uint

	err = nvmlErr(C.nvmlDeviceGetHandleByIndex(C.uint(idx), &dev))
	if err != nil {
		return
	}
	err = nvmlErr(C.nvmlDeviceGetMinorNumber(dev, &minor))
	path = fmt.Sprintf("/dev/nvidia%d", uint(minor))
	return
}
func initNVMLLibrary() {
	C.nvmlInit()
}
