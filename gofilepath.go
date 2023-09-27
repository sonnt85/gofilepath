// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package filepath implements utility routines for manipulating filename paths
// in a way compatible with the target operating system-defined file paths.
//
// The filepath package uses either forward slashes or backslashes,
// depending on the operating system. To process paths such as URLs
// that always use forward slashes regardless of the operating
// system, see the path package.
package gofilepath

import (
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/sonnt85/gosutils/sregexp"
)

func Clean(path string) string {
	return filepath.Clean(filepath.FromSlash(path))
}

func _toSlashSmart(path string) string {
	path = sregexp.New("^/(.)\\*").ReplaceAllString(path, "${1}/*") // /C*
	path = sregexp.New("^/(.)/").ReplaceAllString(path, "${1}/")    // /C/Users -> C/Users
	path = sregexp.New("^(.)/").ReplaceAllString(path, "${1}:/")    //for C/ -> C:/
	return path
}

// same ToSlash but process path for windows:
// /C* => C/*
// /C/Users -> C/Users
// C:/ -> C/
func ToSlashSmart(path string, isFullPaths ...bool) string {
	isFullPath := false
	if len(isFullPaths) != 0 {
		isFullPath = isFullPaths[0]
	}
	if runtime.GOOS == "windows" && strings.Contains(path, "/") {
		if isFullPath || strings.HasPrefix(path, "/") {
			path = _toSlashSmart(path)
		}
		// path = sregexp.New("^(.):/").ReplaceAllString(path, "${1}/")     //for C:/ -> C/
		// return path
	}
	// return filepath.FromSlash(path)
	return filepath.ToSlash(path)
}

// ToSlash returns the result of replacing each separator character
// in path with a slash ('/') character. Multiple separators are
// replaced by multiple slashes.
func ToSlash(path string) string {
	return filepath.ToSlash(path)
}

func FromSlashSmart(path string, isFullPath bool) string {
	path = ToSlashSmart(path, isFullPath)
	return filepath.FromSlash(path)
}

// FromSlash returns the result of replacing each slash ('/') character
// in path with a separator character. Multiple slashes are replaced
// by multiple separators.
func FromSlash(path string) string {
	// if runtime.GOOS == "windows" && strings.Contains(path, "/") {
	// 	path = sregexp.New("^(.)/").ReplaceAllString(path, "${1}:/") //for C/ -> C:/
	// }
	return filepath.FromSlash(path)
}

// SplitList splits a list of paths joined by the OS-specific ListSeparator,
// usually found in PATH or GOPATH environment variables.
// Unlike strings.Split, SplitList returns an empty slice when passed an empty
// string.
func SplitList(path string) []string {
	return filepath.SplitList(filepath.FromSlash(path))
}

// Split splits path immediately following the final Separator,
// separating it into a directory and file name component.
// If there is no Separator in path, Split returns an empty dir
// and file set to path.
// The returned values have the property that path = dir+file.
func Split(path string) (dir, file string) {
	return filepath.Split(filepath.FromSlash(path))
}

// Join joins any number of path elements into a single path,
// separating them with an OS specific Separator. Empty elements
// are ignored. The result is Cleaned. However, if the argument
// list is empty or all its elements are empty, Join returns
// an empty string.
// On Windows, the result will only be a UNC path if the first
// non-empty element is a UNC path.
func Join(elem ...string) string {
	paths := []string{}
	for _, path := range elem {
		paths = append(paths, filepath.FromSlash(path))
	}
	return filepath.Join(paths...)
}

// Ext returns the file name extension used by path.
// The extension is the suffix beginning at the final dot
// in the final element of path; it is empty if there is
// no dot.
func Ext(path string) string {
	return filepath.Ext(filepath.FromSlash(path))
}

// EvalSymlinks returns the path name after the evaluation of any symbolic
// links.
// If path is relative the result will be relative to the current directory,
// unless one of the components is an absolute symbolic link.
// EvalSymlinks calls Clean on the result.
func EvalSymlinks(path string) (string, error) {
	return filepath.EvalSymlinks(filepath.FromSlash(path))
}

// Abs returns an absolute representation of path.
// If the path is not absolute it will be joined with the current
// working directory to turn it into an absolute path. The absolute
// path name for a given file is not guaranteed to be unique.
// Abs calls Clean on the result.
func Abs(path string) (string, error) {
	return filepath.Abs(filepath.FromSlash(path))
}

func IsAbs(path string) bool {
	return filepath.IsAbs(filepath.FromSlash(path))
}

