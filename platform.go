package main

import (
	"fmt"
	"log"

	"github.com/veandco/go-sdl2/sdl"
)

func newPlatformRect(pos relativeRectPosition) *sdl.Rect {
	return &sdl.Rect{tileSourceWidth * int32(pos.xIndex), tileSourceHeight * int32(pos.yIndex), tileSourceWidth, tileSourceHeight}
}

type platformRects struct {
	topLeftRect   *sdl.Rect
	topMiddleRect *sdl.Rect
	topRightRect  *sdl.Rect
	midLeftRect   *sdl.Rect
	midMiddleRect *sdl.Rect
	midRightRect  *sdl.Rect
}

type platformDecoration struct {
	texture *sdl.Texture
	srcRect *sdl.Rect
	dstRect *sdl.Rect
}

func (pd *platformDecoration) draw(renderer *sdl.Renderer) {
	err := renderer.Copy(pd.texture, pd.srcRect, pd.dstRect)
	if err != nil {
		log.Fatalf("could not copy platform decoration texture: %v", err)
	}
}

type platform struct {
	x           int32
	y           int32
	w           int32
	h           int32
	texture     *sdl.Texture
	sourceRects platformRects
	decorations []platformDecoration
}

func newWalkablePlatform(x, y, w, h int32, texture *sdl.Texture) (platform, error) {
	walkablePlatformRects := platformRects{
		topLeftRect:   newPlatformRect(relativeRectPosition{10, 0}),
		topMiddleRect: newPlatformRect(relativeRectPosition{11, 0}),
		topRightRect:  newPlatformRect(relativeRectPosition{12, 0}),
		midLeftRect:   newPlatformRect(relativeRectPosition{10, 1}),
		midMiddleRect: newPlatformRect(relativeRectPosition{11, 1}),
		midRightRect:  newPlatformRect(relativeRectPosition{12, 1}),
	}
	return newPlatform(x, y, w, h, texture, walkablePlatformRects)
}

func newPlatform(x, y, w, h int32, texture *sdl.Texture, sourceRects platformRects) (platform, error) {
	if w < tileDestWidth*3 {
		return platform{}, fmt.Errorf("width value: %v must be higher (at least %v)", w, tileDestWidth*3)
	}
	return platform{x, y, w, h, texture, sourceRects, []platformDecoration{}}, nil
}

func (p *platform) addUpperLeftDecoration(x, y int32) error {
	topLeftDecorationRect := &sdl.Rect{tileSourceWidth*7 + 1, 0, tileSourceWidth, tileSourceHeight - 1}
	return p.addDecoration(topLeftDecorationRect, x, y)
}

func (p *platform) addUpperMiddleDecoration(x, y int32) error {
	topMiddleDecorationRect := &sdl.Rect{tileSourceWidth * 8, 0, tileSourceWidth, tileSourceHeight - 1}
	return p.addDecoration(topMiddleDecorationRect, x, y)
}

func (p *platform) addUpperRightDecoration(x, y int32) error {
	topRightDecorationRect := &sdl.Rect{tileSourceWidth * 9, 0, tileSourceWidth - 1, tileSourceHeight - 1}
	return p.addDecoration(topRightDecorationRect, x, y)
}

func (p *platform) addLowerMiddleDecoration(x, y int32) error {
	midMiddleDecorationRect := &sdl.Rect{tileSourceWidth*7 + 1, tileSourceHeight, tileSourceWidth - 2, tileSourceHeight - 1}
	return p.addDecoration(midMiddleDecorationRect, x, y)
}

// addDecoration adds a decoration tile from src of the platform texture to the position (relative to the platform)
func (p *platform) addDecoration(srcRect *sdl.Rect, x, y int32) error {
	if p.x-p.w/2+x+tileDestWidth > p.x+p.w/2 || p.x-p.w/2+x < p.x-p.w/2 {
		return fmt.Errorf("invalid decoration position x: %v. Decoration width exceeds platform width (%v)", x, p.w)
	}
	if p.y-p.h/2+y+tileDestHeight > p.y+p.h/2 || p.y-p.h/2+y < p.y-p.h/2 {
		return fmt.Errorf("invalid decoration position y: %v. Decoration height exceeds platform height (%v)", y, p.h)
	}
	dstRect := &sdl.Rect{p.x - p.w/2 + x, p.y - p.h/2 + y, tileDestWidth, tileDestHeight}
	pd := platformDecoration{p.texture, srcRect, dstRect}
	p.decorations = append(p.decorations, pd)
	return nil
}

func (p *platform) draw(renderer *sdl.Renderer) {
	// Top row
	p.drawRow(renderer, p.sourceRects.topLeftRect, p.sourceRects.topMiddleRect, p.sourceRects.topRightRect, 0)
	// Other rows
	for y := tileDestHeight; y < p.h; y += tileDestHeight - 1 {
		p.drawRow(renderer, p.sourceRects.midLeftRect, p.sourceRects.midMiddleRect, p.sourceRects.midRightRect, y)
	}
	for _, pd := range p.decorations {
		pd.draw(renderer)
	}
}

func (p *platform) drawRow(renderer *sdl.Renderer, tileLeftRect, tileMiddleRect, tileRightRect *sdl.Rect, y int32) {
	err := renderer.Copy(p.texture, tileLeftRect, &sdl.Rect{p.x - p.w/2, p.y - p.h/2 + y, tileDestWidth, tileDestHeight})
	if err != nil {
		log.Fatalf("could not copy platform left texture: %v", err)
	}
	for x := tileDestWidth; x < p.w-tileDestWidth; x += tileDestWidth {
		err = renderer.Copy(p.texture, tileMiddleRect, &sdl.Rect{p.x - p.w/2 + x, p.y - p.h/2 + y, tileDestWidth, tileDestHeight})
		if err != nil {
			log.Fatalf("could not copy platform middle texture: %v", err)
		}
	}
	err = renderer.Copy(p.texture, tileRightRect, &sdl.Rect{p.x + p.w/2 - tileDestWidth, p.y - p.h/2 + y, tileDestWidth, tileDestHeight})
	if err != nil {
		log.Fatalf("could not copy platform right texture: %v", err)
	}
}
