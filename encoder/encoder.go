package encoder

import (
	"io"
	"time"
)

type fakeEncoder struct {
	singleFrameEncodingTimeInMs int32
	totalFrames                 int32
}

func NewEncoder(singleFrameEncodingTimeInMs int32, totalFrames int32) *fakeEncoder {
	return &fakeEncoder{
		singleFrameEncodingTimeInMs: singleFrameEncodingTimeInMs,
		totalFrames:                 totalFrames,
	}
}

func (f *fakeEncoder) Encode(frameNo int32) error {
	time.Sleep(time.Millisecond * time.Duration(f.singleFrameEncodingTimeInMs))
	if frameNo == f.totalFrames {
		return io.EOF
	}

	return nil
}
