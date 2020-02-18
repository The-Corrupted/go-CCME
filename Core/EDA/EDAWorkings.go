package EDA

import (
	"fmt"
	handler "github.com/The-Corrupted/go-CCME/Core/CCMEHandlers"
	sql "github.com/The-Corrupted/go-CCME/Core/CCMESql"
	"github.com/The-Corrupted/go-CCME/Helpers/OsHandler"
	"github.com/The-Corrupted/gozbar"
	"io/ioutil"
	"os"
	"os/exec"
	"runtime"
	"sort"
	"strconv"
	"sync"
	timef "time"
)

func NewScanner() *zbar.Scanner {
	r := zbar.NewScanner()
	if r == nil {
		fmt.Println("Error: Scanner not made.")
	}
	return r
}

func NewEDA(Args handler.ArgumentLink) *EDAClass {
	switch Args.Name {
	case "Deconstruct":
		var UserDir = OsHandler.SetUserDir()
		Paths := []string{fmt.Sprintf("%s/Videos/StampedVideos/%s", UserDir, *Args.Values["DeconstructVideoName"]),
			fmt.Sprintf("%s/Pictures/EDA/%s/", UserDir, *Args.Values["DeconstructFrameNames"])}
		return &EDAClass{
			Function: Args.Name,
			Deconstruct: Deconstruct{
				VideoName:  *Args.Values["DeconstructVideoName"],
				FramesName: *Args.Values["DeconstructFrameNames"],
			},
			Paths: Paths,
		}
	case "Analyze":
		var UserDir = OsHandler.SetUserDir()
		expNumber, _ := strconv.ParseUint(*Args.Values["AnalyzeExpFrames"], 10, 64)
		Paths := []string{fmt.Sprintf("%s/Pictures/EDA/%s/", UserDir, *Args.Values["AnalyzeFrameNames"])}
		return &EDAClass{
			Function: Args.Name,
			Analyze: Analyze{
				FramesName:        *Args.Values["AnalyzeFrameNames"],
				ExpNumberOfFrames: expNumber,
				Delete:            *Args.Values["AnalyzeDelete"],
			},
			Paths: Paths,
		}
	}
	return &EDAClass{Function: "Fail"}
}

func (EDA *EDAClass) DeconstructVideo(c chan string) {
	if !EDA.CheckFunction("Deconstruct") {
		c <- "This eda instance was not created to use deconstruct, check the instances function argument."
	}
	EDA.MakeDirs()
	cmd := exec.Command("ffmpeg", "-i", EDA.Paths[0], "-qscale:v", "2", "-start_number",
		"0", "-vsync", "0", fmt.Sprintf("%s%s%%d.jpeg", EDA.Paths[1], EDA.Deconstruct.FramesName))
	stdout, stderr := cmd.CombinedOutput()
	if stderr != nil {
		c <- fmt.Sprintf("ERROR: %s\n%v", BytesToString(stdout), stderr)
	}
	fmt.Printf("%s\n", BytesToString(stdout))
	c <- fmt.Sprintf("%s", "Success")
}

/*

TODO:

Refactor the  main loop. Possible means could be a Tracker struct to contain the needed information
to analyze the video and the checks can be broken up into struct methods. Would need its own file.

*/

