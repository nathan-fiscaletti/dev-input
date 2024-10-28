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
	bitOffset := 0

	for _, field := range cap.KeyTypes {
		bitsInField := field.BitLen()

		if keyCode < bitOffset+bitsInField {
			return field.Bit(keyCode-bitOffset) == 1
		}

		bitOffset += bitsInField
	}

	return false
}

// SupportsEvent returns true if the device supports the given event.
func (cap DeviceCapabilities) SupportsEvent(eventType EventFlag) bool {
	et := big.NewInt(int64(eventType))
	return new(big.Int).And(cap.EventTypes, et).Cmp(et) == 0
}
