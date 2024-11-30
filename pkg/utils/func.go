package utils

func ExecuteFunc(function func() error) error {
	var err error = nil
	if function != nil {
		err = function()
	}
	return err
}
