package cmd

import (
	"fmt"
	"math/rand"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegisterAndLogin(t *testing.T) {
	client := &http.Client{}

	n := rand.Intn(1000000)
	username := fmt.Sprintf("user%d", n)
	password := fmt.Sprintf("password%d", n)

	user, err := PostItem[User](t, client, "/api/register", User{Username: username, Password: password})
	if err != nil {
		t.Error(err)
		return
	}
	assert.Equal(t, user.Username, username)

	_, err = login(t, &user)
	if err != nil {
		t.Error(err)
		return
	}
}
