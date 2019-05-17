package writer

import (
	"time"
)

type fakeWriter struct {
	singleFrameWritingTimeInMs time.Duration
}

func NewWriter(singleFrameWritingTimeInMs time.Duration) *fakeWriter {
	return &fakeWriter{
		singleFrameWritingTimeInMs: singleFrameWritingTimeInMs,
	}
}

func (f *fakeWriter) Write(frameNo int) error {
	time.Sleep(f.singleFrameWritingTimeInMs)

	return nil
}
