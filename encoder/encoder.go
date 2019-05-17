package encoder

import (
	"io"
	"time"
)

type fakeEncoder struct {
	singleFrameEncodingTimeInMs time.Duration
	totalFrames                 int
	curFrameNo                  int
}

func NewEncoder(singleFrameEncodingTimeInMs time.Duration, totalFrames int) *fakeEncoder {
	return &fakeEncoder{
		singleFrameEncodingTimeInMs: singleFrameEncodingTimeInMs,
		totalFrames:                 totalFrames,
	}
}

func (f *fakeEncoder) Encode() (int, error) {
	time.Sleep(f.singleFrameEncodingTimeInMs)
	if f.curFrameNo == f.totalFrames {
		return f.curFrameNo, io.EOF
	}
	f.curFrameNo++

	return f.curFrameNo, nil
}
