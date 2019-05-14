package EDA

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/draw"
	"os"
	"strconv"
	"strings"
	"sync"
	CCMEExcept "github.com/The-Corrupted/go-CCME/Helpers/CCMEExcept"
	"github.com/The-Corrupted/gozbar"
	"github.com/The-Corrupted/go-CCME/Helpers/OsHandler"
)

func setCropZone(width int, height int) (int, int) {
	//Largest supported to least: Leave range to avoid
	//excessive checks.
	if width >= 3840 && height >= 1080 {
		if height >= 2160 {
			return width/12, height/8
		}
		return width/12, height/4
	}
	if width >= 1920 && height > 720 {
		if height >= 1080 {
			return width/7, height/4
		}
		return width/7, height/2
	}
	if width >= 1280 && height >= 720 {
		if height >= 1080 {
			return width/5, height/4
		}
		return width/5, height/2
	}
	if width >= 960 && height >= 768 {
		if height >= 960 {
			return width/3, height/3
		}
		return width/3, height/3
	}
	if width >= 720 && height >= 480 {
		if height == 600 {
			return int(float32(width)/2.5), height/2
		}
		return int(float32(width)/2.5), int(float32(height)/1.5)
	}
	if width == 640 && height >= 360 {
		if height >= 480 {
			return int(float32(width)/2.5), int(float32(height)/1.5)
		} 
		return int(float32(width)/2.5), int(float32(height)/1.3)
	} else {
		return width/1, height/1
	}
}

func getImageDimension(imagePath string) (int, int) {
	file, err := os.Open(imagePath)
	if err != nil {
		fmt.Println("Unable to open image file.")
	}
	image, _, err := image.DecodeConfig(file)
	if err != nil {
		fmt.Println("Unable to decode image file.")
	}
	return image.Width, image.Height
}

func (cs *ConcurrentSlice) ReadImage(c chan string, file string, filename string, wg *sync.WaitGroup) {
	defer wg.Done()
	CCMEExcept.Catcher {
		Try: func() {
			globalGrab := ""
			scanner := NewScanner()
			scanner.SetConfig(zbar.QRCODE, zbar.CFG_ADD_CHECK, 1)
			f, err := os.Open(file)
			defer f.Close()
			if err != nil {
				CCMEExcept.Throw(err)
			}

			i, _ := jpeg.Decode(f)
			width, height := getImageDimension(file)
			cropw, croph := setCropZone(width, height)

			subi := i.(SubImager).SubImage(image.Rect(0,0,cropw, croph))
			img := zbar.FromImage(subi)
			res := scanner.Scan(img)
			if res == 0 {
				panic("No symbols found.")
				// cs.Lock()
				// cs.ReturnedItems = append(cs.ReturnedItems, ReturnVal{"None", 0, "No qr symbols found."})
				// defer cs.Unlock()
				// return
			}
			if res == -1 {
				panic("Error occured.")
				// cs.Lock()
				// cs.ReturnedItems = append(cs.ReturnedItems, ReturnVal{"None", 0, "An error occured while scanning the image."})
				// defer cs.Unlock()
			}
			img.First().Each(func(str string) {
				globalGrab = str
			})
			fmt.Println(img)
			fmt.Println(scanner)
			img.Destroy()
			scanner.Destroy()
			fmt.Printf("GlobalGrab: %s", globalGrab)
			split := strings.Split(globalGrab, "\n")
			time, frame := split[0], split[1]
		//	info := fmt.Sprintf("Time: " + time + " frame#: " + frame)
			frameInt, _ := strconv.ParseUint(frame, 10, 64)
			cs.Lock()
			cs.ReturnedItems = append(cs.ReturnedItems, ReturnVal{time, frameInt, nil})
			defer cs.Unlock()
		},
		Catch: func(e CCMEExcept.Exception) {
			fmt.Println(e)
			fmt.Println("Failed to read the image.")
			ret, _ := os.Create(fmt.Sprintf("%d/Pictures/BadImages/%s.jpeg", OsHandler.SetUserDir(), filename))
			f, _ := os.Open(file)
			i, _ := jpeg.Decode(f)
			width, height := getImageDimension(file)
			cropw, croph := setCropZone(width, height)
			subi := i.(SubImager).SubImage(image.Rect(0,0,cropw,croph))
			bounds := subi.Bounds()
			imgGray := image.NewGray(bounds)
			draw.Draw(imgGray, bounds, subi, image.ZP, draw.Src)
			jpeg.Encode(ret, imgGray, nil)
			return
		},
		Finally: func() {
			return
		},
	}.Do()
}

func ( EDA *EDAClass ) CheckFunction(funct string) bool {
	if funct == EDA.Function {
		return true
	}
	return false
}

func BytesToString(byteArray []byte) string {
	return strings.Trim(string(byteArray[:]), "\n")
}

func (EDA *EDAClass) MakeDirs() {
	var i = 0
	for i < len(EDA.Paths) {
		if _, err := os.Stat(EDA.Paths[i]); os.IsNotExist(err) {
			os.MkdirAll(EDA.Paths[i], os.ModePerm)
		}
		i++
	}
}