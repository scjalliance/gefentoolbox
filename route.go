package gefentoolbox

import (
	"fmt"
	"regexp"
	"strconv"
)

// Route is a route
type Route struct {
	Input      int
	InputName  string
	Output     int
	OutputName string
}

var routeMatch = regexp.MustCompile(`^\s*OUTPUT (\d+)(\((.*)\))? IS ROUTED TO INPUT (\d+)(\((.*)\))?\s*`)

// Route sets the input for a particular output
func (s *Switch) Route(output, input int) (string, error) {
	result, err := s.RawCommand(fmt.Sprintf("r %d %d", output, input))
	if i, ok := err.(errParameterOutOfRange); ok {
		if i == 1 {
			return "", ErrOutputOutOfRange
		}
		if i == 2 {
			return "", ErrInputOutOfRange
		}
	}
	return result, err
}

// GetRoute gets the input for a particular output
func (s *Switch) GetRoute(output int) (*Route, error) {
	result, err := s.RawCommand(fmt.Sprintf("#show_r %d", output))
	if i, ok := err.(errParameterOutOfRange); ok {
		if i == 1 {
			return nil, ErrOutputOutOfRange
		}
	}
	strs := routeMatch.FindStringSubmatch(result)
	if len(strs) > 0 {
		input, err := strconv.Atoi(strs[4])
		if err != nil {
			return nil, err
		}
		output, err := strconv.Atoi(strs[1])
		if err != nil {
			return nil, err
		}
		return &Route{
			Input:      input,
			InputName:  strs[6],
			Output:     output,
			OutputName: strs[3],
		}, nil
	}
	return nil, ErrUnexpected
}
