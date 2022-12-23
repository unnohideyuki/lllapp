package handler

import (
	"encoding/json"
	"io/fs"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/labstack/echo/v4"

	"lllapp-backend/config"
)

func GetBookList(c echo.Context) error {
	id := c.Param("id")
	dirName := config.DataDir + "/" + id + "/books"

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
