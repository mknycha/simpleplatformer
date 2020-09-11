package platforms

import (
	"fmt"
	"log"
	"simpleplatformer/common"
	"simpleplatformer/constants"

	"github.com/veandco/go-sdl2/sdl"
)

func newPlatformRect(pos common.RelativeRectPosition) *sdl.Rect {
	return &sdl.Rect{
		constants.TileSourceWidth * int32(pos.XIndex),
		constants.TileSourceHeight * int32(pos.YIndex),
		constants.TileSourceWidth,
		constants.TileSourceHeight,
	}
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

type Platform struct {
	X           int32
	Y           int32
	W           int32
	H           int32
	texture     *sdl.Texture
	sourceRects platformRects
	decorations []platformDecoration
}

func NewWalkablePlatform(x, y, w, h int32, texture *sdl.Texture) (Platform, error) {
	walkablePlatformRects := platformRects{
		topLeftRect:   newPlatformRect(common.RelativeRectPosition{10, 0}),
		topMiddleRect: newPlatformRect(common.RelativeRectPosition{11, 0}),
		topRightRect:  newPlatformRect(common.RelativeRectPosition{12, 0}),
		midLeftRect:   newPlatformRect(common.RelativeRectPosition{10, 1}),
		midMiddleRect: newPlatformRect(common.RelativeRectPosition{11, 1}),
		midRightRect:  newPlatformRect(common.RelativeRectPosition{12, 1}),
	}
	return newPlatform(x, y, w, h, texture, walkablePlatformRects)
}

func newPlatform(x, y, w, h int32, texture *sdl.Texture, sourceRects platformRects) (Platform, error) {
	if w < constants.TileDestWidth*3 {
		return Platform{}, fmt.Errorf("width value: %v must be higher (at least %v)", w, constants.TileDestWidth*3)
	}
	return Platform{x, y, w, h, texture, sourceRects, []platformDecoration{}}, nil
}

func (p *Platform) AddUpperLeftDecoration(x, y int32) error {
	topLeftDecorationRect := &sdl.Rect{constants.TileSourceWidth*7 + 1, 0, constants.TileSourceWidth, constants.TileSourceHeight - 1}
	return p.addDecoration(topLeftDecorationRect, x, y)
}

func (p *Platform) AddUpperMiddleDecoration(x, y int32) error {
	topMiddleDecorationRect := &sdl.Rect{constants.TileSourceWidth * 8, 0, constants.TileSourceWidth, constants.TileSourceHeight - 1}
	return p.addDecoration(topMiddleDecorationRect, x, y)
}

func (p *Platform) AddUpperRightDecoration(x, y int32) error {
	topRightDecorationRect := &sdl.Rect{constants.TileSourceWidth * 9, 0, constants.TileSourceWidth - 1, constants.TileSourceHeight - 1}
	return p.addDecoration(topRightDecorationRect, x, y)
}

func (p *Platform) AddLowerMiddleDecoration(x, y int32) error {
	midMiddleDecorationRect := &sdl.Rect{constants.TileSourceWidth*7 + 1, constants.TileSourceHeight, constants.TileSourceWidth - 2, constants.TileSourceHeight - 1}
	return p.addDecoration(midMiddleDecorationRect, x, y)
}

// addDecoration adds a decoration tile from src of the platform texture to the position (relative to the platform)
func (p *Platform) addDecoration(srcRect *sdl.Rect, x, y int32) error {
	if p.X-p.W/2+x+constants.TileDestWidth > p.X+p.W/2 || p.X-p.W/2+x < p.X-p.W/2 {
		return fmt.Errorf("invalid decoration position x: %v. Decoration width exceeds platform width (%v)", x, p.W)
	}
	if p.Y-p.H/2+y+constants.TileDestHeight > p.Y+p.H/2 || p.Y-p.H/2+y < p.Y-p.H/2 {
		return fmt.Errorf("invalid decoration position y: %v. Decoration height exceeds platform height (%v)", y, p.H)
	}
	dstRect := &sdl.Rect{p.X - p.W/2 + x, p.Y - p.H/2 + y, constants.TileDestWidth, constants.TileDestHeight}
	pd := platformDecoration{p.texture, srcRect, dstRect}
	p.decorations = append(p.decorations, pd)
	return nil
}

func (p *Platform) Draw(renderer *sdl.Renderer) {
	// Top row
	p.drawRow(renderer, p.sourceRects.topLeftRect, p.sourceRects.topMiddleRect, p.sourceRects.topRightRect, 0)
	// Other rows
	for y := constants.TileDestHeight; y < p.H; y += constants.TileDestHeight - 1 {
		p.drawRow(renderer, p.sourceRects.midLeftRect, p.sourceRects.midMiddleRect, p.sourceRects.midRightRect, y)
	}
	for _, pd := range p.decorations {
		pd.draw(renderer)
	}
}

func (p *Platform) drawRow(renderer *sdl.Renderer, tileLeftRect, tileMiddleRect, tileRightRect *sdl.Rect, y int32) {
	err := renderer.Copy(p.texture, tileLeftRect, &sdl.Rect{p.X - p.W/2, p.Y - p.H/2 + y, constants.TileDestWidth, constants.TileDestHeight})
	if err != nil {
		log.Fatalf("could not copy platform left texture: %v", err)
	}
	tileDestWidth := constants.TileDestWidth
	tileDestHeight := constants.TileDestHeight
	for x := tileDestWidth; x < p.W-tileDestWidth; x += tileDestWidth {
		err = renderer.Copy(p.texture, tileMiddleRect, &sdl.Rect{p.X - p.W/2 + x, p.Y - p.H/2 + y, tileDestWidth, tileDestHeight})
		if err != nil {
			log.Fatalf("could not copy platform middle texture: %v", err)
		}
	}
	err = renderer.Copy(p.texture, tileRightRect, &sdl.Rect{p.X + p.W/2 - tileDestWidth, p.Y - p.H/2 + y, tileDestWidth, tileDestHeight})
	if err != nil {
		log.Fatalf("could not copy platform right texture: %v", err)
	}
}
