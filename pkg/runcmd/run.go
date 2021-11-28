package runcmd

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/davecgh/go-spew/spew"
	"github.com/peterbourgon/ff/v3/ffcli"
	"gopkg.in/yaml.v2"
)

var (
	schemaFilename = "schema.graphql"
	modelFilename  = "models.go"
	genFilename    = "generated.go"
	cfgFiles       = []string{"quokka.yml", "quokka.yaml"}
)

type StringList []string

type PackageConfig struct {
	Filename string `yaml:"filename,omitempty"`
	Package  string `yaml:"package,omitempty"`
}

type DirectiveConfig struct {
	SkipRuntime bool `yaml:"skip_runtime"`
}

type TypeMapField struct {
	Resolver  bool   `yaml:"resolver"`
	FieldName string `yaml:"fieldName"`
	Method    string `yaml:"method"`
}

type TypeMapEntry struct {
	Model  StringList              `yaml:"model"`
	Fields map[string]TypeMapField `yaml:"fields,omitempty"`
}

type TypeMap map[string]TypeMapEntry

type Config struct {
	SchemaFilename StringList                 `yaml:"schema,omitempty"`
	Model          PackageConfig              `yaml:"model,omitempty"`
	Models         TypeMap                    `yaml:"models,omitempty"`
	Exec           PackageConfig              `yaml:"exec"`
	Directives     map[string]DirectiveConfig `yaml:"directives,omitempty"`
}

// TODO: --dryrun command
func New() *ffcli.Command {
	flagSet := flag.NewFlagSet("quokka run", flag.ExitOnError)
	return &ffcli.Command{
		Name:       "run",
		ShortUsage: "quokka run",
		ShortHelp:  "Create migrations, migrate database and generate graphql server",
		FlagSet:    flagSet,
		Exec: func(ctx context.Context, args []string) error {
			if len(args) > 0 {
				return errors.New("command does not accept args")
			}

			cfg, err := getConfig()

			fmt.Println(cfg, err, "<=== CONTINUE IN MAIN")

			return nil
		},
	}
}

func getConfig() (*Config, error) {
	file, err := findCfg()
	if err != nil {
		return nil, err
	}

	err = os.Chdir(filepath.Dir(file))
	if err != nil {
		fmtError := fmt.Errorf("unable to change working directory: %w", err)
		return nil, fmtError
	}

	config := &Config{
		SchemaFilename: StringList{schemaFilename},
		Model:          PackageConfig{Filename: modelFilename},
		Models:         TypeMap{},
		Exec:           PackageConfig{Filename: genFilename},
		Directives:     map[string]DirectiveConfig{},
	}

	b, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, fmt.Errorf("unable to read config: %w", err)
	}

	if err := yaml.UnmarshalStrict(b, config); err != nil {
		return nil, fmt.Errorf("unable to parse config: %w", err)
	}

	if err := populateConfig(config); err != nil {
		return nil, err
	}

	return config, nil
}

func populateConfig(config *Config) error {
	defaultDirectives := map[string]DirectiveConfig{
		// skip and include are required by GQL language spec.
		"skip":       {SkipRuntime: true},
		"include":    {SkipRuntime: true},
		"deprecated": {SkipRuntime: true},
	}

	for key, value := range defaultDirectives {
		if _, defined := config.Directives[key]; !defined {
			config.Directives[key] = value
		}
	}

	spew.Dump(defaultDirectives, "XXXXXXX")
	return nil
}

func findCfg() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("os.Getwd() failed to get current directory")
	}

	cfg, err := getFilePath(dir)
	if err != nil {
		return "", err
	}
	return cfg, nil
}

func getFilePath(dir string) (string, error) {
	for _, file := range cfgFiles {
		path := filepath.Join(dir, file)
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}
	}
	return "", os.ErrNotExist // TODO HERE LETS SAY SPECIFIC FILE NOT FOUND
}
