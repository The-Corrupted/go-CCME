package OsHandler

import (
	"fmt"
	"os"
	"runtime"
)

func SetUserDir() string {
	if runtime.GOOS == "windows" {
		home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		if home == "" {
			home = os.Getenv("USERPROFILE")
		}
		return home
	}
	return os.Getenv("HOME")
}

func MakeTempDirComp(videoName string) string {
	TempDir := fmt.Sprintf("%s/TempFiles/Temp%s/", SetUserDir(), videoName)
	os.MkdirAll(TempDir, os.ModePerm)
	return TempDir
}

func MakeTempDirOverlay(VideoName string) []string {
	slice := make([]string, 2)
	TempDir := fmt.Sprintf("%s/TempFiles/Temp%s/", SetUserDir(), VideoName)
	slice[0] = TempDir
	os.MkdirAll(TempDir, os.ModePerm)
	TempDir2 := fmt.Sprintf("%s/TempFiles/Temp%s/OverlayResults/", SetUserDir(), VideoName)
	slice[1] = TempDir2
	os.MkdirAll(TempDir2, os.ModePerm)
	return slice
}
