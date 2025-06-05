package exception

type CustomError struct {
	Code    int
	Message string
}

func (c CustomError) Error() string {
	return c.Message
}
