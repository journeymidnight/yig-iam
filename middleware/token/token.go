package token

import (
	"github.com/journeymidnight/yig-iam/helper"
	"github.com/journeymidnight/yig-iam/db"
	"gopkg.in/iris.v4"
	. "github.com/journeymidnight/yig-iam/api/datatype"
)

type Middleware struct {
	Config Config
}

func New(cfg ...Config) *Middleware {

	var c Config
	if len(cfg) == 0 {
		c = Config{}
	} else {
		c = cfg[0]
	}

	if c.ContextKey == "" {
		c.ContextKey = DefaultContextKey
	}

	if c.CookieKey == "" {
		c.CookieKey = DefaultCookieKey
	}

	return &Middleware{Config: c}
}

func (m *Middleware) Get(ctx *iris.Context) *string {
	return ctx.Get(m.Config.ContextKey).(*string)
}

// Serve the middleware's action
func (m *Middleware) Serve(ctx *iris.Context) {
	query := QueryRequest{}
	if err := ctx.ReadJSON(&query); err != nil {
		ctx.JSON(iris.StatusOK, QueryResponse{RetCode:4000,Message:"bad request"})
		return
	}
	ctx.Set("queryRequest", query)
	helper.Logger.Println(5, "enter middleware", query)
	if query.Action == ACTION_ConnectService || query.Action == ACTION_DescribeAccessKeys{
		ctx.Next()
	} else if query.Token == "" {
		ctx.JSON(iris.StatusOK, QueryResponse{RetCode:4000,Message:"Not authorized, Need login first"})
		return
	} else {
		err := m.CheckToken(ctx, query.Token)
		if err != nil {
			ctx.JSON(iris.StatusOK, QueryResponse{RetCode:4000,Message:"Token invalid"})
			return
		}
	}
	ctx.Next()
}

func (m *Middleware) CheckToken(ctx *iris.Context, token string) error {
	record, err:= db.GetTokenRecord(token)
	if err != nil {
		return err
	}
	ctx.Set(m.Config.ContextKey, record)
	helper.Logger.Println(5, "token set success", record)
	return err
}
