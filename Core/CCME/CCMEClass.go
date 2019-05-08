package CCME

type CCMEClass struct {
	//structs to use for different class types
	Function string
	Create
	Overlay
	GetFrames
	EditVideo
	DeleteDownload
	Paths []string
}

type Create struct {
	PhotoName string
	GenNumber uint64
	FPS float64
	FormattedTime string
}

type Overlay struct {
	UnderlayVideo string
	OverlayName string
	OverlayPosition string
	SaveOverlay string
	FinalVideoName string
}

type GetFrames struct {
	VideoName string
}

type EditVideo struct {
	VideoPath string
	NewFormat string
	Quality string
	EncodingSpeed string
	Width string
	Height string
	OrgVidName string
	MaxBV string
}

type DeleteDownload struct {
	DeleteOrDownload string
	File string
}

type TimeCounter struct {
	TimeValues []uint8
	ZPaddedTime string
}