package ImageComp

import (
	"bytes"
	"compress/zlib"
	"encoding/gob"
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"io"
	"os"
	"sync"
)

// ---------------------------------VideoFrameData member functions------------------------------

func newSet(Frame uint64) *ImageData {
	Set := make(map[string]*PixelSet, 4)
	return &ImageData{Set, Frame}
}

func (vfg *VideoFrameData) SetName(vfgName string) {
	vfg.vfgName = vfgName
}

func (vfg *VideoFrameData) ReturnName() string {
	return vfg.vfgName
}

func (vfg *VideoFrameData) ReturnFrameData(frame uint64) (*ImageData, error) {
	if frame == 0 || frame-1 > uint64(len(vfg.PixelSets)) {
		return nil, errors.New("Bad frame selected")
	}
	if vfg.PixelSets[frame-1].Frame == frame {
		return vfg.PixelSets[frame-1], nil
	}
	return nil, errors.New("Frame not found")
}

func (vfg *VideoFrameData) ReturnPixelSet(frame uint64, key string) (*PixelSet, error) {
	frameData, err := vfg.ReturnFrameData(frame)
	if err != nil {
		return nil, err
	}
	return frameData.Set[key], nil
}

func (vfg *VideoFrameData) ReturnPixelValue(PixelSet *PixelSet, index uint16) ([]uint8, error) {
	index -= 1
	if index-1 > uint16(len(PixelSet.Pixels)) || index-1 < 0 {
		return []uint8{}, errors.New("Index out of range")
	}
	return []uint8{PixelSet.Pixels[index].R, PixelSet.Pixels[index].G, PixelSet.Pixels[index].B, PixelSet.Pixels[index].A}, nil
}

func (vfg *VideoFrameData) ReturnPixelLocation(PixelSet *PixelSet, index uint16) ([]uint16, error) {
	index -= 1
	if index > uint16(len(PixelSet.Pixels)) || index < 0 {
		return []uint16{}, errors.New("Index out of range")
	}
	return []uint16{PixelSet.Pixels[index].X, PixelSet.Pixels[index].Y}, nil
}

func (vfg *VideoFrameData) ReadFrame(c chan string, videoName string, frameNumber uint64, wg *sync.WaitGroup) error {
	defer wg.Done()
	image.RegisterFormat("jpeg", "jpeg", jpeg.Decode, jpeg.DecodeConfig)
	file, err := os.Open(videoName)
	if err != nil {
		fmt.Println("Failed to open image.")
		return err
	}
	defer file.Close()
	pixels, err := getPixels(file, frameNumber)
	if err != nil {
		fmt.Println("Failed to retrieve pixel information.")
		return err
	}
	vfg.Lock()
	vfg.PixelSets = append(vfg.PixelSets, pixels)
	vfg.Unlock()
	return nil
}

func (vfg *VideoFrameData) StoreFrame(frameNumber uint64) bool {
	var fileToCreate = fmt.Sprintf("%s%d.txt", vfg.vfgName, frameNumber)
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	enc.Encode(vfg.PixelSets)
	encodedStruct := buf.Bytes()
	buf.Reset()
	w := zlib.NewWriter(&buf)
	w.Write(encodedStruct)
	defer w.Close()
	_, err := os.Create(fileToCreate)
	if err != nil {
		fmt.Println(err)
	}
	file, err := os.OpenFile(fileToCreate, os.O_WRONLY, os.ModeAppend)
	if err != nil {
		fmt.Println(err)
	}
	if _, err := file.Write(buf.Bytes()); err != nil {
		fmt.Println(err)
	}
	return true
}

func (vfg *VideoFrameData) ReadStoredFrameData(frameNumber uint64) bool {
	var fileToOpen = fmt.Sprintf("%s%d.txt", vfg.vfgName, frameNumber)
	file, err := os.OpenFile(fileToOpen, os.O_RDONLY, 04)
	if err != nil {
		fmt.Println("Failed to open file. Does it exist?")
		return false
	}
	
}

//----------------------------Video Frame Data Private Functions---------------------

func getPixels(file io.Reader, frameNumber uint64) (*ImageData, error) {
	img, _, err := image.Decode(file)
	imageData := newSet(frameNumber)
	if err != nil {
		fmt.Println("Failed to decode the image.")
		return imageData, err
	}
	bounds := img.Bounds()
	width, height := uint16(bounds.Max.X), uint16(bounds.Max.Y)
	var xVals = [4]uint16{width / 2, width / 4, width - width/4, width - width/6}
	var yVals = [4]uint16{height / 2, height / 4, height - height/4, height - height/6}
	var x uint8 = 0
	for ; x < 4; x++ {
		imageData.Set[fmt.Sprintf("imgX%d", x)] = scroll(img, yVals[x], width, "h")
	}
	var y uint8 = 0
	for ; y < 4; y++ {
		imageData.Set[fmt.Sprintf("imgY%d", y)] = scroll(img, xVals[y], height, "v")
	}
	return imageData, nil
}

func scroll(img image.Image, startPos uint16, end uint16, direction string) *PixelSet {
	var Pixels PixelSet
	var x uint16 = 1
	switch direction {
	case "v":
		for ; x < end; x++ {
			var r, g, b, a uint32 = img.At(int(startPos), int(x)).RGBA()
			r = r / 257
			g = g / 257
			b = b / 257
			a = a / 257
			Pixel := &PixelValue{
				X: startPos,
				Y: x,
				R: uint8(r),
				G: uint8(g),
				B: uint8(b),
				A: uint8(a),
			}
			Pixels.Pixels = append(Pixels.Pixels, Pixel)
		}
		break
	case "h":
		for ; x < end; x++ {
			var r, g, b, a uint32 = img.At(int(x), int(startPos)).RGBA()
			r = r / 257
			g = g / 257
			b = b / 257
			a = a / 257
			Pixel := &PixelValue{
				X: x,
				Y: startPos,
				R: uint8(r),
				G: uint8(g),
				B: uint8(b),
				A: uint8(a),
			}
			Pixels.Pixels = append(Pixels.Pixels, Pixel)
		}
		break
	}
	return &Pixels
}
