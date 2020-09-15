package constants

const (
	WindowWidth         = 860
	WindowHeight        = 510
	Gravity             = 0.05
	JumpSpeed           = 4
	CharacterXSpeed     = 1
	CharacterStaminaMax = 30
	SwooshXSpeed        = 1
	SwooshXShift        = 10
)

const (
	scaleX                = WindowWidth / 288
	scaleY                = WindowHeight / 172
	TileSourceWidth       = int32(16)
	TileSourceHeight      = int32(128 / 8)
	TileDestWidth         = int32(TileSourceWidth * scaleX)
	TileDestHeight        = int32(TileSourceHeight * scaleY)
	CharacterSourceWidth  = int32(32)
	CharacterSourceHeight = int32(32)
	CharacterDestWidth    = int32(CharacterSourceWidth * scaleX)
	CharacterDestHeight   = int32(CharacterSourceHeight * scaleY)
)
