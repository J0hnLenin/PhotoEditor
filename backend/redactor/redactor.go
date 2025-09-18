package redactor

import (
	imatix "github.com/J0hnLenin/ComputerVision/imatrix"
	"github.com/J0hnLenin/ComputerVision/processors"
)

func Redact(imputImage imatix.Image, parameters imatix.Parameters) imatix.Image {
	if parameters.RedBrightness != 1.0 {
		processors.ChangeBrightness(imputImage, 0, parameters.RedBrightness)
	}
	if parameters.GreenBrightness != 1.0 {
		processors.ChangeBrightness(imputImage, 1, parameters.RedBrightness)
	}
	if parameters.BlueBrightness != 1.0 {
		processors.ChangeBrightness(imputImage, 1, parameters.RedBrightness)
	}
	if parameters.Negative {
		processors.Negative(imputImage)
	}
	if parameters.VertivalMirror {
		processors.VertivalMirror(imputImage)
	}
	if parameters.HorisontalMirror {
		processors.HorisontalMirror(imputImage)
	}
	if parameters.Magic != 0.0 {
		processors.Magic(imputImage, parameters.Magic)
	}
	return imputImage
}
