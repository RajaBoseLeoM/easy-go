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
	exitEncodeCh := make(chan bool)
	encryptCh := make(chan int32, *totFrames)
	frameNo := int32(1)
	for i := 0; i < *totFrames; i++ {
		go func(frameNo int32) {
			fmt.Printf("Processing frame %v for encoding\n", frameNo)
			err := e.Encode(frameNo)
			encryptCh <- frameNo
			if err == io.EOF {
				fmt.Println("Encoding complete")
				exitEncodeCh <- true
				close(encryptCh)
			}
		}(frameNo)
		frameNo++
	}

	// Encrypting
	en := encryptor.NewEncryptor(int32(*encryptDur), int32(*totFrames))
	exitEncryptCh := make(chan bool)
	ok := false
	exitEncryption := false
	for {
		select {
		case frameNo, ok = <-encryptCh:
			fmt.Println("")
			if ok {
				go func(frameNo int32) {
					fmt.Printf("Processing frame %v for encryption\n", frameNo)
					err := en.Encrypt(frameNo)
					if err == io.EOF {
						fmt.Println("Encryption complete")
						exitEncryption = true
						exitEncryptCh <- true
					}
				}(frameNo)
			} else {
				fmt.Println("Processed all frames for encryption!")
				exitEncryption = true
				exitEncryptCh <- true
			}
		default:
			fmt.Print("Waiting for frames to encrypt\r")
			time.Sleep(1 * time.Millisecond)
		}
		if exitEncryption {
			break
		}
	}

	<-exitEncodeCh
	<-exitEncryptCh

	elapsed := time.Since(start)
	fmt.Printf("Encoding/Encryption took %v\n", elapsed)
}
