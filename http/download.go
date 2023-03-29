package http

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/thoriqadillah/gown/config"
	"github.com/thoriqadillah/gown/worker"
)

type Chunk struct {
	*response
	wg *sync.WaitGroup
	*config.Config
	index int
	start int64
	end   int64
	size  int64
	Data  []byte
}

func Download(res *response, index int, wg *sync.WaitGroup, config *config.Config) worker.Job {
	// get the range part that we want to download
	totalpart := int64(res.totalpart)
	partsize := res.size / totalpart

	start := int64(index * int(partsize))
	end := start + int64(int(partsize)-1)

	if index == int(totalpart)-1 {
		end = res.size
	}

	return &Chunk{
		response: res,
		wg:       wg,
		Config:   config,
		index:    index,
		start:    start,
		end:      end,
		size:     partsize,
	}
}

func (d *Chunk) Struct() interface{} {
	return d
}

// TODO: implement retry
// TODO: handle wether to download file entirely or split
func (d *Chunk) Execute() error {
	defer d.wg.Done()

	httpclient := &http.Client{}

	part := fmt.Sprintf("bytes=%d-%d", d.start, d.end)

	if d.size == -1 {
		log.Printf("Downloading part %d with size unknown", d.index+1)
	} else {
		log.Printf("Downloading part %d from %d to %d", d.index+1, d.start, d.end)
	}

	req, err := http.NewRequest("GET", d.url, nil)
	if err != nil {
		return err
	}

	start := time.Now()
	req.Header.Add("Range", part)
	res, err := httpclient.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	d.Data, err = io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	elapsed := time.Since(start)
	log.Printf("Downloading done for worker with id %d in %v s\n", d.index, elapsed.Seconds())
	return nil
}

func (d *Chunk) HandleError(err error) {
	// TODO: handle error
	log.Printf("Error downloading the file: %v", err)
}
