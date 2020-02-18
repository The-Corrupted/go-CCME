package CCME

import (
	"fmt"
	"github.com/The-Corrupted/go-CCME/Helpers/OsHandler"
	"io"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

func DownloadFile(APIFunction string, r *http.Request) []interface{} {
	UserDir := OsHandler.SetUserDir()
	data := make([]interface{}, 0)
	switch APIFunction {
	case "UploadFile":
		fmt.Println("UploadFile")
		fmt.Println("Parsing multipart form.")
		UserDir := OsHandler.SetUserDir()
		err := r.ParseMultipartForm(18446744079)
		if err != nil {
			fmt.Printf("ERROR: %v", err)
		}
		fmt.Println("File parsed.")
		fmt.Printf("%v", r)
		//Get File, Header and err from Parsed Multiform file
		file, header, err := r.FormFile("UploadVideo")
		if err != nil {
			fmt.Printf("%v", err)
		}
		fmt.Println("File and header retrieved.")
		//Get File Type Based on File's first 512 bytes
		buffer := make([]byte, 512)
		_, err = file.Read(buffer)
		if err != nil {
			fmt.Printf("%v", err)
		}
		fmt.Println("File type found.")
		//Construct File Name
		file.Seek(0, 0)
		split := strings.Split(http.DetectContentType(buffer), "/")
		name := strings.Split(header.Filename, ".")
		fullFile := fmt.Sprintf("%s.%s", name[0], split[1])
		//Check For Proper File Type
		var fileerror uint8 = 0
		if split[0] != "video" {
			fmt.Println("Recieved Wrong File Type")
			fileerror = 1
		}
		fmt.Println(fullFile)

		//--------------------Save File To Disk If File Error and Error hasn't occured-------------------

		if err == nil && fileerror == 0 {
			out, err := os.Create(fmt.Sprintf("%s/Videos/UnderlayVids/%s", UserDir, fullFile))
			if err != nil {
				fmt.Printf("%v", err)
			}
			defer out.Close()
			n, err := io.Copy(out, file)
			if fileerror != 0 {
				fmt.Printf("ERROR DOWNLOADING: %v", err)
			}
			fmt.Println(n, "Bytes Downloaded")
			data = append(data, "Upload Success")
			data = append(data, int64(n))
			data = append(data, fullFile)
			data = append(data, fmt.Sprintf("%s/%s", split[0], split[1]))
			return data
		} else {
			data = append(data, "Upload Failed!")
			data = append(data, int64(0))
			data = append(data, "NULL")
			data = append(data, "NULL")
			return data
		}
		data = append(data, "Upload Failed!")
		data = append(data, 0)
		data = append(data, "NULL")
		data = append(data, "NULL")
		return data
	case "EditVideo":
		r.ParseMultipartForm(0)
		file, header, err := r.FormFile("UploadVideo")
		if err != nil {
			fmt.Printf("%v\n", err)
		}
		fmt.Println("Got Video")
		buffer := make([]byte, 512)
		_, err = file.Read(buffer)
		if err != nil {
			fmt.Printf("%v\n", err)
		}
		fmt.Println("Construct video name.")
		name := strings.Split(header.Filename, ".")
		fmt.Println(name)
		file.Seek(0, 0)
		split := strings.Split(http.DetectContentType(buffer), "/")
		fullFile := fmt.Sprintf("%s.%s", name[0], name[1])
		fullFile = strings.Trim(fullFile, "\n")
		fmt.Println(fullFile)
		fmt.Println(split)
		fmt.Println("Checking file type.")
		var fileerror uint8 = 0
		if name[1] != "mp4" && name[1] != "webm" && name[1] != "ts" {
			fmt.Println("Recieved Wrong File Type")
			fileerror = 1
		}
		//Make Video Path and Save Video
		fmt.Println("Saving video")
		NewDir := fmt.Sprintf("%s/TempFiles/%s/", UserDir, name[0])
		os.MkdirAll(NewDir, os.ModePerm)
		VideoPath := fmt.Sprintf("%s%s", NewDir, fullFile)
		if err == nil && fileerror == 0 {
			out, err := os.Create(VideoPath)
			if err != nil {
				fmt.Printf("%v", err)
			}
			defer out.Close()
			_, err = io.Copy(out, file)
			if fileerror != 0 && err != nil {
				fmt.Printf("ERROR DOWNLOADING: %v", err)
			}
		}
		data = append(data, VideoPath)
		data = append(data, name[0])
		return data
	case "QuickAnalyze":
		fmt.Println("QuickAnalyze")
		r.ParseMultipartForm(0)
		fmt.Println(r.FormValue("UploadVideo"))
		file, header, err := r.FormFile("UploadVideo")
		if err != nil {
			fmt.Printf("%v\n", err)
		}
		fmt.Println("Got video")
		buffer := make([]byte, 512)
		_, err = file.Read(buffer)
		if err != nil {
			fmt.Printf("%v\n", err)
		}
		//Construct File Name
		fmt.Println("Contructing video file name.")
		name := strings.Split(header.Filename, ".")
		fmt.Println(name)
		file.Seek(0, 0)
		split := strings.Split(http.DetectContentType(buffer), "/")
		fullFile := fmt.Sprintf("%s.%s", name[0], name[1])
		fullFile = strings.Trim(fullFile, "\n")
		fmt.Println(fullFile)
		fmt.Println(split)

		//Check For Proper File Type
		fmt.Println("Checking file type.")
		var fileerror uint8 = 0
		if name[1] != "mp4" && name[1] != "webm" && name[1] != "ts" {
			fmt.Println("Recieved Wrong File Type")
			fileerror = 1
		}
		//Make Video Path and Save Video
		fmt.Println("Saving video")
		VideoPath := fmt.Sprintf("%s/Videos/StampedVideos/%s", UserDir, fullFile)
		if err == nil && fileerror == 0 {
			out, err := os.Create(VideoPath)
			if err != nil {
				fmt.Printf("%v", err)
			}
			defer out.Close()
			_, err = io.Copy(out, file)
			if fileerror != 0 && err != nil {
				fmt.Printf("ERROR DOWNLOADING: %v", err)
			}
		}
		DeconPath := fmt.Sprintf("%s/Pictures/EDA/Frames/%s/", UserDir, name[0])
		data = append(data, DeconPath)
		data = append(data, name)
		data = append(data, fullFile)
		data = append(data, VideoPath)
		return data
	}
	return data
}

// ---------------------------------------CCME Helpers------------------------------------------

func (TimeCounter *TimeCounter) IncrementCounter() {
	TimeCounter.TimeValues[2] += 1
	//Deal With Seconds Overflow
	if TimeCounter.TimeValues[2] == 60 {
		TimeCounter.TimeValues[1] += 1
		TimeCounter.TimeValues[2] = 0
		//Deal With Minutes Overflow
		if TimeCounter.TimeValues[1] == 60 {
			TimeCounter.TimeValues[0] += 1
			TimeCounter.TimeValues[1] = 0
		}
	}
	TimeCounter.ZPaddedTime = fmt.Sprintf("%02d:%02d:%02d", TimeCounter.TimeValues[0], TimeCounter.TimeValues[1], TimeCounter.TimeValues[2])
}

func GetTimeInSeconds(formatedTime string) uint64 {
	TimeSplitters := strings.Split(formatedTime, ":")
	Seconds, _ := strconv.ParseUint(TimeSplitters[1], 10, 32)
	Minutes, _ := strconv.ParseUint(TimeSplitters[0], 10, 32)
	Seconds = uint64(Seconds) + (uint64(Minutes) * 60)
	return Seconds
}

func SplitAndConvert(formatedFR string) float64 {
	split := strings.Split(formatedFR, "/")
	Numerator, _ := strconv.ParseFloat(split[0], 64)
	Denominator, _ := strconv.ParseFloat(split[1], 64)
	var Result = Numerator / Denominator
	return Result
}

func GetAudioAndFramerate(VideoName string) []string {
	Slice := make([]string, 2)
	var Format = ""
	fmt.Printf("")
	AudioFormat := regexp.MustCompile("(?m)Audio:\\s\\w{1,}")
	Fps := regexp.MustCompile("(?m)fps,\\s\\d{2,}\\.\\d{2,}|(?m)fps,\\s\\d{2,}")
	cmdAudioType := exec.Command("ffmpeg", "-i", VideoName)
	Output, Error := cmdAudioType.CombinedOutput()
	if Error != nil {
		fmt.Println("ERROR: ", Error)
	}
	NonByteOut := BytesToString(Output)
	fmt.Println(AudioFormat.MatchString(NonByteOut))
	fmt.Println(Fps.MatchString(NonByteOut))

	AudioFormatFound := AudioFormat.FindStringSubmatch(NonByteOut)
	FpsFound := Fps.FindStringSubmatch(NonByteOut)
	if len(AudioFormatFound) != 0 {
		if AudioFormatFound[0] != "" {
			Slice[0] = strings.Split(AudioFormatFound[0], ": ")[1]
		} else {
			fmt.Println(Format)
			Slice[0] = "None"
		}
	} else {
		Slice[0] = "None"
	}
	if FpsFound[0] != "" {
		Slice[1] = strings.Split(FpsFound[0], ", ")[1]
	} else {
		Slice[1] = "None"
	}

	switch Slice[0] {
	case "vorbis":
		Slice[0] = "ogg"
	case "aac":
		Slice[0] = "aac"
	}
	fmt.Printf("%v\n", Slice)
	return Slice
}

func (CCME *CCMEClass) MakeDirs() {
	var i = 0
	for i < len(CCME.Paths) {
		if _, err := os.Stat(CCME.Paths[i]); os.IsNotExist(err) {
			os.MkdirAll(CCME.Paths[i], os.ModePerm)
		}
		i++
	}
}

func BytesToString(byteArray []byte) string {
	return strings.Trim(string(byteArray[:]), "\n")
}

func GetVideoCodec(s string, extension string) string {
	if extension == "ts" {
		return "mpeg2video"
	}
	cmd := exec.Command("ffprobe", "-v", "error", "-select_streams", "v:0", "-show_entries",
		"stream=codec_name", "-of", "default=noprint_wrappers=1:nokey=1", s)

	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Sprintf("%v: FAILED", err)
	}
	n := len(out)
	byteString := string(out[:n-1])
	return byteString
}