// Rel returns a relative path that is lexically equivalent to targpath when
// joined to basepath with an intervening separator. That is,
// Join(basepath, Rel(basepath, targpath)) is equivalent to targpath itself.
// On success, the returned path will always be relative to basepath,
// even if basepath and targpath share no elements.
// An error is returned if targpath can't be made relative to basepath or if
// knowing the current working directory would be necessary to compute it.
// Rel calls Clean on the result.
func Rel(basepath, targpath string) (string, error) {
	return filepath.Rel(filepath.FromSlash(basepath), filepath.FromSlash(targpath))
}

// WalkFunc is the type of the function called by Walk to visit each
// file or directory.
//
// The path argument contains the argument to Walk as a prefix.
// That is, if Walk is called with root argument "dir" and finds a file
// named "a" in that directory, the walk function will be called with
// argument "dir/a".
//
// The directory and file are joined with Join, which may clean the
// directory name: if Walk is called with the root argument "x/../dir"
// and finds a file named "a" in that directory, the walk function will
// be called with argument "dir/a", not "x/../dir/a".
//
// The info argument is the fs.FileInfo for the named path.
//
// The error result returned by the function controls how Walk continues.
// If the function returns the special value SkipDir, Walk skips the
// current directory (path if info.IsDir() is true, otherwise path's
// parent directory). Otherwise, if the function returns a non-nil error,
// Walk stops entirely and returns that error.
//
// The err argument reports an error related to path, signaling that Walk
// will not walk into that directory. The function can decide how to
// handle that error; as described earlier, returning the error will
// cause Walk to stop walking the entire tree.
//
// Walk calls the function with a non-nil err argument in two cases.
//
// First, if an os.Lstat on the root directory or any directory or file
// in the tree fails, Walk calls the function with path set to that
// directory or file's path, info set to nil, and err set to the error
// from os.Lstat.
//
// Second, if a directory's Readdirnames method fails, Walk calls the
// function with path set to the directory's path, info, set to an
// fs.FileInfo describing the directory, and err set to the error from
// Readdirnames.
type WalkFunc filepath.WalkFunc

// type WalkFunc func(path string, info fs.FileInfo, err error) error

// WalkDir walks the file tree rooted at root, calling fn for each file or
// directory in the tree, including root.
//
// All errors that arise visiting files and directories are filtered by fn:
// see the fs.WalkDirFunc documentation for details.
//
// The files are walked in lexical order, which makes the output deterministic
// but requires WalkDir to read an entire directory into memory before proceeding
// to walk that directory.
//
// WalkDir does not follow symbolic links.
func WalkDir(root string, fn fs.WalkDirFunc) error {
	return filepath.WalkDir(filepath.FromSlash(root), fn)
}

// Walk walks the file tree rooted at root, calling fn for each file or
// directory in the tree, including root.
//
// All errors that arise visiting files and directories are filtered by fn:
// see the WalkFunc documentation for details.
//
// The files are walked in lexical order, which makes the output deterministic
// but requires Walk to read an entire directory into memory before proceeding
// to walk that directory.
//
// Walk does not follow symbolic links.
//
// Walk is less efficient than WalkDir, introduced in Go 1.16,
// which avoids calling os.Lstat on every visited file or directory.
func Walk(root string, fn WalkFunc) error {
	return filepath.Walk(filepath.FromSlash(root), any(fn).(filepath.WalkFunc))

}

// Base returns the last element of path.
// Trailing path separators are removed before extracting the last element.
// If the path is empty, Base returns ".".
// If the path consists entirely of separators, Base returns a single separator.
func Base(path string) string {
	return filepath.Base(filepath.FromSlash(path))
}

// Dir returns all but the last element of path, typically the path's directory.
// After dropping the final element, Dir calls Clean on the path and trailing
// slashes are removed.
// If the path is empty, Dir returns ".".
// If the path consists entirely of separators, Dir returns a single separator.
// The returned path does not end in a separator unless it is the root directory.
func Dir(path string) string {
	return filepath.Dir(filepath.FromSlash(path))
}

// VolumeName returns leading volume name.
// Given "C:\foo\bar" it returns "C:" on Windows.
// Given "\\host\share\foo" it returns "\\host\share".
// On other platforms it returns "".
func VolumeName(path string) string {
	return filepath.VolumeName(filepath.FromSlash(path))
}

// func Match(pattern, name string) (matched bool, err error) {
// 	name =filepath.FromSlash(name)
// 	return filepath.Match(pattern, name)
// }

type WarkdirFunc func(root string, fn fs.WalkDirFunc) error

type StatDirEntry struct {
	info fs.FileInfo
}

