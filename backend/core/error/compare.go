package error

func IsNotFoundError(err error) bool {
	e, converr := ToError(err)
	if converr != nil {
		return false
	}

	return e.ErrorCode == ErrorCodeNotFound
}

func IsUnauthenticatedError(err error) bool {
	e, converr := ToError(err)
	if converr != nil {
		return false
	}

	return e.ErrorCode == ErrorCodeUnauthenticated
}

func IsUnavailableError(err error) bool {
	e, converr := ToError(err)
	if converr != nil {
		return false
	}

	return e.ErrorCode == ErrorCodeUnavailable
}
