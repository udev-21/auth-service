package util

func GetErrorIfExist(methods ...func() error) error {
	for _, method := range methods {
		if err := method(); err != nil {
			return err
		}
	}
	return nil
}
