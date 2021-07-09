package config

//number of ships
const (
	PATROL_CRAFT_COUNT = 4
	SUBMARINE_COUNT    = 3
	DESTROYER_COUNT    = 2
	BATTLESHIP_COUNT   = 1
)

//ship lives
const (
	PATROL_CRAFT = iota + 1
	SUBMARINE    = 2
	DESTROYER    = 3
	BATTLESHIP   = 4
)

//fielld marks (ascii)
const (
	SHIP                 = "#"
	MISSED               = "0"
	HIT                  = "X"
	KILL                 = "X"
	FIELD                = "."
	COLUMNS              = 10
	ROWS                 = 10
	COORDINATE_SEPARATOR = "x"
	COLUMN_START         = "A"
	COLUMN_END           = "J"
	ROW_START            = 1
	ROW_END              = 10
)
