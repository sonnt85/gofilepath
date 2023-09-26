package gofilepath

import (
	"fmt"
	"path/filepath"
	"testing"
)

func TestFromSlashSmart(t *testing.T) {
	fmt.Println(filepath.ToSlash(_toSlashSmart("C/Windows/System32/AcGenral.dll")))
}

// "C:/Windows/System32/AcGenral.dll"
