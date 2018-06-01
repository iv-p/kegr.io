package util

import (
	"bytes"
	"compress/gzip"
	"crypto/sha1"
	"errors"
	"io"

	"github.com/rs/xid"
)

// ID returns a new unique ID
func ID() string {
	return xid.New().String()
}

// GetFileHash reads a file from the FS and
// computes the md5 hash based on the file contents
func GetFileHash(file io.Reader) []byte {
	hash := sha1.New()
	io.TeeReader(file, hash)
	return hash.Sum(nil)
}

// GetBitFromByteArray returns the bit found FROM THE BACK in that position
// from the byte array.
func GetBitFromByteArray(pos int, array []byte) (byte, error) {
	arrayIndex := pos / 8

	if arrayIndex < 0 {
		return 0, errors.New("Invalid arguments")
	}

	return (array[arrayIndex] & byte(1<<uint(pos%8))) >> uint(pos%8), nil
}

// GzipBytes takes a byte array and gzip compresses it and returns
// the resulting byte array
func GzipBytes(content []byte) []byte {
	var compressed bytes.Buffer
	w := gzip.NewWriter(&compressed)
	w.Write(content)
	w.Close()
	return compressed.Bytes()
}
