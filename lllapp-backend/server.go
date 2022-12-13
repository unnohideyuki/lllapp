package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

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

func getLastRead(c echo.Context) error {
	id := c.Param("id")
	num := c.Param("num")

	fileName := data_dir + "/" + id + "/books/" + num + "/pghistory.txt"

	fp, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	defer fp.Close()

	scanner := bufio.NewScanner(fp)

	var p0 string
	var p1 string
	var d string

	for scanner.Scan() {
		t := strings.TrimSpace(string(scanner.Text()))
		if len(t) == 0 {
			continue
		}

		re := regexp.MustCompile(`^(\d+)-(\d+),\s+([\d\-]+)`)
		a := re.FindStringSubmatch(t)

		p0 = a[1]
		p1 = a[2]
		d = a[3]
	}

	return c.JSON(http.StatusOK, []string{p0, p1, d})
}

func getPageProgress(c echo.Context) error {
	id := c.Param("id")
	num := c.Param("num")

	fileName := data_dir + "/" + id + "/books/" + num + "/pghistory.txt"

	fp, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	defer fp.Close()

	mp := map[int]int{}
	scanner := bufio.NewScanner(fp)
	pagemx := 0

	for scanner.Scan() {
		t := strings.TrimSpace(string(scanner.Text()))
		if len(t) == 0 {
			continue
		}

		re := regexp.MustCompile(`^(\d+)-(\d+)`)
		a := re.FindStringSubmatch(t)

		p0, _ := strconv.Atoi(a[1])
		p1, _ := strconv.Atoi(a[2])

		if p1 > pagemx {
			pagemx = p1
		}

		for i := p0; i <= p1; i++ {
			ct, _ := mp[i]
			mp[i] = ct + 1
		}
	}

	pghis := make([]int, pagemx+1)
	for i := 0; i <= pagemx; i++ {
		ct, _ := mp[i]
		pghis[i] = ct
	}

	return c.JSON(http.StatusOK, pghis)
}

type BookInfo struct {
	Num     string
	Title   string
	Author  string
	Isbn    string
	ModTime string
}

func getBookInfo(c echo.Context) error {
	id := c.Param("id")
	num := c.Param("num")
	fileName := data_dir + "/" + id + "/books/" + num + "/info.json"
	bytes, _ := ioutil.ReadFile(fileName)
	var jsonObj BookInfo
	_ = json.Unmarshal(bytes, &jsonObj)
	jsonObj.Num = num
	// ModTime is not set
	return c.JSON(http.StatusOK, jsonObj)
}

func postPages(c echo.Context) error {
	id := c.Param("id")
	num := c.Param("num")
	fileName := data_dir + "/" + id + "/books/" + num + "/pghistory.txt"

	pgFrom, _ := strconv.Atoi(c.FormValue("pgfrom"))
	pgTo, _ := strconv.Atoi(c.FormValue("pgto"))
	if pgTo < pgFrom {
		pgTo = pgFrom
	}

	// TODO : file lock, syscall 系を使うので pghistory とは別にロックファイルを作るか？

	fp, err := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		log.Fatal(err)
	}
	defer fp.Close()

	fmt.Fprintf(fp, "%d-%d, %s\n", pgFrom, pgTo, time.Now().Format(time.RFC3339))

	return c.String(http.StatusOK, "postPages ok :"+id+", "+num+", ("+c.FormValue("pgfrom")+")-("+c.FormValue("pgto")+")")
}

func getBookList(c echo.Context) error {
	id := c.Param("id")
	dirName := data_dir + "/" + id + "/books"

	dirs, err := ioutil.ReadDir(dirName)
	if err != nil {
		log.Fatal(err)
	}

	var data []interface{}
	var files []string

	for _, dir := range dirs {
		if dir.IsDir() {
			fileName := filepath.Join(dirName, dir.Name(), "/info.json")
			bytes, _ := ioutil.ReadFile(fileName)
			var i BookInfo
			err = json.Unmarshal(bytes, &i)
			if err != nil {
				log.Fatal(err)
			}

			fileSystem := os.DirFS(filepath.Join(dirName, dir.Name()))
			rlogFileName := "pghistory.txt"
			fi, err := fs.Stat(fileSystem, rlogFileName)
			if err != nil {
				log.Fatal(err)
			}

			i.ModTime = fi.ModTime().Format(time.RFC3339)
			i.Num = dir.Name()
			data = append(data, i)
			files = append(files, fileName)
		}
	}

	return c.JSON(http.StatusOK, data)
}

func main() {
	e := echo.New()

	e.Use(middleware.CORSWithConfig(middleware.DefaultCORSConfig))

	e.GET("/users/:id/books/list", getBookList)

	e.GET("/users/:id/books/:num/info", getBookInfo)
	e.GET("/users/:id/books/:num/toc", getBookToc)
	e.GET("/users/:id/books/:num/pghistory", getPageProgress)
	e.GET("/users/:id/books/:num/lastread", getLastRead)
	e.POST("/users/:id/books/:num/postpages", postPages)

	e.Logger.Fatal(e.Start(":8080"))
}
