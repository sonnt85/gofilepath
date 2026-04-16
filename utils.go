package gofilepath

import (
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func DirIsEmpty(name string) (bool, error) {
	f, err := os.Open(name)
	if err != nil {
		return false, err
	}
	defer f.Close()

	// read in ONLY one file
	_, err = f.Readdir(1)

	// and if the file is EOF... well, the dir is empty.
	if err == io.EOF {
		return true, nil
	}
	return false, err
}

var RunePathSeparators = []string{"/", "\\"}

func GetPathSeparator(inputpath string) string {
	for i := 0; i < len(RunePathSeparators); i++ {
		if strings.Contains(inputpath, RunePathSeparators[i]) {
			return RunePathSeparators[i]
		}
	}
	return ""
}

func HasEndPathSeparators(inputpath string) bool {
	for i := 0; i < len(RunePathSeparators); i++ {
		if strings.HasSuffix(inputpath, RunePathSeparators[i]) {
			return true
			// return RunePathSeparators[i]
		}
	}
	return false
}

func ConvertPathSeparators(frompath, reference string) string {
	for i := 0; i < len(RunePathSeparators); i++ {
		if strings.Contains(reference, RunePathSeparators[i]) {
			namePathSeparators := ""
			for j := 0; j < len(RunePathSeparators); j++ {
				if strings.Contains(frompath, RunePathSeparators[j]) {
					namePathSeparators = RunePathSeparators[j]
					break
				}
			}
			if len(namePathSeparators) != 0 && namePathSeparators != RunePathSeparators[i] {
				frompath = strings.ReplaceAll(frompath, namePathSeparators, RunePathSeparators[i])
			}
			return frompath
		}
	}
	// strings.Join()
	return frompath
}

// JoinSmart joins path elements using the detected path separator from existing elements.
func JoinSmart(fallbackPathSeparators string, elem ...string) string {
	retpath := ""
	firstPathSeparators := ""
	func() {
		for i := 0; i < len(elem); i++ {
			for j := 0; j < len(RunePathSeparators); j++ {
				firstPathSeparators = GetPathSeparator(elem[i])
				if len(firstPathSeparators) != 0 {
					return
				}
			}
		}
	}()
	if len(firstPathSeparators) == 0 {
		if fallbackPathSeparators != "" {
			firstPathSeparators = fallbackPathSeparators
		} else {
			firstPathSeparators = string(os.PathSeparator)
		}
	}
	retpath = Join(elem...)
	if string(os.PathSeparator) != firstPathSeparators {
		retpath = strings.ReplaceAll(retpath, string(os.PathSeparator), firstPathSeparators)
	}
	return retpath
}

func RelSmart(basepath, targpath string, fallbackPathSeparators ...string) (retpath string, err error) {
	basePathSeparators := GetPathSeparator(basepath)
	if len(basePathSeparators) == 0 {
		if len(fallbackPathSeparators) != 0 {
			basePathSeparators = fallbackPathSeparators[0]
		} else {
			basePathSeparators = string(os.PathSeparator)
		}
	}

	retpath, err = Rel(filepath.FromSlash(basepath), filepath.FromSlash(targpath))
	if len(basePathSeparators) != 0 && basePathSeparators != string(os.PathSeparator) {
		retpath = strings.ReplaceAll(retpath, string(os.PathSeparator), basePathSeparators)
	}
	return
}

func CountPathSeparator(namepath string) int {
	pathSeparators := GetPathSeparator(namepath)
	if len(pathSeparators) == 0 {
		return 0
	}
	return strings.Count(namepath, pathSeparators)
}

func PathIsUnixSocket(addr string) bool {
	fileInfo, err := os.Stat(addr)
	if err != nil {
		return false // Assume addr is a TCP socket address
	}
	return fileInfo.Mode()&os.ModeSocket != 0
}

func PathIsChildOf(path, parentDir string) (b bool, err error) {
	var absParentDir, absPath string
	if filepath.IsAbs(parentDir) {
		absParentDir = parentDir
	} else {
		if absParentDir, err = filepath.Abs(parentDir); err != nil {
			return false, err
		}
	}

	if filepath.IsAbs(path) {
		absPath = path
	} else {
		if absPath, err = filepath.Abs(path); err != nil {
			return false, err
		}
	}
	absPath = filepath.Clean(absPath)
	absParentDir = filepath.Clean(absParentDir)

	if absPath == absParentDir {
		return false, nil
	}
	return strings.HasPrefix(absPath, absParentDir), nil
}

// This function returns the first existing path in the given list of paths.
func FirstExistPath(path string) string {
	for _, v := range filepath.SplitList(path) {
		if _, err := os.Stat(v); err == nil {
			return v
		}
	}
	return ""
}

func GetPathInPaths(pathToCheck, paths string) string {
	for _, p := range strings.Split(paths, string(os.PathListSeparator)) {
		if b, _ := PathIsChildOf(pathToCheck, p); b {
			return p
		}
	}
	return ""
}

func PathHasSubpath(subpath, PATH string) bool {
	pashListSeparator := GetPathSeparator(PATH)
	for _, val := range strings.Split(PATH, string(pashListSeparator)) {
		if _, err := os.Stat(filepath.Join(val, subpath)); err == nil {
			return true
		}
	}
	return false
}

// --- Cross-platform path functions ---
// These functions handle both '/' and '\' as separators regardless of the
// current OS. Useful for manipulating remote paths where the CLI OS may
// differ from the target OS (e.g., Linux CLI → Windows agent or vice versa).
// All output uses '/' as separator, which is accepted by Go's os package
// on every platform including Windows.

// NormalizeSeparators replaces all backslashes with forward slashes.
// This is the foundation for cross-platform path manipulation:
// both "/" and "\" are treated as separators, output always uses "/".
func NormalizeSeparators(p string) string {
	return strings.ReplaceAll(p, "\\", "/")
}

// Deprecated: Use JoinSmart instead.
func JointSmart(fallbackPathSeparators string, elem ...string) string {
	return JoinSmart(fallbackPathSeparators, elem...)
}

// BaseSmart returns the last element of the path, handling both '/' and '\'
// as separators regardless of the current OS.
//
//	BaseSmart("C:\\Users\\file.txt")  → "file.txt"  (even on Linux)
//	BaseSmart("/home/user/file.txt")  → "file.txt"  (even on Windows)
//	BaseSmart("")                     → "."
func BaseSmart(p string) string {
	return path.Base(NormalizeSeparators(p))
}

// DirSmart returns all but the last element of the path, handling both
// '/' and '\' as separators. Output uses '/'.
//
//	DirSmart("C:\\Users\\file.txt")  → "C:/Users"
//	DirSmart("/home/user/file.txt")  → "/home/user"
func DirSmart(p string) string {
	return path.Dir(NormalizeSeparators(p))
}

// SplitSmart splits path into directory and file components, handling both
// '/' and '\' as separators. Output uses '/'.
func SplitSmart(p string) (dir, file string) {
	return path.Split(NormalizeSeparators(p))
}

// ExtSmart returns the file extension, handling both separators.
func ExtSmart(p string) string {
	return path.Ext(NormalizeSeparators(p))
}

// BaseNoExtSmart returns the filename without extension, handling both separators.
//
//	BaseNoExtSmart("C:\\Users\\archive.tar.gz") → "archive.tar"
func BaseNoExtSmart(p string) string {
	base := BaseSmart(p)
	return strings.TrimSuffix(base, ExtSmart(p))
}

// JoinSlash joins path elements using '/' separator, normalizing any '\' in
// the input elements. Suitable for constructing remote/network paths.
//
//	JoinSlash("/remote/dir", "sub\\folder", "file.txt") → "/remote/dir/sub/folder/file.txt"
//	JoinSlash("C:\\Users", "docs")                      → "C:/Users/docs"
func JoinSlash(elem ...string) string {
	normalized := make([]string, len(elem))
	for i, e := range elem {
		normalized[i] = NormalizeSeparators(e)
	}
	return path.Join(normalized...)
}

// RelSlash returns a relative path from basepath to targpath, handling both
// separators in input. Output uses '/'.
//
//	RelSlash("C:\\Users\\home", "C:\\Users\\home\\docs\\file.txt") → "docs/file.txt"
func RelSlash(basepath, targpath string) (string, error) {
	rel, err := filepath.Rel(
		filepath.FromSlash(NormalizeSeparators(basepath)),
		filepath.FromSlash(NormalizeSeparators(targpath)),
	)
	if err != nil {
		return "", err
	}
	return filepath.ToSlash(rel), nil
}

// CleanSmart cleans a path handling both separators. Output uses '/'.
func CleanSmart(p string) string {
	return path.Clean(NormalizeSeparators(p))
}

// IsAbsSmart checks if a path is absolute, handling both separator styles
// and Windows drive letters (e.g., "C:/", "C:\").
func IsAbsSmart(p string) bool {
	normalized := NormalizeSeparators(p)
	if path.IsAbs(normalized) {
		return true
	}
	// Check Windows drive letter: "C:/" or "C:"
	if len(normalized) >= 2 && normalized[1] == ':' {
		c := normalized[0]
		if (c >= 'A' && c <= 'Z') || (c >= 'a' && c <= 'z') {
			return true
		}
	}
	return false
}

func PathsPointToSameFile(path1, path2 string) (bool, error) {
	realPath1, err := filepath.EvalSymlinks(path1)
	if err != nil {
		return false, fmt.Errorf("error resolving symlink for %s: %w", path1, err)
	}

	realPath2, err := filepath.EvalSymlinks(path2)
	if err != nil {
		return false, fmt.Errorf("error resolving symlink for %s: %w", path2, err)
	}

	info1, err := os.Lstat(realPath1)
	if err != nil {
		return false, fmt.Errorf("error getting file info for %s: %w", realPath1, err)
	}

	info2, err := os.Lstat(realPath2)
	if err != nil {
		return false, fmt.Errorf("error getting file info for %s: %w", realPath2, err)
	}

	return os.SameFile(info1, info2), nil
}
