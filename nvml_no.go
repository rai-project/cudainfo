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
func (d *nvmlDevice) Status() (*DeviceStatus, error) {
	return nil, ErrNVMLUnavailable
}

// GetP2PLink ...
func GetP2PLink(dev1, dev2 *nvmlDevice) (P2PLinkType, error) {
	return P2PLinkType(0), ErrNVMLUnavailable
}

// GetDevicePath ...
func GetDevicePath(idx uint) (string, error) {
	return "", ErrNVMLUnavailable
}

// NewNvmlDevice ...
func NewNvmlDevice(idx uint) (device *nvmlDevice, err error) {
	return nil, ErrNVMLUnavailable
}

func initNVMLLibrary() {

}
