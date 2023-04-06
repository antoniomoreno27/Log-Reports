package generator

type Generator interface {
	Generate() error
	String() string
}
