package main

import (
	"log"
	"sync"
	"time"

	"github.com/thoriqadillah/gown/fs"
	"github.com/thoriqadillah/gown/http"
	"github.com/thoriqadillah/gown/setting"
	pool "github.com/thoriqadillah/gown/worker"
)

func main() {
	setting := setting.New()

	response, err := http.Fetch("https://storage.googleapis.com/flutter_infra_release/releases/stable/linux/flutter_linux_3.7.9-stable.tar.xz", &setting)
	if err != nil {
		log.Fatal(err)
	}

	worker, err := pool.New(setting.Concurrency, setting.SimmultanousNum)
	if err != nil {
		log.Fatal("Error creating worker")
	}
	worker.Start()
	defer worker.Stop()

	var wg sync.WaitGroup

	start := time.Now()

	file := fs.New(response.Parts(), &setting)
	downloadjobs := make([]*http.Chunk, response.Parts())
	for part := range downloadjobs {
		downloadjobs[part] = http.NewChunk(response, part, &wg)
	}

	for _, job := range downloadjobs {
		wg.Add(1)
		worker.Add(job)
	}

	wg.Wait()

	for part, chunk := range downloadjobs {
		file.Combine(chunk.Data(), part)
	}

	if err := file.Save(response.Filename()); err != nil {
		log.Fatal("Error writing file", err)
	}

	elapsed := time.Since(start)
	log.Printf("Took %v s to download %s\n", elapsed.Seconds(), response.Filename())
}
