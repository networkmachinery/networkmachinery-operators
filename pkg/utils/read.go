package utils

import (
	"fmt"
	"io"
)

func ReadIn(reader io.Reader, p []byte) error {
	for {
		n, err := reader.Read(p)
		if err != nil{
			if err == io.EOF {
				fmt.Println(string(p[:n])) //should handle any remainding bytes.
				break
			}
			return err
		}
		fmt.Println(string(p[:n]))
	}
	return nil
}
