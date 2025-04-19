package main

import (
	"io"
	"mime/multipart"
	"net/http"
	"net/http/cgi"
	"os"
	"path"
	"regexp"
)

func main() {
	cgi.Serve(http.HandlerFunc(uploadHandler))
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := NewTmpl()
	err := tmpl.Init()
	if err != nil {
		panic(err)
	}

	if r.Method == "GET" {
		tmpl.RenderPage(w, "upload.html")
	} else if r.Method == "POST" {
		link, err := processUpload(w, r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		tmpl.Link = link
		tmpl.RenderPage(w, "done.html")
	}
}

func processUpload(w http.ResponseWriter, r *http.Request) (string, error) {
	err := r.ParseMultipartForm(10 << 20) // 10 MB limit
	if err != nil {
		return "", err
	}

	files := r.MultipartForm.File["files"]
	dir, err := uploadDir()
	if err != nil {
		return "", err
	}

	link := regexp.MustCompile(`/[^/]+/[^/]+$`).FindString(dir)
	names := []string{}

	for _, fh := range files {
		name, err := processFile(fh, dir)
		if err != nil {
			return "", err
		}
		names = append(names, name)
	}

	idx, err := os.OpenFile(path.Join(dir, "index.html"), os.O_CREATE | os.O_WRONLY | os.O_TRUNC, 0644)
	if err != nil {
		return "", err
	}
	defer idx.Close()

	tpl := NewTmpl()
	if err := tpl.Init(); err != nil {
		return "", err
	}

	for _, name := range names {
		tpl.AddPhoto(name)
	}

	tpl.RenderPage(idx, "photos.html")
	return link, nil
}

func processFile(fh *multipart.FileHeader, dir string) (string, error) {
	file, err := fh.Open()
	if err != nil {
		return "", err
	}
	defer file.Close()

	fileName := fh.Filename
	savePath := path.Join(dir, sanitizeFileName(fileName))
	outFile, err := os.Create(savePath)
	if err != nil {
		return "", err
	}
	defer outFile.Close()

	if _, err = io.Copy(outFile, file); err != nil {
		return "", err
	}
	if err = scaleJPG(savePath); err != nil {
		return "", err
	}

	return fh.Filename, nil
}

func sanitizeFileName(fileName string) string {
	allowed := regexp.MustCompile(`[^a-zA-Z0-9_\-\.]+`)
	return allowed.ReplaceAllString(fileName, "_")
}