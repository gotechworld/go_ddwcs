package mapper

type Mapper interface {
	MapLoginToken([]byte) string
	MapAwbInfo([]byte, string) interface{}
}
