package constants

const (
	WindowWidth         = 860
	WindowHeight        = 510
	Gravity             = 0.05
	JumpSpeed           = 4
	CharacterVX         = 1
	CharacterStaminaMax = 30
	SwooshVX            = float32(1.0)
	SwooshXShift        = 10
	DefaultEnemyHealth  = 1
	DefaultPlayerHealth = 1
	HitStateLength      = 70
	CharacterVYWhenHit  = -2
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
