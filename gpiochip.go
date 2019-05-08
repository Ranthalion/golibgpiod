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
	Lines      **GpioLine
}

type GpioLine struct {
	Offset      int
	Direction   int
	ActiveState int
	Used        bool
	OpenSource  bool
	OpenDrain   bool

	State    int
	UpToDate bool

	//struct gpiod_chip *chip;
	//struct line_fd_handle *fd_handle;

	Name     string
	consumer string
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

type gpioChipInfo struct {
	name  [32]byte
	label [32]byte
	lines uint32
}

func getChip(name string) (gpio GpioChip) {

	fd, err := os.OpenFile("/dev/"+name, os.O_RDWR, 0755)
	if err != nil {
		log.Fatal(err)
	}
	defer fd.Close()

	stat, _ := fd.Stat()
	log.Println("stat:", stat)
	log.Println("name:", stat.Name())

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

	// /* Get the basename. */
	// pathcpy = strdup(path);
	// if (!pathcpy)
	// 	goto out;

	// name = basename(pathcpy);

	// /* Do we have a corresponding sysfs attribute? */
	// rv = asprintf(&sysfsp, "/sys/bus/gpio/devices/%s/dev", name);
	// if (rv < 0)
	// 	goto out_free_pathcpy;

	// if (access(sysfsp, R_OK) != 0) {
	// 	/*
	// 	 * This is a character device but not the one we're after.
	// 	 * Before the introduction of this function, we'd fail with
	// 	 * ENOTTY on the first GPIO ioctl() call for this file
	// 	 * descriptor. Let's stay compatible here and keep returning
	// 	 * the same error code.
	// 	 */
	// 	errno = ENOTTY;
	// 	goto out_free_sysfsp;
	// }

	// /*
	//  * Make sure the major and minor numbers of the character device
	//  * correspond with the ones in the dev attribute in sysfs.
	//  */
	// snprintf(devstr, sizeof(devstr), "%u:%u",
	// 	 major(statbuf.st_rdev), minor(statbuf.st_rdev));

	// fd = open(sysfsp, O_RDONLY);
	// if (fd < 0)
	// 	goto out_free_sysfsp;

	// memset(sysfsdev, 0, sizeof(sysfsdev));
	// rd = read(fd, sysfsdev, strlen(devstr));
	// close(fd);
	// if (rd < 0)
	// 	goto out_free_sysfsp;

	// if (strcmp(sysfsdev, devstr) != 0) {
	// 	errno = ENODEV;
	// 	goto out_free_sysfsp;
	// }

	ret = true
	return
}
