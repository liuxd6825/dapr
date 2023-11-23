package utils

import "errors"

func GetRecoverError(err error, recErr any) (resErr error) {
	if err != nil {
		return err
	}
	if recErr != nil {
		switch recErr.(type) {
		case string:
			{
				msg, _ := recErr.(string)
				resErr = errors.New(msg)
			}
		case error:
			{
				if e, ok := recErr.(error); ok {
					resErr = e
				}
			}
		}
	}
	return resErr
}
