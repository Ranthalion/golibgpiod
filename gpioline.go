package golibgpiod

import (
	"log"
	"reflect"
	"unsafe"
)

type GpioLine struct {
	Offset      uint32
	Direction   int
	ActiveState int
	Used        bool
	OpenSource  bool
	OpenDrain   bool

	State    int
	UpToDate bool

	chip *GpioChip
	//struct gpiod_chip *chip;
	//struct line_fd_handle *fd_handle;

	Name     string
	Consumer string
}

/**
 * struct gpioline_info - Information about a certain GPIO line
 * @line_offset: the local offset on this GPIO device, fill this in when
 * requesting the line information from the kernel
 * @flags: various flags for this line
 * @name: the name of this GPIO line, such as the output pin of the line on the
 * chip, a rail or a pin header name on a board, as specified by the gpio
 * chip, may be NULL
 * @consumer: a functional name for the consumer of this GPIO line as set by
 * whatever is using it, will be NULL if there is no current user but may
 * also be NULL if the consumer doesn't set this up
 */
type GpioLineInfo struct {
	LineOffset uint32
	Flags      uint32
	Name       [32]byte
	Consumer   [32]byte
}

func newGpioLine(chip *GpioChip, lineNumber uint32) GpioLine {

	line := new(GpioLine)

	line.chip = chip
	line.Offset = lineNumber

	line.refresh()

	return *line
}

func (l *GpioLine) refresh() {

	var info GpioLineInfo
	info.LineOffset = l.Offset

	sizeptr := reflect.TypeOf(info).Size()

	GPIO_GET_LINEINFO_IOCTL := IoRW(0xB4, 0x02, sizeptr)

	err := Ioctl(
		uintptr(l.chip.FileHandle.Fd()),
		GPIO_GET_LINEINFO_IOCTL,
		uintptr(unsafe.Pointer(&info)),
	)

	if err != nil {
		log.Println("err = ", err)
		return
	}

	l.Direction = GPIOD_LINE_DIRECTION_INPUT
	if (info.Flags & GPIOLINE_FLAG_IS_OUT) != 0 {
		l.Direction = GPIOD_LINE_DIRECTION_OUTPUT
	}

	l.ActiveState = GPIOD_LINE_ACTIVE_STATE_HIGH
	if (info.Flags & GPIOLINE_FLAG_ACTIVE_LOW) != 0 {
		l.ActiveState = GPIOD_LINE_ACTIVE_STATE_LOW
	}

	l.Used = (info.Flags & GPIOLINE_FLAG_KERNEL) != 0
	l.OpenDrain = (info.Flags & GPIOLINE_FLAG_OPEN_DRAIN) != 0
	l.OpenSource = (info.Flags & GPIOLINE_FLAG_OPEN_SOURCE) != 0
	l.Name = string(info.Name[:])
	l.Consumer = string(info.Consumer[:])
	l.UpToDate = true
}
