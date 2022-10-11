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

var RunePathSeparators = []string{"/", `\`}

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

func ConvertPathSeparators(fompath, reference string) string {
	for i := 0; i < len(RunePathSeparators); i++ {
		if strings.Contains(reference, RunePathSeparators[i]) {
			namePathSeparators := ""
			for j := 0; j < len(RunePathSeparators); j++ {
				if strings.Contains(fompath, RunePathSeparators[j]) {
					namePathSeparators = RunePathSeparators[j]
					break
				}
			}
			if len(namePathSeparators) != 0 && namePathSeparators != RunePathSeparators[i] {
				fompath = strings.ReplaceAll(fompath, namePathSeparators, RunePathSeparators[i])
			}
			return fompath
		}
	}
	// strings.Join()
	return fompath
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
