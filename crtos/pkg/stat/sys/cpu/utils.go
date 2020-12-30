package cpu

import (
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

func readFile(path string) (string, error) {
	contents, err := ioutil.ReadFile(path)
	if err != nil {
		return "", errors.Wrapf(err, "cpu/stat read file failed!", path)
	}
	return strings.TrimSpace(string(contents)), err
}

func parseUint(s string) (uint64, error) {
	v, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		intValue, intErr := strconv.ParseInt(s, 10, 64)

		if intErr == nil && intValue < 0 {
			return 0, nil
		} else if intErr != nil &&
			intErr.(*strconv.NumError).Err == strconv.ErrRange &&
			intValue < 0 {
			return 0, nil
		}
		return 0, errors.Wrapf(err, "parseUint failed(%s)", s)
	}
	return v, nil
}
