package CCMEHandlers

import (
	"fmt"
	"github.com/The-Corrupted/go-CCME/Helpers/OsHandler"
	"io"
	"net/http"
	"os"
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
