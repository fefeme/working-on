package workingon

import (
	"fmt"
	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
	"reflect"
	"strings"
	"time"
)

type Config struct {
	CreatedWith string                 `yaml:"-"`
	Settings    Settings               `yaml:"settings" mapstructure:"settings"`
	Projects    []ProjectMapping       `mapstructure:"mappings"`
	Templates   []TemplateConfig       `yaml:"templates" mapstructure:"templates"`
	Sources     map[string]interface{} `yaml:"sources" mapstructure:"sources"`
}

type Settings struct {
	Location          time.Location `yaml:"location" mapstructure:"location"`
	DayFirst          bool          `mapstructure:"day_first" yaml:"day_first"`
	DateLayout        string        `mapstructure:"date_layout" yaml:"date_layout"`
	DateTimeLayout    string        `mapstructure:"date_time_layout" yaml:"date_time_layout"`
	ToggleApiToken    string        `mapstructure:"toggl_api_token" yaml:"toggle_api_token"`
	ToggleWid         int           `mapstructure:"toggl_wid" yaml:"toggl_wid"`
	TogglePidRequired bool          `mapstructure:"toggl_pid_required" yaml:"toggl_pid_required"`
	DefaultTaskSource string	    `mapstructure:"default_task_source" yaml:"default_task_source"`
}

type TemplateConfig struct {
	Alias       string `mapstructure:"alias"`
	Description string `mapstructure:"description"`
	Start       string `mapstructure:"start"`
	Stop        string `mapstructure:"stop"`
	Project     int    `mapstructure:"project"`
	TogglTask   int    `mapstructure:"toggl_task"`
}

type ProjectMapping struct {
	Name      string `yaml:"name" mapstructure:"name"`
	TogglePid int    `yaml:"toggl_pid"  mapstructure:"toggl_pid"`
	Git       string `yaml:"git" mapstructure:"git"`
	Jira      string `yaml:"jira" mapstructure:"jira"`
}

var (
	Configuration Config
)

func (c *Config) GetTemplate(alias string) (*TemplateConfig, error) {
	templates := c.Templates
	for _, template := range templates {
		if strings.EqualFold(template.Alias, alias) {
			return &template, nil
		}
	}
	return nil, nil
}

func (c *Config) GetMapping(key string) (*ProjectMapping, error) {
	var projectMapping *ProjectMapping
	for _, n := range c.Projects {
		if strings.EqualFold(n.Name, key) {
			projectMapping = &n
			break
		}
	}
	if projectMapping == nil {
		return nil, fmt.Errorf("project mapping not found for key %s", key)
	}
	return projectMapping, nil
}

func newStringToLocationHookFunc() mapstructure.DecodeHookFunc {
	return func(from reflect.Type, to reflect.Type, data interface{}) (interface{}, error) {
		if from.Kind() == reflect.String && to.Name() == "Location" {
			return time.LoadLocation(data.(string))
		}
		return data, nil
	}
}

func InitConfig() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	viper.AddConfigPath(".")
	viper.AddConfigPath("$HOME/.config/working_on")

	viper.SetEnvPrefix("WO")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	Configuration.CreatedWith = "working_on"

	err := viper.Unmarshal(&Configuration,
		viper.DecodeHook(mapstructure.ComposeDecodeHookFunc(newStringToLocationHookFunc())))

	if err != nil {
		return nil, err
	}
	return &Configuration, nil
}
