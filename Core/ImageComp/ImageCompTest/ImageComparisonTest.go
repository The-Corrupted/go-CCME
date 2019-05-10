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
	_, err := vfd.ReturnFrameData(1)
	if err != nil {
		fmt.Printf("Fail ( ReturnFrameData ): %v\n", err)
	} else {
		fmt.Printf("Success ( ReturnFrameData )\n")
	}
	PixelSet, err := vfd.ReturnPixelSet(1, "imgX0")
	if err != nil {
		fmt.Printf("Fail ( ReturnPixelSet )\n")
	} else {
		fmt.Println("Success ( ReturnPixelSet )")
	}
	var pixelToFind uint16 = 2
	pixelColor, err := vfd.ReturnPixelValue(PixelSet, pixelToFind)
	if err != nil {
		fmt.Printf("Fail ( ReturnPixelValue ): %v\n", err)
	} else {
		fmt.Printf("Success ( ReturnPixelValue ): R: %d G: %d B: %d A: %d\n", pixelColor[0], pixelColor[1], pixelColor[2], pixelColor[3])
	}
	pixelLocation, err := vfd.ReturnPixelLocation(PixelSet, pixelToFind)
	if err != nil {
		fmt.Printf("Fail ( ReturnPixelLocation ): %v\n", err)
	} else {
		fmt.Printf("Success ( ReturnPixelLocation ): X: %d Y: %d\n", pixelLocation[0], pixelLocation[1])
	}
	Success := vfd.StoreFrame(1)
	if !Success {
		fmt.Println("Failed to store frame.")
	}
	fmt.Println("Success...?")
}
