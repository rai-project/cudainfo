// +build !cgo !linux !amd64,!ppc64le

package cudainfo

type nvmlDeviceHandle struct{}

// GetDeviceCount ...
func GetDeviceCount() (uint, error) {
	return 0, ErrNVMLUnavailable
}

// GetDriverVersion ...
func GetDriverVersion() (string, error) {
	return "", ErrNVMLUnavailable
}

// Status ...
func (d *NVMLDevice) Status() (*DeviceStatus, error) {
	return nil, ErrNVMLUnavailable
}

// GetP2PLink ...
func GetP2PLink(dev1, dev2 *NVMLDevice) (P2PLinkType, error) {
	return P2PLinkType(0), ErrNVMLUnavailable
}

// GetDevicePath ...
func GetDevicePath(idx uint) (string, error) {
	return "", ErrNVMLUnavailable
}

// NewNvmlDevice ...
func NewNvmlDevice(idx uint) (device *NVMLDevice, err error) {
	return nil, ErrNVMLUnavailable
}

func initNVMLLibrary() {

}
