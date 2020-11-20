package toggl

import (
	"encoding/json"
	"fmt"
)

type WorkspaceList struct {
	Count      int
	Workspaces []Workspace
}


type ProjectList struct {
	Count    int
	Projects []Project
}

type WorkspaceClient struct {
	client *Client
}

func (w *WorkspaceClient) GetWorkspaces() (*WorkspaceList, error) {
	message, err := w.client.NewMessage("GET", "workspaces", nil)
	if err != nil {
		return nil, err
	}

	data, err := w.client.SendRequest(message)
	if err != nil {
		return nil, err
	}

	var workspaces []Workspace
	if err := json.Unmarshal(*data, &workspaces); err != nil {
		return nil, err
	}

	return &WorkspaceList{
		Workspaces: workspaces,
		Count:      len(workspaces),
	}, nil
}

func (w *WorkspaceClient) ListProjects(wid int) (*ProjectList, error) {
	message, err := w.client.NewMessage("GET", fmt.Sprintf("workspaces/%d/projects", wid), nil)
	if err != nil {
		return nil, err
	}

	data, err := w.client.SendRequest(message)
	if err != nil {
		return nil, err
	}

	var projects []Project
	err = json.Unmarshal(*data, &projects)

	return &ProjectList{
		Projects: projects,
		Count:    len(projects),
	}, nil

}
