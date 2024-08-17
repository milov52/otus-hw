package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"time"
)

type TelnetClient interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
}

type telnetClient struct {
	address string
	timeout time.Duration
	conn    net.Conn
	in      io.ReadCloser
	out     io.Writer
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &telnetClient{
		address: address,
		timeout: timeout,
		in:      in,
		out:     out,
	}
}

func (tc *telnetClient) Connect() error {
	conn, err := net.DialTimeout("tcp", tc.address, tc.timeout)
	if err != nil {
		return err
	}
	tc.conn = conn
	fmt.Fprintf(os.Stderr, "...Connected to %s\n", tc.address)
	return nil
}

func (tc *telnetClient) Send() error {
	if tc.conn == nil {
		return fmt.Errorf("connection is not established")
	}
	_, err := io.Copy(tc.conn, tc.in)
	return err
}

func (tc *telnetClient) Receive() error {
	if tc.conn == nil {
		return fmt.Errorf("connection is not established")
	}
	_, err := io.Copy(tc.out, tc.conn)
	return err
}

func (tc *telnetClient) Close() error {
	if tc.conn != nil {
		fmt.Fprintln(os.Stderr, "...Connection closed")
		return tc.conn.Close()
	}
	return nil
}