func (d *StatDirEntry) Name() string               { return d.info.Name() }
func (d *StatDirEntry) IsDir() bool                { return d.info.IsDir() }
func (d *StatDirEntry) Type() fs.FileMode          { return d.info.Mode().Type() }
func (d *StatDirEntry) Info() (fs.FileInfo, error) { return d.info, nil }

func isSymlinkToDir(path string, d fs.DirEntry) bool {
	if (d.Type() & os.ModeSymlink) == 0 {
		return false
	}

	realPath, err := filepath.EvalSymlinks(path)
	if err != nil {
		return false
	}

	info, err := os.Stat(realPath)
	if err != nil {
		return false
	}

	return info.IsDir()
}

func isSymlinkToFile(path string, d fs.DirEntry) bool {
	if (d.Type() & os.ModeSymlink) == 0 {
		return false
	}

	realPath, err := filepath.EvalSymlinks(path)
	if err != nil {
		return false
	}

	info, err := os.Stat(realPath)
	if err != nil {
		return false
	}

	return !info.IsDir()
}

func isSymlinkToDirectory(path string) (bool, error) {
	fileInfo, err := os.Lstat(path)
	if err != nil {
		return false, err
	}
	if fileInfo.Mode()&os.ModeSymlink == 0 {
		return false, nil
	}

	targetPath, err := os.Readlink(path)
	if err != nil {
		return false, err
	}

	targetFileInfo, err := os.Stat(targetPath)
	if err != nil {
		return false, err
	}

	return targetFileInfo.IsDir(), nil
}

// func FindFilesMatchPathFromRoot(root, pattern string, maxdeep int, matchfile, matchdir bool, matchFunc func(pattern, relpath string) bool, warkdirs ...WarkdirFunc) (matches []string) {
// FindFilesMatchPathFromRoot finds files and directories that match a specified pattern
// within a given root directory and up to a maximum depth (valid from 0).
//
// Parameters:
//   - root: The starting directory path from which the search begins.
//   - pattern: The pattern used for matching file and directory paths.
//   - maxdeep: The maximum depth (distance from the root) to search for matching files and directories.
//     A value of 0 means only the root directory will be considered.
//   - matchfile: Set to true to include matching files in the results.
//   - matchdir: Set to true to include matching directories in the results.
//   - matchFunc: A function that takes a pattern and a relative path as arguments and returns true if the path matches the pattern.
//   - warkdirs: Optional functions that can be applied to each directory encountered during the search.
//
// Returns:
//   - matches: A slice of strings containing the paths of the matching files and directories.
func FindFilesMatchPathFromRoot(root, pattern string, maxdeep int, matchfile, matchdir bool, matchFunc func(pattern, relpath string) bool, warkdirs ...WarkdirFunc) (matches []string) {

	matches = make([]string, 0)
	if matchFunc == nil {
		return
	}
	// if ok, _ := isSymlinkToDirectory(root); ok && !strings.HasSuffix(root, string(os.PathSeparator)) {
	// 	root += string(os.PathSeparator)
	// }
	rootPath := filepath.FromSlash(root)
	if finfo, err := os.Stat(root); err == nil {
		if !finfo.IsDir() { //is file
			if matchFunc(pattern, root) {
				matches = []string{root}
			}
			return
		}
	}
	pattern = filepath.FromSlash(pattern)
	warkdir := filepath.WalkDir
	if len(warkdirs) != 0 {
		warkdir = warkdirs[0]
	}
	var relpath string
	var deep int
	if nil != warkdir(rootPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil { //signaling that Walk will not walk into this directory.
			// return err
			return nil
		}

		relpath, err = filepath.Rel(rootPath, path)
		if err != nil {
			return nil
		}
		if maxdeep >= 0 {
			deep = strings.Count(relpath, string(os.PathSeparator))
			if deep > maxdeep {
				if d.IsDir() {
					return filepath.SkipDir
				} else {
					return nil
				}
			}
		}
		if (matchdir && (d.IsDir() || isSymlinkToDir(path, d))) || (matchfile && (!d.IsDir() || isSymlinkToFile(path, d))) {
			if matchFunc(pattern, relpath) {
				matches = append(matches, path)
			}
			if isSymlinkToDir(path, d) {
				newMaxdeep := maxdeep
				if maxdeep >= 1 {
					newMaxdeep = maxdeep - deep - 1
					if newMaxdeep < 0 {
						return nil
					}
				}

				newRoot := path + string(os.PathSeparator)
				if submatches := FindFilesMatchPathFromRoot(newRoot, pattern, newMaxdeep, matchfile, matchdir, matchFunc, warkdir); len(submatches) != 0 {
					matches = append(matches, submatches...)
				}
				return nil
			}
		}
		return nil
	}) {
		return nil
	}
	return matches
}

