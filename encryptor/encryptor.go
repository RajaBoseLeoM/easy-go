package encryptor

import (
	"time"
)

type fakeEncryptor struct {
	singleFrameEncryptionTimeInMs time.Duration
}

func NewEncryptor(singleFrameEncryptionTimeInMs time.Duration) *fakeEncryptor {
	return &fakeEncryptor{
		singleFrameEncryptionTimeInMs: singleFrameEncryptionTimeInMs,
	}
}

func (f *fakeEncryptor) Encrypt(frameNo int) (int, error) {
	time.Sleep(f.singleFrameEncryptionTimeInMs)

	return frameNo, nil
}
