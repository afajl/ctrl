package config

import (
	"bytes"
	"encoding/json"
	"io"
)

func FromJson(c *Config, r io.Reader) error {
	buf := new(bytes.Buffer)
	if _, err := buf.ReadFrom(r); err != nil {
		return err
	}
	return FromJsonBytes(c, buf.Bytes())
}

func FromJsonBytes(c *Config, b []byte) error {
	if err := json.Unmarshal(b, &c); err != nil {
		return err
	}
	return nil
}
