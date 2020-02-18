package CCMEHandlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	_ "io/ioutil"
	"log"
	"net/http"
	_ "strconv"
)

func ParseJson(APIFunction string, r *http.Request) ArgumentLink {
	var buf bytes.Buffer
	Logger := log.New(&buf, "logger: ", log.Llongfile)
	r.ParseForm()
	var args ArgumentLink
	args.Values = make(map[string]*string)
	switch APIFunction {
	case "Deconstruct":
		type Storage struct {
			Video, SaveName string
		}
		var s Storage
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&s)
		if err != nil {
			fmt.Println("Failed to decode json data")
			return args
		}
		fmt.Printf("Struct: %v\n", s)
		args.Name = APIFunction
		args.Values["DeconstructVideoName"] = &s.Video
		args.Values["DeconstructFrameNames"] = &s.SaveName
		return args
	case "Create":
		/*
			This uses a map of strings as compared to
			the structs used for others. It is preferable to do it
			this way to avoid parsing r twice, and allows for less confusing
			flow control.
			When possible, use structs instead.
		*/
		data := make(map[string]string)
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&data)
		if err != nil {
			fmt.Printf("ERROR IN CREATE: %v\n", err)
			return args
		}
		args.Name = data["Type"]
		if args.Name == "CreateFixed" {
			formDat := [...]string{data["FrameSolidName"], data["FPSSolid"],
				data["FramesSolid"]}
			args.Values["CreateVideoName"] = &formDat[0]       //data["FrameSolidName"]
			args.Values["CreateFrameNumberRate"] = &formDat[1] //data["FPSSolid"]
			args.Values["CreateFrameNumber"] = &formDat[2]     //data["FramesSolid"]
		} else {
			formDat := [...]string{data["calculatedVidLength"], data["calculatedFPS"],
				data["FrameTimedName"]}
			args.Values["CreateVideoTime"] = &formDat[0] //data["calculatedVidLength"]
			args.Values["CreateFrameRate"] = &formDat[1] //data["calculatedFPS"]
			args.Values["CreateVideoName"] = &formDat[2] //data["FrameTimedName"]
		}
		return args
	case "Overlay":
		type Storage struct {
			OverlayFrames string
			UnderlayVideo string
			SaveOverlay   string
			VideoName     string
		}
		var s Storage
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&s)
		if err != nil {
			fmt.Println("Failed to decode json data")
			return args
		}
		args.Name = APIFunction
		args.Values["UnderlayVideo"] = &s.UnderlayVideo
		args.Values["OverlayVideo"] = &s.OverlayFrames
		args.Values["OverlayFinalVideoName"] = &s.VideoName
		args.Values["SaveOverlay"] = &s.SaveOverlay
		return args
	case "GetVidInfo":
		type Storage struct {
			VideoName string
		}
		var s Storage
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&s)
		if err != nil {
			fmt.Println("Failed to decode json data")
			return args
		}
		args.Name = APIFunction
		args.Values["Video"] = &s.VideoName
		Logger.Printf("args.Values %v", s)
		fmt.Println(&buf)
	case "Analyze":
		type Storage struct {
			FramesName        string
			ExpNumberOfFrames string
			Delete            string
		}
		var s Storage
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&s)
		if err != nil {
			fmt.Println("Failed to decode json data.")
			return args
		}
		args.Name = APIFunction
		args.Values["AnalyzeFrameNames"] = &s.FramesName
		args.Values["AnalyzeExpFrames"] = &s.ExpNumberOfFrames
		args.Values["AnalyzeDelete"] = &s.Delete
	case "FullOutput":
		type Storage struct {
			Id string
		}
		var s Storage
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&s)
		if err != nil {
			fmt.Println("Failed to decode json data.")
			return args
		}
		args.Values["id"] = &s.Id
	}
	return args
}
