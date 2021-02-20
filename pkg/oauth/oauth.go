package oauth

import (
	"context"
	"crypto/sha256"
	"fmt"
)

const (
	tokenOne = "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJpc3MiOiJPbmxpbmUgSldUIEJ1aWxkZXIiLCJpYXQiOjE2MTM4NDMwNjksImV4cCI6MTY0NTM3OTA2OSwiYXVkIjoid3d3LmV4YW1wbGUuY29tIiwic3ViIjoib25lQHVzZXIuY29tIn0.LCRkH09YVi2BrnjXKPalLP2aNwn3lGDUiuhi5sx4tGY"
	tokenTwo = "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJpc3MiOiJPbmxpbmUgSldUIEJ1aWxkZXIiLCJpYXQiOjE2MTM4NDMwNjksImV4cCI6MTY0NTM3OTA2OSwiYXVkIjoid3d3LmV4YW1wbGUuY29tIiwic3ViIjoidHdvQHVzZXIuY29tIn0._0TE7qepR1o9R_COobgwNnfeBGeOLXGM6qPQRj2R-t8"
)

type Insecure struct {
	token string
	sub   string
}

func NewInsecureOne() *Insecure {
	return &Insecure{
		sub:   mustHash("one@user.com"),
		token: tokenOne,
	}
}

func NewInsecureTwo() *Insecure {
	return &Insecure{
		sub:   mustHash("two@user.com"),
		token: tokenTwo,
	}
}

func (i *Insecure) GetRequestMetadata(ctx context.Context, in ...string) (map[string]string, error) {
	m := map[string]string{
		"authorization": "bearer " + i.token,
	}

	return m, nil
}

func (i *Insecure) RequireTransportSecurity() bool {
	return false
}

func (i *Insecure) User() string {
	return i.sub
}

func mustHash(s string) string {
	h := sha256.New()
	_, err := h.Write([]byte(s))
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("%x", h.Sum(nil))
}
