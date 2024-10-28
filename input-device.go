package input

import (
	"bytes"
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"strings"
	"unsafe"
)

// Device holds metadata for an input device.
type Device struct {
	DeviceCapabilities

	ID        int
	Path      string
	SysFsPath string
	Name      string

	handle *os.File
}

// ListKeyboards returns a list of Devices that support the EV_KEY event type, and support common keyboard keys.
func ListKeyboards() ([]*Device, error) {
	var keyboards []*Device
	devices, err := ListDevices()
	if err != nil {
		return nil, err
	}

keyboardCheck:
	for _, device := range devices {
		if !device.SupportsEvent(EV_KEY) || device.SupportsEvent(EV_REL) || device.SupportsEvent(EV_ABS) {
			continue
		}

		// Check for common key codes, Escape, A, B, C and Enter
		common := []int{1, 30, 46, 48, 28}

		for _, key := range common {
			if !device.SupportsKey(key) {
				continue keyboardCheck
			}
		}

		keyboards = append(keyboards, device)
	}

	return keyboards, nil
}

// ListPointerDevices returns a list of Devices that support the EV_REL or EV_ABS event type. These are normally pointer devices like mice and touch-pads/track-pads.
func ListPointerDevices() ([]*Device, error) {
	var pointers []*Device
	devices, err := ListDevices()
	if err != nil {
		return nil, err
	}

	for _, device := range devices {
		if device.SupportsEvent(EV_REL) || device.SupportsEvent(EV_ABS) {
			pointers = append(pointers, device)
		}
	}

	return pointers, nil
}

// ListMice returns a list of Devices that support the EV_REL event type. These are normally mice.
func ListMice() ([]*Device, error) {
	devices, err := ListDevices()
	if err != nil {
		return nil, err
	}

	mice := []*Device{}

	for _, device := range devices {
		if device.SupportsEvent(EV_REL) {
			mice = append(mice, device)
		}
	}

	return mice, nil
}

// ListTouchDevices returns a list of Devices that support the EV_ABS event type. These are normally touch-pads/track-pads.
func ListTouchPads() ([]*Device, error) {
	devices, err := ListDevices()
	if err != nil {
		return nil, err
	}

	touchPads := []*Device{}

	for _, device := range devices {
		if device.SupportsEvent(EV_ABS) {
			touchPads = append(touchPads, device)
		}
	}

	return touchPads, nil
}

// GetDevice reads and returns metadata for a specific input device.
func GetDevice(eventID int) (*Device, error) {
	basePath := fmt.Sprintf("/sys/class/input/event%d/device", eventID)

	// Read the device name
	data, err := os.ReadFile(filepath.Join(basePath, "name"))
	if err != nil {
		return nil, err
	}
	name := strings.TrimSpace(string(data))

	// Read capabilities
	evCapabilityMasks, err := readBitMasks(filepath.Join(basePath, "capabilities/ev"))
	if err != nil {
		return nil, err
	}
	if len(evCapabilityMasks) != 1 {
		return nil, fmt.Errorf("expected 1 bit mask for %s, got %d", filepath.Join(basePath, "capabilities/ev"), len(evCapabilityMasks))
	}

	keyCapabilityMasks, err := readBitMasks(filepath.Join(basePath, "capabilities/key"))
	if err != nil {
		return nil, err
	}

	capabilities := DeviceCapabilities{
		EventTypes: evCapabilityMasks[0],
		KeyTypes:   keyCapabilityMasks,
	}

	return &Device{
		ID:                 eventID,
		SysFsPath:          basePath,
		Path:               fmt.Sprintf("/dev/input/event%d", eventID),
		Name:               name,
		DeviceCapabilities: capabilities,
	}, nil
}

// ListDevices returns a list of Devices from /sys/class/input.
func ListDevices() ([]*Device, error) {
	var devices []*Device

	// Loop through event directories in /sys/class/input
	for i := 0; ; i++ {
		device, err := GetDevice(i)
		if err != nil {
			if !errors.Is(err, os.ErrNotExist) {
				return nil, fmt.Errorf("Failed to get input device: %v\n", err)
			}

			// Stop if we run out of devices (no such file/directory)
			break
		}
		devices = append(devices, device)
	}

	return devices, nil
}

// Read reads an event from the device.
func (d *Device) Read(e *Event) error {
	if d.handle == nil {
		return fmt.Errorf("device not open")
	}

	buffer := make([]byte, unsafe.Sizeof(*e))
	_, err := d.handle.Read(buffer)
	if err != nil {
		return err
	}

	if isLittleEndian() {
		err = binary.Read(bytes.NewBuffer(buffer), binary.LittleEndian, e)
	} else {
		err = binary.Read(bytes.NewBuffer(buffer), binary.BigEndian, e)
	}
	if err != nil {
		return err
	}

	return nil
}

// Open opens the device for reading.
func (d *Device) Open() error {
	if d.handle != nil {
		return errors.New("device already open")
	}

	handle, err := os.Open(d.Path)
	if err != nil {
		return err
	}

	d.handle = handle
	return nil
}

// Close closes the device.
func (d *Device) Close() error {
	if d.handle == nil {
		return errors.New("device not open")
	}
	d.handle.Close()
	d.handle = nil
	return nil
}

// ListenWithChannel listens for events on the device and writes them to the channel.
func (d *Device) ListenWithChannel(ctx context.Context, eventChan *chan Event) error {
	err := d.Open()
	if err != nil {
		return fmt.Errorf("Error opening device: %v\n", err)
	}

	go func() {
		defer d.Close()

		for {
			if ctx.Err() != nil {
				return
			}

			var event Event
			err := d.Read(&event)
			if err != nil {
				close(*eventChan)
				*eventChan = nil
				return
			}
			*eventChan <- event
		}
	}()

	return nil
}

// Listen listens for events on the device and calls the handler for each event. This function will block until the context is canceled.
func (d *Device) Listen(ctx context.Context, handler Handler) error {
	eventChan := make(chan Event)
	err := d.ListenWithChannel(ctx, &eventChan)
	if err != nil {
		return err
	}

	go func() {
		for {
			event := <-eventChan
			handler(event)
		}
	}()

	<-ctx.Done()
	return ctx.Err()
}

func isLittleEndian() bool {
	var i int32 = 0x01020304
	u := unsafe.Pointer(&i)
	pb := (*byte)(u)
	b := *pb
	return b == 0x04
}

func readBitMasks(path string) ([]*big.Int, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	parts := strings.Fields(string(data))

	segments := []*big.Int{}

	for _, part := range parts {
		segment := new(big.Int)
		_, ok := segment.SetString(part, 16)
		if !ok {
			return nil, fmt.Errorf("failed to parse hex segment: %s", part)
		}

		segments = append(segments, segment)
	}

	if isLittleEndian() {
		reversedSegments := []*big.Int{}
		for i := len(segments) - 1; i >= 0; i-- {
			reversedSegments = append(reversedSegments, segments[i])
		}
		segments = reversedSegments
	}

	return segments, nil
}
