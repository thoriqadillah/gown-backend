package http

import (
	"log"
	"net/http"
	"strings"

	"github.com/thoriqadillah/gown/setting"
)

type response struct {
	url         string
	filename    string
	size        int64
	contentType string
	cansplit    bool
	totalpart   int
}

func Fetch(url string, setting *setting.Setting) (*response, error) {
	log.Printf("Fetching file data")

	// get the redirected url
	res, err := http.Head(url)
	if err != nil {
		log.Printf("Error fetching url %v", err)
		return nil, err
	}

	newurl := res.Request.URL.String()
	if url != newurl {
		log.Printf("Following link to %s", newurl[:50]+"...")
	}

	url = newurl

	// get content-length (size in bytes) of a file
	res, err = http.Head(url)
	if err != nil {
		log.Printf("Error fetching file url %v", err)
		return nil, err
	}

	size := res.ContentLength

	contentType := res.Header.Get("Content-Type")
	split := strings.Split(contentType, "/")
	contentType = "." + split[len(split)-1]

	split = strings.Split(url, "/")
	filename := split[len(split)-1]
	filename = strings.Split(filename, "?")[0]
	if filename == "" {
		filename = "file" + contentType
	} else {
		filename += contentType
	}

	// check if the file support cansplit download
	cansplit := res.Header.Get("Accept-Ranges") == "bytes"

	totalpart := int(size / setting.Partsize)
	if size == -1 || !cansplit {
		totalpart = 1
		log.Println("File does not support download in chunks. Downloading the file entirely")
	}
	setting.Concurrency = totalpart

	response := &response{
		url:         url,
		filename:    filename,
		size:        size,
		contentType: contentType,
		cansplit:    cansplit,
		totalpart:   totalpart,
	}

	return response, nil
}

func (r *response) Parts() int {
	return r.totalpart
}

func (r *response) Filename() string {
	return r.filename
}
