package socks

import (
	"fmt"
	"github.com/armon/go-socks5"
	"github.com/kolesa-team/go-socks-ldap/storage"
	"github.com/sirupsen/logrus"
	"io"
)

const (
	// бит с версией клиента
	clientAuthVersion uint8 = 1

	// биты с кодами ответа
	clientAuthSuccess uint8 = 0
	clientAuthFailure uint8 = 1
)

// структура аутентификатора
type LdapAuthenticator struct {
	storageClient *storage.Client
}

func NewLdapAuthenticator(storage *storage.Client) *LdapAuthenticator {
	return &LdapAuthenticator{storage}
}

// Возращаем socks клиенту ответ, что авторизация идет через пару логин:пасс
func (a LdapAuthenticator) GetCode() uint8 {
	return socks5.UserPassAuth
}

// Функция авторизации socks клиента
func (a LdapAuthenticator) Authenticate(reader io.Reader, writer io.Writer) (*socks5.AuthContext, error) {
	if _, err := writer.Write([]byte{uint8(5), socks5.UserPassAuth}); err != nil {
		return nil, err
	}

	header := []byte{0, 0}
	if _, err := io.ReadAtLeast(reader, header, 2); err != nil {
		return nil, err
	}

	if header[0] != clientAuthVersion {
		return nil, fmt.Errorf("unsupported auth version: %v", header[0])
	}

	userLen := int(header[1])
	user := make([]byte, userLen)
	if _, err := io.ReadAtLeast(reader, user, userLen); err != nil {
		return nil, err
	}

	if _, err := reader.Read(header[:1]); err != nil {
		return nil, err
	}

	passLen := int(header[0])
	pass := make([]byte, passLen)

	if _, err := io.ReadAtLeast(reader, pass, passLen); err != nil {
		return nil, err
	}

	logrus.WithFields(logrus.Fields{
		"user": string(user),
	}).Info("user connected")

	if entry, err := a.storageClient.GetEntry(string(user)); err == nil && entry.Validate(string(pass)) {
		if _, err := writer.Write([]byte{clientAuthVersion, clientAuthSuccess}); err != nil {
			logrus.WithFields(logrus.Fields{
				"user": string(user),
			}).Info(err)

			return nil, err
		}

		logrus.WithFields(logrus.Fields{
			"user": string(user),
		}).Info("user credentials is valid")

		return &socks5.AuthContext{
			Method: socks5.UserPassAuth,
			Payload: map[string]string{
				"Username": string(user),
			},
		}, nil
	}

	if _, err := writer.Write([]byte{clientAuthVersion, clientAuthFailure}); err != nil {
		logrus.WithFields(logrus.Fields{
			"user": string(user),
		}).Info(err)

		return nil, err
	}

	logrus.WithFields(logrus.Fields{
		"user": string(user),
	}).Info("credential is invalid")

	return nil, socks5.UserAuthFailed
}

// клозер для закрытия корректного
func (a *LdapAuthenticator) Close() error {
	return a.storageClient.Close()
}
