package main

import (
	"flag"
	"fmt"
	"io"
	"time"

	"github.com/RealImage/easy-go/encoder"
	"github.com/RealImage/easy-go/encryptor"
	"github.com/RealImage/easy-go/writer"
)

func main() {
	start := time.Now()

	encodeDur := flag.Int("s", 250, "encoding time taken for single frame")
	encryptDur := flag.Int("r", 100, "encryption time taken for single frame")
	writerDur := flag.Int("w", 25, "write time taken for single frame")
	totFrames := flag.Int("d", 10, "total no of frames")

	// Encoding
	e := encoder.NewEncoder(int32(*encodeDur), int32(*totFrames))
	frameNo := int32(1)
	for {
		fmt.Printf("Processing frame %v for encoding\n", frameNo)
		err := e.Encode(frameNo)
		if err == io.EOF {
			fmt.Println("Encoding complete")
			break
		}
		frameNo++
	}

	// Encrypting
	en := encryptor.NewEncryptor(int32(*encryptDur), int32(*totFrames))
	frameNo = int32(1)
	for {
		fmt.Printf("Processing frame %v for encryption\n", frameNo)
		err := en.Encrypt(frameNo)
		if err == io.EOF {
			fmt.Println("Encryption complete")
			break
		}
		frameNo++
	}

	// Writing
	w := writer.NewWriter(int32(*writerDur), int32(*totFrames))
	frameNo = int32(1)
	for {
		fmt.Printf("Processing frame %v for writing\n", frameNo)
		err := w.Write(frameNo)
		if err == io.EOF {
			fmt.Println("Writing complete")
			break
		}
		frameNo++
	}

	elapsed := time.Since(start)
	fmt.Printf("Encoding took %v\n", elapsed)
}
