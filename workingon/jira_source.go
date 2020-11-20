package workingon

import (
	jira "github.com/andygrunwald/go-jira"
	"github.com/mitchellh/mapstructure"
)

type jiraConfig struct {
	Username string
	Password string
	URL      string `mapstructure:url`
}

// JiraSource A source for Jira Tasks
type JiraSource struct {
	client *jira.Client
}

func init() {
	err := Registry.Register(&JiraSource{})
	if err != nil {
		panic(err)
	}
}

func (j *JiraSource) Configure(cfg *Config) error {
	var config jiraConfig
	err := mapstructure.Decode(cfg.Sources["jira"], &config)
	if err != nil {
		return err
	}

	tp := jira.BasicAuthTransport{
		Username: config.Username,
		Password: config.Password,
	}

	j.client, err = jira.NewClient(tp.Client(), config.URL)
	if err != nil {
		return err
	}

	return nil
}

// GetName returns the name of the source
func (j *JiraSource) GetName() string {
	return "Jira"
}

// GetTasks returns all Jira issues assigned to the current jira user
func (j *JiraSource) GetTasks() ([]Task, error) {
	jql := "assignee=currentUser() and statusCategory != Done order by key"

	issues, resp, err := j.client.Issue.Search(jql, nil)

	if err != nil {
		return nil, err
	}

	var tasks []Task

	if resp.StatusCode == 200 {
		for _, issue := range issues {
			tasks = append(tasks, Task{
				Key:     issue.Key,
				Summary: issue.Fields.Summary,
				Project: Project{
					Key:  issue.Fields.Project.Key,
					Name: issue.Fields.Project.Name,
				},
			},
			)
		}
	}

	return tasks, nil
}

// GetProjects returns all public projects
func (j *JiraSource) GetProjects() ([]Project, error) {
	jiraProjects, resp, err := j.client.Project.GetList()

	if err != nil {
		return nil, err
	}

	var projects []Project

	if resp.StatusCode == 200 {
		for _, project := range *jiraProjects {
			projects = append(projects, Project{Name: project.Name})
		}
	}

	return projects, nil

}

func (j *JiraSource) GetTask(key string) (*Task, error) {
	issue, _, err := j.client.Issue.Get(key, nil)

	if err != nil {
		return nil, err
	}
	t := &Task{
		Key:     issue.Key,
		Summary: issue.Fields.Summary,
		Project: Project{
			Key:  issue.Fields.Project.Key,
			Name: issue.Fields.Project.Name,
		},
	}
	return t, nil

}
