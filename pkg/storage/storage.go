package storage

import (
	"os"
	"path/filepath"
)

type Storage struct {
	HomeDir    string
	ProjectDir string
}

func NewStorage() (*Storage, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	taoHome := filepath.Join(home, ".tao")
	if err := os.MkdirAll(taoHome, 0755); err != nil {
		return nil, err
	}

	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	return &Storage{
		HomeDir:    taoHome,
		ProjectDir: cwd,
	}, nil
}

func (s *Storage) LogLesson(dir, lesson string) error {
	path := filepath.Join(s.ProjectDir, dir, "README.md")
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString("\n## Lesson Learned\n" + lesson + "\n")
	return err
}
