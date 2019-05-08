package sca

import (
	"errors"
	"io"
	"strings"
)

func Kill(closers ...io.Closer) (err error) {
	var errs []string
	for i := range closers {
		if c := closers[i]; c != nil {
			if err := c.Close(); err != nil {
				errs = append(errs, err.Error())
			}
		}
	}

	if len(errs) != 0 {
		err = errors.New(strings.Join(errs, "\n"))
	}
	return
}
