package http

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/thoriqadillah/gown/setting"
)

type Chunk struct {
	*response
	wg *sync.WaitGroup
	*setting.Setting
	index int
	start int64
	end   int64
	size  int64
	data  []byte
}

func Download(res *response, index int, wg *sync.WaitGroup, setting *setting.Setting) *Chunk {
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
		Setting:  setting,
		index:    index,
		start:    start,
		end:      end,
		size:     partsize,
	}
}

func (c *Chunk) download() error {
	defer c.wg.Done()

	httpclient := &http.Client{}

	part := fmt.Sprintf("bytes=%d-%d", c.start, c.end)

	if c.size == -1 {
		log.Printf("Downloading part %d with size unknown", c.index+1)
	} else {
		log.Printf("Downloading part %d from %d to %d", c.index+1, c.start, c.end)
	}

	req, err := http.NewRequest("GET", c.url, nil)
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

	c.data, err = io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	elapsed := time.Since(start)
	log.Printf("Downloading done for worker with id %d in %v s\n", c.index, elapsed.Seconds())

	return nil
}

// TODO: handle wether to download file entirely or split
func (c *Chunk) Execute() error {
	var err error

	for retry := 0; retry < c.Maxtries; retry++ {
		if err = c.download(); err == nil {
			return nil
		}

		log.Printf("Error while downloading chunk %d: %v. Retrying....\n", c.index+1, err)
	}

	return err
}

func (c *Chunk) HandleError(err error) {
	// TODO: handle error
	log.Printf("Error downloading the file: %v", err)
}

func (c *Chunk) Data() []byte {
	return c.data
}
