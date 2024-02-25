package typ

type BodyErrPair struct {
	Body []byte
	Err  error
}

func NewRcvBuf(body []byte, err error) *BodyErrPair {
	return &BodyErrPair{
		Err:  err,
		Body: body,
	}
}
