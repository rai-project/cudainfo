package cudainfo

import (
	"errors"
	"os"
	"os/exec"
	"runtime"
)

type cuMemoryInfo struct {
	ECC       bool
	Global    uint
	Shared    uint // includes L1 cache
	Constant  uint
	L2Cache   uint
	Bandwidth uint
}

type cuDevice struct {
	handle cuDeviceHandle
	Family string
	Arch   string
	Cores  uint
	HyperQ uint
	Memory cuMemoryInfo
}

type nvmlDevice struct {
	handle nvmlDeviceHandle

	Model       string
	UUID        string
	Path        string
	Power       uint
	CPUAffinity uint
	PCI         PCIInfo
	Clocks      ClockInfo
	Topology    []P2PLink
}

type Device struct {
	Idx int
	*cuDevice
	*nvmlDevice
}

type P2PLinkType uint

type P2PLink struct {
	BusID string
	Link  P2PLinkType
}

type ClockInfo struct {
	Core   uint
	Memory uint
}

type PCIInfo struct {
	BusID     string
	BAR1      uint64
	Bandwidth uint
}

type UtilizationInfo struct {
	GPU     uint
	Encoder uint
	Decoder uint
}

type PCIThroughputInfo struct {
	RX uint
	TX uint
}

type PCIStatusInfo struct {
	BAR1Used   uint64
	Throughput PCIThroughputInfo
}

type ECCErrorsInfo struct {
	L1Cache uint64
	L2Cache uint64
	Global  uint64
}

type nvmlMemoryInfo struct {
	GlobalUsed uint64
	ECCErrors  ECCErrorsInfo
}

type ProcessInfo struct {
	PID  uint
	Name string
}

type DeviceStatus struct {
	Power       uint
	Temperature uint
	Utilization UtilizationInfo
	Memory      nvmlMemoryInfo
	Clocks      ClockInfo
	PCI         PCIStatusInfo
	Processes   []ProcessInfo
}

const (
	P2PLinkUnknown P2PLinkType = iota
	P2PLinkCrossCPU
	P2PLinkSameCPU
	P2PLinkHostBridge
	P2PLinkMultiSwitch
	P2PLinkSingleSwitch
	P2PLinkSameBoard
)

var (
	ErrCGODisabled        = errors.New("Cannot get device information since CGO is disabled.")
	ErrCPUAffinity        = errors.New("failed to retrieve CPU affinity")
	ErrUnsupportedP2PLink = errors.New("unsupported P2P link type")
	ErrNVMLUnavailable    = errors.New("NVML is unavailable on the system. It's only available on linux with CGO enabled.")
	archToFamily          = map[string]string{
		"1": "Tesla",
		"2": "Fermi",
		"3": "Kepler",
		"5": "Maxwell",
		"6": "Pascal",
	}
	archToCoresPerSM = map[string]uint{
		"1.0": 8,   // Tesla Generation (SM 1.0) G80 class
		"1.1": 8,   // Tesla Generation (SM 1.1) G8x G9x class
		"1.2": 8,   // Tesla Generation (SM 1.2) GT21x class
		"1.3": 8,   // Tesla Generation (SM 1.3) GT20x class
		"2.0": 32,  // Fermi Generation (SM 2.0) GF100 GF110 class
		"2.1": 48,  // Fermi Generation (SM 2.1) GF10x GF11x class
		"3.0": 192, // Kepler Generation (SM 3.0) GK10x class
		"3.2": 192, // Kepler Generation (SM 3.2) TK1 class
		"3.5": 192, // Kepler Generation (SM 3.5) GK11x GK20x class
		"3.7": 192, // Kepler Generation (SM 3.7) GK21x class
		"5.0": 128, // Maxwell Generation (SM 5.0) GM10x class
		"5.2": 128, // Maxwell Generation (SM 5.2) GM20x class
		"5.3": 128, // Maxwell Generation (SM 5.3) TX1 class
		"6.0": 64,  // Pascal Generation (SM 6.0) GP100 class
		"6.1": 128, // Pascal Generation (SM 6.1) GP10x class
		"6.2": 128, // Pascal Generation (SM 6.2) GP10x class
	}

	archToHyperQ = map[string]uint{
		"1.0": 1,  // Tesla Generation (SM 1.0) G80 class
		"1.1": 1,  // Tesla Generation (SM 1.1) G8x G9x class
		"1.2": 1,  // Tesla Generation (SM 1.2) GT21x class
		"1.3": 1,  // Tesla Generation (SM 1.3) GT20x class
		"2.0": 1,  // Fermi Generation (SM 2.0) GF100 GF110 class
		"2.1": 1,  // Fermi Generation (SM 2.1) GF10x GF11x class
		"3.0": 1,  // Kepler Generation (SM 3.0) GK10x class
		"3.2": 1,  // Kepler Generation (SM 3.2) TK1 class
		"3.5": 32, // Kepler Generation (SM 3.5) GK11x GK20x class
		"3.7": 32, // Kepler Generation (SM 3.7) GK21x class
		"5.0": 32, // Maxwell Generation (SM 5.0) GM10x class
		"5.2": 32, // Maxwell Generation (SM 5.2) GM20x class
		"5.3": 32, // Maxwell Generation (SM 5.3) TX1 class
		"6.0": 32, // Pascal Generation (SM 6.0) GP100 class
		"6.1": 32, // Pascal Generation (SM 6.1) GP10x class
		"6.2": 32, // Pascal Generation (SM 6.2) GP10x class
	}
)

func (t P2PLinkType) String() string {
	switch t {
	case P2PLinkCrossCPU:
		return "Cross CPU socket"
	case P2PLinkSameCPU:
		return "Same CPU socket"
	case P2PLinkHostBridge:
		return "Host PCI bridge"
	case P2PLinkMultiSwitch:
		return "Multiple PCI switches"
	case P2PLinkSingleSwitch:
		return "Single PCI switch"
	case P2PLinkSameBoard:
		return "Same board"
	case P2PLinkUnknown:
	}
	return "???"
}

func GetCUDADevice(id int) (*Device, error) {
	if runtime.GOOS != "linux" {
		return nil, errors.New("Unsupported OS")
	}
	cuDev, err := NewCUDADeviceByIdx(id)
	if err != nil {
		return nil, err
	}
	nvmlDev, err := NewNvmlDevice(uint(id))
	if err != nil {
		return nil, err
	}
	return &Device{
		Idx:        id,
		cuDevice:   cuDev,
		nvmlDevice: nvmlDev,
	}, nil
}

func LoadUVM() error {
	if runtime.GOOS != "linux" {
		return nil
	}
	if _, err := os.Stat("/dev/nvidia-uvm"); err == nil {
		return nil
	}
	return exec.Command("nvidia-modprobe", "-u", "-c=0").Run()
}
