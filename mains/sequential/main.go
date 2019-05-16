package main

import (
	"flag"
	"fmt"
	"io"
	"time"

	"github.com/RealImage/easy-go/encoder"
	"github.com/RealImage/easy-go/encryptor"
)

func main() {
	start := time.Now()

	encDur := flag.Int("s", 250, "single frame encoding duration")
	encryptDur := flag.Int("r", 100, "single frame encryption duration")
	totFrames := flag.Int("d", 10, "total no of frames")

	// Encoding
	e := encoder.NewEncoder(int32(*encDur), int32(*totFrames))
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
	elapsed := time.Since(start)
	fmt.Printf("Encoding took %v\n", elapsed)
}
