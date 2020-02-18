package CCME

import (
	"bytes"
	"fmt"
	handler "github.com/The-Corrupted/go-CCME/Core/CCMEHandlers"
	"github.com/The-Corrupted/go-CCME/Helpers/OsHandler"
	"github.com/skip2/go-qrcode"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	timef "time"
)

// buf bytes.Buffer
// Logger = log.New(&buf, "logger: ", log.Ldate, log.Ltime, log.Llongfile)

func NewCCME(Args handler.ArgumentLink) *CCMEClass {
	fmt.Println(Args.Name)
	switch Args.Name {
	case "CreateFixed":
		gen, _ := strconv.ParseUint(*Args.Values["CreateFrameNumber"], 10, 32)
		fps, _ := strconv.ParseFloat(*Args.Values["CreateFrameNumberRate"], 64)
		Paths := []string{fmt.Sprintf("%s/Pictures/QRFrames/%s/", OsHandler.SetUserDir(), *Args.Values["CreateVideoName"])}
		return &CCMEClass{
			Function: "Create",
			Create: Create{
				PhotoName: *Args.Values["CreateVideoName"],
				GenNumber: gen,
				FPS:       fps,
			},
			Paths: Paths,
		}
	case "CreateTimed":
		fps, _ := strconv.ParseFloat(*Args.Values["CreateFrameRate"], 64)
		Paths := []string{fmt.Sprintf("%s/Pictures/QRFrames/%s/", OsHandler.SetUserDir(), *Args.Values["CreateVideoName"])}
		return &CCMEClass{
			Function: "Create",
			Create: Create{
				PhotoName:     *Args.Values["CreateVideoName"],
				FPS:           fps,
				FormattedTime: *Args.Values["CreateVideoTime"],
			},
			Paths: Paths,
		}
	case "Overlay":
		UserDir := OsHandler.SetUserDir()
		Paths := []string{fmt.Sprintf("%s/Videos/UnderlayVids/", UserDir), //Underlay Video Path
			fmt.Sprintf("%s/Pictures/QRFrames/%s/", UserDir, *Args.Values["OverlayVideo"]), //Overlay Frames Path
			fmt.Sprintf("%s/Videos/StampedVideos/", UserDir)}                               //Save directory for new video
		fmt.Sprintf("%s/Pictures/OriginalStamped/%s/", UserDir, *Args.Values["OverlayVideo"]) // Original Frames save directory
		return &CCMEClass{
			Function: Args.Name,
			Overlay: Overlay{
				UnderlayVideo:   *Args.Values["UnderlayVideo"],
				OverlayName:     *Args.Values["OverlayVideo"],
				OverlayPosition: "0:0",
				SaveOverlay:     *Args.Values["SaveOverlay"],
				FinalVideoName:  *Args.Values["OverlayFinalVideoName"],
			},
			Paths: Paths,
		}
	case "GetVidInfo":
		Paths := []string{fmt.Sprintf("%s/Videos/UnderlayVids/", OsHandler.SetUserDir())}
		return &CCMEClass{
			Function: Args.Name,
			GetFrames: GetFrames{
				VideoName: *Args.Values["Video"],
			},
			Paths: Paths,
		}
	case "EditVideo":
		Paths := []string{fmt.Sprintf("%s/Videos/Transcoded/", OsHandler.SetUserDir())}
		return &CCMEClass{
			Function: Args.Name,
			EditVideo: EditVideo{
				VideoPath:     *Args.Values["VideoPath"],
				NewFormat:     *Args.Values["VidFormat"],
				Quality:       *Args.Values["Quality"],
				EncodingSpeed: *Args.Values["EncodingSpeed"],
				OrgVidName:    *Args.Values["VideoName"],
				Width:         *Args.Values["Width"],
				Height:        *Args.Values["Height"],
				MaxBV:         *Args.Values["BVRate"],
			},
			Paths: Paths,
		}
	}
	return &CCMEClass{Function: "Fail"}
}

func newTimeCounter(StartValues []uint8) *TimeCounter {
	if len(StartValues) < 3 {
		return &TimeCounter{
			TimeValues:  []uint8{0, 0, 0},
			ZPaddedTime: "00:00:00",
		}
	} else {
		return &TimeCounter{
			TimeValues:  StartValues,
			ZPaddedTime: fmt.Sprintf("%02d:%02d:%02d", StartValues[0], StartValues[1], StartValues[2]),
		}
	}
}

