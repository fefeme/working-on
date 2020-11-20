package workingon

type Project struct {
	Key  string
	Name string
}


type Task struct {
	Key       string
	Summary   string
	Project   Project
	TogglTask int
}

