package imatix

type Image struct {
	Matrix [][][3]uint8
	Height int
	Width  int
}

type Parameters struct {
	RedBrightness    float64
	GreenBrightness  float64
	BlueBrightness   float64
	Contrast         float64
	Negative         bool
	Order            string // RGB RBG GRB
	VertivalMirror   bool
	HorisontalMirror bool
	Magic            int
}
