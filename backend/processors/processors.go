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

var core = map[string]imatix.ProccessorCore{
	"logarithmicBrightness": logarithmicBrightnessCore,
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

func ChangeContrast(img imatix.Image, contrast float64) {
	for i := 0; i < img.Height; i++ {
		for j := 0; j < img.Width; j++ {
			for k := 0; k < 3; k++ {
				if contrast > 1.0 {
					gamma := 30.0
					oldValue := float64(img.Matrix[i][j][k]) / 255
					newValue := 1.0 / (1.0 + math.Exp(gamma*(0.5-oldValue)))
					average := ((contrast-1.0)*newValue + (2.0-contrast)*oldValue)
					img.Matrix[i][j][k] = transform(average * 255.0)
				} else if contrast < 1.0 {
					gamma := 1 / 30.0
					oldValue := float64(img.Matrix[i][j][k]) / 255
					newValue := 1.0 / (1.0 + math.Exp(gamma*(0.5-oldValue)))
					average := ((1.0-contrast)*newValue + (contrast)*oldValue)
					img.Matrix[i][j][k] = transform(average * 255.0)
				}
			}

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

func VerticalMirror(img imatix.Image) {
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

func HorizontalMirror(img imatix.Image) {
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
				if int(img.Matrix[i][j][k])+t <= 255 && int(img.Matrix[i][j][k])-t >= 0 {
					img.Matrix[i][j][k] = 255 - img.Matrix[i][j][k]
				}
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

func ChangeOrder(img imatix.Image, order string) imatix.Image {
	newImage := imatix.Image{
		Matrix: make([][][3]uint8, img.Height),
		Height: img.Height,
		Width:  img.Width,
	}

	channelMap := make([]int, 3)
	for i, char := range order {
		switch char {
		case 'R':
			channelMap[i] = 0
		case 'G':
			channelMap[i] = 1
		case 'B':
			channelMap[i] = 2
		}
	}

	for y := 0; y < img.Height; y++ {
		newImage.Matrix[y] = make([][3]uint8, img.Width)
		for x := 0; x < img.Width; x++ {
			original := img.Matrix[y][x]
			newImage.Matrix[y][x] = [3]uint8{
				original[channelMap[0]],
				original[channelMap[1]],
				original[channelMap[2]],
			}
		}
	}

	return newImage
}

func ApplyCore(img imatix.Image, parameterName string, parameterValue float64) {
	img.Apply(core[parameterName], parameterValue)
}

func logarithmicBrightnessCore(value float64, c float64) float64 {
	return c * math.Log(1.0+value)
}