func (EDA *EDAClass) AnalyzeImages(c chan string) {

	//-------------------------------Setting CPU Usage------------------------------------------
	var cpus = runtime.NumCPU()
	runtime.GOMAXPROCS(cpus)

	//-------------------------------Create Objects---------------------------------------------
	var cs ConcurrentSlice
	scanner := NewScanner()
	defer scanner.Destroy()
	scanner.SetConfig(0, zbar.CFG_ENABLE, 1)

	//------------------------Initialize needed loop variables----------------------------------
	var loopDetectedMessage, BufferMessage,
		FullMessage, LastFrameData = "", "", "", ""
	var goodFrames, Frames, MangledOrMissing, LastFrame, ExtraFrames,
		CurrentFrame, NewNumber, offSet uint64 = 0, 0, 0, 0, 0, 0, 0, 0
	var ConsecutiveRepeat, ExtraFound, MissingMangledFound,
		finalFrameFound = false, false, false, false
	TimesRepeated := 0
	var ConsecutiveMissing uint16 = 0
	dirItems, _ := ioutil.ReadDir(EDA.Paths[0])
	dirItemNums := len(dirItems)
	var EntryName = fmt.Sprintf("Manual Analysis: %s", EDA.Analyze.FramesName)
	var wg sync.WaitGroup
	now := timef.Now()
	var x uint64 = 0

	//-----------------------Begin Analysis (Async File read)----------------------------------
	for x = 0; x < uint64(dirItemNums); x++ {
		if x%uint64(cpus) == 0 {
			wg.Wait()
		}
		wg.Add(1)
		num := strconv.FormatUint(x, 10)
		var file = EDA.Paths[0] + EDA.Analyze.FramesName + num + ".jpeg"
		var filename = EDA.Analyze.FramesName + num + ".jpeg"
		fmt.Println(file)
		ch := make(chan string)
		defer close(ch)
		go cs.ReadImage(ch, file, filename, &wg)
	}
	wg.Wait()

	//----------------------Sort then read collected data--------------------------------------
	sort.Slice(cs.ReturnedItems[:], func(i, j int) bool {
		return cs.ReturnedItems[i].Frame < cs.ReturnedItems[j].Frame
	})
	for x = 0; x < uint64(len(cs.ReturnedItems)); x++ {
		var time = cs.ReturnedItems[x].Time
		var frame = cs.ReturnedItems[x].Frame
		if ExtraFound == true {
			CurrentFrame += 1
		} else {
			ExtraFound = false
		}
		if MissingMangledFound == true {
			CurrentFrame = NewNumber
			offSet += 1
			MissingMangledFound = false
		}
		if ConsecutiveMissing == 240 {
			MangledOrMissing = MangledOrMissing - 240
			break
		}
		if finalFrameFound == true {
			fmt.Println("Final Frame found. Quiting.")
			break
		}
		CurrentFrame = x + 1 + offSet
		if LastFrame == 0 {
			if frame > 1 {
				fmt.Printf("First frame(s) missing. Frame counter increased by %d\n", frame)
				MangledOrMissing += frame
				LastFrame = frame
				continue
			}
		}
		fmt.Printf("%d\t%d\n", CurrentFrame, frame)
		if frame == CurrentFrame && x >= 0 {
			LastFrame = CurrentFrame
			LastFrameData = strconv.FormatUint(frame, 10)
			goodFrames += 1
			fmt.Println("QRCode Retrieved: " + LastFrameData)
			if frame == Frames {
				finalFrameFound = true
			}
			if TimesRepeated != 0 {
				if ConsecutiveRepeat != false {
					BufferMessage += fmt.Sprintf("Buffering/Stuttering found at frame %d: Repeated %d times.\n", LastFrame-1, TimesRepeated)
				}
				TimesRepeated = 0
				ConsecutiveRepeat = false
			}
			continue
		} else if LastFrameData == strconv.FormatUint(frame, 10) {
			LastFrameData = strconv.FormatUint(frame, 10)
			ExtraFrames += 1
			TimesRepeated += 1
			if TimesRepeated >= 20 {
				fmt.Println("Frame stutter/buffering found.")
				ConsecutiveRepeat = true
			}
			ExtraFound = true
			continue
		} else if LastFrame > frame {
			loopDetectedMessage += fmt.Sprintf("Video Loop detected. Last frame: %d Current frame: %d")
			break
		} else {
			if LastFrame == 0 && frame > 1 {
				FullMessage += fmt.Sprintln("First Frame Mangled or Missing.\n")
				FullMessage += fmt.Sprintf("%s\t%d\n", time, frame)
				continue
			} else if frame > 0 {
				FullMessage += fmt.Sprintf("Missing or mangled frame(s) after frame %d\n", LastFrame)
				LastFrameData = strconv.FormatUint(frame, 10)
				NewNumber = frame
				MissingMangledFound = true
				FullMessage += fmt.Sprintf("QR DATA: %d\n", frame)
				FullMessage += fmt.Sprintf("Occurence at file %d\n", x)
			} else {
				//First numbered frame not found yet.
				continue
			}
			MangledOrMissing += 1
			ConsecutiveMissing += 1
			continue
		}
	}
	then := timef.Now()
	elapsed := then.Sub(now)
	fmt.Printf("Elapsed time: %v\n", elapsed)
	percent := (float64(MangledOrMissing) / float64(EDA.Analyze.ExpNumberOfFrames)) * 100
	finalResult := fmt.Sprintf("%.3f%% of frames mangled or missing\n", percent)
	lastID := sql.Last_Id()
	ID := lastID + 1
	FullMessage += finalResult
	FullMessage += fmt.Sprintf(`%d frames repeated and or added`, ExtraFrames)
	_, err := sql.UpdateDB(ID, EntryName, MangledOrMissing, ExtraFrames, EDA.Analyze.ExpNumberOfFrames, finalResult, FullMessage)
	if err != nil {
		fmt.Printf("Error updating database: %v\n", err)
		c <- "FAIL"
	}
	fmt.Println(EDA.Delete)
	if EDA.Delete == "ON" {
		err := os.RemoveAll(EDA.Paths[0])
		if err != nil {
			fmt.Printf("ERROR: %v\n", err)
			c <- "FAIL"
		}
	}
	fmt.Println(FullMessage)
	c <- "Success"
}
