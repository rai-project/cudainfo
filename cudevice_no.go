// +build !cgo !linux !amd64

package cudainfo

type cuDeviceHandle struct{}

// GetCUDAVersion ...
func GetCUDAVersion() (string, error) {
	return "", ErrCGODisabled
} 
// NewCUDADevice ...
func NewCUDADevice(busID string) (*cuDevice, error) {
	return nil, ErrCGODisabled
} 
// NewCUDADeviceByIdx ...
func NewCUDADeviceByIdx(id int) (*cuDevice, error) {
	return nil, ErrCGODisabled
} 
// CanAccessPeer ...
func CanAccessPeer(dev1, dev2 *cuDevice) (bool, error) {
	return false, ErrCGODisabled
} 
// DeviceGetCount ...
func DeviceGetCount() (int, error) {
	return 0, ErrCGODisabled
}
