package file

import (
	"io"
	"os"
	"time"
)

func GetFileName(template string) string {
	now := time.Now()
	return now.Format(template)
}

func CopyFile(filePath string, fileContent io.ReadCloser) error {
	out, err := os.Create(filePath)
	if err != nil {
		return err
	}

	defer out.Close()

	_, err = io.Copy(out, fileContent)

	return err
}
