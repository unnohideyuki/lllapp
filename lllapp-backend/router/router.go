package router

import (
	h "lllapp-backend/handler"

	"github.com/labstack/echo/v4"
)

func Register(g *echo.Group) {
	users := g.Group("/users")
	// for User Page
	//users.GET("/:id/summary", h.GetUserSummary)
	users.GET("/:id/books/list", h.GetBookList)

	// for Reading Record
	users.GET("/:id/books/:num/info", h.GetBookInfo)
	users.GET("/:id/books/:num/toc", h.GetBookToc)
	users.GET("/:id/books/:num/pghistory", h.GetPageProgress)
	users.GET("/:id/books/:num/lastread", h.GetLastRead)
	users.POST("/:id/books/:num/postpages", h.PostPages)
}
