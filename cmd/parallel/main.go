package main

import (
	"flag"
	"fmt"
	"io"
	"log"
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

	flag.Parse()

	// Encoding
	e := encoder.NewEncoder(time.Millisecond*time.Duration(*encodeDur), *totFrames)
	encryptCh := make(chan int, *totFrames)
	for i := 0; i < *totFrames; i++ {
		go func() {
			frameNo, err := e.Encode()
			fmt.Printf("Processing frame %v for encoding\n", frameNo)
			if err == io.EOF {
				fmt.Printf("\nEncoding complete\n")
			} else if err != nil {
				log.Fatalf("failed to encode for frame: %v\n", frameNo)
			}
			encryptCh <- frameNo
		}()
	}

	// Encrypting
	en := encryptor.NewEncryptor(time.Millisecond * time.Duration(*encryptDur))
	writerCh := make(chan int, *totFrames)
	currentEncryptionCount := 0
	for {
		select {
		case frameNo, ok := <-encryptCh:
			if ok {
				go func(frameNo int) {
					fmt.Printf("\nProcessing frame %v for encryption\n", frameNo)
					eFrame, err := en.Encrypt(frameNo)
					if err != nil {
						log.Fatalf("failed to encrypt frame: %v\n", frameNo)
					}
					writerCh <- eFrame
				}(frameNo)
				currentEncryptionCount++
				if currentEncryptionCount == *totFrames {
					close(encryptCh)
				}
			} else {
				fmt.Printf("\nProcessed all frames for encryption!\n")
			}
		default:
			fmt.Print("Waiting for frames to encrypt\r")
			time.Sleep(1 * time.Millisecond)
		}
		if currentEncryptionCount == *totFrames {
			break
		}
	}
	fmt.Printf("\nEncryption complete\n")

	// Writing
	w := writer.NewWriter(time.Millisecond * time.Duration(*writerDur))
	exitWriter := false
	expectedFrameNo := 1
	for {
		select {
		case frameNo, ok := <-writerCh:
			if ok {
				// FIXME: Instead sorting might be a better approach
				if expectedFrameNo != frameNo {
					writerCh <- frameNo
					break
				}
				fmt.Printf("\nProcessing frame %v for writing\n", frameNo)
				err := w.Write(frameNo)
				if err != nil {
					log.Fatalf("failed to write frame: %v\n", frameNo)
				}
				if expectedFrameNo == *totFrames {
					close(writerCh)
					break
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
	fmt.Printf("\nWrite complete\n")

	elapsed := time.Since(start)
	fmt.Printf("Time taken to Encoding & Encryption in parallel along with sequencial writing is %v\n", elapsed)
}
