/*
An API client for opensubtitles.org

This is a client for the OSDb protocol. Currently the package only allows movie
identification, subtitles search, and download.
*/
package osdb

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"os"

	"github.com/kolo/xmlrpc"
)

const (
	ChunkSize = 65536 // 64k
)

// Allocate a new OSDB client
func NewClient() (*Client, error) {
	osdbServer := os.Getenv("OSDB_SERVER")
	if osdbServer == "" {
		osdbServer = DefaultOSDBServer
	}
	rpc, err := xmlrpc.NewClient(osdbServer, nil)
	if err != nil {
		return nil, err
	}

	c := &Client{
		UserAgent: DefaultUserAgent,
		Client:    rpc, // xmlrpc.Client
	}

	return c, nil
}

// Generate an OSDB hash for an *os.File.
func HashFile(file *os.File) (hash uint64, err error) {
	fi, err := file.Stat()
	if err != nil {
		return
	}
	if fi.Size() < ChunkSize {
		return 0, fmt.Errorf("File is too small")
	}

	// Read head and tail blocks.
	buf := make([]byte, ChunkSize*2)
	err = readChunk(file, 0, buf[:ChunkSize])
	if err != nil {
		return
	}
	err = readChunk(file, fi.Size()-ChunkSize, buf[ChunkSize:])
	if err != nil {
		return
	}

	// Convert to uint64, and sum.
	var nums [(ChunkSize * 2) / 8]uint64
	reader := bytes.NewReader(buf)
	err = binary.Read(reader, binary.LittleEndian, &nums)
	if err != nil {
		return 0, err
	}
	for _, num := range nums {
		hash += num
	}

	return hash + uint64(fi.Size()), nil
}

// Generate a OSDB hash for a file.
func Hash(path string) (uint64, error) {
	file, err := os.Open(path)
	if err != nil {
		return 0, err
	}
	defer file.Close()
	return HashFile(file)
}

// Read a chunk of a file at `offset` so as to fill `buf`.
func readChunk(file *os.File, offset int64, buf []byte) (err error) {
	n, err := file.ReadAt(buf, offset)
	if err != nil {
		return
	}
	if n != ChunkSize {
		return fmt.Errorf("Invalid read %v", n)
	}
	return
}
