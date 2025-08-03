package input

import (
	"math/big"
)

// DeviceCapabilities holds information on the event and key capabilities of the device.
type DeviceCapabilities struct {
	EventTypes *big.Int
	KeyTypes   []*big.Int
}

// SupportsKey returns true if the device supports the given key code.
func (cap DeviceCapabilities) SupportsKey(keyCode int) bool {
	bitsPerField := 8 * len(cap.KeyTypes[0].Bytes())
	for i, field := range cap.KeyTypes {
		start := i * bitsPerField
		end := start + bitsPerField
		if keyCode >= start && keyCode < end {
			return field.Bit(keyCode-start) == 1
		}
	}
	return false
}

// SupportsEvent returns true if the device supports the given event.
func (cap DeviceCapabilities) SupportsEvent(eventType EventFlag) bool {
	return cap.EventTypes.Bit(int(eventType)) == 1
}
