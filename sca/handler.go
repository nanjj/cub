package sca

import (
	"context"
)

type Handler func(context.Context, []DataObject) ([]DataObject, error)
