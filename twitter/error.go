package twitter

// NotFoundError not found error
type NotFoundError struct {
	Msg string
}

func (n NotFoundError) Error() string {
	return n.Msg
}
