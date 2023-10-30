package utils

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

func GetErrorWithStack(err error) string {
	if err == nil {
		return ""
	}

	causeErr := errors.Cause(err)
	lastErr := err
	for {
		st := getStackTrace(err)
		if err.Error() == causeErr.Error() && len(st) == 0 {
			return fmt.Sprintf("%+v", lastErr)
		}
		if len(st) != 0 {
			lastErr = err
		}
		err = errors.Unwrap(err)
		if err == nil {
			return fmt.Sprintf("%+v", lastErr)
		}
	}
}

func StringToLines(s string) []string {
	lines := strings.Split(s, "\n")
	res := []string{}
	for _, line := range lines {
		line = strings.ReplaceAll(line, "\t", "    ")
		res = append(res, line)
	}
	return res
}

func getStackTrace(err error) []errors.Frame {
	type stackTracer interface {
		StackTrace() errors.StackTrace
	}
	//nolint:errorlint // We want to check if the error implements the stackTracer interface
	if err, ok := err.(stackTracer); ok {
		frames := []errors.Frame{}
		for _, f := range err.StackTrace() {
			frames = append(frames, f)
		}
		return frames
	}
	return nil
}
