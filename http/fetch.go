package http

import (
	"log"
	"math"
	"net/http"
	"strings"

	"github.com/thoriqadillah/gown/setting"
)

type response struct {
	url         string
	size        int64
	contentType string
	cansplit    bool
	totalpart   int
	*setting.Setting
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

	contentType := res.Header.Get("Content-Type")

	// check if the file support cansplit download
	cansplit := res.Header.Get("Accept-Ranges") == "bytes"
	size := res.ContentLength

	totalpart := dynamicPartition(size, setting.Partsize)
	if size == -1 || !cansplit {
		totalpart = 1
		log.Println("File does not support download in chunks. Downloading the file entirely")
	}

	setting.Concurrency = totalpart

	response := &response{
		url:         url,
		size:        size,
		contentType: contentType,
		cansplit:    cansplit,
		totalpart:   totalpart,
		Setting:     setting,
	}

	return response, nil
}

func dynamicPartition(size int64, defaultParitionSize int64) int {
	num := math.Log10(float64(size / (1024 * 1024)))
	partsize := defaultParitionSize
	for i := 0; i < int(num); i++ {
		partsize *= 3 // 3 is just author's self configured number
	}

	return int(size / partsize)
}

func (r *response) Parts() int {
	return r.totalpart
}

func (r *response) Filename() string {
	split := strings.Split(r.contentType, "/")
	type_ := "." + split[len(split)-1]

	split = strings.Split(r.url, "/")
	filename := split[len(split)-1]
	filename = strings.Split(filename, "?")[0]
	if filename == "" {
		return "file" + type_
	}

	split = strings.Split(filename, ".")
	if len(split) != 0 {
		return filename
	}

	return filename + type_
}
