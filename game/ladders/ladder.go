package ladders

import (
	"fmt"
	"log"
	"simpleplatformer/common"
	"simpleplatformer/constants"

	"github.com/veandco/go-sdl2/sdl"
)

type ladderRects struct {
	topRect *sdl.Rect
	midRect *sdl.Rect
	botRect *sdl.Rect
}

func NewLadder(x, y, w, h int32, texture *sdl.Texture) (Ladder, error) {
	if w < constants.TileDestWidth {
		return Ladder{}, fmt.Errorf("invalid ladder width provided: %v. Must be at least %v", w, constants.TileDestWidth)
	}
	if h < 2*constants.TileDestHeight {
		return Ladder{}, fmt.Errorf("invalid ladder height provided: %v. Must be at least %v", h, 2*constants.TileDestHeight)
	}
	rects := ladderRects{
		topRect: newLadderRect(common.RelativeRectPosition{7, 4}),
		midRect: newLadderRect(common.RelativeRectPosition{7, 5}),
		botRect: newLadderRect(common.RelativeRectPosition{7, 6}),
	}
	return Ladder{x, y, w, h, texture, rects}, nil
}

// TODO: Duplication from platforms package, refactor.
func newLadderRect(pos common.RelativeRectPosition) *sdl.Rect {
	return &sdl.Rect{
		constants.TileSourceWidth*int32(pos.XIndex) - 1,
		constants.TileSourceHeight * int32(pos.YIndex),
		constants.TileSourceWidth,
		constants.TileSourceHeight,
	}
}

type Ladder struct {
	X           int32
	Y           int32
	W           int32
	H           int32
	texture     *sdl.Texture
	sourceRects ladderRects
}

func (l *Ladder) Draw(renderer *sdl.Renderer) {
	dst := &sdl.Rect{l.X - l.W/2, l.Y - l.H/2, constants.TileDestWidth, constants.TileDestHeight}
	// Draw top
	err := renderer.Copy(l.texture, l.sourceRects.topRect, dst)
	if err != nil {
		log.Fatalf("could not copy ladder texture (top): %v", err)
	}
	// Draw bottom
	dst = &sdl.Rect{l.X - l.W/2, l.Y + l.H/2 - constants.TileDestHeight, constants.TileDestWidth, constants.TileDestHeight}
	err = renderer.Copy(l.texture, l.sourceRects.botRect, dst)
	if err != nil {
		log.Fatalf("could not copy ladder texture (bottom): %v", err)
	}
	// Draw the rest
	for tempY := l.Y - l.H/2 + constants.TileDestHeight; tempY < l.Y+l.H/2-constants.TileDestHeight; tempY += constants.TileDestHeight {
		dst.Y = tempY
		err = renderer.Copy(l.texture, l.sourceRects.midRect, dst)
		if err != nil {
			log.Fatalf("could not copy ladder texture (middle): %v", err)
		}
	}
}