func (CCME *CCMEClass) CreateQR(c chan string) {
	fmt.Println(CCME.Function)
	if !CCME.CheckFunction("Create") {
		c <- "This class was not constructed to use create. Check its function argument."

	}
	CCME.MakeDirs()
	TimeCounter := newTimeCounter([]uint8{0, 0, 0})
	var passes uint64 = 0
	var FramesToSeconds uint8 = 0
	if CCME.Create.GenNumber > 0 {
		for passes < CCME.Create.GenNumber {
			FramesToSeconds += 1
			if float64(FramesToSeconds)-CCME.Create.FPS >= 0 {
				fmt.Printf("%f, %f", float64(FramesToSeconds), CCME.Create.FPS)
				FramesToSeconds = 0
				TimeCounter.IncrementCounter()
			}
			err := qrcode.WriteFile(fmt.Sprintf("%s\n%d", TimeCounter.ZPaddedTime, passes+1), qrcode.Highest, 256, fmt.Sprintf("%s%s%d.png", CCME.Paths[0], CCME.Create.PhotoName, passes))
			if err == nil {
				passes++
				continue
			} else {
				ReturnString := fmt.Sprintf("Error: %v", err)
				c <- ReturnString
			}
		}
	} else {
		VideoTime := GetTimeInSeconds(CCME.Create.FormattedTime)
		GenNumber := VideoTime * uint64(CCME.Create.FPS)
		for passes < GenNumber {
			FramesToSeconds += 1
			if float64(FramesToSeconds)-CCME.Create.FPS <= 0 {
				fmt.Printf("%f, %f", float64(FramesToSeconds), CCME.Create.FPS)
				FramesToSeconds = 0
				TimeCounter.IncrementCounter()
			}
			err := qrcode.WriteFile(fmt.Sprintf("%s\n%d", TimeCounter.ZPaddedTime, passes+1), qrcode.Highest, 256, fmt.Sprintf("%s%s%d.png", CCME.Paths[0], CCME.Create.PhotoName, passes))
			if err == nil {
				passes++
				continue
			} else {
				ReturnString := fmt.Sprintf("Error: %v", err)
				c <- ReturnString
			}
		}
	}
	c <- "Codes Successfully Created!"
}

