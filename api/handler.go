package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image/jpeg"
	"image/png"
	"net/http"
	"strconv"
)

func CollatzHandler(w http.ResponseWriter, r *http.Request) {
	var params Params
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		http.Error(w, fmt.Errorf("error parsing Params: %w", err).Error(), http.StatusBadRequest)
		return
	}

	img, err := Draw(&params)
	if err != nil {
		http.Error(w, fmt.Errorf("error drawing image: %w", err).Error(), http.StatusInternalServerError)
		return
	}

	buffer := new(bytes.Buffer)
	switch params.Format {

	case ImageFormatJPG:
		err := jpeg.Encode(buffer, img, &jpeg.Options{Quality: 100})
		if err != nil {
			http.Error(w, fmt.Errorf("error encoding jpeg image: %w", err).Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "image/jpeg")

	case ImageFormatPNG:
		err := png.Encode(buffer, img)
		if err != nil {
			http.Error(w, fmt.Errorf("error encoding png image: %w", err).Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "image/png")
	}

	w.Header().Set("Content-Length", strconv.Itoa(buffer.Len()))

	if _, err := w.Write(buffer.Bytes()); err != nil {
		http.Error(w, fmt.Errorf("error writing buffer to response: %w", err).Error(), http.StatusInternalServerError)
		return
	}
}
