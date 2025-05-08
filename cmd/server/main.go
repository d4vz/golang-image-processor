package main

import (
	"fmt"
	"image"
	"net/http"
	"os"
	"os/signal"

	_ "image/jpeg"
	_ "image/png"

	"github.com/d4vz/golang-image-processor/packages/api"
	"github.com/d4vz/golang-image-processor/packages/logger"
	"github.com/d4vz/golang-image-processor/packages/processor"
)

func main() {
	api := api.NewApi(api.NewApiConfig(":8080"))
	logger.Info("Starting server...")
	processor := processor.NewProcessor(processor.NewProcessorConfig("/tmp"))

	api.Handle("/process", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.Body = http.MaxBytesReader(w, r.Body, 10<<20)
		err := r.ParseMultipartForm(10 << 20)

		if err != nil {
			http.Error(w, "The uploaded file is too large. Please choose an image less than 10MB in size", http.StatusBadRequest)
			return
		}

		file, _, err := r.FormFile("image")

		if err != nil {
			logger.Error(fmt.Sprintf("Failed to get form file: %v", err))
			http.Error(w, "Invalid form data", http.StatusBadRequest)
			return
		}
		defer file.Close()

		img, _, err := image.Decode(file)

		if err != nil {
			logger.Error(fmt.Sprintf("Failed to decode image: %v", err))
			http.Error(w, "Failed to decode image", http.StatusInternalServerError)
			return
		}

		logger.Info("Starting image processing...")
		err = processor.Process(r.Context(), img)

		if err != nil {
			logger.Error(fmt.Sprintf("Failed to process image: %v", err))
			http.Error(w, "Failed to process image", http.StatusInternalServerError)
			return
		}

		logger.Info("Image processed successfully")
		w.WriteHeader(http.StatusOK)
	}))

	go func() {
		err := api.Start()

		if err != nil {
			logger.Error("Failed to start server")
			os.Exit(1)
		}
	}()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	<-sig
}
