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

func TestNormalizeSeparators(t *testing.T) {
	tests := []struct {
		input, want string
	}{
		{`C:\Users\file.txt`, "C:/Users/file.txt"},
		{"/home/user/file.txt", "/home/user/file.txt"},
		{`C:\mixed/path\file.txt`, "C:/mixed/path/file.txt"},
		{"", ""},
		{"nopath", "nopath"},
		{`\\server\share\file`, "//server/share/file"},
	}
	for _, tt := range tests {
		got := NormalizeSeparators(tt.input)
		if got != tt.want {
			t.Errorf("NormalizeSeparators(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestBaseSmart(t *testing.T) {
	tests := []struct {
		input, want string
	}{
		// Windows paths on any OS
		{`C:\Users\file.txt`, "file.txt"},
		{`C:\Users\docs\`, "docs"},
		{`C:\file.txt`, "file.txt"},
		// Unix paths
		{"/home/user/file.txt", "file.txt"},
		{"/home/user/docs/", "docs"},
		// Mixed
		{`C:\mixed/path\file.txt`, "file.txt"},
		// Edge cases
		{"file.txt", "file.txt"},
		{"", "."},
		{"/", "/"},
		{`\`, "/"},
		{`C:\`, "C:"},
	}
	for _, tt := range tests {
		got := BaseSmart(tt.input)
		if got != tt.want {
			t.Errorf("BaseSmart(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestDirSmart(t *testing.T) {
	tests := []struct {
		input, want string
	}{
		{`C:\Users\file.txt`, "C:/Users"},
		{"/home/user/file.txt", "/home/user"},
		{`C:\Users\docs\report.pdf`, "C:/Users/docs"},
		{"file.txt", "."},
		{"/file.txt", "/"},
	}
	for _, tt := range tests {
		got := DirSmart(tt.input)
		if got != tt.want {
			t.Errorf("DirSmart(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestJoinSlash(t *testing.T) {
	tests := []struct {
		elems []string
		want  string
	}{
		{[]string{"/remote/dir", "file.txt"}, "/remote/dir/file.txt"},
		{[]string{`C:\Users`, "docs", "file.txt"}, "C:/Users/docs/file.txt"},
		{[]string{"/remote", `sub\folder`, "file.txt"}, "/remote/sub/folder/file.txt"},
		{[]string{`C:\Users\home`, `sub\dir`}, "C:/Users/home/sub/dir"},
		{[]string{"/a/", "/b/"}, "/a/b"},
		{[]string{"", "file.txt"}, "file.txt"},
	}
	for _, tt := range tests {
		got := JoinSlash(tt.elems...)
		if got != tt.want {
			t.Errorf("JoinSlash(%v) = %q, want %q", tt.elems, got, tt.want)
		}
	}
}

func TestRelSlash(t *testing.T) {
	tests := []struct {
		base, targ, want string
	}{
		{`C:\Users\home`, `C:\Users\home\docs\file.txt`, "docs/file.txt"},
		{"/home/user", "/home/user/docs/file.txt", "docs/file.txt"},
		{`C:\a\b`, `C:\a\b\c\d`, "c/d"},
		// Mixed separators
		{`C:\Users\home`, "C:/Users/home/docs/file.txt", "docs/file.txt"},
	}
	for _, tt := range tests {
		got, err := RelSlash(tt.base, tt.targ)
		if err != nil {
			t.Errorf("RelSlash(%q, %q) error: %v", tt.base, tt.targ, err)
			continue
		}
		if got != tt.want {
			t.Errorf("RelSlash(%q, %q) = %q, want %q", tt.base, tt.targ, got, tt.want)
		}
	}
}

func TestExtSmart(t *testing.T) {
	tests := []struct {
		input, want string
	}{
		{`C:\Users\file.txt`, ".txt"},
		{`C:\Users\archive.tar.gz`, ".gz"},
		{"/home/user/noext", ""},
		{`C:\Users\.hidden`, ".hidden"},
	}
	for _, tt := range tests {
		got := ExtSmart(tt.input)
		if got != tt.want {
			t.Errorf("ExtSmart(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestBaseNoExtSmart(t *testing.T) {
	tests := []struct {
		input, want string
	}{
		{`C:\Users\file.txt`, "file"},
		{`C:\Users\archive.tar.gz`, "archive.tar"},
		{"/home/user/noext", "noext"},
	}
	for _, tt := range tests {
		got := BaseNoExtSmart(tt.input)
		if got != tt.want {
			t.Errorf("BaseNoExtSmart(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestIsAbsSmart(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{"/home/user", true},
		{`C:\Users`, true},
		{"C:/Users", true},
		{`D:\`, true},
		{"relative/path", false},
		{"file.txt", false},
		{"", false},
	}
	for _, tt := range tests {
		got := IsAbsSmart(tt.input)
		if got != tt.want {
			t.Errorf("IsAbsSmart(%q) = %v, want %v", tt.input, got, tt.want)
		}
	}
}

func TestCleanSmart(t *testing.T) {
	tests := []struct {
		input, want string
	}{
		{`C:\Users\..\Users\file.txt`, "C:/Users/file.txt"},
		{"/home//user/./file.txt", "/home/user/file.txt"},
		{`C:\a\b\..\c`, "C:/a/c"},
	}
	for _, tt := range tests {
		got := CleanSmart(tt.input)
		if got != tt.want {
			t.Errorf("CleanSmart(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

// TestCrossPlatformScenario simulates CLI on Linux manipulating Windows remote paths
func TestCrossPlatformScenario(t *testing.T) {
	// Scenario: Linux CLI uploads directory to Windows agent
	// localPath = "/home/user/mydir" (Linux)
	// remotePath = "C:\Users\Public\uploads" (Windows)

	remotePath := `C:\Users\Public\uploads`
	localBasename := "mydir" // filepath.Base on local path works fine

	// Build remote path for COPY_DIR mode
	remoteDir := JoinSlash(remotePath, localBasename)
	if remoteDir != "C:/Users/Public/uploads/mydir" {
		t.Errorf("JoinSlash remote dir = %q, want %q", remoteDir, "C:/Users/Public/uploads/mydir")
	}

	// Build remote path for individual files in directory
	relPath := "sub/file.txt" // from filepath.Rel on local, already forward slash on Linux
	remoteFile := JoinSlash(remoteDir, relPath)
	if remoteFile != "C:/Users/Public/uploads/mydir/sub/file.txt" {
		t.Errorf("JoinSlash remote file = %q, want %q", remoteFile, "C:/Users/Public/uploads/mydir/sub/file.txt")
	}

	// Scenario: Linux CLI downloads from Windows agent
	// Extract basename from Windows remote path for local dir name
	windowsRemote := `C:\Users\Public\mydir`
	basename := BaseSmart(windowsRemote)
	if basename != "mydir" {
		t.Errorf("BaseSmart(%q) = %q, want %q", windowsRemote, basename, "mydir")
	}

	// Scenario: Windows CLI uploads to Linux agent
	// relPath from filepath.Rel on Windows = "sub\file.txt"
	winRelPath := `sub\file.txt`
	linuxRemote := "/opt/data"
	remoteDest := JoinSlash(linuxRemote, winRelPath)
	if remoteDest != "/opt/data/sub/file.txt" {
		t.Errorf("JoinSlash with Windows relPath = %q, want %q", remoteDest, "/opt/data/sub/file.txt")
	}
}