// FindFilesMatchRegexpPathFromRoot finds files and directories that match a regular expression pattern
// starting from the specified root directory and considering a maximum depth (valid from 0).
//
// Parameters:
//   - root: The starting directory path from which the search begins.
//   - pattern: The regular expression pattern used for matching file and directory names.
//   - maxdeep: The maximum depth (distance from the root) to search for matching files and directories.
//     A value of 0 means only the root directory will be considered.
//   - matchfile: Set to true to include matching files in the results.
//   - matchdir: Set to true to include matching directories in the results.
//   - warkdirs: Optional functions that can be applied to each directory encountered during the search.
//
// Returns:
//   - matches: A slice of strings containing the paths of the matching files and directories.
func FindFilesMatchRegexpPathFromRoot(root, pattern string, maxdeep int, matchfile, matchdir bool, warkdirs ...WarkdirFunc) (matches []string) {
	matchFunc := func(pattern, relpath string) bool {
		return sregexp.New(pattern).MatchString(relpath)
	}
	return FindFilesMatchPathFromRoot(root, pattern, maxdeep, matchfile, matchdir, matchFunc, warkdirs...)
}

// FindFilesMatchRegexpName finds files and directories that match a regular expression pattern
// by considering their names within a specified root directory and within a maximum depth (valid from 0).
//
// Parameters:
//   - root: The starting directory path from which the search begins.
//   - pattern: The regular expression pattern used for matching file and directory names.
//   - maxdeep: The maximum depth (distance from the root) to search for matching files and directories.
//     A value of 0 means only the root directory will be considered.
//   - matchfile: Set to true to include matching files in the results.
//   - matchdir: Set to true to include matching directories in the results.
//   - warkdirs: Optional functions that can be applied to each directory encountered during the search.
//
// Returns:
//   - matches: A slice of strings containing the paths of the matching files and directories.
func FindFilesMatchRegexpName(root, pattern string, maxdeep int, matchfile, matchdir bool, warkdirs ...WarkdirFunc) (matches []string) {
	matchFunc := func(pattern, relpath string) bool {
		return sregexp.New(pattern).MatchString(filepath.Base(relpath))
	}
	return FindFilesMatchPathFromRoot(root, pattern, maxdeep, matchfile, matchdir, matchFunc, warkdirs...)
}

// FindFilesMatchName finds files and directories whose names match the specified pattern
// within a given root directory and up to a maximum depth (valid from 0).
//
// Parameters:
//   - root: The starting directory path from which the search begins.
//   - pattern: The pattern used for matching file and directory names.
//   - maxdeep: The maximum depth (distance from the root) to search for matching files and directories.
//     A value of 0 means only the root directory will be considered.
//   - matchfile: Set to true to include matching files in the results.
//   - matchdir: Set to true to include matching directories in the results.
//   - warkdirs: Optional functions that can be applied to each directory encountered during the search.
//
// Returns:
//   - matches: A slice of strings containing the paths of the matching files and directories.
func FindFilesMatchName(root, pattern string, maxdeep int, matchfile, matchdir bool, warkdirs ...WarkdirFunc) (matches []string) {
	matchFunc := func(pattern, relpath string) bool {
		if match, err := filepath.Match(pattern, filepath.Base(relpath)); err == nil && match {
			return true
		}
		return false
	}
	return FindFilesMatchPathFromRoot(root, pattern, maxdeep, matchfile, matchdir, matchFunc, warkdirs...)
}

func GetDrives() ([]string, error) {
	return getDrives()
}

func TempFileCreateWithContent(data []byte, filename ...string) (fpath string) {
	if len(filename) != 0 && len(filename[0]) != 0 { //has name
		wdir, err := os.MkdirTemp("", "systempath")
		if err != nil {
			return
		}
		tmpfile := filepath.Join(wdir, filename[0])
		if f, err := os.Create(tmpfile); err == nil {
			if _, err = f.Write(data); err != nil {
				os.Remove(f.Name())
			} else {
				fpath = tmpfile
			}
			f.Close()
		}
	} else {
		if f, err := os.CreateTemp("", ""); err == nil {
			fpath = f.Name()
			f.Close()
		}
	}

	return
}

func Cat(files ...string) (contents string, err error) {
	var tmpbytes []byte
	for _, fname := range files {
		tmpbytes, err = os.ReadFile(fname)
		if err != nil {
			return
		}
		if len(contents) != 0 {
			contents += "\n" + string(tmpbytes)
		} else {
			contents += string(tmpbytes)
		}
	}
	return
}

func BaseNoExt(fpath string) string {
	return strings.TrimSuffix(Base(fpath), Ext(fpath))
}
