 package main

import (
	"github.com/kbinani/screenshot"
	"image/png"
	"image"
	"os"
	"fmt"
	"flag"
	"time"
	term "github.com/nsf/termbox-go"
)

func reset() {
	term.Sync()
}

func main() {

	screenPtr := flag.Int("numScreens",screenshot.NumActiveDisplays(), "Number of screens to use")
	pollingPtr := flag.Int("pollingInt", 1, "How many seconds between screen grabs")
	bufferSizePtr := flag.Int("bufferSize", 60, "How many pictures to store in memory per screen")
	flag.Parse()

	var pics [][]*image.RGBA
	n := *screenPtr
	ticker := time.NewTicker(time.Duration(*pollingPtr) * time.Second)
	defer ticker.Stop()
	finished := make(chan bool)
	pause := false
	output := false
	previousNow := time.Now()
	err := term.Init()
	if err != nil {
		panic(err)
	}
	defer term.Close()
	go func() {
		for {
			switch ev := term.PollEvent(); ev.Type {
				case term.EventKey:
					switch ev.Key {
						case term.KeyEsc:
							reset()
							fmt.Println("Exiting")
							finished <- true
							return
						case term.KeyF1:
							reset()
							fmt.Println("Compiling Images")
							output = true
						case term.KeyF2:
							reset()
							if pause {
								pause = false
								fmt.Println("Resuming")
							} else {
								pause = true
								fmt.Println("Pausing")
							}
						default:
							reset()
					}
			}
		}
	}()
	for {
		if output {
			for _, picArray := range pics {
				for index, pic := range picArray {
					fileName := fmt.Sprintf("%d_%d-%02d-%02dT%02d.%02d.%02d.png",index, 
        previousNow.Year(), previousNow.Month(), previousNow.Day(),
        previousNow.Hour(), previousNow.Minute(), previousNow.Second())
					file, err := os.Create(fileName)
					if err != nil {
						panic(err)
					}
					defer file.Close()
					png.Encode(file, pic)
				}
			}
			output = false
		} else if !pause {
			select {
				case <-finished:
					ticker.Stop()
					return
				case t := <-ticker.C:
					previousNow = t
					fmt.Println("Current time: ", t.UTC().Format(time.UnixDate))
					var screenPics []*image.RGBA
					for i := 0; i < n; i++ {
			 			bounds := screenshot.GetDisplayBounds(i)

				 		img, err := screenshot.CaptureRect(bounds)
				 		if err != nil {
			 				panic(err)
			 			}
			 			screenPics = append(screenPics,img)
				 	}
				 	if len(pics) == *bufferSizePtr{
			 			pics = pics[1:]
			 		}
				 	pics = append(pics,screenPics)
			}
		}
	}
}