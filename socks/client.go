package socks

import (
	"github.com/armon/go-socks5"
	"github.com/kolesa-team/go-socks-ldap/storage"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"log"
	"net"
	"syscall"
	"time"
)

type (
	// client options
	Options struct {
		Listen    string
		SoftLimit uint64
	}
	// client struct
	Client struct {
		options     Options
		authClient  *LdapAuthenticator
		socksClient *socks5.Server
		listener    net.Listener
	}
)

func NewClient(options Options, storage *storage.Client) *Client {
	return &Client{
		options:    options,
		authClient: NewLdapAuthenticator(storage),
	}
}

// client initialization
func (c *Client) Init() error {
	var err error

	c.socksClient, err = socks5.New(&socks5.Config{
		AuthMethods: []socks5.Authenticator{c.authClient},
		Logger:      log.New(ioutil.Discard, "", log.LstdFlags),
	})

	if err != nil {
		return err
	}

	if c.listener, err = net.Listen("tcp", c.options.Listen); err != nil {
		return err
	}

	if err := c.setLimits(); err != nil {
		logrus.WithError(err).Error("cannot set the soft limit")
	}

	return nil
}

// set system limits
func (c *Client) setLimits() error {
	var rLimit syscall.Rlimit

	if err := syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rLimit); err != nil {
		return err
	}

	if rLimit.Cur < c.options.SoftLimit {
		rLimit.Cur = c.options.SoftLimit

		if err := syscall.Setrlimit(syscall.RLIMIT_NOFILE, &rLimit); err != nil {
			return err
		}
	}

	if err := syscall.Getrlimit(syscall.RLIMIT_NOFILE, &rLimit); err != nil {
		return err
	}

	return nil
}

// listen and serve
func (c *Client) Serve() error {
	for {
		conn, err := c.listener.Accept()

		if err != nil {
			logrus.WithField("remote", conn.RemoteAddr().String()).
				WithError(err).
				Error("cannot accept connection")

			if conn != nil {
				//noinspection GoUnhandledErrorResult
				conn.Close()
			}

			// подождем чтобы все прошло
			time.Sleep(time.Second)
			continue
		}

		logrus.WithFields(logrus.Fields{
			"remote": conn.RemoteAddr().String(),
		}).Info("new connection")

		go c.serveConn(conn)
	}
}

// handle the connection
func (c *Client) serveConn(conn net.Conn) {
	if err := c.socksClient.ServeConn(conn); err != nil {
		logrus.WithField("remote", conn.RemoteAddr().String()).Error(err)
	}
}

// client closer
func (c *Client) Close() error {
	if err := c.authClient.Close(); err != nil {
		return err
	}

	if err := c.listener.Close(); err != nil {
		return err
	}

	return nil
}
