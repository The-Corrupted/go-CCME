package CCMEHandlers

import (
	"fmt"
	"net/http"
	"strconv"
)

type ArgumentLink struct {
	Name string
	Values map[string]*string
}

func ParseForm(APIFunction string, r *http.Request)  ArgumentLink {
	r.ParseForm()
	var args ArgumentLink
	args.Values = make(map[string]*string)
	switch APIFunction {
	case "Create":
		timedFR, _ := strconv.ParseInt(r.Form["FR"][0], 10, 64)
		if timedFR <= 0  {
			args.Name = "CreateFixed"
			args.Values["CreateVideoName"] = &r.Form["QRName"][0]
			args.Values["CreateFrameNumber"] = &r.Form["Frames"][0]
			args.Values["CreateFrameNumberRate"] = &r.Form["videoFrameRate"][0]
			fmt.Printf("FR: %v\n", r.Form["FR"])
		} else {
			args.Name = "CreateTimed"
			args.Values["CreateVideoTime"] = &r.Form["lengthOfVid"][0]
			args.Values["CreateFrameRate"] = &r.Form["FR"][0]
			args.Values["CreateVideoName"] = &r.Form["timedVidName"][0]
		}
		return args
	case "Compact":
		args.Name = APIFunction
		args.Values["CompactName"] = &r.Form["FrameNames"][0]
		args.Values["CompactFormat"] = &r.Form["VidType"][0]
		args.Values["CompactFrameRate"] = &r.Form["Framerate"][0]
		args.Values["CompactVideoName"] = &r.Form["videoName"][0]
		return args
	case "Overlay":
		args.Name = APIFunction
		args.Values["UnderlayVideo"] = &r.Form["underlayVideo"][0]
		args.Values["OverlayVideo"] = &r.Form["overlayVideo"][0]
		args.Values["SaveOverlay"] = &r.Form["SaveOverlay"][0]
		args.Values["OverlayFinalVideoName"] = &r.Form["genVidName"][0]
		return args
	case "Deconstruct":
		fmt.Printf("Form: %v", r.Form)
		args.Name = APIFunction
		args.Values["DeconstructVideoName"] = &r.Form["deconstructVid"][0]
		args.Values["DeconstructFrameNames"] = &r.Form["DeconFrameNames"][0]
		return args
	case "Analyze":
		args.Name = APIFunction
		args.Values["AnalyzeFrameNames"] = &r.Form["analyzeFrameNames"][0]
		args.Values["AnalyzeExpFrames"] = &r.Form["expectedFrames"][0]
		args.Values["AnalyzeDelete"] = &r.Form["delete"][0]
		return args
	case "GetVidInfo":
		args.Name = APIFunction
		args.Values["Video"] = &r.Form["videoToGetFR"][0]
		return args
	case "EditVideo":
		data := DownloadFile(APIFunction, r)
		var BV = SelectBV(r.FormValue("device"))
		var VideoPath = data[0].(string)
		var Name = data[1].(string)
		formdata := [...]string{r.FormValue("formatForConversion"),
							  r.FormValue("Quality"), r.FormValue("EncodingSpeed"),
							  r.FormValue("vidWidth"), r.FormValue("vidHeight")}
		args.Name = APIFunction
		args.Values["VideoPath"] = &VideoPath
		args.Values["VidFormat"] = &formdata[0]
		args.Values["Quality"] = &formdata[1]
		args.Values["EncodingSpeed"] = &formdata[2]
		args.Values["Width"] = &formdata[3]
		args.Values["Height"] = &formdata[4]
		args.Values["VideoName"] = &Name
		args.Values["BVRate"] = &BV
		return args
	case "DeleteDownload":
		args.Name = APIFunction
		args.Values["TransOption"] = &r.Form["transOption"][0]
		args.Values["DelDown"] = &r.Form["delDown"][0]
		return args
	}
	return args
}