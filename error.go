package gefentoolbox

import (
	"fmt"
	"regexp"
	"strconv"
)

func asError(in string) error {
	if err := errParameterOutOfRangeFromString(in); err != nil {
		return err
	}
	return nil
}

type errParameterOutOfRange int

var matchErrParameterOutOfRange = regexp.MustCompile(`^\s*THE PARAMETER ((\d+) )?IS OUT OF RANGE\s*$`)

func (e errParameterOutOfRange) Error() string {
	return fmt.Sprintf("The parameter %d is out of range.", e)
}

func errParameterOutOfRangeFromString(in string) error {
	strs := matchErrParameterOutOfRange.FindStringSubmatch(in)
	if len(strs) > 0 {
		if strs[2] == "" {
			return errParameterOutOfRange(1)
		}
		i, err := strconv.Atoi(strs[2])
		if err != nil {
			return err
		}
		return errParameterOutOfRange(i)
	}
	return nil
}
