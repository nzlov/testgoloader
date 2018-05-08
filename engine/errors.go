package engine

import (
	"fmt"
)

type ErrInitParams struct {
	Params string
}

func NewErrInitParams(params string) *ErrInitParams {
	return &ErrInitParams{
		Params: params,
	}
}

func (e *ErrInitParams) Error() string {
	return fmt.Sprintf("Plugin Init Params Error:%+v", e.Params)

}
