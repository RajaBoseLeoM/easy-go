package encryptor

import (
	"io"
	"time"
)

type fakeEncryptor struct {
	singleFrameEncryptionTimeInMs int32
	totalFrames                   int32
}

func NewEncryptor(singleFrameEncryptionTimeInMs int32, totalFrames int32) *fakeEncryptor {
	return &fakeEncryptor{
		singleFrameEncryptionTimeInMs: singleFrameEncryptionTimeInMs,
		totalFrames:                   totalFrames,
	}
}

func (f *fakeEncryptor) Encrypt(frameNo int32) error {
	time.Sleep(time.Millisecond * time.Duration(f.singleFrameEncryptionTimeInMs))
	if frameNo == f.totalFrames {
		return io.EOF
	}

	return nil
}
