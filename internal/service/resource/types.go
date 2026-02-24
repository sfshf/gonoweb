package resource

import (
	. "github.com/sfshf/gonoweb/internal/model"
)

type TypedResource struct {
	Menus   []TResource `json:"menus"`
	Widgets []TResource `json:"widgets"`
	APIs    []TResource `json:"apis"`
}
