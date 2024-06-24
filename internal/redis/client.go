package redis

import (
	"context"
	"fmt"
	"io"
	"net"
	"time"
)

type client struct {
	Addr string // redis addr
}

type Result struct {
	Data string
	Err  error
}

type Options struct {
	Addr string
}

func NewClient(options Options) *client {
	addr := "127.0.0.1:8080"
	if len(options.Addr) > 2 {
		addr = options.Addr
	}
	return &client{
		Addr: addr,
	}
}

func (c *client) Dial() (net.Conn, error) {
	return net.Dial("tcp", c.Addr)
}

func (c *client) Set(ctx context.Context, k, v string) (string, error) {
	resChn := make(chan Result, 1)
	start := time.Now()
	defer func() {
		fmt.Printf("took %v\n", time.Since(start))
	}()
	go c.set(k, v, resChn)
	select {
	case res := <-resChn:
		return res.Data, res.Err
	case <-ctx.Done():
		return "", fmt.Errorf("conn time out")
	}
}

func (c *client) Get(ctx context.Context, k string) (string, error) {
	resChn := make(chan Result, 1)
	go c.get(k, resChn)
	select {
	case res := <-resChn:
		return res.Data, res.Err
	case <-ctx.Done():
		return "", fmt.Errorf("conn time out")
	}
}

func (c *client) get(k string, resChn chan Result) {
	conn, err := c.Dial()
	if err != nil {
		resChn <- Result{Err: err}
		return
	}
	defer conn.Close()

	// Setting a deadline for the connection
	conn.SetDeadline(time.Now().Add(2 * time.Second))

	cmd := fmt.Sprintf("get %v", k)
	encodedCmd, err := Marshal(cmd)
	if err != nil {
		resChn <- Result{Err: err}
		return
	}

	// Write the cmd "get key"
	if _, err = conn.Write(encodedCmd); err != nil {
		resChn <- Result{Err: err}
		return
	}

	// Read server response
	buf, err := io.ReadAll(conn)
	if err != nil {
		resChn <- Result{Err: err}
		return
	}
	//unmarshal server response
	s, err := Unmarshal(string(buf))
	if err != nil {
		resChn <- Result{Err: err}
	}

	if str, ok := s.(string); ok {
		resChn <- Result{Err: nil, Data: str}
	}

	resChn <- Result{Err: fmt.Errorf("unknow server error")}
}

func (c *client) set(k, v string, resChn chan Result) {
	conn, err := c.Dial()
	if err != nil {
		resChn <- Result{Err: err}
		return
	}
	defer conn.Close()

	// Setting a deadline for the connection
	conn.SetDeadline(time.Now().Add(2 * time.Second))

	cmd := fmt.Sprintf("set %v %v", k, v)
	encodedCmd, err := Marshal(cmd)
	if err != nil {
		resChn <- Result{Err: err}
		return
	}

	// Write the cmd "set key value"
	if _, err = conn.Write(encodedCmd); err != nil {
		resChn <- Result{Err: err}
		return
	}

	// Read server response
	buf, err := io.ReadAll(conn)
	if err != nil {
		resChn <- Result{Err: err}
		return
	}

	resChn <- Result{Err: nil, Data: string(buf)}
}
