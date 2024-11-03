package core

type Encoder[T any] interface {
	//Encode(io.Writer, T) error
	Encode(T) error
}

type Decoder[T any] interface {
	//Decode(io.Reader, T) error
	Decode(T) error
}
