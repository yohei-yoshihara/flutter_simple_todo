package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"testing"
)

const BASE_URL = "http://localhost:8000"

func GetItems[T any](t *testing.T, client *http.Client, path string) ([]T, error) {
	resp, err := client.Get(BASE_URL + path)
	if err != nil {
		t.Error(err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
		return nil, err
	}

	items := []T{}
	err = json.Unmarshal(body, &items)
	if err != nil {
		log.Println(err)
		t.Error(err)
		return nil, err
	}
	return items, nil
}

func GetItem[T Task | Folder | User](t *testing.T, client *http.Client, path string) (T, error) {
	resp, err := client.Get(BASE_URL + path)
	if err != nil {
		t.Error(err)
		var item T
		return item, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var item T
		return item, fmt.Errorf(resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
		var item T
		return item, err
	}

	var item T
	err = json.Unmarshal(body, &item)
	if err != nil {
		log.Println(err)
		t.Error(err)
		var item T
		return item, err
	}
	return item, nil
}

func PostItem[T Task | Folder | User](t *testing.T, client *http.Client, path string, item T) (T, error) {
	itemJson, err := json.Marshal(item)
	if err != nil {
		t.Error(err)
		var item T
		return item, err
	}
	resp, err := client.Post(BASE_URL+path, "application/json; charset=UTF-8", bytes.NewBuffer(itemJson))
	if err != nil {
		t.Error(err)
		var item T
		return item, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var item T
		return item, fmt.Errorf(resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Error(err)
		var item T
		return item, err
	}

	var retItem T
	err = json.Unmarshal(body, &retItem)
	if err != nil {
		t.Error(err)
		var item T
		return item, err
	}

	return retItem, nil
}

func login(t *testing.T, user *User) (*http.Client, error) {
	if user == nil {
		user = &User{
			Username: "user1",
			Password: "password",
		}
	}

	jar, err := cookiejar.New(&cookiejar.Options{})
	if err != nil {
		t.Error(err)
		return nil, err
	}

	client := &http.Client{
		Jar: jar,
	}

	userJson, err := json.Marshal(*user)
	if err != nil {
		t.Error(err)
		return nil, err
	}

	_, err = client.Post(BASE_URL+"/api/login", "application/json; charset=UTF-8", bytes.NewBuffer(userJson))
	if err != nil {
		t.Error(err)
		return nil, err
	}

	return client, nil
}
