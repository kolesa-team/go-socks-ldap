package validator

import (
	"fmt"
	"github.com/GehirnInc/crypt"
	_ "github.com/GehirnInc/crypt/md5_crypt"
	"github.com/sirupsen/logrus"
	"regexp"
)

var re = regexp.MustCompile(`^.*?(\$.*)\$`)

// cryptMD5 validator struct
type CryptMd5Validator struct {
	crypt crypt.Crypter
}

func NewCryptMd5Validator() Client {
	return &CryptMd5Validator{
		crypt: crypt.MD5.New(),
	}
}

// get the field related with the current validator
func (v *CryptMd5Validator) GetFieldName() string {
	return "userPassword"
}

// password validation
func (v *CryptMd5Validator) Validate(hash string, password string) bool {
	return v.validatePassword(hash, password)
}

// password validation
func (v *CryptMd5Validator) validatePassword(hash string, password string) bool {
	r := re.FindStringSubmatch(hash)

	if len(r) != 2 {
		logrus.WithFields(logrus.Fields{
			"hash":      hash,
			"validator": "CryptMD5",
			"r":         r,
		}).Error("invalid hash")

		return false
	}

	userHash, err := v.crypt.Generate([]byte(password), []byte(r[1]))

	if err != nil {
		logrus.WithFields(logrus.Fields{
			"hash":      hash,
			"validator": "CryptMD5",
			"r":         r,
			"password":  password,
		}).Error("generate user hash failed")

		return false
	}

	return fmt.Sprintf("{CRYPT}%s", userHash) == hash
}
