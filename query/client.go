package query

import (
	"bufio"
	"fmt"
	"net"
)

type Request struct {
}

type Client interface {
	Connect(host string) error
	IsConnected() bool
	Query() (string, error)
}

type queryClient struct {
	Client
	c net.Conn
}

func NewClient() Client {
	return &queryClient{}
}

func (c *queryClient) Connect(host string) error {
	conn, err := net.Dial("udp", host)
	if err != nil {
		err = fmt.Errorf("creating connection to %s failed: %w", host, err)
		return err
	}

	c.c = conn

	return nil
}

func (c *queryClient) IsConnected() bool {
	return c.c != nil
}

func (c *queryClient) Query() (string, error) {
	p := make([]byte, 2048)
	if _, err := fmt.Fprintf(c.c, "Hi UDP Server, How are you doing?"); err != nil {
		return "", err
	}
	if _, err := bufio.NewReader(c.c).Read(p); err == nil {
		fmt.Printf("%s\n", p)
	} else {
		fmt.Printf("Some error %v\n", err)
	}
	fmt.Printf("%s\n", p)
	return "", nil
}
