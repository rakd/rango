package rangolib

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
)

type File struct {
	Name    string `json:"name"`
	Path    string `json:"path"`
	IsDir   bool   `json:"isDir"`
	Size    int64  `json:"size"`
	ModTime int64  `json:"mtime"`
}

func (f *File) Load(info os.FileInfo) {
	f.Name = info.Name()
	f.IsDir = info.IsDir()
	f.Size = info.Size()
	f.ModTime = info.ModTime().Unix()
}

// NewFile constructs a new File based on a path and file info
func NewFile(path string, info os.FileInfo) *File {
	file := &File{Path: path}
	file.Load(info)
	return file
}

// ReadDir lists the contents of a directory
func ReadDir(dirname string) ([]*File, error) {
	files := make([]*File, 0)
	contents, err := ioutil.ReadDir(dirname)
	if err != nil {
		return nil, err
	}

	for _, info := range contents {
		file := NewFile(filepath.Join(dirname, info.Name()), info)
		files = append(files, file)
	}

	return files, nil
}

// CreateDir creates a new directory
func CreateDir(dirname string) (*File, error) {

	// make directory
	if err := os.MkdirAll(dirname, 0755); err != nil {
		return nil, err
	}

	// check that directory was created
	info, err := os.Stat(dirname)
	if err != nil {
		return nil, err
	}

	// convert fileinfo into something we can print
	return NewFile(dirname, info), nil
}

// UpdateDir renames an existing directory
func UpdateDir(src string, dest string) (*File, error) {

	// check that destination doesn't exist
	info, err := os.Stat(dest)
	if info != nil {
		return nil, errors.New("Cannot overwrite destination")
	}

	// move directory including it's contents
	if err := moveDir(src, dest); err != nil {
		return nil, err
	}

	// check that directory was created
	info, err = os.Stat(dest)
	if err != nil {
		return nil, err
	}

	// convert fileinfo into something we can print
	return NewFile(dest, info), nil
}

// DeleteDir will delete a directory and it's contents
func DeleteDir(dirname string) error {

	dir, err := os.Stat(dirname)
	if err != nil {
		return err
	}

	if dir.IsDir() == false {
		return errors.New("DeleteDir can only delete directories")
	}

	// remove the directory
	return os.RemoveAll(dirname)
}