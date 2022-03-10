package hash

import (
	"crypto/md5"
	"io"
	"os"
)

func MD5File(filename string) ([]byte, error) {

	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	hash := md5.New()
	if _, err = io.Copy(hash, file); err != nil {
		return nil, err
	}

	return hash.Sum(nil), nil
}
