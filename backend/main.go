package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"net/http"
)

func imageHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Use POST method", http.StatusMethodNotAllowed)
		return
	}

	fmt.Println("Body:", r.Body)
	fmt.Println("Header:", r.Header.Get("Content-Type"))

	img, err := png.Decode(r.Body)
	if err != nil {
		http.Error(w, "Can't decode image", http.StatusBadRequest)
		return
	}

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

	fmt.Printf("Matrix created: %dx%dx%d\n", height, width, 3)

	pixel := imageMatrix[0][0]
	fmt.Printf("First pixel: R=%d, G=%d, B=%d\n", pixel[0], pixel[1], pixel[2])

	newImg := image.NewRGBA(image.Rect(0, 0, width, height))

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			pixel := imageMatrix[y][x]
			color := color.RGBA{
				R: pixel[0],
				G: pixel[1],
				B: pixel[2],
				A: 255,
			}
			newImg.Set(x, y, color)
		}
	}

	w.Header().Set("Content-Type", "image/png")

	if err := png.Encode(w, newImg); err != nil {
		http.Error(w, "Failed to encode image", http.StatusInternalServerError)
		return
	}
}

func main() {
	fmt.Println("Server started")

	http.HandleFunc("/api/v1/image", imageHandler)
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		fmt.Println("Error:", err)
	}

}
