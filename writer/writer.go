package writer

import (
	"io"
	"time"
)

type fakeWriter struct {
	singleFrameWritingTimeInMs int32
	totalFrames                int32
}

func NewWriter(singleFrameWritingTimeInMs int32, totalFrames int32) *fakeWriter {
	return &fakeWriter{
		singleFrameWritingTimeInMs: singleFrameWritingTimeInMs,
		totalFrames:                totalFrames,
	}
}

func (f *fakeWriter) Write(frameNo int32) error {
	time.Sleep(time.Millisecond * time.Duration(f.singleFrameWritingTimeInMs))
	if frameNo == f.totalFrames {
		return io.EOF
	}

	return nil
}
