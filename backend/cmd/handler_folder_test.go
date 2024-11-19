package cmd

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFolders(t *testing.T) {
	client, err := login(t, nil)
	if err != nil {
		t.Error(err)
		return
	}

	folders, err := GetItems[Folder](t, client, "/api/folders")
	if err != nil {
		t.Error(err)
		return
	}
	assert.True(t, len(folders) > 0, "")
}

func TestFolder(t *testing.T) {
	client, err := login(t, nil)
	if err != nil {
		t.Error(err)
		return
	}

	folder, err := GetItem[Folder](t, client, "/api/folders/1")
	if err != nil {
		t.Error(err)
		return
	}
	assert.True(t, folder.ID > 0)
	assert.True(t, len(folder.Name) > 0)
}

func TestCreateAndDeleteFolder(t *testing.T) {
	client, err := login(t, nil)
	if err != nil {
		t.Error(err)
		return
	}

	var folderId int64
	{ // create
		folder := Folder{
			Name: "テスト",
		}

		retFolder, err := PostItem[Folder](t, client, "/api/folders/create", folder)
		if err != nil {
			t.Error(err)
			return
		}

		assert.Equal(t, retFolder.Name, "テスト")
		folderId = retFolder.ID
	}

	{ // delete
		folder := Folder{
			ID: folderId,
		}

		retFolder, err := PostItem[Folder](t, client, "/api/folders/delete", folder)
		if err != nil {
			t.Error(err)
			return
		}

		assert.Equal(t, retFolder.ID, folderId)
	}
}

func TestUpdateFolder(t *testing.T) {
	client, err := login(t, nil)
	if err != nil {
		t.Error(err)
		return
	}

	items, err := GetItems[Folder](t, client, "/api/folders")
	if err != nil {
		t.Error(err)
		return
	}

	n := rand.Intn(100000)
	folderName := fmt.Sprintf("folder-%d", n)
	folder := Folder{
		ID:   items[0].ID,
		Name: folderName,
	}
	retFolder, err := PostItem[Folder](t, client, "/api/folders/update", folder)

	assert.Equal(t, retFolder.Name, folderName)

	folder = Folder{
		ID:   items[0].ID,
		Name: "プライベート",
	}
	_, err = PostItem[Folder](t, client, "/api/folders/update", folder)
}
