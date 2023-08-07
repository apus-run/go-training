package ginx

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Context a wrapper of gin.Context
type Context struct {
	*gin.Context
}

// HandlerFunc defines the handler to wrap gin.Context
type HandlerFunc func(c *Context)

// Handle convert HandlerFunc to gin.HandlerFunc
func Handle(h HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := &Context{
			c,
		}
		h(ctx)
	}
}

type ginKey struct{}

// NewGinContext returns a new Context that carries gin.Context value.
func NewGinContext(ctx context.Context, c *gin.Context) context.Context {
	return context.WithValue(ctx, ginKey{}, c)
}

// FromGinContext returns the gin.Context value stored in ctx, if any.
func FromGinContext(ctx context.Context) (c *gin.Context, ok bool) {
	c, ok = ctx.Value(ginKey{}).(*gin.Context)
	return
}

const (
	// CodeOK means a successful response
	CodeOK = 0
	// CodeErr means a failure response
	CodeErr = 1
)

// Response defines HTTP JSON response
type Response struct {
	Code    int         `json:"code"`
	Msg     string      `json:"msg"`
	Data    interface{} `json:"data"`
	Details []string    `json:"details,omitempty"`
}

type Pagination struct {
	// pageNumber 当前页
	pageNumber int `json:"page_number"`
	// pageSize 分页数
	pageSize int `json:"page_size"`
	// total means total page count
	total int `json:"total"`
	// data 数据
	data []interface{} `json:"data"`
	// totalPages 总页数
	totalPages int  `json:"total_pages,omitempty"`
	hasNext    bool `json:"has_next,omitempty"`
}

// JSON returns JSON response
// e.x. {"code":<code>, "msg":<msg>, "data":<data>, "details":<details>}
func (c *Context) JSON(httpStatus int, resp Response) {
	c.Context.JSON(httpStatus, resp)
}

// JSONOK returns JSON response with successful business code and data
// e.x. {"code":0, "msg":"成功", "data":<data>}
func (c *Context) JSONOK(data interface{}) {
	j := new(Response)
	j.Code = CodeOK
	j.Msg = "ok"

	if data == nil {
		j.Data = gin.H{}
	} else {
		j.Data = data
	}
	c.Context.JSON(http.StatusOK, j)
	return
}

// JSONE returns JSON response with failure business code ,msg and data
// e.x. {"code":<code>, "msg":<msg>, "data":<data>}
func (c *Context) JSONE(code int, msg string, data interface{}) {
	j := new(Response)
	j.Code = code
	j.Msg = msg
	switch d := data.(type) {
	case error:
		j.Data = d.Error()
	default:
		j.Data = data
	}

	c.Context.JSON(http.StatusOK, j)
	return
}

// JSONPage returns JSON response with pagination
// e.x. {"code":<code>, "msg":<msg>, "data":<data>, "pagination":<pagination>}
// <pagination> { "pageNumber":1, "pageSize":20, "total": 9 }
func (c *Context) JSONPage(data interface{}, p pagination.Pagination) {
	j := new(ResponsePagination)
	j.Code = CodeOK
	j.Data = data
	j.Pagination = p
	c.Context.JSON(http.StatusOK, j)
}

// NotFound 未找到相关路由
func (c *Context) NotFound() {
	c.String(http.StatusNotFound, "the route not found")
}
