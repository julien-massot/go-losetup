package losetup

import (
	"fmt"
	"os"
	"golang.org/x/sys/unix"
)

// Device represents a loop device /dev/loop#
type Device struct {
	// device number (i.e. 7 if /dev/loop7)
	number uint64

	// flags with which to open the device with
	flags int
}

// New creates a reference to a specific loop device, if you know which one you
// want to reference.
func New(number uint64, flags int) Device {
	return Device{number, flags}
}

// open returns a file handle to /dev/loop# and returns an error if it cannot
// be opened.
func (device Device) open() (*os.File, error) {
	return os.OpenFile(device.Path(), device.flags, 0660)
}

// Path returns the path to the loopback device
func (device Device) Path() string {
	return fmt.Sprintf(DeviceFormatString, device.number)
}

// Change the default block size (512 bytes) for a loop device,
// returns an error if the blocksize can't be set.
func (device Device) SetBlockSize(blockSize uint32) error {
	f, err := device.open()
	if err != nil {
		return fmt.Errorf("couldn't open loop device: %v", err)
	}

	defer f.Close()

	_, _, err = unix.Syscall(unix.SYS_IOCTL, f.Fd(), SetBlockSize, uintptr(blockSize))
	if err != 0 {
		return fmt.Errorf("Failed to set block size: %w", err)
	}

	return nil
}
