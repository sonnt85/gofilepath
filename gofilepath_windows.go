package gofilepath

import (
	"golang.org/x/sys/windows"
)

func bitsToDrives(bitMap uint32) (drives []string) {
	availableDrives := []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z"}

	for i := range availableDrives {
		if bitMap&1 == 1 {
			drives = append(drives, availableDrives[i])
		}
		bitMap >>= 1
	}
	return
}

func getDrives() ([]string, error) {
	narkDrives, e := windows.GetLogicalDrives()
	if e != nil {
		return nil, e
	}
	return bitsToDrives(narkDrives), nil
}
