package utils

import (
	"fmt"
	"io"
)

func BytesToInt64(b []byte) int64 {
	var val int64
	fmt.Sscanf(string(b), "%d", &val)
	return val
}

func ReadDockerOutput(reader io.Reader, w io.Writer) error {
	header := make([]byte, 8)
	for {
		_, err := io.ReadFull(reader, header)
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		size := int(header[4])<<24 | int(header[5])<<16 | int(header[6])<<8 | int(header[7])
		if size == 0 {
			continue
		}
		buf := make([]byte, size)
		_, err = io.ReadFull(reader, buf)
		if err != nil {
			return err
		}
		w.Write(buf)
	}
}
