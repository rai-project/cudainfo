// +build !cgo !linux !amd64

package cudainfo

type cuDeviceHandle struct{}

func GetCUDAVersion() (string, error) {
	return "", ErrCGODisabled
}
func NewCUDADevice(busID string) (*cuDevice, error) {
	return nil, ErrCGODisabled
}
func NewCUDADeviceByIdx(id int) (*cuDevice, error) {
	return nil, ErrCGODisabled
}
func CanAccessPeer(dev1, dev2 *cuDevice) (bool, error) {
	return false, ErrCGODisabled
}
func DeviceGetCount() (int, error) {
	return 0, ErrCGODisabled
}
