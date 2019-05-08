package main

import (
	"sync"
	_ "testing"
	_ "os"
	"fmt"
	imgComp "github.com/The-Corrupted/go-CCME/Core/ImageComp"
);
func main() {
	var vfd imgComp.VideoFrameData
	vfd.SetName("Dillpickle")
	fmt.Println(vfd.ReturnName())
	ch := make(chan string)
	defer close(ch)
	var wg sync.WaitGroup
	wg.Add(1)
	vfd.ReadFrame(ch, "./DS.jpeg", 1, &wg)
	frameData, err := vfd.ReturnFrameData(1)
	if err != nil {
		fmt.Printf("Fail ( ReturnFrameData ): %v\n", err)
	} else {
		fmt.Printf("Success ( ReturnFrameData ): %v\n", frameData)
	}
	PixelSet, err := vfd.ReturnPixelSet(1, "imgX0")
	if err != nil {
		fmt.Printf("Fail ( ReturnPixelSet ): %v\n", err)
	} else {
		fmt.Println("Success ( ReturnPixelSet )")
	}
	var pixelToFind uint16 = 2
	pixelColor, err := vfd.ReturnPixelValue(PixelSet, pixelToFind)
	if err != nil {
		fmt.Printf("Fail ( ReturnPixelValue ): %v\n", err)
	} else {
		for x := 0; x < len(PixelSet.Pixels); x++ {
			fmt.Println(PixelSet.Pixels[x])
		}
		fmt.Printf("Success ( ReturnPixelValue ): R: %d G: %d B: %d A: %d\n", pixelColor[0], pixelColor[1], pixelColor[2], pixelColor[3])
	}
	pixelLocation, err := vfd.ReturnPixelLocation(PixelSet, pixelToFind)
	if err != nil {
		fmt.Printf("Fail ( ReturnPixelLocation ): %v\n", err)
	} else {
		fmt.Printf("Success ( ReturnPixelLocation ): X: %d Y: %d\n", pixelLocation[0], pixelLocation[1])
	}
}
