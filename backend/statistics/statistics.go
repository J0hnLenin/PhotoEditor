package statistics

import (
	imatix "github.com/J0hnLenin/ComputerVision/imatrix"
)

type BrightnessHistogram struct {
	Red   [256]uint16
	Green [256]uint16
	Blue  [256]uint16
	Gray  [256]uint16
}

type ImageStatistics struct {
	Brightness BrightnessHistogram
}

func GetBrightnessHistogram(img imatix.Image) BrightnessHistogram {
	hist := BrightnessHistogram{}

	for i := 0; i < img.Height; i++ {
		for j := 0; j < img.Width; j++ {
			redValue := img.Matrix[i][j][0]
			greenValue := img.Matrix[i][j][1]
			blueValue := img.Matrix[i][j][2]

			grayValue := uint8(float32(redValue)*0.299) +
				uint8(float32(greenValue)*0.587) +
				uint8(float32(blueValue)*0.114)

			hist.Red[redValue]++
			hist.Green[greenValue]++
			hist.Blue[blueValue]++
			hist.Gray[grayValue]++
		}
	}

	return hist
}

func GetStatistics(img imatix.Image) ImageStatistics {
	return ImageStatistics{
		Brightness: GetBrightnessHistogram(img),
	}
}
