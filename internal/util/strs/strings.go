package strs

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

type ResourceIdentifier struct {
	Type int32  `json:"type"`
	Obj  string `json:"obj"`
	Act  string `json:"act"`
}

func ValidateResourceIdentifier(typ int32, identifier string) (*ResourceIdentifier, error) {
	switch typ {
	case 1: // menu新增情况，typ=1
		matched, err := regexp.MatchString(`^/([A-Za-z0-9_-]+)(/[A-Za-z0-9_-]+)*$`, identifier)
		if err != nil {
			return nil, err
		}
		if !matched {
			return nil, errors.New("menu identifier must be `^/([A-Za-z0-9_-]+)(/[A-Za-z0-9_-]+)*$`")
		}
		splits := strings.Split(identifier, " ")
		return &ResourceIdentifier{
			Type: typ,
			Obj:  splits[1],
			Act:  splits[0],
		}, nil
	case 2: // widget新增情况，typ=2
		matched, err := regexp.MatchString(`^[A-Za-z0-9]+(_[A-Za-z0-9]+)*$`, identifier)
		if err != nil {
			return nil, err
		}
		if !matched {
			return nil, errors.New("widget identifier must be `^[A-Za-z0-9]+(_[A-Za-z0-9]+)*$`")
		}
		splits := strings.Split(identifier, " ")
		return &ResourceIdentifier{
			Type: typ,
			Obj:  splits[1],
			Act:  splits[0],
		}, nil
	case 3: // api新增情况，typ=3
		matched, err := regexp.MatchString(`^(GET|POST|PUT|DELETE|PATCH|OPTIONS|HEAD) /([A-Za-z0-9_-]+)(/[A-Za-z0-9_-]+)*$`, identifier)
		if err != nil {
			return nil, err
		}
		if !matched {
			return nil, errors.New("API identifier must be `^(GET|POST|PUT|DELETE|PATCH|OPTIONS|HEAD) /([A-Za-z0-9_-]+)(/[A-Za-z0-9_-]+)*$`")
		}
		splits := strings.Split(identifier, " ")
		return &ResourceIdentifier{
			Type: typ,
			Obj:  splits[1],
			Act:  splits[0],
		}, nil
	default: // 资源授权情况，typ=-1
		matched, err := regexp.MatchString(`^(read|write) /([A-Za-z0-9_-]+)(/[A-Za-z0-9_-]+)*$`, identifier)
		if err != nil {
			return nil, err
		}
		if matched {
			splits := strings.Split(identifier, " ")
			return &ResourceIdentifier{
				Type: 1,
				Obj:  splits[1],
				Act:  splits[0],
			}, nil
		}
		matched, err = regexp.MatchString(`^(read|write) [A-Za-z0-9]+(_[A-Za-z0-9]+)*$`, identifier)
		if err != nil {
			return nil, err
		}
		if matched {
			splits := strings.Split(identifier, " ")
			return &ResourceIdentifier{
				Type: 2,
				Obj:  splits[1],
				Act:  splits[0],
			}, nil
		}
		matched, err = regexp.MatchString(`^(GET|POST|PUT|DELETE|PATCH|OPTIONS|HEAD) /([A-Za-z0-9_-]+)(/:?[A-Za-z0-9_-]+)*$`, identifier)
		if err != nil {
			return nil, err
		}
		if matched {
			splits := strings.Split(identifier, " ")
			return &ResourceIdentifier{
				Type: 3,
				Obj:  splits[1],
				Act:  splits[0],
			}, nil
		}
		return nil, fmt.Errorf("unsupported identifier format: %s", identifier)
	}
}
