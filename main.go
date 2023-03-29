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
	setting := setting.Default()

	response, err := http.Fetch("https://rr2---sn-5go7ynld.googlevideo.com/videoplayback?expire=1680116117&ei=NTUkZLuoIpHR7ASC4azABA&ip=188.126.94.102&id=o-ALhmOJ3oieUYAXE3mPY_Sj47Iwt-1BF9NXyINPRlWa5w&itag=22&source=youtube&requiressl=yes&mh=_T&mm=31%2C26&mn=sn-5go7ynld%2Csn-i5heen7z&ms=au%2Conr&mv=m&mvi=2&pl=23&initcwndbps=1851250&spc=99c5CaJolB8T_Xr99VByK4brgVJYJjBhBOpyYyc6HwT9UFGO9A&vprv=1&mime=video%2Fmp4&ns=4HUKe8I22PeU_NUng8TX0KIM&cnr=14&ratebypass=yes&dur=757.713&lmt=1680089430107467&mt=1680094186&fvip=2&fexp=24007246&c=WEB&txp=4432434&n=iPK0LMPnHKXVqQ&sparams=expire%2Cei%2Cip%2Cid%2Citag%2Csource%2Crequiressl%2Cspc%2Cvprv%2Cmime%2Cns%2Ccnr%2Cratebypass%2Cdur%2Clmt&sig=AOq0QJ8wRAIgUk8xVxYBWX6IqD3z4b-xrSbkjumytmAYEYHbrE78R7YCICgKgDP7En9USRd-zzrkCA6uWXmKzkrCvgnKsWPypt6R&lsparams=mh%2Cmm%2Cmn%2Cms%2Cmv%2Cmvi%2Cpl%2Cinitcwndbps&lsig=AG3C_xAwRgIhALF882h-ODFIyT1OWKNAJLuUgxcUptn7Rbm7gSRyhikZAiEA0iqZV6beENOBQvjq_OIYAva1DH_KLJdZ5ksUoFWlmQM%3D&title=He's%20Doing%20It%20Again", &setting)
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
		downloadjobs[part] = http.NewChunk(response, part, &wg, &setting)
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
