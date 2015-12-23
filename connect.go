package gefentoolbox

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"unicode"
)

// Switch is the Gefen Toolbox Matrix Switcher
type Switch struct {
	conn    net.Conn
	buffer  *bufio.ReadWriter
	welcome string
}

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

	fmt.Println("Waiting for Username prompt...")

	for {
		prompt, err := s.buffer.ReadString(':')
		if err != nil {
			return nil, err
		}
		prompt = strings.TrimFunc(prompt, func(x rune) bool { return !unicode.IsLetter(x) })
		fmt.Printf("username[%s]\n", prompt)
		if prompt == "UserID" {
			fmt.Println("Wants UserID.")
			fmt.Fprintln(s.buffer, username)
			s.buffer.Flush()
			break
		}
	}

	fmt.Println("Waiting for Password prompt...")

	for {
		prompt, err := s.buffer.ReadString(':')
		if err != nil {
			return nil, err
		}
		prompt = strings.TrimFunc(prompt, func(x rune) bool { return !unicode.IsLetter(x) })
		fmt.Printf("password[%s]\n", prompt)
		if prompt == "Password" {
			fmt.Println("Wants Password.")
			fmt.Fprintln(s.buffer, password)
			s.buffer.Flush()
			break
		}
	}

	line, err := s.buffer.ReadString('>')
	if err != nil {
		return nil, err
	}
	fmt.Println(line)

	return s, nil
}

// Welcome shows the text that was presented during device login
func (s *Switch) Welcome() string {
	return s.welcome
}

// RawCommand sends a raw command to the switcher
func (s *Switch) RawCommand(cmd string) (string, error) {
	fmt.Println(cmd)
	fmt.Fprintf(s.buffer, "%s\r\n", cmd)
	s.buffer.Flush()
	return "", nil
}
