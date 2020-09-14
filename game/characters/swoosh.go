package characters

import (
	"log"
	"simpleplatformer/common"
	"simpleplatformer/constants"

	"github.com/veandco/go-sdl2/sdl"
)

func updateSwooshes(sws []*swoosh) []*swoosh {
	result := []*swoosh{}
	for _, sw := range sws {
		if sw.destroyed == false {
			sw.update()
			result = append(result, sw)
		}
	}
	return result
}

type swoosh struct {
	time       int
	texture    *sdl.Texture
	rects      []*sdl.Rect
	x          int32
	y          int32
	w          int32
	h          int32
	facedRight bool
	destroyed  bool
}

func newSwoosh(tex *sdl.Texture, x, y int32, facedRight bool) *swoosh {
	rects := newCharacterAnimationRects([]common.RelativeRectPosition{
		{0, 0},
		{1, 0},
		{2, 0},
		{3, 0},
	})
	return &swoosh{
		time:       0,
		texture:    tex,
		x:          x,
		y:          y,
		w:          32,
		h:          32,
		rects:      rects,
		facedRight: facedRight,
		destroyed:  false,
	}
}

func (s *swoosh) update() {
	s.time++
	if s.time > len(s.rects)*10 {
		s.destroyed = true
	}
}

func (s *swoosh) draw(r *sdl.Renderer) {
	displayedFrame := s.time / 10 % len(s.rects)
	src := s.rects[displayedFrame]
	characterDestWidth := constants.CharacterDestWidth
	characterDestHeight := constants.CharacterDestHeight
	dst := &sdl.Rect{s.x - characterDestWidth/2, s.y - characterDestHeight/2, characterDestWidth, characterDestHeight}
	var flip sdl.RendererFlip
	if s.facedRight {
		flip = sdl.FLIP_NONE
	} else {
		flip = sdl.FLIP_HORIZONTAL
	}
	err := r.CopyEx(s.texture, src, dst, 0, nil, flip)
	if err != nil {
		log.Fatalf("could not copy Swoosh texture: %v", err)
	}
}
