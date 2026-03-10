package utils

import (
	"mime/multipart"
)

func FileHeaderToBytes(fileHeader *multipart.FileHeader) ([]byte, error) {
	file, err := fileHeader.Open()
	if err != nil {
		return nil, err
	}
	defer file.Close()

	data := make([]byte, fileHeader.Size)
	_, err = file.Read(data)
	if err != nil {
		return nil, err
	}

	return data, nil
}
