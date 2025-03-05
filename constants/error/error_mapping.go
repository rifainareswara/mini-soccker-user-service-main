package error

func ErrMapping(err error) bool {
	allErrors := make([]error, 0)
	allErrors = append(append(GeneralErrors[:], UserErrors[:]...))

	for _, item := range allErrors {
		if err.Error() == item.Error() {
			return true
		}
	}

	return false
}
