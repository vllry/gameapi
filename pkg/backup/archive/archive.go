package archive

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)

func ArchiveFilesToBuffer(sourceDirectory string) (bytes.Buffer, error) {
	var tarBuf bytes.Buffer
	tw := tar.NewWriter(&tarBuf)
	var gzipBuf bytes.Buffer
	gzipWriter := gzip.NewWriter(&gzipBuf)

	walkErr := filepath.Walk(
		sourceDirectory,
		func(path string, info os.FileInfo, err error) error {
			if err != nil { // Why is this  check necessary...?
				return err
			}

			header, err := tar.FileInfoHeader(info, path)
			if err != nil {
				return err
			}

			header.Name = filepath.Join(sourceDirectory, strings.TrimPrefix(path, sourceDirectory))
			if err := tw.WriteHeader(header); err != nil {
				return err
			}

			if !info.Mode().IsRegular() { //nothing more to do for non-regular
				return nil
			}

			f, openErr := os.Open(path)
			if openErr != nil {
				return errors.Wrap(openErr, "couldn't open file")
			}
			defer f.Close()

			_, writeErr := io.Copy(tw, f)
			return errors.Wrap(writeErr, "couldn't write file")
		},
	)
	if walkErr != nil {
		return bytes.Buffer{}, errors.Wrap(walkErr, "couldn't add files to ArchiveFilesToFile")
	}

	if err := tw.Close(); err != nil {
		return bytes.Buffer{}, err
	}

	gzipWriter.Write(tarBuf.Bytes())
	return tarBuf, nil
}

func ArchiveFilesToFile(sourceDirectory string, outputFile string) error {
	aFile, err := os.Create(outputFile)
	if err != nil {
		return errors.Wrap(err, "couldn't create archive file")
	}
	defer aFile.Close()

	buf, err := ArchiveFilesToBuffer(sourceDirectory)
	if err != nil {
		return err
	}

	_, err = buf.WriteTo(aFile)
	return err
}

/*

func ArchiveFilesToFile(sourceDirectory string, outputFile string) error {
	aFile, err := os.Create(outputFile)
	if err != nil {
		return errors.Wrap(err, "couldn't create archive file")
	}
	defer aFile.Close()

	var buf bytes.Buffer
	tw := tar.NewWriter(&buf)
	gzipWriter := gzip.NewWriter(aFile)

	walkErr := filepath.Walk(
		sourceDirectory,
		func(path string, info os.FileInfo, err error) error {
			if err != nil { // Why is this  check necessary...?
				return err
			}

			header, err := tar.FileInfoHeader(info, path)
			if err != nil {
				return err
			}

			header.Name = filepath.Join(sourceDirectory, strings.TrimPrefix(path, sourceDirectory))
			if err := tw.WriteHeader(header); err != nil {
				return err
			}

			if !info.Mode().IsRegular() { //nothing more to do for non-regular
				return nil
			}

			f, openErr := os.Open(path)
			if openErr != nil {
				return errors.Wrap(openErr, "couldn't open file")
			}
			defer f.Close()

			_, writeErr := io.Copy(tw, f)
			return errors.Wrap(writeErr, "couldn't write file")
		},
	)
	if walkErr != nil {
		return errors.Wrap(walkErr, "couldn't add files to ArchiveFilesToFile")
	}

	if err := tw.Close(); err != nil {
		return err
	}

	gzipWriter.Write(buf.Bytes())

	err = gzipWriter.Close()
	if err != nil {
		return errors.Wrap(err, "couldn't flush archive writes")
	}
	return nil
}
*/

// UnArchiveToBytes takes the path of a .tar.gz file, and returns the contents.
func UnArchiveToBytes(filePath string) ([]byte, error) {
	archiveF, err := os.Open(filePath)
	if err != nil {
		return []byte{}, errors.Wrap(err, "failed to open archive")
	}
	defer archiveF.Close()

	gzipReader, err := gzip.NewReader(archiveF)
	if err != nil {
		return []byte{}, nil
	}
	defer gzipReader.Close()

	b, err := ioutil.ReadAll(gzipReader)
	if err != nil {
		return []byte{}, err
	}
	return b, nil
}
