package utils

import (
	"mime/multipart"

	"github.com/pkg/errors"
)

func FileHeaderToBytes(fileHeader *multipart.FileHeader) ([]byte, error) {
	file, err := fileHeader.Open()
	if err != nil {
		return nil, errors.Wrapf(err, "FileHeaderToBytes failed, headerSize=%d, name=%s", fileHeader.Size, fileHeader.Filename)
	}
	defer file.Close()

	data := make([]byte, fileHeader.Size)
	_, err = file.Read(data)
	if err != nil {
		return nil, errors.Wrapf(err, "FileHeaderToBytes failed, headerSize=%d, name=%s", fileHeader.Size, fileHeader.Filename)
	}

	return data, nil
}