func (CCME *CCMEClass) OverlayQR(c chan string) {
	fmt.Println("Setting up needed vars")
	fmt.Printf("OverlayVideo: %v\nUnderlayVideo: %v\n", CCME.Overlay.OverlayName, CCME.Overlay.UnderlayVideo)
	if !CCME.CheckFunction("Overlay") {
		c <- "This CCME instance was not created to use overlay. Check its function argument."
	}
	start := timef.Now()
	split := strings.Split(CCME.Overlay.UnderlayVideo, ".")
	VideoName := split[0]
	VideoFormat := split[1]
	AudioFS := GetAudioAndFramerate(fmt.Sprintf("%s%s", CCME.Paths[0], CCME.Overlay.UnderlayVideo))
	Bitrate, err := EstimateBitrate(fmt.Sprintf("%s%s", CCME.Paths[0], CCME.Overlay.UnderlayVideo))
	if err != nil {
		c <- fmt.Sprintf("Failed to get the videos bitrate. %v\n", err)
	}
	var BufSize = SetBufSize(Bitrate)
	TempDirPath := OsHandler.MakeTempDirOverlay(VideoName)
	if AudioFS[0] != "None" {
		cmdExtractAudio := exec.Command("ffmpeg", "-i", fmt.Sprintf("%s%s", CCME.Paths[0], CCME.Overlay.UnderlayVideo), "-vn",
			"-acodec", "copy", fmt.Sprintf("%s%s.%s", TempDirPath[0], VideoName, AudioFS[0]), "-y")

		out, err := cmdExtractAudio.CombinedOutput()
		if err != nil {
			c <- "Failed to extract audio."
		}
		fmt.Println(BytesToString(out))
	}

	fmt.Println("Bitrate: " + Bitrate)
	fmt.Println("BufSize: " + BufSize)

	//Dismantle Video

	fmt.Println("Dismantling video")
	cmdDismantle := exec.Command("ffmpeg", "-i", fmt.Sprintf("%s%s", CCME.Paths[0], CCME.Overlay.UnderlayVideo),
		"-start_number", "0", "-qscale:v", "2", fmt.Sprintf("%s%s%%d.jpeg", TempDirPath[0], VideoName))
	Output, Error := cmdDismantle.CombinedOutput()
	if Error != nil {
		fmt.Println(BytesToString(Output))
		c <- "Failed to dismantle video frames"
	}

	//Manually overlay frames
	fmt.Println("Overlaying frames")
	files, err := ioutil.ReadDir(CCME.Paths[1])
	if err != nil {
		fmt.Println(CCME.Paths[1])
		c <- "Failed to read directory."
		return
	}
	size := len(files)
	files, err = ioutil.ReadDir(TempDirPath[0])
	if err != nil {
		c <- "Failed to read directory"
		return
	}
	size2 := len(files)
	//Make sure a minimum size is found ( if one goes over the other the program fails )
	if size2-2 < size && AudioFS[0] != "None" {
		fmt.Printf("Size set to TempDirPath: %d\n", size2-2)
		size = size2 - 2
	} else if size2-1 < size && AudioFS[0] == "None" {
		fmt.Printf("Size set to TempDirPath: %d\n", size2-1)
		size = size2 - 1
	} else {
		fmt.Printf("Size set to QRFramesPath: %d\n", size)
	}

	for x := 0; x < size; x++ {
		fmt.Printf(fmt.Sprintf("%s%s%d.jpeg", TempDirPath[0], VideoName, x) + "\n" + fmt.Sprintf("%s%s%d.png", CCME.Paths[1], CCME.Overlay.OverlayName, x) + "\n")
		overlayFile := exec.Command("ffmpeg", "-i", fmt.Sprintf("%s%s%d.jpeg", TempDirPath[0], VideoName, x),
			"-i", fmt.Sprintf("%s%s%d.png", CCME.Paths[1], CCME.Overlay.OverlayName, x),
			"-filter_complex", "overlay", fmt.Sprintf("%s%s%d.jpeg", TempDirPath[1], VideoName, x), "-y")
		out, err := overlayFile.CombinedOutput()
		if err != nil {
			fmt.Printf("%s\n%v\n", BytesToString(out), err)
			c <- "Failed to overlay frames"
			return
		}
		error := os.Remove(fmt.Sprintf("%s%s%d.jpeg", TempDirPath[0], VideoName, x))
		if error != nil {
			fmt.Printf("Error: %v", error)
			c <- "Failed to remove temporary frames"
			return
		}
	}

	//Reconstruct The Video

	codec := GetVideoCodec(fmt.Sprintf("%s%s", CCME.Paths[0], CCME.Overlay.UnderlayVideo), VideoFormat)
	var codecToUse string
	fmt.Printf("Framerate: %s\n", AudioFS[1])
	switch codec {
	case "h264":
		cmdConstruct := exec.Command("ffmpeg", "-framerate", AudioFS[1], "-i", fmt.Sprintf("%s%s%%d.jpeg", TempDirPath[1], VideoName),
			"-c:v", "libx264", "-b:v", Bitrate, "-maxrate", Bitrate, "-bufsize", BufSize, "-minrate", BufSize, "-preset", "ultrafast", "-profile", "high", "-pix_fmt", "yuv420p", "-y",
			"-strict", "-2", fmt.Sprintf("%s%s.%s", CCME.Paths[2], CCME.Overlay.FinalVideoName, VideoFormat))
		codecToUse = "libx264"
		Output, Error = cmdConstruct.CombinedOutput()
		if Error != nil {
			fmt.Println(BytesToString(Output))
			c <- "Failed to construct back into a video."
		}
	case "hevc":
		cmdConstruct := exec.Command("ffmpeg", "-i", fmt.Sprintf("%s%s%%d.jpeg", TempDirPath[1], VideoName), "-framerate", AudioFS[1],
			"-c:v", "libx265", "-b:v", Bitrate, "-maxrate", Bitrate, "-bufsize", BufSize, "-minrate", BufSize, "-preset", "ultrafast", "-pix_fmt", "yuv420p", "-y",
			"-vsync", "0", "-strict", "-2", fmt.Sprintf("%s%s.%s", CCME.Paths[2], CCME.Overlay.FinalVideoName, VideoFormat))
		codecToUse = "libx265"
		Output, Error = cmdConstruct.CombinedOutput()
		if Error != nil {
			fmt.Println(BytesToString(Output))
			c <- "Failed to construct back into a video."
		}
	case "vp9":
		cmdConstruct := exec.Command("ffmpeg", "-i", fmt.Sprintf("%s%s%%d.jpeg", TempDirPath[1], VideoName), "-framerate", AudioFS[1],
			"-c:v", "libvpx-vp9", "-b:v", Bitrate, "-maxrate", Bitrate, "-bufsize", BufSize, "-minrate", BufSize, "-pix_fmt", "yuv420p", "-y", "-deadline", "realtime",
			"-vsync", "0", "-strict", "-2", fmt.Sprintf("%s%s.%s", CCME.Paths[2], CCME.Overlay.FinalVideoName, VideoFormat))
		codecToUse = "libvpx-vp9"
		Output, Error = cmdConstruct.CombinedOutput()
		if Error != nil {
			fmt.Println(BytesToString(Output))
			c <- "Failed to construct back into a video."
		}
	case "vp8":
		cmdConstruct := exec.Command("ffmpeg", "-i", fmt.Sprintf("%s%s%%d.jpeg", TempDirPath[1], VideoName), "-framerate", AudioFS[1],
			"-c:v", "libvpx", "-b:v", Bitrate, "-maxrate", Bitrate, "-bufsize", BufSize, "-minrate", BufSize, "-pix_fmt", "yuv420p", "-y", "deadline", "realtime",
			"-vsync", "0", "-strict", "-2", fmt.Sprintf("%s%s.%s", CCME.Paths[2], CCME.Overlay.FinalVideoName, VideoFormat))
		codecToUse = "libvpx"
		Output, Error = cmdConstruct.CombinedOutput()
		if Error != nil {
			fmt.Println(BytesToString(Output))
			c <- "Failed to construct back into a video."
		}
	default:
		errorr := os.RemoveAll(TempDirPath[0])
		if errorr != nil {
			fmt.Printf("%v\n", errorr)
		}
		c <- codec + " not currently supported."
	}

	if AudioFS[0] != "None" {
		switch AudioFS[0] {
		case "aac":
			cmdAudioAdd := exec.Command("ffmpeg", "-i", fmt.Sprintf("%s%s.%s", CCME.Paths[2], CCME.Overlay.FinalVideoName, VideoFormat),
				"-i", fmt.Sprintf("%s%s.%s", TempDirPath[0], VideoName, AudioFS[0]), "-c:v", "copy", "-c:a", "libfdk_aac",
				"-y", fmt.Sprintf("%s%s%d.%s", CCME.Paths[2], CCME.Overlay.FinalVideoName, 2, VideoFormat))
			Output, Error = cmdAudioAdd.CombinedOutput()
			if Error != nil {
				fmt.Println(BytesToString(Output))
				c <- "Failed To Re-add Audio"
			}
		case "ogg":
			cmdAudioAdd := exec.Command("ffmpeg", "-i", fmt.Sprintf("%s%s.%s", CCME.Paths[2], CCME.Overlay.FinalVideoName, VideoFormat),
				"-i", fmt.Sprintf("%s%s.%s", TempDirPath[0], VideoName, AudioFS[0]), "-c:v", "copy", "-c:a", "libvorbis",
				"-y", fmt.Sprintf("%s%s%d.%s", CCME.Paths[2], CCME.Overlay.FinalVideoName, 2, VideoFormat))
			Output, Error = cmdAudioAdd.CombinedOutput()
			if Error != nil {
				fmt.Println(BytesToString(Output))
				c <- "Failed To Re-add Audio"
			}
		case "mp3":
			cmdAudioAdd := exec.Command("ffmpeg", "-i", fmt.Sprintf("%s%s.%s", CCME.Paths[2], CCME.Overlay.FinalVideoName, VideoFormat),
				"-i", fmt.Sprintf("%s%s.%s", TempDirPath[0], VideoName, AudioFS[0]), "-c:v", "copy", "-c:a", "libmp3lame",
				"-y", fmt.Sprintf("%s%s%d.%s", CCME.Paths[2], CCME.Overlay.FinalVideoName, 2, VideoFormat))
			Output, Error = cmdAudioAdd.CombinedOutput()
			if Error != nil {
				fmt.Println(BytesToString(Output))
				c <- "Failed To Re-add Audio"
			}
		default:
			errorr := os.RemoveAll(TempDirPath[0])
			if errorr != nil {
				fmt.Printf("%v\n", errorr)
			}
			c <- AudioFS[0] + " currently not a supported audio format."
		}
	}

	//Add timer to the bottom of the video
	if AudioFS[0] != "None" {
		cmdTimer := exec.Command("ffmpeg", "-i", fmt.Sprintf("%s%s%d.%s", CCME.Paths[2], CCME.Overlay.FinalVideoName, 2, VideoFormat), "-vf",
			"drawtext=fontfile=/usr/share/fonts/truetype/freefont/FreeSerif.ttf:text='%{pts \\: hms}':x=0:y=h-th:fontsize=16:fontcolor=white@0.9:box=1:boxcolor=black@0.6",
			"-vsync", "0", "-y", "-c:a", "copy", "-vcodec", codecToUse, "-threads", "3", "-strict", "-2", fmt.Sprintf("%s%s%d.%s", CCME.Paths[2], CCME.Overlay.FinalVideoName, 3, VideoFormat))

		Output, Error = cmdTimer.CombinedOutput()
		if Error != nil {
			fmt.Println(BytesToString(Output))
			c <- "Failed To Add Timer"
		}
	} else {
		cmdTimer := exec.Command("ffmpeg", "-i", fmt.Sprintf("%s%s.%s", CCME.Paths[2], CCME.Overlay.FinalVideoName, VideoFormat), "-vf",
			"drawtext=fontfile=/usr/share/fonts/truetype/freefont/FreeSerif.ttf:text='%{pts \\: hms}':x=0:y=h-th:fontsize=16:fontcolor=white@0.9:box=1:boxcolor=black@0.6",
			"-vsync", "0", "-y", "-c:a", "copy", "-vcodec", codecToUse, "-threads", "3", "-strict", "-2", fmt.Sprintf("%s%s%d.%s", CCME.Paths[2], CCME.Overlay.FinalVideoName, 3, VideoFormat))

		Output, Error = cmdTimer.CombinedOutput()
		if Error != nil {
			fmt.Println(BytesToString(Output))
			c <- "Failed To Add Timer"
		}
	}

	//Remove and rename

	os.Remove(fmt.Sprintf("%s%s.%s", CCME.Paths[2], CCME.Overlay.FinalVideoName, VideoFormat))
	if AudioFS[0] != "None" {
		os.Remove(fmt.Sprintf("%s%s%d.%s", CCME.Paths[2], CCME.Overlay.FinalVideoName, 2, VideoFormat))
	}
	os.Rename(fmt.Sprintf("%s%s%d.%s", CCME.Paths[2], CCME.Overlay.FinalVideoName, 3, VideoFormat),
		fmt.Sprintf("%s%s.%s", CCME.Paths[2], CCME.Overlay.FinalVideoName, VideoFormat))

	fmt.Printf("TempDirPath: %v\n", TempDirPath[0])
	errorr := os.RemoveAll(TempDirPath[0])
	if errorr != nil {
		fmt.Printf("%v\n", errorr)
	}
	end := timef.Now()
	elapsed := end.Sub(start)
	fmt.Printf("It took %v to overlay the video.\n", elapsed)

	c <- fmt.Sprintf("%s.%s", CCME.Overlay.FinalVideoName, VideoFormat)
}

