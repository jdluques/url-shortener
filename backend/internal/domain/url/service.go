package url

type IDGenerator interface {
	NextID() (int64, error)
}

type ShortCodeGenerator interface {
	Generate(id int64) (string, error)
}
