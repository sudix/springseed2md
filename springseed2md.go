package main

import (
	"encoding/json"
	"fmt"
	"github.com/codegangsta/cli"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type Article struct {
	Name     string `json:"name"`
	Content  string `json:"content"`
	Notebook string `json:"notebook"`
	Id       string `json:"id"`
	Date     int64  `json:"date"`
}

func parseFile(filePath string) (Article, error) {
	var article Article
	b, err := ioutil.ReadFile(filePath)
	if err != nil {
		return article, err
	}
	if err := json.Unmarshal(b, &article); err != nil {
		return article, err
	}
	return article, nil
}

func writeMarkDown(filePath string, article Article) error {
	err := ioutil.WriteFile(filePath, []byte(article.Content), 0644)
	if err != nil {
		return err
	}
	return nil
}

func convert(fromDir, distDir string) error {
	fileInfos, err := ioutil.ReadDir(fromDir)
	if err != nil {
		return err
	}

	for _, fileInfo := range fileInfos {
		if fileInfo.IsDir() {
			continue
		}
		name := fileInfo.Name()
		if !strings.Contains(name, ".note") {
			continue
		}
		article, err := parseFile(filepath.Join(fromDir, name))
		if err != nil {
			return err
		}
		snakeFileName := strings.Replace(article.Name, " ", "_", -1)
		fileBase := strings.ToLower(snakeFileName) + ".md"

		fmt.Println(fileBase)
		mdFilePath := filepath.Join(distDir, fileBase)
		nil := writeMarkDown(mdFilePath, article)
		if err != nil {
			return err
		}

	}

	return nil
}

func main() {
	app := cli.NewApp()
	app.Name = "springseed2md notes converter"
	app.Usage = "notes_directory markdown_output_directory"
	app.Action = func(c *cli.Context) {
		args := c.Args()
		if len(args) < 2 {
			fmt.Println("Error, you need to specify both notes_directory and markdown_output_directory.")
			os.Exit(1)
		}
		fromDir := args[0]
		distDir := args[1]
		err := convert(fromDir, distDir)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}
	app.Run(os.Args)
}
