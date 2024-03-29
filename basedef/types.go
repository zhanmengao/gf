package basedef

type Int interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64
}
type Uint interface {
	~uint | ~uint8 | ~uint16 | ~uint32
}

type Ptr[T any] interface {
	*T
}
