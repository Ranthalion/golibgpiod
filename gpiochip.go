package golibgpiod

import (
	"log"
	"os"
	"reflect"
	"strings"
	"unsafe"
)

type GpioChip struct {
	Name       string
	Label      string
	NumLines   uint32
	FileHandle *os.File
	Lines      []GpioLine
}

type gpioChipInfo struct {
	name  [32]byte
	label [32]byte
	lines uint32
}

func GetGpioChips() []GpioChip {

	file, err := os.Open("/dev")
	if err != nil {
		log.Fatalf("failed opening directory: %s", err)
	}
	defer file.Close()

	gpioCheck := func(s string) bool {
		return strings.HasPrefix(s, "gpiochip")
	}

	list, _ := file.Readdirnames(0)
	chipNames := filter(list, gpioCheck)

	var chips []GpioChip

	for _, name := range chipNames {
		chips = append(chips, getChip(name))
	}

	return chips
}

func filter(dirs []string, test func(string) bool) (ret []string) {
	for _, dir := range dirs {
		if test(dir) {
			ret = append(ret, dir)
		}
	}
	return
}

func getChip(name string) (gpio GpioChip) {

	fd, err := os.OpenFile("/dev/"+name, os.O_RDWR, 0755)
	if err != nil {
		log.Fatal(err)
	}
	//defer fd.Close()

	//stat, _ := fd.Stat()

	if isGpioChipCdev("/dev/"+name) == false {
		log.Println(name, "is not cdev")
		return
	}

	var info gpioChipInfo
	sizeptr := reflect.TypeOf(info).Size()
	GPIO_GET_CHIPINFO_IOCTL := IoR(0xB4, 0x01, sizeptr)

	err = Ioctl(
		uintptr(fd.Fd()),
		GPIO_GET_CHIPINFO_IOCTL,
		uintptr(unsafe.Pointer(&info)),
	)

	if err != nil {
		log.Println("err = ", err)
		return
	}

	gpio.Name = string(info.name[:])
	gpio.FileHandle = fd
	gpio.NumLines = info.lines
	gpio.Label = string(info.label[:])

	//The kernel sets the label of a GPIO device to "unknown" if it
	//hasn't been defined in DT, board file etc. On the off-chance that
	//we got an empty string, do the same.
	if gpio.Label == "" {
		gpio.Label = "unknown"
	}

	return
}

func isGpioChipCdev(path string) (ret bool) {
	ret = false

	stat, err := os.Lstat(path)
	if err != nil {
		log.Fatal("unable to lstat", err)
		return
	}

	/* Is it a character device? */
	if stat.Mode()&os.ModeCharDevice == 0 {
		log.Fatal("Not a character device", stat.Mode())
		return
	}

	ret = true
	return
}

func (c *GpioChip) GetLines() {

	c.Lines = make([]GpioLine, c.NumLines)
	for i := uint32(0); i < c.NumLines; i++ {
		c.Lines[i] = newGpioLine(c, i)
	}
}

func (c *GpioChip) Close() {
	c.FileHandle.Close()
}
