package gofilepath

import (
	"io"
	"os"
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

func JointSmart(fallbackPathSeparators string, elem ...string) string {
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
		if b, _ := PathIsChildOf(pathToCheck, paths); b {
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
