package database

type UTINYINT uint8 //

type MEDUIMINT int32   //
type UMEDUIMINT uint32 // 16777215 (unit32 is bigger so validate first)

type TINYTEXT string // TinyText (255)

type TIMESTAMP uint64

const MAX_UMEDUIMINT UMEDUIMINT = 16777215

func Valid_UMEDUIMINT(value UMEDUIMINT) bool {
	return UMEDUIMINT(value) <= MAX_UMEDUIMINT
}
