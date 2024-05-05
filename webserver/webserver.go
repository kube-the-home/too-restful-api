package webserver

import (
	"encoding/json"
	"errors"
	"fmt"
	"kube-the-home/too-restful-api/config"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
)

func Execute() {
	files, err := os.ReadDir(config.CONFIG.Path)
	if err != nil {
		slog.Error("failed to read directory", "error", err.Error())
	}

	filesMap := make(map[string][]byte)
	for _, file := range files {
		if !file.IsDir() {
			content, err := os.ReadFile(filepath.Join(config.CONFIG.Path, file.Name()))

			if err != nil {
				slog.Error("failed to read file", "error", err.Error())
				continue
			}

			filesMap[file.Name()] = content
		}
	}

	// Start the webserver
	Init(filesMap)
}

func Init(files map[string][]byte) {
	mux := http.NewServeMux()

	for file, content := range files {
		path := "/items/" + file
		slog.Info("http://localhost:" + config.CONFIG.Port + path)
		mux.HandleFunc("GET "+path, func(w http.ResponseWriter, r *http.Request) {
			getData(w, r, content)
		})
	}

	mux.HandleFunc("GET /list", func(w http.ResponseWriter, r *http.Request) {
		list(w, r, files)
	})

	mux.HandleFunc("GET /metrics", func(w http.ResponseWriter, r *http.Request) {
		metrics(w, r, files)
	})

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", config.CONFIG.Port),
		Handler: mux,
	}

	slog.Info("Started webserver", "port", config.CONFIG.Port)
	if err := server.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
		slog.Info("listen", "error", err.Error())
	}
	slog.Warn("Goodbye!")
}

func getData(w http.ResponseWriter, r *http.Request, data []byte) {
	w.Header().Set("Content-Type", "application/javascript")

	if data == nil {
		http.Error(w, "No data available", http.StatusInternalServerError)
		return
	} else {
		_, err := w.Write(data)
		if err != nil {
			slog.Error("writing data", "error", err.Error())
			http.Error(w, "Error writing data", http.StatusInternalServerError)
			return
		}
	}
}

func metrics(w http.ResponseWriter, r *http.Request, files map[string][]byte) {
	w.Header().Set("Content-Type", "application/json")

	data := fmt.Sprint("too_restful_api_file_count ", len(files))

	byteData := []byte(data)
	_, err := w.Write(byteData)
	if err != nil {
		slog.Error("writing data", "error", err.Error())
		http.Error(w, "Error writing data", http.StatusInternalServerError)
		return
	}

}

func list(w http.ResponseWriter, r *http.Request, files map[string][]byte) {
	w.Header().Set("Content-Type", "application/json")

	var keys []string
	for k := range files { // Main code to get keys
		keys = append(keys, k)
	}

	jsonString, err := json.Marshal(keys)

	if err != nil {
		slog.Error("parsing list data", "error", err.Error())
		http.Error(w, "Error parsing list data", http.StatusInternalServerError)
		return
	}

	_, err = w.Write(jsonString)
	if err != nil {
		slog.Error("writing data", "error", err.Error())
		http.Error(w, "Error writing data", http.StatusInternalServerError)
		return
	}

}
