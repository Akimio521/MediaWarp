package constants

import (
	"encoding/json"
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

type FliterMode uint8

const (
	WHITELIST FliterMode = iota
	BLACKLIST
)

func (f FliterMode) String() string {
	switch f {
	case WHITELIST:
		return "WhiteList"
	case BLACKLIST:
		return "BlackList"
	default:
		return "Unknown"
	}
}

func (f FliterMode) MarshalJSON() ([]byte, error) {
	return json.Marshal(f.String())
}

func (f *FliterMode) UnmarshalJSON(data []byte) error {
	var s string
	err := json.Unmarshal(data, &s)
	if err != nil {
		return err
	}
	switch strings.ToLower(s) {
	case "whitelist":
		*f = WHITELIST
	case "blacklist":
		*f = BLACKLIST
	default:
		return fmt.Errorf("unknown FliterMode: %s", s)
	}
	return nil
}

func (f FliterMode) MarshalYAML() (any, error) {
	return f.String(), nil
}

func (f *FliterMode) UnmarshalYAML(value *yaml.Node) error {
	var s string
	if err := value.Decode(&s); err != nil {
		return err
	}
	switch strings.ToLower(s) {
	case "whitelist":
		*f = WHITELIST
	case "blacklist":
		*f = BLACKLIST
	default:
		return fmt.Errorf("unknown FliterMode: %s", s)
	}
	return nil
}

func (f FliterMode) MarshalText() ([]byte, error) {
	return []byte(f.String()), nil
}

func (f *FliterMode) UnmarshalText(text []byte) error {
	s := string(text)
	switch strings.ToLower(s) {
	case "whitelist":
		*f = WHITELIST
	case "blacklist":
		*f = BLACKLIST
	default:
		return fmt.Errorf("unknown FliterMode: %s", s)
	}
	return nil
}
