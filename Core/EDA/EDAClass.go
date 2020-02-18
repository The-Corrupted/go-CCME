package EDA

import (
	"image"
	"sync"
)

type SubImager interface {
	SubImage(r image.Rectangle) image.Image
}

type EDAClass struct {
	Function string
	Deconstruct
	Analyze
	Paths []string
}

type Deconstruct struct {
	VideoName  string
	FramesName string
}

type Analyze struct {
	FramesName        string
	ExpNumberOfFrames uint64
	Delete            string
}

type ReturnVal struct {
	Time  string
	Frame uint64
	Error interface{}
}

type ConcurrentSlice struct {
	sync.RWMutex
	ReturnedItems []ReturnVal
}
