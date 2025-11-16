package client

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
)

type Client struct {
	conn net.Conn
	addr string
}

// NewClient creates a new Redis client
func NewClient(addr string) (*Client, error) {
	conn, err := net.DialTimeout("tcp", addr, 5*time.Second)
	if err != nil {
		return nil, fmt.Errorf("failed to connect: %w", err)
	}

	return &Client{
		conn: conn,
		addr: addr,
	}, nil
}

// Close closes the connection
func (c *Client) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

// Set sets a key-value pair
func (c *Client) Set(key, value string) error {
	cmd := encodeArray([]string{"SET", key, value})
	_, err := c.conn.Write(cmd)
	if err != nil {
		return err
	}

	// Read OK response
	_, err = c.readResponse()
	return err
}

// SetEx sets a key-value pair with expiration in seconds
func (c *Client) SetEx(key, value string, seconds int64) error {
	cmd := encodeArray([]string{"SET", key, value, "EX", fmt.Sprintf("%d", seconds)})
	_, err := c.conn.Write(cmd)
	if err != nil {
		return err
	}

	// Read OK response
	_, err = c.readResponse()
	return err
}

// Get retrieves a value by key
func (c *Client) Get(key string) (string, error) {
	cmd := encodeArray([]string{"GET", key})
	_, err := c.conn.Write(cmd)
	if err != nil {
		return "", err
	}

	resp, err := c.readResponse()
	if err != nil {
		return "", err
	}

	if resp == nil {
		return "", fmt.Errorf("key not found")
	}

	if str, ok := resp.(string); ok {
		return str, nil
	}

	return fmt.Sprintf("%v", resp), nil
}

// Ping checks if the server is alive
func (c *Client) Ping() (string, error) {
	cmd := encodeArray([]string{"PING"})
	_, err := c.conn.Write(cmd)
	if err != nil {
		return "", err
	}

	resp, err := c.readResponse()
	if err != nil {
		return "", err
	}

	if str, ok := resp.(string); ok {
		return str, nil
	}

	return "PONG", nil
}

// readResponse reads and decodes RESP response
func (c *Client) readResponse() (interface{}, error) {
	reader := bufio.NewReader(c.conn)

	line, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}

	if len(line) < 3 {
		return nil, fmt.Errorf("invalid response")
	}

	line = strings.TrimSuffix(line, "\r\n")

	switch line[0] {
	case '+': // Simple String
		return line[1:], nil
	case '-': // Error
		return nil, fmt.Errorf(line[1:])
	case ':': // Integer
		val, err := strconv.ParseInt(line[1:], 10, 64)
		if err != nil {
			return nil, err
		}
		return val, nil
	case '$': // Bulk String
		length, err := strconv.Atoi(line[1:])
		if err != nil {
			return nil, err
		}
		if length == -1 {
			return nil, nil // Null bulk string
		}

		buf := make([]byte, length+2) // +2 for \r\n
		_, err = reader.Read(buf)
		if err != nil {
			return nil, err
		}
		return string(buf[:length]), nil
	case '*': // Array
		count, err := strconv.Atoi(line[1:])
		if err != nil {
			return nil, err
		}
		if count == -1 {
			return nil, nil
		}

		result := make([]interface{}, count)
		for i := 0; i < count; i++ {
			val, err := c.readResponse()
			if err != nil {
				return nil, err
			}
			result[i] = val
		}
		return result, nil
	default:
		return nil, fmt.Errorf("unknown response type: %c", line[0])
	}
}

// encodeArray encodes string array to RESP array format
func encodeArray(args []string) []byte {
	var result strings.Builder

	result.WriteString(fmt.Sprintf("*%d\r\n", len(args)))
	for _, arg := range args {
		result.WriteString(fmt.Sprintf("$%d\r\n%s\r\n", len(arg), arg))
	}

	return []byte(result.String())
}
