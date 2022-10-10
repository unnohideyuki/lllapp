package main

import (
	"bufio"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

const data_dir = "./data"

func getBookToc(c echo.Context) error {
	id := c.Param("id")
	num := c.Param("num")

	fileName := data_dir + "/" + id + "/books/" + num + "/toc.txt"

	fp, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	defer fp.Close()

	scanner := bufio.NewScanner(fp)

	toc := [](map[string]string){}

	titleName := ""
	pageFrom := -1

	for scanner.Scan() {
		t := strings.TrimSpace(string(scanner.Text()))
		if len(t) == 0 {
			continue
		}
		if c := t[0:1]; c == "%" {
			continue
		}
		a := strings.Split(t, ",")
		pg, _ := strconv.Atoi(strings.Trim(a[0], "()"))

		if titleName != "" {
			pageTo := pg - 1
			if pageFrom < 0 {
				pageTo = -1
			}

			mp := map[string]string{
				"name":      titleName,
				"page_from": strconv.Itoa(pageFrom),
				"page_to":   strconv.Itoa(pageTo),
			}
			toc = append(toc, mp)
		}

		if c := a[0][0:1]; c == "(" {
			pageFrom = -1
		} else {
			pageFrom = pg
		}
		if len(a) >= 2 {
			titleName = strings.Trim(a[1], " \t")
		} else {
			titleName = ""
		}
	}

	if err = scanner.Err(); err != nil {
		panic(err)
	}

	return c.JSON(http.StatusOK, toc)
}

func main() {
	e := echo.New()

	e.Use(middleware.CORSWithConfig(middleware.DefaultCORSConfig))

	e.GET("/users/:id/books/:num/toc", getBookToc)
	e.Logger.Fatal(e.Start(":8080"))
}
