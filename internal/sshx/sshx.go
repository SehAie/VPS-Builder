package sshx

import (
	"bytes"
	"fmt"
	"net"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"
)

type Client struct {
	addr   string
	user   string
	config *ssh.ClientConfig
}

func NewPassword(user, host string, port int, password string, timeout time.Duration) *Client {
	return &Client{addr: fmt.Sprintf("%s:%d", host, port), user: user, config: &ssh.ClientConfig{
		User:            user,
		Auth:            []ssh.AuthMethod{ssh.Password(password)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         timeout,
	}}
}

func NewKey(user, host string, port int, signer ssh.Signer, timeout time.Duration) *Client {
	return &Client{addr: fmt.Sprintf("%s:%d", host, port), user: user, config: &ssh.ClientConfig{
		User:            user,
		Auth:            []ssh.AuthMethod{ssh.PublicKeys(signer)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         timeout,
	}}
}

func (c *Client) Run(script string) (string, error) {
	conn, err := ssh.Dial("tcp", c.addr, c.config)
	if err != nil {
		return "", err
	}
	defer conn.Close()
	sess, err := conn.NewSession()
	if err != nil {
		return "", err
	}
	defer sess.Close()
	var out bytes.Buffer
	sess.Stdout = &out
	sess.Stderr = &out
	cmd := "bash -lc " + shellQuote(script)
	if err := sess.Run(cmd); err != nil {
		return out.String(), err
	}
	return out.String(), nil
}

func (c *Client) Ping() error {
	conn, err := ssh.Dial("tcp", c.addr, c.config)
	if err != nil {
		return err
	}
	_ = conn.Close()
	return nil
}

func CanDial(host string, port int, timeout time.Duration) error {
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", host, port), timeout)
	if err != nil {
		return err
	}
	return conn.Close()
}

func shellQuote(s string) string {
	return "'" + strings.ReplaceAll(s, "'", "'\\''") + "'"
}
