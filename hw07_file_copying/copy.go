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

func createProgressBar(srcFile *os.File, offset, limit int64) *pb.ProgressBar {
	srcFileInfo, _ := srcFile.Stat()
	// Вычисление ширины полосы прогресса
	var maxBarWidth int64
	if limit > 0 {
		if offset+limit > srcFileInfo.Size() {
			maxBarWidth = srcFileInfo.Size() - offset
		} else {
			maxBarWidth = limit
		}
	} else {
		maxBarWidth = srcFileInfo.Size() - offset
	}

	// Создание и запуск полосы прогресса
	return pb.Full.Start64(maxBarWidth)
}

func handleSrcFileErrors(srcFile *os.File, offset int64) error {
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
	return nil
}

func Copy(fromPath, toPath string, offset, limit int64) error {
	if fromPath == "" || toPath == "" {
		return ErrUnsupportedFile
	}

	srcFile, err := os.Open(fromPath)
	if err != nil {
		return err
	}

	if err := handleSrcFileErrors(srcFile, offset); err != nil {
		return err
	}

	// Установка смещения (offset)
	_, err = srcFile.Seek(offset, io.SeekStart)
	if err != nil {
		return err
	}

	// Получаем текущий рабочий каталог
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	tmpDstFile, err := os.CreateTemp(cwd, "dstFileTmp")
	if err != nil {
		return err
	}
	defer tmpDstFile.Close()

	bar := createProgressBar(srcFile, offset, limit)
	defer bar.Finish()
	barReader := bar.NewProxyReader(srcFile)

	// Копирование данных
	srcFileInfo, _ := srcFile.Stat()
	if limit > 0 && offset+limit <= srcFileInfo.Size() {
		_, err = io.CopyN(tmpDstFile, barReader, limit)
	} else {
		_, err = io.Copy(tmpDstFile, barReader)
	}

	if err != nil && !errors.Is(err, io.EOF) {
		return err
	}

	srcFile.Close()
	err = os.Rename(tmpDstFile.Name(), toPath)
	if err != nil {
		return err
	}
	return nil
}
