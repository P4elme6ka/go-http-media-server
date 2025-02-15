package vhostHandler

import (
	"github.com/P4elme6ka/go-http-media-server/src/param"
	"github.com/P4elme6ka/go-http-media-server/src/serverErrHandler"
	"github.com/P4elme6ka/go-http-media-server/src/serverHandler"
	"github.com/P4elme6ka/go-http-media-server/src/serverLog"
	"github.com/P4elme6ka/go-http-media-server/src/tpl"
	"github.com/P4elme6ka/go-http-media-server/src/user"
	"net/http"
)

type VhostHandler struct {
	p            *param.Param
	logger       *serverLog.Logger
	errorHandler *serverErrHandler.ErrHandler
	theme        tpl.Theme
	Handler      http.Handler
}

func NewHandler(
	p *param.Param,
	logger *serverLog.Logger,
	errorHandler *serverErrHandler.ErrHandler,
	theme tpl.Theme,
) *VhostHandler {
	users := user.NewList(p.UserMatchCase)
	for _, u := range p.UsersPlain {
		errorHandler.LogError(users.AddPlain(u.Username, u.Password))
	}
	for _, u := range p.UsersBase64 {
		errorHandler.LogError(users.AddBase64(u.Username, u.Password))
	}
	for _, u := range p.UsersMd5 {
		errorHandler.LogError(users.AddMd5(u.Username, u.Password))
	}
	for _, u := range p.UsersSha1 {
		errorHandler.LogError(users.AddSha1(u.Username, u.Password))
	}
	for _, u := range p.UsersSha256 {
		errorHandler.LogError(users.AddSha256(u.Username, u.Password))
	}
	for _, u := range p.UsersSha512 {
		errorHandler.LogError(users.AddSha512(u.Username, u.Password))
	}

	muxHandler := serverHandler.NewMultiplexer(p, *users, theme, logger, errorHandler)
	pathTransformHandler := serverHandler.NewPathTransformer(p.PrefixUrls, p.BaseUrls, muxHandler)

	vhostHandler := &VhostHandler{
		p:            p,
		logger:       logger,
		errorHandler: errorHandler,
		Handler:      pathTransformHandler,
	}

	return vhostHandler
}
