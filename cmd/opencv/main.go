package main

import (
	"image/color"

	"gocv.io/x/gocv"
)

func main() {
	webcam, err := gocv.VideoCaptureDevice(0)
	if err != nil {
		println(err)
	}
	defer webcam.Close()

	window := gocv.NewWindow("Detector")
	defer window.Close()

	classifier := gocv.NewCascadeClassifier()
	defer classifier.Close()

	classifier.Load("haarcascade_frontalface_default.xml")
	for {
		img := gocv.NewMat()
		if ok := webcam.Read(&img); !ok {
			println("Cannot read the camera")
			return
		}
		if img.Empty() {
			continue
		}

		myFace := classifier.DetectMultiScale(img)
		for _, r := range myFace {
			gocv.Rectangle(&img, r, color.RGBA{0, 255, 0, 0}, 2)
		}

		window.IMShow(img)
		if window.WaitKey(1) >= 0 {
			break
		}
	}
}
