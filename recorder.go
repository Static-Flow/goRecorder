package main

import (
	"flag"
	"fmt"
	"github.com/kbinani/screenshot"
	term "github.com/nsf/termbox-go"
	"image"
	"image/png"
	"os"
	"strconv"
	"strings"
	"time"
)

func reset() {
	term.Sync()
}

func splitScreens(ids string) []int {
	idList := strings.Split(ids, ",")
	idIntList := []int{}
	for _, i := range idList {
		j, err := strconv.Atoi(i)
		if err != nil {
			panic(err)
		}
		idIntList = append(idIntList, j)
	}
	return idIntList
}

func getScreenInfo() {
	screensInfo := ""
	n := screenshot.NumActiveDisplays()
	for i := 0; i < n; i++ {
		bounds := screenshot.GetDisplayBounds(i)
		screensInfo += fmt.Sprintf("%d: %d by %d Monitor\n", i, bounds.Dx(), bounds.Dy())
	}
	fmt.Println(screensInfo)
}

func main() {

	screenPtr := flag.String("screenIds", "", "list of screen id to use, e.g. 0,1,n. Use -listScreens to see available ids")
	pollingPtr := flag.Int("pollingInt", 1, "How many seconds between screen grabs")
	bufferSizePtr := flag.Int("bufferSize", 60, "How many pictures to store in memory per screen")
	listScreensPtr := flag.Bool("listScreens", false, "List Screen Info")
	flag.Parse()
	if *listScreensPtr {
		getScreenInfo()
	} else {
		var pics [][]*image.RGBA
		numberOfScreens := splitScreens(*screenPtr)
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
						fileName := fmt.Sprintf("%d_%d-%02d-%02dT%02d.%02d.%02d.png", numberOfScreens[index],
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
					for i := 0; i < len(numberOfScreens); i++ {
						bounds := screenshot.GetDisplayBounds(numberOfScreens[i])

						img, err := screenshot.CaptureRect(bounds)
						if err != nil {
							panic(err)
						}
						screenPics = append(screenPics, img)
					}
					if len(pics) == *bufferSizePtr {
						pics = pics[1:]
					}
					pics = append(pics, screenPics)
				}
			}
		}
	}
}
