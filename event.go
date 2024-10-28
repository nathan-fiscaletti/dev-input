package input

type Handler func(Event)

type Event struct {
	Time  [2]uint64 // Timestamp (seconds and microseconds)
	Type  EventType // Event type
	Code  uint16    // Event code
	Value int32     // Event value (1 for key down, 0 for key up)
}

// EventFlag is a flag that is used to determine if a particular device supports a particular event.
type EventFlag int

const (
	EV_SYN       EventFlag = 0x0001
	EV_KEY       EventFlag = 0x0002
	EV_REL       EventFlag = 0x0004
	EV_ABS       EventFlag = 0x0008
	EV_MSC       EventFlag = 0x0010
	EV_SW        EventFlag = 0x0020
	EV_LED       EventFlag = 0x0100
	EV_SND       EventFlag = 0x0200
	EV_REP       EventFlag = 0x0400
	EV_FF        EventFlag = 0x0800
	EV_PWR       EventFlag = 0x1000
	EV_FF_STATUS EventFlag = 0x2000
)

// EventType is a type of event that is sent by the device.
type EventType uint16

const (
	EV_TYPE_SYN       EventType = 0x00
	EV_TYPE_KEY       EventType = 0x01
	EV_TYPE_REL       EventType = 0x02
	EV_TYPE_ABS       EventType = 0x03
	EV_TYPE_MSC       EventType = 0x04
	EV_TYPE_SW        EventType = 0x05
	EV_TYPE_LED       EventType = 0x11
	EV_TYPE_SND       EventType = 0x12
	EV_TYPE_REP       EventType = 0x14
	EV_TYPE_FF        EventType = 0x15
	EV_TYPE_PWR       EventType = 0x16
	EV_TYPE_FF_STATUS EventType = 0x17
)
