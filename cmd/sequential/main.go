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
	for {
		frameNo, err := e.Encode()
		fmt.Printf("Processing frame %v for encoding\n", frameNo)
		if err == io.EOF {
			fmt.Println("Encoding complete")
			close(encryptCh)
			break
		} else if err != nil {
			log.Fatalf("failed to encode frame: %v\n", frameNo)
		}
		encryptCh <- frameNo
	}

	// Encrypting
	en := encryptor.NewEncryptor(time.Millisecond * time.Duration(*encryptDur))
	writerCh := make(chan int, *totFrames)
	for frameNo := range encryptCh {
		fmt.Printf("Processing frame %v for encryption\n", frameNo)
		eFrame, err := en.Encrypt(frameNo)
		if err != nil {
			log.Fatalf("failed to encrypt frame: %v\n", eFrame)
		}
		writerCh <- eFrame
	}
	close(writerCh)

	// Writing
	w := writer.NewWriter(time.Millisecond * time.Duration(*writerDur))
	for frameNo := range writerCh {
		fmt.Printf("Processing frame %v for writing\n", frameNo)
		err := w.Write(frameNo)
		if err != nil {
			log.Fatalf("failed to write frame: %v\n", frameNo)
		}
	}

	elapsed := time.Since(start)
	fmt.Printf("Sequencial encoding, encryption and writing took %v\n", elapsed)
}
