package main

import (
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"net/http"
	"net/url"
	"strconv"

	imatix "github.com/J0hnLenin/ComputerVision/imatrix"
	"github.com/J0hnLenin/ComputerVision/processors"
	"github.com/J0hnLenin/ComputerVision/redactor"
	"github.com/J0hnLenin/ComputerVision/statistics"
	"github.com/gorilla/mux"
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

func imageApplyHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	funcName := vars["func"]
	processor, exists := processors.Functions[funcName]
	if !exists {
		http.Error(w, "Invalid name of function", http.StatusNotFound)
		return
	}

	img, err := png.Decode(r.Body)
	if err != nil {
		http.Error(w, "Can't decode image", http.StatusBadRequest)
		return
	}

	formatedImage := imageDeserializer(img)
	processor(formatedImage)
	newImg := imageSerializer(formatedImage)

	w.Header().Set("Content-Type", "image/png")

	if err := png.Encode(w, newImg); err != nil {
		http.Error(w, "Failed to encode image", http.StatusInternalServerError)
		return
	}
}

func imageStatisticsHandler(w http.ResponseWriter, r *http.Request) {
	img, err := png.Decode(r.Body)
	if err != nil {
		http.Error(w, "Can't decode image", http.StatusBadRequest)
		return
	}

	formatedImage := imageDeserializer(img)
	stats := statistics.GetStatistics(formatedImage)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(stats); err != nil {
		http.Error(w, "JSON encoding error", http.StatusInternalServerError)
		return
	}
}

func parseQuery(queryParameters url.Values) imatix.Parameters {
	redBrightness, _ := strconv.Atoi(queryParameters.Get("redBrightness"))
	if redBrightness == 0 {
		redBrightness = 100
	}
	greenBrightness, _ := strconv.Atoi(queryParameters.Get("greenBrightness"))
	if greenBrightness == 0 {
		greenBrightness = 100
	}
	blueBrightness, _ := strconv.Atoi(queryParameters.Get("blueBrightness"))
	if blueBrightness == 0 {
		blueBrightness = 100
	}
	contrast, _ := strconv.Atoi(queryParameters.Get("contrast"))
	if contrast == 0 {
		contrast = 100
	}
	negative, _ := strconv.ParseBool(queryParameters.Get("negative"))
	order := queryParameters.Get("order")
	vertivalMirror, _ := strconv.ParseBool(queryParameters.Get("vertivalMirror"))
	horisontalMirror, _ := strconv.ParseBool(queryParameters.Get("horisontalMirror"))
	magic, _ := strconv.Atoi(queryParameters.Get("magic"))
	if contrast == 0 {
		contrast = 100
	}

	return imatix.Parameters{
		RedBrightness:    float64(redBrightness) / 100,
		GreenBrightness:  float64(greenBrightness) / 100,
		BlueBrightness:   float64(blueBrightness) / 100,
		Contrast:         float64(contrast) / 100,
		Negative:         negative,
		Order:            order,
		VertivalMirror:   vertivalMirror,
		HorisontalMirror: horisontalMirror,
		Magic:            magic,
	}

}

func imageRedactorHandler(w http.ResponseWriter, r *http.Request) {
	img, err := png.Decode(r.Body)
	if err != nil {
		http.Error(w, "Can't decode image", http.StatusBadRequest)
		return
	}

	parameters := parseQuery(r.URL.Query())

	formatedImage := imageDeserializer(img)

	redactedImage := redactor.Redact(formatedImage, parameters)

	newImg := imageSerializer(redactedImage)

	w.Header().Set("Content-Type", "image/png")

	if err := png.Encode(w, newImg); err != nil {
		http.Error(w, "Failed to encode image", http.StatusInternalServerError)
		return
	}
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/api/v1/image/apply/{func}", imageApplyHandler).Methods("POST")
	r.HandleFunc("/api/v1/image/statistics", imageStatisticsHandler).Methods("POST")
	r.HandleFunc("/api/v1/image/redactor", imageRedactorHandler).Methods("POST")
	fmt.Println("Server starting on :8000")
	log.Fatal(http.ListenAndServe(":8000", r))
}
