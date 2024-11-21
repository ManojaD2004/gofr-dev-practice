package route

import (
	"gofr.dev/pkg/gofr"
	t "github.com/ManojaD2004/types"
)

func UserGetRoute(ctx *gofr.Context) (interface{}, error) {
	req := t.UserRequestType{}
	ctx.Bind(&req)
	res := t.UserResponseType{}
	// Your logic

	return res, nil
}
