package cmd

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTasks(t *testing.T) {
	client, err := login(t, nil)
	if err != nil {
		t.Error(err)
		return
	}

	tasks, err := GetItems[Task](t, client, "/api/tasks")
	if err != nil {
		t.Error(err)
		return
	}
	assert.True(t, len(tasks) > 0, "")
}

func TestTask(t *testing.T) {
	client, err := login(t, nil)
	if err != nil {
		t.Error(err)
		return
	}

	task, err := GetItem[Task](t, client, "/api/tasks/1")
	if err != nil {
		t.Error(err)
		return
	}
	assert.True(t, task.ID > 0)
	assert.True(t, len(task.Title) > 0)
	assert.True(t, task.FolderId > 0)
}

func TestCreateAndDeleteTask(t *testing.T) {
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
	assert.True(t, len(folders) > 0)

	var taskId int64
	{ // create
		task := Task{
			Title:    "テスト",
			FolderId: folders[0].ID,
		}

		returnedTask, err := PostItem[Task](t, client, "/api/tasks/create", task)
		if err != nil {
			t.Error(err)
			return
		}

		assert.Equal(t, returnedTask.Title, "テスト")
		taskId = returnedTask.ID
	}

	{ // delete
		task := Task{
			ID: taskId,
		}

		returnedTask, err := PostItem[Task](t, client, "/api/tasks/delete", task)
		if err != nil {
			t.Error(err)
			return
		}

		assert.Equal(t, returnedTask.ID, taskId)
	}
}

func TestUpdateTask(t *testing.T) {
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

	tasks, err := GetItems[Task](t, client, "/api/tasks")
	if err != nil {
		t.Error(err)
		return
	}

	n := rand.Intn(100000)
	taskTitle := fmt.Sprintf("task-%d", n)
	task := Task{
		ID:       tasks[0].ID,
		Title:    taskTitle,
		FolderId: folders[0].ID,
	}
	returnedTask, err := PostItem[Task](t, client, "/api/tasks/update", task)
	if err != nil {
		t.Error(err)
		return
	}

	assert.Equal(t, returnedTask.Title, taskTitle)

	_, err = PostItem[Task](t, client, "/api/tasks/update", tasks[0])
	if err != nil {
		t.Error(err)
		return
	}
}
