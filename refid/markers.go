package refid

type Marker interface {
	Generate() any
	IsValid(any) bool
}

var markers map[string]Marker

func init() {
	markers = make(map[string]Marker)
}
