package main

import (
	"log"
	"sync"
	"time"

	"github.com/thoriqadillah/gown/config"
	"github.com/thoriqadillah/gown/fs"
	"github.com/thoriqadillah/gown/http"
	pool "github.com/thoriqadillah/gown/worker"
)

func main() {
	config := config.Default()

	response, err := http.Fetch("https://rr5---sn-h5q7knes.googlevideo.com/videoplayback?expire=1680096599&ei=9-gjZM2JN9nS7QSNko7QAw&ip=185.147.214.82&id=o-AAN6M3DCUsvRlfA8S_jiuBTRYAMP-me6vVLKnz0v1tf7&itag=22&source=youtube&requiressl=yes&mh=ja&mm=31%2C29&mn=sn-h5q7knes%2Csn-h5qzen76&ms=au%2Crdu&mv=m&mvi=5&pl=24&initcwndbps=1935000&spc=99c5Ce-NrzRBCuzbs6VasHOsb3BaLdeZnCJaa-AlgO-vUmQEgA&vprv=1&mime=video%2Fmp4&ns=xuw1Mg_XsEWTO-gWg0CzNZQM&cnr=14&ratebypass=yes&dur=324.220&lmt=1679976939892667&mt=1680074740&fvip=4&fexp=24007246&c=WEB&txp=5532434&n=DK7-G3FlRhyC1Q&sparams=expire%2Cei%2Cip%2Cid%2Citag%2Csource%2Crequiressl%2Cspc%2Cvprv%2Cmime%2Cns%2Ccnr%2Cratebypass%2Cdur%2Clmt&sig=AOq0QJ8wRQIgCHUXwSWAlAT0GEp-IqvQof3tkt2JWMLfGuac7YbK0XYCIQDmEnqOrUrm2ffsraoRPOJCsgjcUm9HKLYiORzXXBI0GQ%3D%3D&lsparams=mh%2Cmm%2Cmn%2Cms%2Cmv%2Cmvi%2Cpl%2Cinitcwndbps&lsig=AG3C_xAwRQIgbOzk8gI-ipxbBvqKh84J7swQl9ZMById1A3DtlDTIrwCIQDp42kzwQjN1X-LSZCq4UVsCyky1E0M2tnr1LRUGugRgg%3D%3D&title=Remnant%202%3A%20The%20First%20Hands-On%20Preview%20-%20IGN%20First", &config)
	if err != nil {
		log.Fatal(err)
	}

	worker, err := pool.New(config.Concurrency, config.SimmultanousNum)
	if err != nil {
		log.Fatal("Error creating worker")
	}
	worker.Start()
	defer worker.Stop()

	var wg sync.WaitGroup

	start := time.Now()

	file := fs.New(response.Parts(), &config)
	downloadjobs := make([]pool.Job, response.Parts())
	for part := range downloadjobs {
		downloadjobs[part] = http.Download(response, part, &wg, &config)
	}

	for _, job := range downloadjobs {
		wg.Add(1)
		worker.Add(job)
	}

	wg.Wait()

	for part, job := range downloadjobs {
		chunk := job.Struct().(*http.Chunk)
		file.Combine(chunk.Data, part)
	}

	if err := file.Save(response.Filename()); err != nil {
		log.Fatal("Error writing file", err)
	}

	elapsed := time.Since(start)
	log.Printf("Took %v s to download %s\n", elapsed.Seconds(), response.Filename())
}
