package zipper

import (
	"archive/zip"
	"fmt"
	"io"
	"os"

	logger "github.com/dm1trypon/easy-logger"
)

/*
Create <Zipper> - init Zipper structure
	Returns <*Zipper>:
		1. Structure pointer
*/
func (z *Zipper) Create() *Zipper {
	z = &Zipper{
		lc: "ZIPPER",
	}

	return z
}

/*
ZipFiles <Zipper> - compresses one or many files into a single zip archive file
	Returns <error>:
		1. error
	Args:
		1. path <string> - filename is the output zip file's name
		2. files <[]string> - files is a list of files to add to the zip
*/
func (z *Zipper) ZipFiles(path string, files []string) error {
	newZipFile, err := os.Create(path)
	if err != nil {
		logger.Error(z.lc, fmt.Sprint("Error creating zip file: ", err.Error()))
		return err
	}
	defer newZipFile.Close()

	zipWriter := zip.NewWriter(newZipFile)
	defer zipWriter.Close()

	for _, file := range files {
		if err = z.addFileToZip(zipWriter, file); err != nil {
			continue
		}
	}
	return nil
}

/*
addFileToZip <Zipper> - adds files to archive
	Returns <error>:
		1. error
	Args:
		1. zipWriter <*zip.Writer> - writer of files to zip archive
		2. path <string> - path of destination
*/
func (z *Zipper) addFileToZip(zipWriter *zip.Writer, path string) error {
	fileToZip, err := os.Open(path)
	if err != nil {
		logger.Error(z.lc, fmt.Sprint("Error opening zip file: ", err.Error()))
		return err
	}
	defer fileToZip.Close()

	info, err := fileToZip.Stat()
	if err != nil {
		logger.Error(z.lc, fmt.Sprint("Error getting zip archive structure: ", err.Error()))
		return err
	}

	header, err := zip.FileInfoHeader(info)
	if err != nil {
		logger.Error(z.lc, fmt.Sprint("Error creating zip header file: ", err.Error()))
		return err
	}

	header.Name = fileToZip.Name()
	header.Method = zip.Deflate

	writer, err := zipWriter.CreateHeader(header)
	if err != nil {
		logger.Error(z.lc, fmt.Sprint("Error adding file to zip archive: ", err.Error()))
		return err
	}

	_, err = io.Copy(writer, fileToZip)
	if err != nil {
		logger.Error(z.lc, fmt.Sprint("Error copying zip archive data: ", err.Error()))
		return err
	}

	return nil
}
