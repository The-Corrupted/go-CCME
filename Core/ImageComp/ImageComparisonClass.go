package ImageComp

import (
	"sync"
)

type PixelValue struct {
	R, G, B, A uint8
	X, Y       uint16
}

type PixelSet struct {
	Pixels []*PixelValue
}

type ImageData struct {
	Set   map[string]*PixelSet
	Frame uint64
}

type VideoFrameData struct {
	sync.RWMutex
	PixelSets []*ImageData
	vfgName   string
}
