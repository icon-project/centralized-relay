package interchaintest

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/icon-project/centralized-relay/test/chains"
)

var ibcConfigPath = filepath.Join(os.Getenv(chains.BASE_PATH), "ibc-config")

func CleanBackupConfig() {
	files, err := filepath.Glob(filepath.Join(ibcConfigPath, "*.json"))
	if err != nil {
		fmt.Println("Error deleting file:", err)
		return
	}

	for _, file := range files {
		err := os.Remove(file)
		if err != nil {
			fmt.Println("Error deleting file:", err)
		}
	}

}

// for saving data in particular format
func BackupConfig(chain chains.Chain) error {
	return nil
}

func GetLocalFileContent(fileName string) ([]byte, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, fmt.Errorf("%s file not found : %w", fileName, err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			fmt.Println("Error closing file:", err)
		}
	}(file)

	fileInfo, err := file.Stat()
	if err != nil {
		return nil, fmt.Errorf("unable to get file: %w", err)
	}
	fileSize := fileInfo.Size()

	// Read the file content into a buffer
	buffer := make([]byte, fileSize)
	_, err = file.Read(buffer)
	if err != nil {
		return nil, fmt.Errorf("unable to get content: %w", err)
	}
	return buffer, nil
}

func RestoreConfig(chain chains.Chain) error {
	return nil
}