func (CCME *CCMEClass) GetVidInfo(c chan map[string]string) {
	var buf bytes.Buffer
	Logger := log.New(&buf, "logger: ", log.Llongfile)
	if !CCME.CheckFunction("GetVidInfo") {
		fmt.Println("This CCME instance was not constructed to use GetVidInfo, check the instances function argument.")
		c <- map[string]string{"Status": "FAIL", "Err": "Wrong instance to use this function."}
	}

	cmd := exec.Command("ffprobe", "-v", "error", "-select_streams", "v:0", "-show_entries",
		"stream=codec_name", "-of", "default=noprint_wrappers=1:nokey=1",
		fmt.Sprintf("%s%s", CCME.Paths[0], CCME.GetFrames.VideoName))

	videoCodec, sterr := cmd.CombinedOutput()
	if sterr != nil {
		Logger.Print(fmt.Printf("ERROR: %s %s \n", sterr, CCME.GetFrames.VideoName))
		fmt.Println(&buf)
		c <- map[string]string{"Status": "FAIL", "Err": fmt.Sprintf("ERROR: %s\n%v", videoCodec, sterr)}
	}

	cmd = exec.Command("ffprobe", "-v", "error", "-select_streams", "a:0", "-show_entries",
		"stream=codec_name", "-of", "default=noprint_wrappers=1:nokey=1",
		fmt.Sprintf("%s%s", CCME.Paths[0], CCME.GetFrames.VideoName))

	audioCodec, sterr := cmd.CombinedOutput()

	if sterr != nil {
		Logger.Print(fmt.Printf("ERROR: %s %s\n", audioCodec, sterr))
		fmt.Println(&buf)
		audioCodec = []byte("No Audio")
	}
	var audioReturn = BytesToString(audioCodec)
	if audioReturn == "" || audioReturn == " " {
		audioCodec = []byte("No Audio")
	}

	cmd = exec.Command("ffprobe", "-loglevel", "fatal", "-v", "0", "-count_frames",
		"-select_streams", "v:0", "-show_entries",
		"stream=nb_read_frames", "-of", "default=nokey=1:noprint_wrappers=1",
		fmt.Sprintf("%s%s", CCME.Paths[0], CCME.GetFrames.VideoName))

	frameNumber, sterr := cmd.CombinedOutput()
	if sterr != nil {
		Logger.Print(fmt.Printf("ERROR: %s %s\n", sterr, CCME.GetFrames.VideoName))
		fmt.Println(&buf)
		c <- map[string]string{"Status": "FAIL", "Err": fmt.Sprintf("ERROR: %s\n%v", frameNumber, sterr)}
	}
	cmd = exec.Command("ffprobe", "-loglevel", "fatal", "-v", "0", "-of", "csv=p=0", "-select_streams", "0",
		"-show_entries", "stream=r_frame_rate",
		fmt.Sprintf("%s%s", CCME.Paths[0], CCME.GetFrames.VideoName))
	frameRate, sterr := cmd.CombinedOutput()
	if sterr != nil {
		Logger.Print(fmt.Printf("ERROR: %s %s\n", sterr, CCME.GetFrames.VideoName))
		fmt.Println(&buf)
		c <- map[string]string{"Status": "FAIL", "Err": fmt.Sprintf("ERROR: %s\n%v\n", frameRate, sterr)}
	}
	cmd = exec.Command("ffprobe", "-v", "error", "-select_streams", "v:0", "-show_entries",
		"stream=profile", "-of", "default=nokey=1:noprint_wrappers=1", fmt.Sprintf("%s%s", CCME.Paths[0], CCME.GetFrames.VideoName))
	videoProfile, sterr := cmd.CombinedOutput()
	if sterr != nil {
		Logger.Print(fmt.Printf("ERROR: %s %s\n", sterr, CCME.GetFrames.VideoName))
		fmt.Println(&buf)
		c <- map[string]string{"Status": "FAIL", "Err": fmt.Sprintf("ERROR: %s\n%v\n", videoProfile, sterr)}
	}
	cmd = exec.Command("ffprobe", "-v", "error", "-select_streams", "v:0", "-show_entries",
		"stream=pix_fmt", "-of", "default=nokey=1:noprint_wrappers=1", fmt.Sprintf("%s%s", CCME.Paths[0], CCME.GetFrames.VideoName))
	pixelFormat, sterr := cmd.CombinedOutput()
	if sterr != nil {
		Logger.Print(fmt.Printf("ERROR: %s %s\n", sterr, CCME.GetFrames.VideoName))
		fmt.Println(&buf)
		c <- map[string]string{"Status": "FAIL", "Err": fmt.Sprintf("ERROR: %s\n%v\n", videoProfile, sterr)}
	}
	c <- map[string]string{"Status": "OK", "VideoCodec": BytesToString(videoCodec), "AudioCodec": BytesToString(audioCodec),
		"FrameNumber": BytesToString(frameNumber), "FrameRate": BytesToString(frameRate),
		"VideoProfile": BytesToString(videoProfile), "PixelFormat": BytesToString(pixelFormat), "Err": "nil"}

}

