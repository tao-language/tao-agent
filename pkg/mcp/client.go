package mcp

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os/exec"
)

type Client struct {
	cmd    *exec.Cmd
	stdin  *bufio.Writer
	stdout *bufio.Scanner
}

func NewClient(command string, args []string) (*Client, error) {
	cmd := exec.Command(command, args...)
	stdinPipe, err := cmd.StdinPipe()
	if err != nil {
		return nil, err
	}
	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	if err := cmd.Start(); err != nil {
		return nil, err
	}

	return &Client{
		cmd:    cmd,
		stdin:  bufio.NewWriter(stdinPipe),
		stdout: bufio.NewScanner(stdoutPipe),
	}, nil
}

func (c *Client) Call(method string, params interface{}) (*JSONRPCResponse, error) {
	req := JSONRPCRequest{
		JSONRPC: "2.0",
		Method:  method,
		Params:  params,
		ID:      1, // Simplified ID management
	}

	data, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	if _, err := fmt.Fprintf(c.stdin, "%s\n", data); err != nil {
		return nil, err
	}
	if err := c.stdin.Flush(); err != nil {
		return nil, err
	}

	if !c.stdout.Scan() {
		return nil, fmt.Errorf("failed to read response: %v", c.stdout.Err())
	}

	var resp JSONRPCResponse
	if err := json.Unmarshal(c.stdout.Bytes(), &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

func (c *Client) Close() error {
	return c.cmd.Process.Kill()
}
