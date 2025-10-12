package processors

import (
	"math"
	"sort"

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
			img.Matrix[i][j][color] = clip(value)
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
					img.Matrix[i][j][k] = clip(average * 255.0)
				} else if contrast < 1.0 {
					gamma := 1 / 30.0
					oldValue := float64(img.Matrix[i][j][k]) / 255
					newValue := 1.0 / (1.0 + math.Exp(gamma*(0.5-oldValue)))
					average := ((1.0-contrast)*newValue + (contrast)*oldValue)
					img.Matrix[i][j][k] = clip(average * 255.0)
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

func clip(v float64) uint8 {
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

func createGaussianKernel(size int, sigma float64) []float64 {
	kernel := make([]float64, size)
	radius := size / 2
	var sum float64

	for i := -radius; i <= radius; i++ {
		x := float64(i)
		value := math.Exp(-(x*x)/(2*sigma*sigma)) / (sigma * math.Sqrt(2*math.Pi))
		kernel[i+radius] = value
		sum += value
	}

	return kernel
}

// оптимизация для одномерных
func applyConvolution(src imatix.Image, dst imatix.Image, kernel []float64, horizontal bool) {
	radius := len(kernel) / 2

	for y := 0; y < src.Height; y++ {
		for x := 0; x < src.Width; x++ {
			var sumR, sumG, sumB float64
			var weightSum float64

			for k := -radius; k <= radius; k++ {
				var px, py int
				//Здесь направление одномерной свёртки
				if horizontal {

					px = x + k
					py = y
				} else {

					px = x
					py = y + k
				}

				if px < 0 || px >= src.Width || py < 0 || py >= src.Height {
					continue
				}

				weight := kernel[k+radius]
				pixel := src.Matrix[py][px]

				sumR += float64(pixel[0]) * weight
				sumG += float64(pixel[1]) * weight
				sumB += float64(pixel[2]) * weight
				weightSum += weight
			}

			if weightSum > 0 {
				dst.Matrix[y][x] = [3]uint8{
					clip(sumR / weightSum),
					clip(sumG / weightSum),
					clip(sumB / weightSum),
				}
			}
		}
	}
}

func GaussianFilter(img imatix.Image, size int, sigma float64) imatix.Image {
	if sigma <= 0 {
		return img
	}

	kernel := createGaussianKernel(size, sigma) //одномерное ядро обязательно для Convolution
	temp := CopyImage(img)

	// Горизонтально
	applyConvolution(img, temp, kernel, true)
	// Вертикально
	applyConvolution(temp, img, kernel, false)
	return img
}

func Changes(a imatix.Image, b imatix.Image) imatix.Image {
	img := CopyImage(b)

	for i := 0; i < b.Height; i++ {
		for j := 0; j < b.Width; j++ {
			for k := 0; k < 3; k++ {
				aValue := imatix.PixelToContinuous(a.Matrix[i][j][k])
				bValue := imatix.PixelToContinuous(b.Matrix[i][j][k])
				img.Matrix[i][j][k] = clip(
					(1 - math.Abs(aValue-bValue)) * 255,
				)
			}
		}
	}
	return img
}

func UnsharpMasking(img imatix.Image, blured imatix.Image, power float64) imatix.Image {
	if power <= 0 {
		return img
	}

	for i := 0; i < img.Height; i++ {
		for j := 0; j < img.Width; j++ {
			for k := 0; k < 3; k++ {
				imgValue := imatix.PixelToContinuous(img.Matrix[i][j][k])
				bluredValue := imatix.PixelToContinuous(blured.Matrix[i][j][k])
				img.Matrix[i][j][k] = clip(
					(imgValue + power*(imgValue-bluredValue)) * 255,
				)
			}
		}
	}
	return img
}

func CopyImage(img imatix.Image) imatix.Image {
	copyImg := imatix.Image{
		Matrix: make([][][3]uint8, img.Height),
		Height: img.Height,
		Width:  img.Width,
	}

	for i := 0; i < img.Height; i++ {
		copyImg.Matrix[i] = make([][3]uint8, img.Width)
		copy(copyImg.Matrix[i], img.Matrix[i])
	}

	return copyImg
}

func SigmaFilter(img imatix.Image, size int, sigma float64, k float64) imatix.Image {
	temp := CopyImage(img)
	radius := (size - 1) / 2

	for y := 0; y < img.Height; y++ {
		for x := 0; x < img.Width; x++ {
			center := temp.Matrix[y][x]

			var sumR, sumG, sumB float64
			var count int

			for dy := -radius; dy <= radius; dy++ {
				for dx := -radius; dx <= radius; dx++ {
					nx, ny := x+dx, y+dy
					if nx < 0 || nx >= img.Width || ny < 0 || ny >= img.Height {
						continue
					}

					neighbor := temp.Matrix[ny][nx]

					if math.Abs(float64(neighbor[0])-float64(center[0])) <= k*sigma &&
						math.Abs(float64(neighbor[1])-float64(center[1])) <= k*sigma &&
						math.Abs(float64(neighbor[2])-float64(center[2])) <= k*sigma {

						sumR += float64(neighbor[0])
						sumG += float64(neighbor[1])
						sumB += float64(neighbor[2])
						count++
					}
				}
			}

			if count > 0 {
				img.Matrix[y][x] = [3]uint8{
					clip(sumR / float64(count)),
					clip(sumG / float64(count)),
					clip(sumB / float64(count)),
				}
			}
		}
	}
	return img
}

func MedianFilter(img imatix.Image, kernelSize int) imatix.Image {
	if kernelSize <= 1 {
		return img
	}

	temp := CopyImage(img)
	radius := kernelSize / 2
	windowSize := kernelSize * kernelSize

	for y := 0; y < img.Height; y++ {
		for x := 0; x < img.Width; x++ {

			windowR := make([]uint8, 0, windowSize)
			windowG := make([]uint8, 0, windowSize)
			windowB := make([]uint8, 0, windowSize)

			for dy := -radius; dy <= radius; dy++ {
				for dx := -radius; dx <= radius; dx++ {
					nx, ny := x+dx, y+dy

					if nx >= 0 && nx < img.Width && ny >= 0 && ny < img.Height {
						pixel := temp.Matrix[ny][nx]
						windowR = append(windowR, pixel[0])
						windowG = append(windowG, pixel[1])
						windowB = append(windowB, pixel[2])
					}
				}
			}

			if len(windowR) > 0 {
				img.Matrix[y][x] = [3]uint8{
					median(windowR),
					median(windowG),
					median(windowB),
				}
			}
		}
	}
	return img
}

func median(data []uint8) uint8 {
	if len(data) == 0 {
		return 0
	}

	sorted := make([]uint8, len(data))
	copy(sorted, data)
	sort.Slice(sorted, func(i, j int) bool { return sorted[i] < sorted[j] })

	middle := len(sorted) / 2
	return sorted[middle]
}

func RectangularFilter(img imatix.Image, kernelSize int) imatix.Image {
	if kernelSize <= 1 {
		return img
	}

	kernel := createRectangularKernel(kernelSize)
	temp := CopyImage(img)

	applyConvolution(img, temp, kernel, true)
	applyConvolution(temp, img, kernel, false)

	return img
}

func createRectangularKernel(size int) []float64 {
	kernel := make([]float64, size)
	value := 1.0 / float64(size)

	for i := 0; i < size; i++ {
		kernel[i] = value
	}

	return kernel
}

func LogarithmicClipAuto(img imatix.Image) imatix.Image {
	result := CopyImage(img)
	maxBrightness := findMaxBrightness(result)

	c := 255.0 / math.Log(1.0+maxBrightness)

	for y := 0; y < result.Height; y++ {
		for x := 0; x < result.Width; x++ {
			for channel := 0; channel < 3; channel++ {
				value := float64(result.Matrix[y][x][channel])
				newValue := c * math.Log(1.0+value)
				result.Matrix[y][x][channel] = clip(newValue)
			}
		}
	}
	return result
}

func PowerClipAuto(img imatix.Image, gamma float64) imatix.Image {
	result := CopyImage(img)
	maxBrightness := findMaxBrightness(result)

	c := 255.0 / math.Pow(maxBrightness/255.0, gamma)

	for y := 0; y < result.Height; y++ {
		for x := 0; x < result.Width; x++ {
			for channel := 0; channel < 3; channel++ {
				normalized := float64(result.Matrix[y][x][channel]) / 255.0
				cliped := c * math.Pow(normalized, gamma)
				result.Matrix[y][x][channel] = clip(cliped)
			}
		}
	}
	return result
}

func Binaryclip(img imatix.Image, threshold uint8) imatix.Image {
	result := CopyImage(img)

	for y := 0; y < result.Height; y++ {
		for x := 0; x < result.Width; x++ {
			for channel := 0; channel < 3; channel++ {

				if result.Matrix[y][x][channel] >= threshold {
					result.Matrix[y][x][channel] = 255
				} else {
					result.Matrix[y][x][channel] = 0
				}
			}
		}
	}
	return result
}

func IntensitySliceConstant(img imatix.Image, lowThreshold, highThreshold, constantValue uint8) imatix.Image {
	result := CopyImage(img)

	for y := 0; y < result.Height; y++ {
		for x := 0; x < result.Width; x++ {

			inRange := true
			for channel := 0; channel < 3; channel++ {
				if result.Matrix[y][x][channel] < lowThreshold || result.Matrix[y][x][channel] > highThreshold {
					inRange = false
					break
				}
			}

			if inRange {

				result.Matrix[y][x] = [3]uint8{constantValue, constantValue, constantValue}
			}

		}
	}
	return result
}

func findMaxBrightness(img imatix.Image) float64 {
	maxVal := 0.0
	for y := 0; y < img.Height; y++ {
		for x := 0; x < img.Width; x++ {
			for channel := 0; channel < 3; channel++ {
				val := float64(img.Matrix[y][x][channel])
				if val > maxVal {
					maxVal = val
				}
			}
		}
	}
	return maxVal
}
