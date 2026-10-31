package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

const (
	UploadDir     = "./uploads"
	MaxUploadSize = 20 << 20 // 20MB
)

func main() {
	// 创建上传目录
	err := os.MkdirAll(UploadDir, os.ModePerm)
	if err != nil {
		panic(err)
	}

	// 静态页面
	http.Handle("/", http.FileServer(http.Dir("./static")))

	// 上传接口
	http.HandleFunc("/upload", uploadHandler)

	// 文件访问
	http.Handle("/files/",
		http.StripPrefix("/files/",
			http.FileServer(http.Dir(UploadDir)),
		),
	)

	fmt.Println("server running at :8080")
	http.ListenAndServe(":8080", nil)
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	// 限制上传大小
	r.Body = http.MaxBytesReader(w, r.Body, MaxUploadSize)

	err := r.ParseMultipartForm(MaxUploadSize)
	if err != nil {
		http.Error(w, "file too large", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "invalid file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	filename := filepath.Base(header.Filename)

	dstPath := filepath.Join(UploadDir, filename)

	dst, err := os.Create(dstPath)
	if err != nil {
		http.Error(w, "save failed", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		http.Error(w, "write failed", http.StatusInternalServerError)
		return
	}

	fileURL := fmt.Sprintf("/files/%s", filename)

	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(fileURL))
}
