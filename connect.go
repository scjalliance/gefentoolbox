package gefentoolbox

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"regexp"
	"strings"
	"sync"
	"unicode"
)

// Switch is the Gefen Toolbox Matrix Switcher
type Switch struct {
	conn    net.Conn
	buffer  *bufio.ReadWriter
	welcome string
	sync.Mutex
}

var (
	prefixStrip = regexp.MustCompile(`^[\0\s]*`)
	promptStrip = regexp.MustCompile(`\s*telnet->\s*$`)
	matchModel  = regexp.MustCompile(`^\s*Welcome to (.*) TELNET\s*$`)
)

// Errors
var (
	ErrUnexpected       = errors.New("unexpected error")
	ErrOutputOutOfRange = errors.New("output value is out of range")
	ErrInputOutOfRange  = errors.New("input value is out of range")
)

// New connects to and then returns a *Switch
func New(addr, username, password string) (*Switch, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}
	s := &Switch{conn: conn}

	s.buffer = bufio.NewReadWriter(bufio.NewReader(conn), bufio.NewWriter(conn))

	for len(s.welcome) == 0 {
		welcome, err := s.buffer.ReadString('\n')
		if err != nil {
			return nil, err
		}
		s.welcome = strings.TrimSpace(welcome)
	}

	for {
		prompt, err := s.buffer.ReadString(':')
		if err != nil {
			return nil, err
		}
		prompt = strings.TrimFunc(prompt, func(x rune) bool { return !unicode.IsLetter(x) })
		if prompt == "UserID" {
			fmt.Fprintln(s.buffer, username)
			s.buffer.Flush()
			break
		}
	}

	for {
		prompt, err := s.buffer.ReadString(':')
		if err != nil {
			return nil, err
		}
		prompt = strings.TrimFunc(prompt, func(x rune) bool { return !unicode.IsLetter(x) })
		if prompt == "Password" {
			fmt.Fprintln(s.buffer, password)
			s.buffer.Flush()
			break
		}
	}

	_, err = s.waitForPrompt()
	if err != nil {
		return nil, err
	}

	return s, nil
}

func (s *Switch) waitForPrompt() (string, error) {
	lines, err := s.buffer.ReadString('>')
	if err != nil {
		return "", err
	}
	lines = prefixStrip.ReplaceAllString(lines, "")
	lines = promptStrip.ReplaceAllString(lines, "")
	return lines, nil
}

// Welcome shows the text that was presented during device login
func (s *Switch) Welcome() string {
	return s.welcome
}

// Model returns the switcher model, based on the Welcome() result
func (s *Switch) Model() string {
	strs := matchModel.FindStringSubmatch(s.welcome)
	if len(strs) > 1 {
		return strs[1]
	}
	return "Unknown"
}

// RawCommand sends a raw command to the switcher
func (s *Switch) RawCommand(cmd string) (string, error) {
	s.Lock()
	defer s.Unlock()
	fmt.Fprintf(s.buffer, "%s\r\n", cmd)
	s.buffer.Flush()

	lines, err := s.waitForPrompt()
	if err != nil {
		return "", err
	}

	if err = asError(lines); err != nil {
		return "", err
	}

	return lines, nil
}
