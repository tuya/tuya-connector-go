package utils

import "os"

// Determine whether the given path file / folder exists
func Exists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		return os.IsExist(err)
	}
	return true
}

// Determine if the given path is a folder
func IsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

// Determine if the given path is a file
func IsFile(path string) bool {
	return !IsDir(path)
}

func Mkdir(path string) error {
	return os.Mkdir(path, 0777)
}
