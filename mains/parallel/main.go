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
	exitEncodeCh := make(chan bool)
	encryptCh := make(chan int32, *totFrames)
	frameNo := int32(1)
	for i := 0; i < *totFrames; i++ {
		go func(frameNo int32) {
			fmt.Printf("Processing frame %v for encoding\n", frameNo)
			err := e.Encode(frameNo)
			encryptCh <- frameNo
			if err == io.EOF {
				fmt.Printf("\nEncoding complete\n")
				exitEncodeCh <- true
				close(encryptCh)
			}
		}(frameNo)
		frameNo++
	}

	// Encrypting
	en := encryptor.NewEncryptor(int32(*encryptDur), int32(*totFrames))
	exitEncryptCh := make(chan bool)
	writerCh := make(chan int32, *totFrames)
	ok := false
	exitEncryption := false
	for {
		select {
		case frameNo, ok = <-encryptCh:
			if ok {
				go func(frameNo int32) {
					fmt.Printf("\nProcessing frame %v for encryption\n", frameNo)
					err := en.Encrypt(frameNo)
					writerCh <- frameNo
					if err == io.EOF {
						fmt.Printf("\nEncryption complete\n")
						exitEncryption = true
						exitEncryptCh <- true
						close(writerCh)
					}
				}(frameNo)
			} else {
				fmt.Printf("\nProcessed all frames for encryption!\n")
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

	// Writing
	w := writer.NewWriter(int32(*writerDur), int32(*totFrames))
	exitWriter := false
	expectedFrameNo := int32(1)
	for {
		select {
		case frameNo, ok = <-writerCh:
			if ok {
				if expectedFrameNo != frameNo {
					writerCh <- frameNo
					break
				}
				fmt.Printf("\nProcessing frame %v for writing\n", frameNo)
				err := w.Write(frameNo)
				if err == io.EOF {
					fmt.Printf("\nWrite complete\n")
					exitWriter = true
				}
				expectedFrameNo++
			} else {
				fmt.Printf("\nProcessed all frames for writing\n")
				exitWriter = true
			}
		default:
			fmt.Print("Waiting for frames to write\r")
			time.Sleep(1 * time.Millisecond)
		}
		if exitWriter {
			break
		}
	}

	<-exitEncodeCh
	<-exitEncryptCh

	elapsed := time.Since(start)
	fmt.Printf("Encoding/Encryption took %v\n", elapsed)
}
