package storage

import (
	"github.com/kolesa-team/go-socks-ldap/ldap"
	"github.com/kolesa-team/go-socks-ldap/validator"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"sync"
	"time"
)

type (
	// local entries storage
	Entries map[string]*ldap.Entry
	// storage client options
	Options struct {
		UpdateInterval time.Duration
		ConnectRetry   int
		LdapOptions    ldap.Options
	}
	// client struct
	Client struct {
		options      Options
		entries      Entries
		closeChannel chan bool
		notifyClose  chan bool
		authClient   validator.Client
		mutex        sync.RWMutex
	}
)

func NewClient(options Options, authClient validator.Client) *Client {
	options.UpdateInterval *= time.Second
	return &Client{
		options:    options,
		authClient: authClient,
		entries:    make(Entries),
	}
}

// get close listener
func (c *Client) ListenClose() <-chan bool {
	return c.closeChannel
}

// graceful closer
func (c *Client) Close() error {
	close(c.closeChannel)
	return nil
}

// initialization the client
func (c *Client) Init() error {
	// update entries
	go func() {
		var cnt = 0

	UpdateCycle:
		for {
			t := time.NewTimer(c.options.UpdateInterval)
			err := c.updateUsers()

			if err != nil {
				logrus.WithFields(logrus.Fields{
					"host":    c.options.LdapOptions.Host,
					"message": err.Error(),
				}).Error("update entries failed")
			}

			if err != nil && cnt > c.options.ConnectRetry {
				logrus.WithFields(logrus.Fields{
					"retry_count": cnt,
					"limit":       c.options.ConnectRetry,
				}).Error(err)

				c.closeChannel <- true
			} else if err != nil {
				cnt++
			}

			select {
			case <-c.notifyClose:
				t.Stop()
				break UpdateCycle
			case <-t.C:
			}
		}
	}()

	return nil
}

// update users from ldap
func (c *Client) updateUsers() error {
	var (
		err     error
		entries []*ldap.Entry
	)

	logrus.WithFields(logrus.Fields{
		"server": c.options.LdapOptions.Host,
	}).Info("fetch all entries form ldap")

	conn := ldap.NewClient(c.options.LdapOptions, c.authClient)
	//noinspection GoUnhandledErrorResult
	defer conn.Close()

	if err := conn.Init(); err != nil {
		return err
	}

	if entries, err = conn.FetchEntries(); err != nil {
		logrus.WithError(err).WithFields(logrus.Fields{
			"server": c.options.LdapOptions.Host,
		}).Error("update failed with error")

		return err
	}

	if len(entries) == 0 {
		return nil
	}

	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.entries = make(Entries)

	for _, entry := range entries {
		if entry == nil {
			continue
		}

		c.entries[entry.GetUID()] = entry
	}

	logrus.WithFields(logrus.Fields{
		"total": len(c.entries),
	}).Info("entries updated")

	return nil
}

// get user entry by username
func (c *Client) GetEntry(username string) (*ldap.Entry, error) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	if user, ok := c.entries[username]; ok {
		return user, nil
	}

	return nil, errors.New("entry not found")
}
