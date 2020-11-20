package toggl

import (
	"encoding/json"
	"fmt"
)

type TaskClient struct {
	client *Client
}

func (t *TaskClient) List(wid int) (*TaskList, error) {
	message, err := t.client.NewMessage("GET", fmt.Sprintf("workspaces/%d/tasks", wid), nil)
	if err != nil {
		return nil, err
	}

	data, err := t.client.SendRequest(message)
	if err != nil {
		return nil, err
	}

	var tasks []Task
	err = json.Unmarshal(*data, &tasks)

	return &TaskList{
		Tasks: tasks,
		Count: len(tasks),
	}, nil
}

func (t *TaskClient) Get(tid int) (*Task, error) {
	message, err := t.client.NewMessage("GET", fmt.Sprintf("tasks/%d", tid), nil)

	if err != nil {
		return nil, err
	}

	data, err := t.client.SendRequest(message)
	if err != nil {
		return nil, err
	}

	var task Task
	err = json.Unmarshal(*data, &task)

	return &task, nil

}
