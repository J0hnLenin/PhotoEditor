package imatix

type Image struct {
	Matrix [][][3]uint8
	Height int
	Width  int
}

type ProccessorCore func(float64, float64) float64

func (img Image) Apply(core ProccessorCore, parameter float64) {
	for i := 0; i < img.Height; i++ {
		for j := 0; j < img.Width; j++ {
			for k := 0; k < 3; k++ {
				continuousValue := PixelToContinuous(img.Matrix[i][j][k])
				img.Matrix[i][j][k] = ContinuousToPixel(core(continuousValue, parameter))
			}
		}
	}
}

func PixelToContinuous(value uint8) float64 {
	return float64(value) / 255.0
}

func ContinuousToPixel(value float64) uint8 {
	return uint8(value * 255.0)
}

type Parameters struct {
	RedBrightness         float64
	GreenBrightness       float64
	BlueBrightness        float64
	LogarithmicBrightness float64
	Contrast              float64
	Negative              bool
	Order                 string // RGB RBG GRB
	VerticalMirror        bool
	HorizontalMirror      bool
	Magic                 int
	Filter                string
	FilterSize            int
	Sigma                 float64
	Interval              int
	UnsharpMasking        float64
	LogarithmicClip 	  bool
	PowerClip 			  float64
	BinaryClip			  int
	ConstantLow			  int
	ConstantHigh		  int
	ConstantValue		  int
}