//Must be redone so that more precise/custom bitrates can be choosen by the user.
//This will require a rewrite of EditVideo as well as the html form.

func (CCME *CCMEClass) TranscodeVideo(c chan string) {
	var stdout []byte
	var stderr error
	var failed = false
	if !CCME.CheckFunction("EditVideo") {
		c <- "This ccme instance was not constructed to use editvideo, check the instances function argument."
	}
	if CCME.EditVideo.Width == "default" || CCME.EditVideo.Height == "default" {
		cmd := exec.Command("ffprobe", "-v", "error", "-select_streams", "v:0", "-show_entries",
			"stream=width,height", "-of", "csv=s=x:p=0", fmt.Sprintf("%s", CCME.EditVideo.VideoPath))
		stdout, stderr = cmd.CombinedOutput()
		if stderr != nil {
			fmt.Println("Probe failed.")
			failed = true
		}
	}
	if CCME.EditVideo.Width == "default" && failed != true && CCME.EditVideo.Height != "default" {
		fmt.Println("Width default")
		fmt.Println(BytesToString(stdout))
		CCME.EditVideo.Width = strings.Split(strings.Trim(BytesToString(stdout), "\n"), "x")[0]
	}
	if CCME.EditVideo.Height == "default" && failed != true && CCME.EditVideo.Width != "default" {
		fmt.Println("Height default")
		fmt.Println(BytesToString(stdout))
		CCME.EditVideo.Height = strings.Split(strings.Trim(BytesToString(stdout), "\n"), "x")[1]
	}
	if CCME.EditVideo.Height == "default" && failed != true && CCME.EditVideo.Width == "default" {
		fmt.Println("All default")
		fmt.Println(BytesToString(stdout))
		var dimensions = strings.Split(strings.Trim(BytesToString(stdout), "\n"), "x")
		CCME.EditVideo.Width = dimensions[0]
		CCME.EditVideo.Height = dimensions[1]
	}
	var encodingInfo = strings.Split(CCME.EditVideo.NewFormat, ":")
	var fileType, encoder = encodingInfo[0], encodingInfo[1]
	var crfRate, webmLossless string
	switch CCME.EditVideo.Quality {
	case "0":
		if encoder == "libvpx" {
			crfRate = "4"
		} else {
			crfRate = "0"
		}
		webmLossless = "1"
	case "1":
		crfRate = "6"
	case "2":
		crfRate = "17"
		webmLossless = "0"
	case "3":
		crfRate = "34"
		webmLossless = "0"
	case "4":
		crfRate = "40"
		webmLossless = "0"
	case "5":
		crfRate = "51"
		webmLossless = "0"
	}
	var webmEncodeRate string
	if fileType == "webm" {
		switch CCME.EditVideo.EncodingSpeed {
		case "ultrafast":
			webmEncodeRate = "realtime"
		case "veryslow":
			webmEncodeRate = "best"
		default:
			webmEncodeRate = "good"
		}
	}
	var command *exec.Cmd
	var CCMENumberString = strings.Trim(CCME.EditVideo.MaxBV, "M")
	CCMENumber, _ := strconv.ParseInt(CCMENumberString, 10, 64)
	if CCMENumber > 0 {
		fmt.Println("Limited Transcoding")
		switch encoder {
		case "libx264":
			command = exec.Command("ffmpeg", "-hide_banner", "-i", CCME.EditVideo.VideoPath, "-c:v", encoder, "-x264-params", "\"nal-hrd=cbr\"", "-b:v", CCME.EditVideo.MaxBV, "-minrate", CCME.EditVideo.MaxBV,
				"-maxrate", CCME.EditVideo.MaxBV, "-BufSize", SetBufSize(CCME.EditVideo.MaxBV), "-vf", fmt.Sprintf("scale=%s:%s", CCME.EditVideo.Width, CCME.EditVideo.Height), "-pix_fmt", "yuv420p", "-c:a", "libfdk_aac",
				"-preset", CCME.EditVideo.EncodingSpeed, fmt.Sprintf("%s%s.%s", CCME.Paths[0], CCME.EditVideo.OrgVidName, fileType), "-y")
		case "libx265":
			fmt.Println(CCME.EditVideo.VideoPath)
			command = exec.Command("ffmpeg", "-hide_banner", "-i", CCME.EditVideo.VideoPath, "-c:v", encoder, "-preset", CCME.EditVideo.EncodingSpeed, "-b:v", CCME.EditVideo.MaxBV,
				"-vf", fmt.Sprintf("scale=%s:%s", CCME.EditVideo.Width, CCME.EditVideo.Height), "-pix_fmt", "yuv420p", "-c:a", "libfdk_aac", fmt.Sprintf("%s%s.%s", CCME.Paths[0], CCME.EditVideo.OrgVidName, fileType), "-y")
		case "libvpx":
			command = exec.Command("ffmpeg", "-hide_banner", "-i", CCME.EditVideo.VideoPath, "-c:v", encoder, "-vf", fmt.Sprintf("scale=%s:%s", CCME.EditVideo.Width, CCME.EditVideo.Height), "-minrate", CCME.EditVideo.MaxBV, "-maxrate", CCME.EditVideo.MaxBV, "-b:v", CCME.EditVideo.MaxBV, "-deadline", webmEncodeRate,
				"-c:a", "libvorbis", fmt.Sprintf("%s%s.%s", CCME.Paths[0], CCME.EditVideo.OrgVidName, fileType), "-y")
		case "libvpx-vp9":
			command = exec.Command("ffmpeg", "-hide_banner", "-i", CCME.EditVideo.VideoPath, "-c:v", encoder, "-vf", fmt.Sprintf("scale=%s:%s", CCME.EditVideo.Width, CCME.EditVideo.Height), "-minrate", CCME.EditVideo.MaxBV, "-maxrate", CCME.EditVideo.MaxBV, "-b:v", CCME.EditVideo.MaxBV, "-deadline", webmEncodeRate,
				"-c:a", "libvorbis", fmt.Sprintf("%s%s.%s", CCME.Paths[0], CCME.EditVideo.OrgVidName, fileType), "-y")
		default:
			command = exec.Command("ffmpeg", "-hide_banner", "-i", "-vf", fmt.Sprintf("scale=%s:%s", CCME.EditVideo.Width, CCME.EditVideo.Height),
				"-b:v", CCME.EditVideo.MaxBV, CCME.EditVideo.VideoPath, fmt.Sprintf("%s%s.%s", CCME.Paths[0], CCME.EditVideo.OrgVidName, fileType), "-y")
		}
		stdout, stderr = command.CombinedOutput()
		if stderr != nil {
			fmt.Printf("ERROR: %v\n", stderr)
		}
	} else {
		fmt.Println("No limit Transcoding")
		switch encoder {
		case "libx264":
			command = exec.Command("ffmpeg", "-hide_banner", "-i", CCME.EditVideo.VideoPath, "-c:v", encoder, "-preset", CCME.EditVideo.EncodingSpeed, "-crf", crfRate,
				"-vf", fmt.Sprintf("scale=%s:%s", CCME.EditVideo.Width, CCME.EditVideo.Height), "-pix_fmt", "yuv420p", "-c:a", "libfdk_aac", fmt.Sprintf("%s%s.%s", CCME.Paths[0], CCME.EditVideo.OrgVidName, fileType), "-y")
		case "libx265":
			command = exec.Command("ffmpeg", "-hide_banner", "-i", CCME.EditVideo.VideoPath, "-c:v", encoder, "-preset", CCME.EditVideo.EncodingSpeed, "-crf", crfRate,
				"-vf", fmt.Sprintf("scale=%s:%s", CCME.EditVideo.Width, CCME.EditVideo.Height), "-pix_fmt", "yuv420p", "-c:a", "libfdk_aac", fmt.Sprintf("%s%s.%s", CCME.Paths[0], CCME.EditVideo.OrgVidName, fileType), "-y")
		case "libvpx-vp9":
			command = exec.Command("ffmpeg", "-hide_banner", "-i", CCME.EditVideo.VideoPath, "-c:v", encoder, "-lossless", webmLossless, "-crf", crfRate,
				"-deadline", webmEncodeRate, "-vf", fmt.Sprintf("scale=%s:%s", CCME.EditVideo.Width, CCME.EditVideo.Height), "-pix_fmt",
				"yuv420p", "-c:a", "libvorbis", fmt.Sprintf("%s%s.%s", CCME.Paths[0], CCME.EditVideo.OrgVidName, fileType), "-y")
		case "libvpx":
			command = exec.Command("ffmpeg", "-hide_banner", "-i", CCME.EditVideo.VideoPath, "-c:v", encoder, "-lossless", webmLossless, "-crf", crfRate,
				"-deadline", webmEncodeRate, "-vf", fmt.Sprintf("scale=%s:%s", CCME.EditVideo.Width, CCME.EditVideo.Height), "-pix_fmt",
				"yuv420p", "-c:a", "libvorbis", fmt.Sprintf("%s%s.%s", CCME.Paths[0], CCME.EditVideo.OrgVidName, fileType), "-y")
		default:
			command = exec.Command("ffmpeg", "-hide_banner", "-i", "-vf", fmt.Sprintf("scale=%s:%s", CCME.EditVideo.Width, CCME.EditVideo.Height),
				"-b:v", CCME.EditVideo.MaxBV, CCME.EditVideo.VideoPath, fmt.Sprintf("%s%s.%s", CCME.Paths[0], CCME.EditVideo.OrgVidName, fileType), "-y")
		}
		stdout, stderr = command.CombinedOutput()
		if stderr != nil {
			fmt.Printf("ERROR: %v\n", stderr)
		}
	}
	fmt.Println("stdout: " + BytesToString(stdout))
	fmt.Println(encoder + " " + fileType)
	fmt.Println("Width: " + CCME.EditVideo.Width + " Height: " + CCME.EditVideo.Height)
	fmt.Println("Webm encode rate: " + webmEncodeRate + " crfRate: " + crfRate + " webmLossless: " + webmLossless)
	fmt.Println(EstimateBitrate(CCME.EditVideo.VideoPath))
	c <- "Done"
}
