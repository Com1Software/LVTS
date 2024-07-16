package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	url := "http://com1software.com"
	url = "http://192.168.1.118:8080"
	tc := 0
	tctl := 0
	loop := false
	switch {
	case tctl == 0:
		resp, err := http.Get(url)
		if err != nil {
			log.Fatal(err)
		}
		reader := bufio.NewReader(resp.Body)
		if loop {
			x := 0
			for {
				line, erra := reader.ReadBytes('\n')
				if erra != nil {
					log.Fatal(err)
				}
				fmt.Printf("%s", string(line))
				if x > 10 {
					err = os.WriteFile("rec.dat", line, 0644)
					if err != nil {
						log.Fatal(err)
					}
					x = 0
				}
			}

		} else {

			line, erra := reader.ReadBytes('\n')
			if erra != nil {
				log.Fatal(err)
			}
			fmt.Printf("%s", string(line))
			err = os.WriteFile("rec.dat", line, 0644)
			if err != nil {
				log.Fatal(err)
			}
		}
	case tctl == 1:

		ctx := cancelCtxOnSigterm(context.Background())
		startWork(ctx, url)
		fmt.Println("test")
		fmt.Println(tc)

	}
}

func cancelCtxOnSigterm(ctx context.Context) context.Context {
	exitCh := make(chan os.Signal, 1)
	signal.Notify(exitCh, os.Interrupt, syscall.SIGTERM)
	ctx, cancel := context.WithCancel(ctx)
	go func() {
		<-exitCh
		cancel()
	}()
	return ctx
}
func startWork(ctx context.Context, url string) {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	reader := bufio.NewReader(resp.Body)
	line, erra := reader.ReadBytes('\n')
	if erra != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s", string(line))

	for {
		fmt.Println("xx")
		if err := work(ctx, url); err != nil {
			fmt.Printf("failed to do work: %s", err)
		}
		select {
		case <-ticker.C:
			fmt.Println("test")
			continue
		case <-ctx.Done():
			fmt.Println("oo")
			return
		}
	}
}

func work(ctx context.Context, url string) error {
	fmt.Println("doing work")

	return nil
}
