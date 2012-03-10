package config

import (
	"fmt"
	"github.com/afajl/ctrl/path"
	"os"
	"path/filepath"
)

const defaultConfigfile = "ctrl.conf"
const defaultLogdir = "ctrlruns"

func FromFile(c *Config, path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()
	if err = FromJson(c, file); err != nil {
		return fmt.Errorf("error parsing %s, %s", path, err)
	}

	c.Rootdir = filepath.Dir(path)

	// set logdir to dirname(path)/defaultLogdir if not set
	if c.Logdir == "" {
		c.Logdir = filepath.Join(c.Rootdir, defaultLogdir)
	}

	return nil
}

func FromFileDefault(c *Config) error {
	path, err := path.FindUpwards(defaultConfigfile)
	if err != nil {
		return err
	}
	if path == "" {
		return fmt.Errorf("could not find ctrl.conf in any parent")
	}
	return FromFile(c, path)
}
