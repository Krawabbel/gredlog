package main

import (
	"fmt"
	"net"
	"strconv"
	"strings"
)

type Client interface {
	Request(string) (string, error)
	Restart() error
	Close() error
}

type client struct {
	conn net.Conn
}

func NewClient(addr string) (Client, error) {
	c := new(client)
	err := c.connect(addr)
	if err != nil {
		return nil, fmt.Errorf("establishing REDIS client connection failed: %s", err)
	}
	return c, nil
}

func (c *client) Close() error {
	return c.conn.Close()
}

func (c *client) Request(cmd string) (string, error) {
	tokens := strings.Split(cmd, " ")
	q := fmt.Sprintf("*%d\r\n", len(tokens))
	for _, tok := range tokens {
		q += fmt.Sprintf("$%d\r\n%s\r\n", len(tok), tok)
	}
	err := c.write(q)
	if err != nil {
		return "", fmt.Errorf("posting request failed: %s", err)
	}
	return c.read()
}

func (c *client) connect(addr string) error {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return fmt.Errorf("error connecting to %s (REDIS server probably not running): %s", addr, err)
	}
	c.conn = conn
	err = c.handshake("PING", "PONG")
	if err != nil {
		return fmt.Errorf("error connecting to %s (probably not a REDIS database): %s", addr, err)
	}
	return nil
}

func (c *client) Restart() error {
	addr := c.conn.RemoteAddr().String()
	_ = c.conn.Close()
	err := c.connect(addr)
	if err != nil {
		return fmt.Errorf("failed to reconnect to %s: %s", addr, err)
	}
	return nil
}

func (c *client) handshake(given, want string) error {
	have, err := c.Request(given)
	if err != nil {
		return fmt.Errorf("handshake failed: %s", err)
	}
	if have != want {
		return fmt.Errorf("handshake failed: given '%s', expected '%s', got '%s'", given, want, have)
	}
	return nil
}

func (c *client) write(s string) error {
	b := []byte(s)
	n, err := c.conn.Write(b)
	switch {
	case err != nil:
		return fmt.Errorf("error writing to REDIS database: %s", err)
	case n != len(b):
		return fmt.Errorf("error writing to REDIS database: only %d of %d bytes sent", n, len(b))
	default:
		return nil
	}
}

func (c *client) read() (string, error) {
	b, err := c.read_byte()
	if err != nil {
		return "", err
	}
	switch b {
	case '+':
		return c.read_simple()
	case '-':
		return c.read_error()
	case ':':
		return c.read_integer()
	case '$':
		return c.read_bulk()
	case '*':
		return c.read_array()
	default:
		return "", fmt.Errorf("expected '(+|-|:|$|*)', got '%s'", string(b))
	}
}

func (c *client) read_byte() (byte, error) {
	buf := make([]byte, 1)
	n, err := c.conn.Read(buf)
	if n != 1 {
		return 0, fmt.Errorf("expected one byte, got zero")
	}
	if err != nil {
		return 0, err
	}
	return buf[0], nil
}

func (c *client) expect_byte(want byte) error {
	have, err := c.read_byte()
	switch {
	case err != nil:
		return err
	case have != want:
		return fmt.Errorf("expected 0x%2x ('%s'), got 0x%2x ('%s')", want, string(want), have, string(have))
	default:
		return nil
	}
}

func (c *client) read_line() (string, error) {
	line := ""
	for {
		b, err := c.read_byte()
		switch {
		case err != nil:
			return line, err
		case b == '\r':
			return line, c.expect_byte('\n')
		default:
			line += string(b)
		}
	}
}

func (c *client) read_simple() (string, error) {
	return c.read_line()
}

func (c *client) read_error() (string, error) {
	r, err := c.read_line()
	if err != nil {
		return "", err
	}
	return "", fmt.Errorf("(error) " + r)
}

func (c *client) read_integer() (string, error) {
	r, err := c.read_line()
	if err != nil {
		return "", err
	}
	return "(integer) " + r, nil
}

func (c *client) read_bulk() (string, error) {
	n, err := c.read_line()
	if err != nil {
		return "", err
	}
	s, err := c.read_line()
	if err != nil {
		return "", err
	}
	if n != fmt.Sprint(len(s)) {
		return "", fmt.Errorf("resp error: expected bulk string length %s, got %d", n, len(s))
	}
	return "\"" + s + "\"", nil
}

func (c *client) read_array() (string, error) {
	n_str, err := c.read_line()
	if err != nil {
		return "", err
	}
	n, err := strconv.Atoi(n_str)
	if err != nil {
		return "", err
	}
	elements := make([]string, n)
	for i := range elements {
		e, err := c.read()
		if err != nil {
			return "", err
		}
		elements[i] = e
	}
	return "[" + strings.Join(elements, ", ") + "]", nil
}
