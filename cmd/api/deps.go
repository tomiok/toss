package api

import uploads "github.com/tomiok/toss/internal/files/handler"

type Deps struct {
	UploadHandler uploads.Handler
}

func NewDeps() Deps {
	handler := uploads.New()

	return Deps{
		UploadHandler: handler,
	}
}
