package redactor

import (
	imatix "github.com/J0hnLenin/ComputerVision/imatrix"
	"github.com/J0hnLenin/ComputerVision/processors"
)

func Redact(inputImage imatix.Image, parameters imatix.Parameters) (imatix.Image, imatix.Image) {

	originalCopy := processors.CopyImage(inputImage)

	if parameters.Order != "RGB" && parameters.Order != "" {
		inputImage = processors.ChangeOrder(inputImage, parameters.Order)
	}
	if parameters.RedBrightness != 1.0 {
		processors.ChangeBrightness(inputImage, 0, parameters.RedBrightness)
	}
	if parameters.GreenBrightness != 1.0 {
		processors.ChangeBrightness(inputImage, 1, parameters.GreenBrightness)
	}
	if parameters.BlueBrightness != 1.0 {
		processors.ChangeBrightness(inputImage, 2, parameters.BlueBrightness)
	}
	if parameters.Contrast != 1.0 {
		processors.ChangeContrast(inputImage, parameters.Contrast)
	}
	if parameters.Negative {
		processors.Negative(inputImage)
	}
	if parameters.VerticalMirror {
		processors.VerticalMirror(inputImage)
	}
	if parameters.HorizontalMirror {
		processors.HorizontalMirror(inputImage)
	}
	if parameters.Magic != 0.0 {
		processors.Magic(inputImage, parameters.Magic)
	}

	temp := processors.CopyImage(inputImage)
	bluredImage := temp

	if parameters.Filter == "gaussian" && parameters.Sigma > 0 {
		bluredImage = processors.GaussianFilter(inputImage, parameters.FilterSize, parameters.Sigma)
	}
	if parameters.Filter == "sigma" && parameters.Sigma > 0 {
		bluredImage = processors.SigmaFilter(inputImage, parameters.FilterSize, parameters.Sigma, float64(parameters.Interval))
	}
	if parameters.Filter == "median" {
		bluredImage = processors.MedianFilter(inputImage, parameters.FilterSize)
	}
	if parameters.Filter == "rectangular" {
		bluredImage = processors.RectangularFilter(inputImage, parameters.FilterSize)
	}

	if parameters.UnsharpMasking != 0.0 {
		inputImage = processors.UnsharpMasking(temp, bluredImage, parameters.UnsharpMasking)
	} else {
		inputImage = bluredImage
	}
	changes := processors.Changes(originalCopy, inputImage)
	// ApplyCore - новая функция чтобы избежать дублирования.
	// Можно код отрефакторить, прическать.
	// Новые обработчики добавлять в этом формате.
	//if parameters.LogarithmicBrightness != 1.0 {
	//	processors.ApplyCore(inputImage, "LogarithmicBrightness", parameters.LogarithmicBrightness)
	//}
	return changes, inputImage
}
