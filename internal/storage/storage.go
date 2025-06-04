package storage

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
)

var storage *Storage

func init() {
	currentUser, err := user.Current()
	if err != nil {
		panic(fmt.Errorf("failed to get current user: %v", err))
	}

	path := filepath.Join(currentUser.HomeDir, ".letstry")

	err = os.MkdirAll(path, 0755)
	if err != nil {
		panic(fmt.Errorf("failed to create storage directory: %v", err))
	}

	storage = &Storage{dir: path}
}

type Storage struct {
	dir string
}

func GetStorage() *Storage {
	return storage
}

func (s *Storage) GetAbsolutePath(value string) string {
	return filepath.Join(s.dir, value)
}

func (s *Storage) OpenFile(name string) (*os.File, error) {
	return s.OpenFileWithDefaultContent(name, []byte{})
}

func (s *Storage) Exists(name string) bool {
	_, err := os.Stat(filepath.Join(s.dir, name))
	return !os.IsNotExist(err)
}

func (s *Storage) OpenFileWithDefaultContent(name string, defaultContent []byte) (*os.File, error) {
	filePath := filepath.Join(s.dir, name)

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		file, err := os.Create(filePath)
		if err != nil {
			return nil, err
		}

		_, err = file.Write(defaultContent)
		if err != nil {
			return nil, err
		}

		_, err = file.Seek(0, 0)
		if err != nil {
			return nil, err
		}

		return file, nil
	}

	return os.OpenFile(filePath, os.O_RDWR, 0644)
}

func (s *Storage) DirectoryExists(name string) bool {
	stat, err := os.Stat(filepath.Join(s.dir, name))
	if err != nil {
		return false
	}

	return stat.IsDir()
}

func (s *Storage) CreateDirectory(name string) error {
	return os.MkdirAll(filepath.Join(s.dir, name), 0755)
}

func (s *Storage) DeleteDirectory(name string) error {
	return os.RemoveAll(filepath.Join(s.dir, name))
}

func (s *Storage) ListDirectories(path string) ([]string, error) {
	// List all directories within the directory specified by the path

	dir, err := os.Open(filepath.Join(s.dir, path))
	if err != nil {
		return []string{}, err
	}

	defer dir.Close()

	files, err := dir.Readdir(-1)
	if err != nil {
		return []string{}, err
	}

	var dirs []string
	for _, file := range files {
		if file.IsDir() {
			dirs = append(dirs, file.Name())
		}
	}

	return dirs, nil
}
