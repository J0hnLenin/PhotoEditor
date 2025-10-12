package main

import (
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	imatix "github.com/J0hnLenin/ComputerVision/imatrix"
	"github.com/J0hnLenin/ComputerVision/redactor"
	"github.com/J0hnLenin/ComputerVision/statistics"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func imageDeserializer(img image.Image) imatix.Image {
	bounds := img.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y
	imageMatrix := make([][][3]uint8, height)
	for y := 0; y < height; y++ {
		imageMatrix[y] = make([][3]uint8, width)
		for x := 0; x < width; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			imageMatrix[y][x] = [3]uint8{
				uint8(r >> 8),
				uint8(g >> 8),
				uint8(b >> 8),
			}
		}
	}
	return imatix.Image{
		Width:  width,
		Height: height,
		Matrix: imageMatrix,
	}
}

func imageSerializer(img imatix.Image) image.Image {
	newImg := image.NewRGBA(image.Rect(0, 0, img.Width, img.Height))

	for y := 0; y < img.Height; y++ {
		for x := 0; x < img.Width; x++ {
			pixel := img.Matrix[y][x]
			color := color.RGBA{
				R: pixel[0],
				G: pixel[1],
				B: pixel[2],
				A: 255,
			}
			newImg.Set(x, y, color)
		}
	}
	return newImg
}

func parseQuery(queryParameters url.Values) imatix.Parameters {
	redBrightness, _ := strconv.Atoi(queryParameters.Get("RedBrightness"))
	if redBrightness == 0 {
		redBrightness = 100
	}
	greenBrightness, _ := strconv.Atoi(queryParameters.Get("GreenBrightness"))
	if greenBrightness == 0 {
		greenBrightness = 100
	}
	blueBrightness, _ := strconv.Atoi(queryParameters.Get("BlueBrightness"))
	if blueBrightness == 0 {
		blueBrightness = 100
	}
	logarithmicBrightness, _ := strconv.Atoi(queryParameters.Get("LogarithmicBrightness"))
	contrast, _ := strconv.Atoi(queryParameters.Get("Contrast"))
	if contrast == 0 {
		contrast = 100
	}
	negative, _ := strconv.ParseBool(queryParameters.Get("Negative"))
	order := queryParameters.Get("Order")
	verticalMirror, _ := strconv.ParseBool(queryParameters.Get("VerticalMirror"))
	horizontalMirror, _ := strconv.ParseBool(queryParameters.Get("HorizontalMirror"))
	magic, _ := strconv.Atoi(queryParameters.Get("Magic"))
	filter := queryParameters.Get("Filter")
	filterSize, _ := strconv.Atoi(queryParameters.Get("FilterSize"))
	sigma, _ := strconv.ParseFloat(queryParameters.Get("Sigma"), 64)
	interval, _ := strconv.Atoi(queryParameters.Get("Interval"))
	unsharpMasking, _ := strconv.ParseFloat(queryParameters.Get("UnsharpMasking"), 64)

	logarithmicClip, _ := strconv.ParseBool(queryParameters.Get("LogarithmicClip"))
	powerClip, _ := strconv.ParseFloat(queryParameters.Get("PowerClip"), 64)
	binaryClip, _ := strconv.Atoi(queryParameters.Get("BinaryClip"))
	constantLow, _ := strconv.Atoi(queryParameters.Get("ConstantLow"))
	constantHigh, _ := strconv.Atoi(queryParameters.Get("ConstantHigh"))
	constantValue, _ := strconv.Atoi(queryParameters.Get("ConstantValue"))

	return imatix.Parameters{
		RedBrightness:         float64(200-redBrightness) / 100,
		GreenBrightness:       float64(200-greenBrightness) / 100,
		BlueBrightness:        float64(200-blueBrightness) / 100,
		LogarithmicBrightness: float64(logarithmicBrightness) / 100,
		Contrast:              float64(contrast) / 100,
		Negative:              negative,
		Order:                 order,
		VerticalMirror:        verticalMirror,
		HorizontalMirror:      horizontalMirror,
		Magic:                 magic,
		Filter:                filter,
		FilterSize:            filterSize,
		Sigma:                 sigma,
		Interval:              interval,
		UnsharpMasking:        unsharpMasking,
		LogarithmicClip: 	   logarithmicClip,
		PowerClip: 			   powerClip,
		BinaryClip:			   binaryClip,
		ConstantLow:		   constantLow,
		ConstantHigh:		   constantHigh,
		ConstantValue:		   constantValue,
	}

}

func processChannel(originalImage imatix.Image, channelName string) imatix.Image {
	channelImage := imatix.Image{
		Width:  originalImage.Width,
		Height: originalImage.Height,
		Matrix: make([][][3]uint8, originalImage.Height),
	}

	for y := 0; y < originalImage.Height; y++ {
		channelImage.Matrix[y] = make([][3]uint8, originalImage.Width)
		for x := 0; x < originalImage.Width; x++ {
			pixel := originalImage.Matrix[y][x]
			switch channelName {
			case "Red":
				channelImage.Matrix[y][x] = [3]uint8{pixel[0], 0, 0}
			case "Green":
				channelImage.Matrix[y][x] = [3]uint8{0, pixel[1], 0}
			case "Blue":
				channelImage.Matrix[y][x] = [3]uint8{0, 0, pixel[2]}
			case "Gray":
				gray := uint8(0.299*float64(pixel[0]) + 0.587*float64(pixel[1]) + 0.114*float64(pixel[2]))
				channelImage.Matrix[y][x] = [3]uint8{gray, gray, gray}
			}
		}
	}
	return channelImage
}

func imageRedactorHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		http.Error(w, "Can't parse form data", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "Can't get image from form data", http.StatusBadRequest)
		return
	}
	defer file.Close()

	fmt.Printf("Processing file: %s, Size: %d bytes\n", header.Filename, header.Size)

	parameters := parseQuery(r.URL.Query())

	img, err := png.Decode(file)
	if err != nil {
		http.Error(w, "Can't decode image", http.StatusBadRequest)
		return
	}

	originalFormatedImage := imageDeserializer(img)

	changesImage, redactedImage := redactor.Redact(originalFormatedImage, parameters)

	redChannel := processChannel(originalFormatedImage, "Red")
	greenChannel := processChannel(originalFormatedImage, "Green")
	blueChannel := processChannel(originalFormatedImage, "Blue")
	grayChannel := processChannel(originalFormatedImage, "Gray")

	stats := statistics.GetStatistics(originalFormatedImage)

	writer := multipart.NewWriter(w)
	w.Header().Set("Content-Type", writer.FormDataContentType())

	redactedPart, err := writer.CreateFormFile("redacted_image", "redacted_"+header.Filename)
	if err != nil {
		http.Error(w, "Can't create form file for redacted image", http.StatusInternalServerError)
		return
	}
	if err := png.Encode(redactedPart, imageSerializer(redactedImage)); err != nil {
		http.Error(w, "Failed to encode redacted image", http.StatusInternalServerError)
		return
	}

	redPart, err := writer.CreateFormFile("red_channel", "red_channel_"+header.Filename)
	if err != nil {
		http.Error(w, "Can't create form file for red channel", http.StatusInternalServerError)
		return
	}
	if err := png.Encode(redPart, imageSerializer(redChannel)); err != nil {
		http.Error(w, "Failed to encode red channel image", http.StatusInternalServerError)
		return
	}

	greenPart, err := writer.CreateFormFile("green_channel", "green_channel_"+header.Filename)
	if err != nil {
		http.Error(w, "Can't create form file for green channel", http.StatusInternalServerError)
		return
	}
	if err := png.Encode(greenPart, imageSerializer(greenChannel)); err != nil {
		http.Error(w, "Failed to encode green channel image", http.StatusInternalServerError)
		return
	}

	bluePart, err := writer.CreateFormFile("blue_channel", "blue_channel_"+header.Filename)
	if err != nil {
		http.Error(w, "Can't create form file for blue channel", http.StatusInternalServerError)
		return
	}
	if err := png.Encode(bluePart, imageSerializer(blueChannel)); err != nil {
		http.Error(w, "Failed to encode blue channel image", http.StatusInternalServerError)
		return
	}

	grayPart, err := writer.CreateFormFile("gray_channel", "gray_channel_"+header.Filename)
	if err != nil {
		http.Error(w, "Can't create form file for gray channel", http.StatusInternalServerError)
		return
	}
	if err := png.Encode(grayPart, imageSerializer(grayChannel)); err != nil {
		http.Error(w, "Failed to encode gray channel image", http.StatusInternalServerError)
		return
	}

	changesPart, err := writer.CreateFormFile("changes_channel", "changes_channel"+header.Filename)
	if err != nil {
		http.Error(w, "Can't create form file for changes channel", http.StatusInternalServerError)
		return
	}
	if err := png.Encode(changesPart, imageSerializer(changesImage)); err != nil {
		http.Error(w, "Failed to encode changes channel image", http.StatusInternalServerError)
		return
	}

	statsPart, err := writer.CreateFormField("statistics")
	if err != nil {
		http.Error(w, "Can't create form field for statistics", http.StatusInternalServerError)
		return
	}
	jsonData, err := json.Marshal(stats)
	if err != nil {
		http.Error(w, "JSON encoding error", http.StatusInternalServerError)
		return
	}
	statsPart.Write(jsonData)

	writer.WriteField("original_filename", header.Filename)
	writer.WriteField("file_size", strconv.FormatInt(header.Size, 10))
	writer.WriteField("image_width", strconv.Itoa(originalFormatedImage.Width))
	writer.WriteField("image_height", strconv.Itoa(originalFormatedImage.Height))
	writer.WriteField("processing_parameters", r.URL.RawQuery)
	writer.WriteField("status", "success")

	writer.Close()
}

func main() {

	allowedOrigins := os.Getenv("CORS_ALLOWED_ORIGINS")
	if allowedOrigins == "" {
		allowedOrigins = "http://localhost,http://127.0.0.1,http://localhost:80,http://127.0.0.1:80"
	}
	fmt.Println(allowedOrigins)
	origins := strings.Split(allowedOrigins, ",")

	c := cors.New(cors.Options{
		AllowedOrigins:   origins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
		Debug:            true,
	})

	r := mux.NewRouter()
	r.HandleFunc("/api/v1/image/redactor", imageRedactorHandler).Methods("POST")

	handler := c.Handler(r)

	fmt.Println("Server starting on :8000")
	log.Fatal(http.ListenAndServe(":8000", handler))
}
