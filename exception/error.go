package exception

//

type NotFoundError struct {
	Entity string
}

func (e *NotFoundError) Error() string {
	return e.Entity + " not found"
}

//

type BadRequestError struct {
	Message string
}

func (e *BadRequestError) Error() string {
	return e.Message
}
