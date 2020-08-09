package main

import (
	"fmt"
	"log"
	"time"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
)

const (
	windowWidth  = 800
	windowHeight = 600
)

func main() {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}
	defer sdl.Quit()

	window, err := sdl.CreateWindow("Platformer", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		int32(windowWidth), int32(windowHeight), sdl.WINDOW_SHOWN)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer window.Destroy()

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer renderer.Destroy()
	sdl.SetHint(sdl.HINT_RENDER_SCALE_QUALITY, "1")

	// explosionBytes, audioSpec := sdl.LoadWAV("balloons/explode.wav")
	// if audioSpec == nil {
	// 	log.Println(sdl.GetError())
	// 	return
	// }
	// audioID, err := sdl.OpenAudioDevice("", false, audioSpec, nil, 0)
	// if err != nil {
	// 	panic(err)
	// }
	// defer sdl.CloseAudioDevice(audioID)
	// defer sdl.FreeWAV(explosionBytes)

	// audioState := audioState{explosionBytes, audioID, audioSpec}

	// go func() {
	// 	sdl.Delay(5000)
	// 	e := sdl.QuitEvent{Type: sdl.QUIT}
	// 	sdl.PushEvent(&e)
	// }()

	var elapsedTime float32

	renderer.SetDrawColor(uint8(66), uint8(135), uint8(245), uint8(0))

	tex, err := img.LoadTexture(renderer, "assets/sheet_9.png")
	if err != nil {
		log.Fatal("could not load texture: %v", err)
	}
	tileSourceHeight := int32(128 / 8)
	tileSourceWidth := int32(16)
	tileStartRect := &sdl.Rect{tileSourceWidth * 7, 0, tileSourceWidth, tileSourceHeight}
	tileMiddleRect := &sdl.Rect{tileSourceWidth * 8, 0, tileSourceWidth, tileSourceHeight}
	tileEndRect := &sdl.Rect{tileSourceWidth * 9, 0, tileSourceWidth, tileSourceHeight}
	tileDestHeight := int32((windowHeight / tileSourceHeight))
	tileDestWidth := int32((windowWidth / tileSourceWidth))

	running := true
	for running {
		frameStart := time.Now()
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				println("Quit")
				running = false
				break
			}
		}
		renderer.Clear()
		renderer.Copy(tex, tileStartRect, &sdl.Rect{0, 100, tileDestWidth, tileDestHeight})
		for i := 1; i < 7; i++ {
			renderer.Copy(tex, tileMiddleRect, &sdl.Rect{int32(i) * tileDestWidth, 100, tileDestWidth, tileDestHeight})
		}
		renderer.Copy(tex, tileEndRect, &sdl.Rect{7 * tileDestWidth, 100, tileDestWidth, tileDestHeight})

		renderer.Present()
		elapsedTime = float32(time.Since(frameStart).Seconds() * 1000)
		// fmt.Println("ms per frame:", elapsedTime)
		if elapsedTime < 5 {
			sdl.Delay(5 - uint32(elapsedTime))
			elapsedTime = float32(time.Since(frameStart).Seconds() * 1000)
		}
	}
}
