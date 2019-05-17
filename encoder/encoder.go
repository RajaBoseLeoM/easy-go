package encoder

import (
	"io"
	"sync"
	"time"
)

type fakeEncoder struct {
	singleFrameEncodingTimeInMs time.Duration
	totalFrames                 int
	curFrameNo                  int
	m                           sync.Mutex
}

func NewEncoder(singleFrameEncodingTimeInMs time.Duration, totalFrames int) *fakeEncoder {
	return &fakeEncoder{
		singleFrameEncodingTimeInMs: singleFrameEncodingTimeInMs,
		totalFrames:                 totalFrames,
	}
}

func (f *fakeEncoder) Encode() (int, error) {
	time.Sleep(f.singleFrameEncodingTimeInMs)

	f.m.Lock()
	defer f.m.Unlock()

	if f.curFrameNo == f.totalFrames {
		return f.curFrameNo, io.EOF
	}
	f.curFrameNo++

	return f.curFrameNo, nil
}
