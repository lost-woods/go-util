package osutil

import (
	"os"
	"path/filepath"
	"strings"
)

func JoinPath(pathChunk ...string) string {
	return filepath.Join(pathChunk...)
}

func PathExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil || !os.IsNotExist(err)
}

func GetLastElementInPath(path string) string {
	return filepath.Base(path)
}

func GetRelativePath(basePath string, fullPath string) string {
	relativePath, err := filepath.Rel(basePath, fullPath)
	if err != nil {
		log.Fatalf("Failed to get relative path from path '%s' for file '%s': %v", basePath, fullPath, err)
	}

	return relativePath
}

func ReadFile(path string) string {
	fileBytes, err := os.ReadFile(path)
	if err != nil {
		log.Errorf("Error reading file: %s", err)
		return ""
	}

	return strings.TrimSuffix(string(fileBytes), "\n")
}

func WriteFile(path string, data []byte) {
	err := os.MkdirAll(filepath.Dir(path), os.ModePerm)
	if err != nil {
		log.Fatalf("Failed to create directories for file '%s': %v", path, err)
	}

	err = os.WriteFile(path, data, 0644)
	if err != nil {
		log.Fatalf("Failed to write file '%s': %v", path, err)
	}
}

func DeleteFile(path string) {
	err := os.Remove(path)
	if err != nil {
		log.Fatalf("Failed to delete file '%s': %v", path, err)
	}
}

func FindFiles(rootDir string, prefix string, suffix string) []string {
	var matches []string

	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		name := info.Name()
		if (prefix == "" || strings.HasPrefix(name, prefix)) &&
			(suffix == "" || strings.HasSuffix(name, suffix)) {
			matches = append(matches, path)
		}

		return nil
	})
	if err != nil {
		log.Fatalf("Failed to find files using path '%s' with prefix '%s' and suffix '%s': %v", rootDir, prefix, suffix, err)
	}

	return matches
}
