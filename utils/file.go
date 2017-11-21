package utils

import (
	"errors"
	"os"
)

var (
	ErrFileHuge = errors.New("file is too huge")
)

func SaveFile(dat []byte, path string, limit int64) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	over, err := IsOverFileSize(file, limit)
	if err != nil {
		return err
	}
	if over {
		return ErrFileHuge
	}
	_, err = file.Write(dat)
	return err
}

func IsOverFileSize(f *os.File, limit int64) (bool, error) {
	info, err := f.Stat()
	if err != nil {
		return false, err
	}
	if info.Size() > limit {
		return true, nil
	}
	return false, nil
}
