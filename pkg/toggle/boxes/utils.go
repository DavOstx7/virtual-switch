package boxes


func extendStopFunc(stopFunc StopFunc, extension func() error) StopFunc {
	if extension == nil {
		return stopFunc
	}

	return func() error {
		if err := stopFunc(); err != nil {
			return err
		}
		return extension()
	}
}
