package util

import (
	"io/fs"
)

func GetAllFiles(efs fs.FS) (files map[string]string, err error) {
	fileNames, err := GetAllFileNames(efs)
	if err != nil {
		return nil, err
	}

	files = make(map[string]string)
	for _, fileName := range fileNames {
		fileBody, err := fs.ReadFile(efs, fileName)
		if err != nil {
			return nil, err
		}
		files[fileName] = string(fileBody)
	}

	return files, nil
}

func GetAllFileNames(efs fs.FS) (fileNames []string, err error) {
	if err := fs.WalkDir(efs, ".", func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}
		fileNames = append(fileNames, path)
		return nil
	}); err != nil {
		return nil, err
	}

	return fileNames, nil
}