func (CCME *CCMEClass) CheckFunction(funct string) bool {
	if funct == CCME.Function {
		return true
	}
	return false
}

func SelectBV(device string) string {
	Terriformity := strings.Split(device, "|")
	var devicesType, quality = Terriformity[0], Terriformity[1]
	if devicesType == "None" {
		return quality
	}
	if quality == "4KV" {
		return HBV4K
	} else {
		if devicesType == "4KDev" {
			return HBVHD
		} else if devicesType == "HDDev" {
			return BVHD
		} else if devicesType == "HD600" {
			fmt.Println(BVHD600)
			return BVHD600
		} else if devicesType == "HD2000" {
			return BVHD2000
		}
	}
	return "0"
}

func SetBufSize(maxbv string) string {
	var bvstring = strings.Trim(maxbv, "k")
	bv, _ := strconv.ParseFloat(bvstring, 64)
	var bufsize = bv * 0.8
	var bufsizeInt = int64(bufsize)
	var bufstring = strconv.FormatInt(bufsizeInt, 10)
	bufstring = bufstring + "k"
	fmt.Println(bufstring)
	return bufstring
}

func EstimateBitrate(video string) (string, error) {
	cmdSize := exec.Command("ffprobe", "-hide_banner", "-v", "error", "-show_entries", "format=size", "-of",
		"default=noprint_wrappers=1:nokey=1", video)
	cmdTime := exec.Command("ffprobe", "-hide_banner", "-v", "error", "-show_entries", "format=duration", "-of",
		"default=noprint_wrappers=1:nokey=1", video)

	byteSize, stderr := cmdSize.CombinedOutput()
	if stderr != nil {
		return "FAILED", stderr
	}
	byteTime, stderr := cmdTime.CombinedOutput()
	if stderr != nil {
		return "FAILED", stderr
	}
	//	fmt.Println(BytesToString(byteSize) + " " + BytesToString(byteTime))
	Size, err := strconv.ParseFloat(BytesToString(byteSize), 64)
	if err != nil {
		return "FAIL", err
	}
	Time, err := strconv.ParseFloat(BytesToString(byteTime), 64)
	if err != nil {
		return "FAIL", err
	}
	var bitrate float64 = 0
	Buf := Size * 8 / 1000
	bitrate = Buf / Time
	bitstring := strconv.FormatFloat(bitrate, 'f', 0, 64)
	fmt.Println(bitstring)
	fmt.Println(bitstring + "k")
	return bitstring + "k", nil
}

func GetVideoProfile(video string) (string, error) {
	cmdProfile := exec.Command("ffprobe", "-hide_banner", "-v", "error", "-select_streams", "v:0",
		"-show_entries", "stream=profile", "-of", "default=noprint_wrappers=1:nokey=1", video)

	byteProfile, stderr := cmdProfile.CombinedOutput()
	if stderr != nil {
		return "FAILED", stderr
	}
	return BytesToString(byteProfile), nil
}

func GetVideoPixelFormat(video string) (string, error) {
	cmdPixelValue := exec.Command("ffprobe", "-hide_banner", "-v", "-error", "-select_streams", "v:0",
		"-show_entries", "stream=pix_fmt", "-of", "default=noprint_wrappers=1:nokey=1", video)

	bytePixVal, stderr := cmdPixelValue.CombinedOutput()
	if stderr != nil {
		return "FAILED", stderr
	}
	return BytesToString(bytePixVal), nil
}

// func ReadImagePixels(imagePath string) [][]PixelValue {
// 	HorizontalPixels := make([]PixelValue, 0)

// }
