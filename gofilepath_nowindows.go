//go:build !windows
// +build !windows

package gofilepath

func getDrives() ([]string, error) {
	return []string{"/"}, nil
}
