package ldap

import (
	"crypto/tls"
	"github.com/kolesa-team/go-socks-ldap/validator"
	"gopkg.in/ldap.v3"
)

type (
	// опции
	Options struct {
		Host         string
		Username     string
		Password     string
		BaseDn       string
		Scope        int
		DerefAliases int
		TimeLimit    int
		Filter       string
		TlsEnabled   bool
	}
	// структура клиента
	Client struct {
		options    Options
		authClient validator.Client
		ldapClient *ldap.Conn
	}
)

func NewClient(options Options, auth validator.Client) *Client {
	return &Client{
		options:    options,
		authClient: auth,
	}
}

// инициализация
func (c *Client) Init() error {
	var err error

	if c.options.TlsEnabled {
		c.ldapClient, err = ldap.DialTLS("tcp", c.options.Host, &tls.Config{InsecureSkipVerify: true})
	} else {
		c.ldapClient, err = ldap.Dial("tcp", c.options.Host)
	}

	if err != nil {
		return err
	}

	if err = c.ldapClient.Bind(c.options.Username, c.options.Password); err != nil {
		c.ldapClient.Close()
		return err
	}

	return nil
}

// клозер для graceful отключения
func (c *Client) Close() error {
	if c.ldapClient != nil {
		c.ldapClient.Close()
	}

	return nil
}

// получаем все записи с ldap
func (c *Client) FetchEntries() ([]*Entry, error) {
	var (
		result *ldap.SearchResult
		err    error
	)

	authField := c.authClient.GetFieldName()
	request := ldap.NewSearchRequest(
		c.options.BaseDn,
		c.options.Scope,
		c.options.DerefAliases,
		0,
		c.options.TimeLimit,
		false,
		c.options.Filter,
		[]string{"uid", authField},
		nil,
	)

	if result, err = c.ldapClient.Search(request); err != nil {
		return nil, err
	}

	entries := make([]*Entry, 0, len(result.Entries))

	for _, entry := range result.Entries {
		var (
			uid, password string
		)

		if uid = entry.GetAttributeValue("uid"); uid == "" {
			continue
		}

		if password = entry.GetAttributeValue(authField); password == "" {
			continue
		}

		entries = append(entries, NewEntry(c.authClient, uid, password))
	}

	return entries, nil
}
