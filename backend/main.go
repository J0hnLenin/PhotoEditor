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
	"os"
	"strconv"
	"strings"

	imatix "github.com/J0hnLenin/ComputerVision/imatrix"
	"github.com/J0hnLenin/ComputerVision/processors"
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
	r.HandleFunc("/api/v1/image/apply/{func}", imageApplyHandler).Methods("POST")
	r.HandleFunc("/api/v1/image/statistics", imageStatisticsHandler).Methods("POST")
	r.HandleFunc("/api/v1/image/redactor", imageRedactorHandler).Methods("POST")

	handler := c.Handler(r)

	fmt.Println("Server starting on :8000")
	log.Fatal(http.ListenAndServe(":8000", handler))
}
