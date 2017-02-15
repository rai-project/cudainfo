// +build !cgo !linux !amd64

package cudainfo

type nvmlDeviceHandle struct{}

func GetDeviceCount() (uint, error) {
	return 0, ErrNVMLUnavailable
}

func GetDriverVersion() (string, error) {
	return "", ErrNVMLUnavailable
}

func (d *nvmlDevice) Status() (*DeviceStatus, error) {
	return nil, ErrNVMLUnavailable
}

func GetP2PLink(dev1, dev2 *nvmlDevice) (P2PLinkType, error) {
	return P2PLinkType(0), ErrNVMLUnavailable
}

func GetDevicePath(idx uint) (string, error) {
	return "", ErrNVMLUnavailable
}

func NewNvmlDevice(idx uint) (device *nvmlDevice, err error) {
	return nil, ErrNVMLUnavailable
}

func initNVMLLibrary() {

}
