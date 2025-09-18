package processors

import (
	"math"

	imatix "github.com/J0hnLenin/ComputerVision/imatrix"
)

type ImageProccessor func(imatix.Image)

var Functions = map[string]ImageProccessor{
	"green": Green,
	"red":   Red,
	"blue":  Blue,
	"gray":  GrayScale,
}

func Red(img imatix.Image) {
	for i := 0; i < img.Height; i++ {
		for j := 0; j < img.Width; j++ {
			img.Matrix[i][j][1] = 0
			img.Matrix[i][j][2] = 0
		}
	}
}

func Green(img imatix.Image) {
	for i := 0; i < img.Height; i++ {
		for j := 0; j < img.Width; j++ {
			img.Matrix[i][j][0] = 0
			img.Matrix[i][j][2] = 0
		}
	}
}

func Blue(img imatix.Image) {
	for i := 0; i < img.Height; i++ {
		for j := 0; j < img.Width; j++ {
			img.Matrix[i][j][0] = 0
			img.Matrix[i][j][1] = 0
		}
	}
}

func GrayScale(img imatix.Image) {
	for i := 0; i < img.Height; i++ {
		for j := 0; j < img.Width; j++ {
			GrayColor := uint8(float32(img.Matrix[i][j][0])*0.299) +
				uint8(float32(img.Matrix[i][j][1])*0.587) +
				uint8(float32(img.Matrix[i][j][2])*0.114)

			img.Matrix[i][j][0] = GrayColor
			img.Matrix[i][j][1] = GrayColor
			img.Matrix[i][j][2] = GrayColor
		}
	}
}

func ChangeBrightness(img imatix.Image, color int, brightness float64) {
	for i := 0; i < img.Height; i++ {
		for j := 0; j < img.Width; j++ {
			value := float64(img.Matrix[i][j][color]) / 255
			value = float64(math.Pow(value, brightness) * 255)
			img.Matrix[i][j][color] = transform(value)
		}
	}
}

func Negative(img imatix.Image) {
	for i := 0; i < img.Height; i++ {
		for j := 0; j < img.Width; j++ {
			for k := 0; k < 3; k++ {
				img.Matrix[i][j][k] = 255 - img.Matrix[i][j][k]
			}
		}
	}
}

func VertivalMirror(img imatix.Image) {
	for i := 0; i < img.Height/2; i++ {
		for j := 0; j < img.Width; j++ {
			for k := 0; k < 3; k++ {
				v := img.Matrix[i][j][k]
				img.Matrix[i][j][k] = img.Matrix[img.Height-i-1][j][k]
				img.Matrix[img.Height-i-1][j][k] = v
			}
		}
	}
}

func HorisontalMirror(img imatix.Image) {
	for i := 0; i < img.Height; i++ {
		for j := 0; j < img.Width/2; j++ {
			for k := 0; k < 3; k++ {
				v := img.Matrix[i][j][k]
				img.Matrix[i][j][k] = img.Matrix[i][img.Width-j-1][k]
				img.Matrix[i][img.Width-j-1][k] = v
			}
		}
	}
}

func Magic(img imatix.Image, t int) {
	for i := 0; i < img.Height; i++ {
		for j := 0; j < img.Width; j++ {
			for k := 0; k < 3; k++ {
				if int(img.Matrix[i][j][k])+t > 255 || int(img.Matrix[i][j][k])-t < 0 {
					img.Matrix[i][j][k] = 255 - img.Matrix[i][j][k]
				}
				img.Matrix[i][j][k] = 255 - img.Matrix[i][j][k]
			}
		}
	}
}

func transform(v float64) uint8 {
	if v > 255 {
		v = 255.0
	} else if v < 0 {
		v = 0.0
	}
	return uint8(v)
}
