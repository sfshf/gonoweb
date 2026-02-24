package strs

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

func ValidateResourceIdentifier(typ int32, identifier string) (string, string, error) {
	switch typ {
	case 1: // menu
		matched, err := regexp.MatchString(`^/([A-Za-z0-9_-]+)(/[A-Za-z0-9_-]+)*$`, identifier)
		if err != nil {
			return "", "", err
		}
		if !matched {
			return "", "", errors.New("menu identifier must be `^/([A-Za-z0-9_-]+)(/[A-Za-z0-9_-]+)*$`")
		}
		return identifier, "", nil
	case 2: // widget
		matched, err := regexp.MatchString(`^[A-Za-z0-9]+(_[A-Za-z0-9]+)*$`, identifier)
		if err != nil {
			return "", "", err
		}
		if !matched {
			return "", "", errors.New("widget identifier must be `^[A-Za-z0-9]+(_[A-Za-z0-9]+)*$`")
		}
		return identifier, "", nil
	case 3: // api
		matched, err := regexp.MatchString(`^(GET|POST|PUT|DELETE|PATCH|OPTIONS|HEAD) /([A-Za-z0-9_-]+)(/[A-Za-z0-9_-]+)*$`, identifier)
		if err != nil {
			return "", "", err
		}
		if !matched {
			return "", "", errors.New("API identifier must be `^(GET|POST|PUT|DELETE|PATCH|OPTIONS|HEAD) /([A-Za-z0-9_-]+)(/[A-Za-z0-9_-]+)*$`")
		}
		splits := strings.Split(identifier, " ")
		return splits[1], splits[0], nil
	default:
		matched, err := regexp.MatchString(`^/([A-Za-z0-9_-]+)(/[A-Za-z0-9_-]+)*$`, identifier)
		if err != nil {
			return "", "", err
		}
		if matched {
			return identifier, "", nil
		}
		matched, err = regexp.MatchString(`^[A-Za-z0-9]+(_[A-Za-z0-9]+)*$`, identifier)
		if err != nil {
			return "", "", err
		}
		if matched {
			return identifier, "", nil
		}
		matched, err = regexp.MatchString(`^(GET|POST|PUT|DELETE|PATCH|OPTIONS|HEAD) /([A-Za-z0-9_-]+)(/:?[A-Za-z0-9_-]+)*$`, identifier)
		if err != nil {
			return "", "", err
		}
		if matched {
			splits := strings.Split(identifier, " ")
			return splits[1], splits[0], nil
		}
		return "", "", fmt.Errorf("unsupported identifier format: %s", identifier)
	}
}
