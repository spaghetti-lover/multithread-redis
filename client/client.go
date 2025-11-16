package client

import (
	"bufio"
	"fmt"
	"net"
	"time"

	"github.com/spaghetti-lover/multithread-redis/internal/core"
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
	cmd := []interface{}{"SET", key, value}
	return c.sendCommand(cmd)
}

// SetEx sets a key-value pair with expiration in seconds
func (c *Client) SetEx(key, value string, seconds int64) error {
	cmd := []interface{}{"SET", key, value, "EX", fmt.Sprintf("%d", seconds)}
	return c.sendCommand(cmd)
}

// Get retrieves a value by key
func (c *Client) Get(key string) (string, error) {
	cmd := []interface{}{"GET", key}
	resp, err := c.sendCommandWithResponse(cmd)
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
	cmd := []interface{}{"PING"}
	resp, err := c.sendCommandWithResponse(cmd)
	if err != nil {
		return "", err
	}

	if str, ok := resp.(string); ok {
		return str, nil
	}

	return "PONG", nil
}

// sendCommand sends a command without waiting for response
func (c *Client) sendCommand(cmd []interface{}) error {
	respData := core.Encode(cmd, false)
	_, err := c.conn.Write(respData)
	return err
}

// sendCommandWithResponse sends a command and waits for response
func (c *Client) sendCommandWithResponse(cmd []interface{}) (interface{}, error) {
	// Send command
	respData := core.Encode(cmd, false)
	_, err := c.conn.Write(respData)
	if err != nil {
		return nil, fmt.Errorf("failed to send command: %w", err)
	}

	// Read response
	reader := bufio.NewReader(c.conn)

	// Read first byte to determine response type
	_, err = reader.ReadByte()
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// Unread the first byte so Decode can read it
	reader.UnreadByte()

	// Read all available data
	buf := make([]byte, 4096)
	n, err := reader.Read(buf)
	if err != nil {
		return nil, fmt.Errorf("failed to read response data: %w", err)
	}

	// Decode RESP response
	decoded, err := core.Decode(buf[:n])
	if err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return decoded, nil
}
