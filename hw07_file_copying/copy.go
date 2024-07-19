package main

import (
	"errors"
	"io"
	"os"

	"github.com/cheggaaa/pb/v3"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	if fromPath == "" || toPath == "" {
		return ErrUnsupportedFile
	}

	srcFile, err := os.Open(fromPath)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	// Полученим информацию о файле
	srcFileInfo, err := srcFile.Stat()
	if err != nil {
		return err
	}

	// Проверка на тип файла
	if srcFileInfo.Mode()&os.ModeType != 0 {
		return ErrUnsupportedFile
	}

	// Проверка, не превышает ли offset размер файла
	if offset > srcFileInfo.Size() {
		return ErrOffsetExceedsFileSize
	}

	// Установка смещения (offset)
	_, err = srcFile.Seek(offset, io.SeekStart)
	if err != nil {
		return err
	}

	dstFile, err := os.Create(toPath)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	bar := pb.Full.Start64(srcFileInfo.Size())
	defer bar.Finish()
	barReader := bar.NewProxyReader(srcFile)

	if limit > 0 {
		_, err = io.CopyN(dstFile, barReader, limit)
	} else {
		_, err = io.Copy(dstFile, barReader)
	}

	if err != nil && errors.Is(err, io.EOF) {
		return err
	}

	return nil
}
