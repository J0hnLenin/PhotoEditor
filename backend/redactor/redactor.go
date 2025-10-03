package redactor

import (
	imatix "github.com/J0hnLenin/ComputerVision/imatrix"
	"github.com/J0hnLenin/ComputerVision/processors"
)

func Redact(imputImage imatix.Image, parameters imatix.Parameters) imatix.Image {
	if parameters.Order != "RGB" && parameters.Order != "" {
		imputImage = processors.ChangeOrder(imputImage, parameters.Order)
	}
	if parameters.RedBrightness != 1.0 {
		processors.ChangeBrightness(imputImage, 0, parameters.RedBrightness)
	}
	if parameters.GreenBrightness != 1.0 {
		processors.ChangeBrightness(imputImage, 1, parameters.GreenBrightness)
	}
	if parameters.BlueBrightness != 1.0 {
		processors.ChangeBrightness(imputImage, 2, parameters.BlueBrightness)
	}
	if parameters.Contrast != 1.0 {
		processors.ChangeContrast(imputImage, parameters.Contrast)
	}
	if parameters.Negative {
		processors.Negative(imputImage)
	}
	if parameters.VerticalMirror {
		processors.VerticalMirror(imputImage)
	}
	if parameters.HorizontalMirror {
		processors.HorizontalMirror(imputImage)
	}
	if parameters.Magic != 0.0 {
		processors.Magic(imputImage, parameters.Magic)
	}
	// ApplyCore - новая функция чтобы избежать дублирования.
	// Можно код отрефакторить, прическать.
	// Новые обработчики добавлять в этом формате.
	if parameters.LogarithmicBrightness != 1.0 {
		processors.ApplyCore(imputImage, "LogarithmicBrightness", parameters.LogarithmicBrightness)
	}
	return imputImage
}
