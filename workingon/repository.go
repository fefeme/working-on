package workingon

import (
	"github.com/tcnksm/go-gitconfig"
)

func FindProjectByGitRepositoryUrl(cfg *Config) int {
	// Check, if this is a git repository
	url, _ := gitconfig.OriginURL()
	if url == "" {
		return 0
	}

	for _, project := range cfg.Projects {
		if project.Git == url {
			return project.TogglePid
		}
	}

	return 0
}
