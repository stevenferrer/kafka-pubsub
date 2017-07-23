package usermanagementsvc

import "errors"

func err2str(err error) string {
	if err == nil {
		return ""
	}

	return err.Error()
}

func str2err(s string) error {
	if s == "" {
		return nil
	}

	return errors.New(s)
}
