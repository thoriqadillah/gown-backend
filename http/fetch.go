package http

import (
	"log"
	"net/http"
	"strings"

	"github.com/thoriqadillah/gown/config"
)

type response struct {
	url         string
	filename    string
	size        int64
	contentType string
	cansplit    bool
	totalpart   int
	config.Config
}

func Fetch(url string, conf *config.Config) (*response, error) {
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
	filename = strings.Split(filename, "?")[0] + contentType

	// check if the file support cansplit download
	cansplit := res.Header.Get("Accept-Ranges") == "bytes"
	if !cansplit {
		log.Println("Does not support split download. Downloading the file entirely")
	}

	totalpart := 1
	if cansplit {
		totalpart = int(size / conf.Partsize)
	}
	conf.Concurrency = totalpart

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
