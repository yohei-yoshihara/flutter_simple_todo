package cmd

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFolderIsProtected(t *testing.T) {
	client := &http.Client{}
	_, err := GetItems[Folder](t, client, "/api/folders")
	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "401 Unauthorized")

	_, err = GetItem[Folder](t, client, "/api/folders/1")
	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "401 Unauthorized")

	_, err = PostItem[Folder](t, client, "/api/folders/create", Folder{Name: "test"})
	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "401 Unauthorized")

	_, err = PostItem[Folder](t, client, "/api/folders/delete", Folder{ID: 1})
	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "401 Unauthorized")

	_, err = PostItem[Folder](t, client, "/api/folders/update", Folder{ID: 1})
	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "401 Unauthorized")
}

func TestTaskIsProtected(t *testing.T) {
	client := &http.Client{}
	_, err := GetItems[Task](t, client, "/api/tasks")
	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "401 Unauthorized")

	_, err = GetItem[Task](t, client, "/api/tasks/1")
	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "401 Unauthorized")

	_, err = PostItem[Task](t, client, "/api/tasks/create", Task{Title: "test", FolderId: 1})
	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "401 Unauthorized")

	_, err = PostItem[Task](t, client, "/api/tasks/delete", Task{ID: 1})
	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "401 Unauthorized")

	_, err = PostItem[Task](t, client, "/api/tasks/update", Task{ID: 1, Title: "test", FolderId: 1})
	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "401 Unauthorized")
}
