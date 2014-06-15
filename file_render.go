package main

import (
	//"bytes"
	"encoding/json"
	//"flag"
	"fmt"
	//"github.com/howeyc/fsnotify"
	//"code.google.com/p/go.exp/fsnotify"
	//"github.com/russross/blackfriday"
	"io/ioutil"
	//"log"
	//"net/http"
	"os"
	//"os/exec"
	"path/filepath"
	//"sort"
	"strings"
	"text/template"
	//"time"
)

// Renders current site to html
func Render_Site() error {
	var err error
	filepath.Walk(config.SourceDir, Walker)
	var pages PagesSlice
	for _, dir := range site.Directories {

		readglob := dir + "/*.md"
		var dirfiles, _ = filepath.Glob(readglob)

		// loop through files in directory
		for _, file := range dirfiles {
			//fmt.Println("  File:", file)
			outfile := filepath.Base(file)
			outfile = strings.Replace(outfile, ".md", ".html", 1)

			// read & parse file for parameters
			page := readParseFile(file)
			page.OutFile = dir + "/" + outfile
			// create array of parsed pages
			pages = append(pages, page)
		}
	}
	//fmt.Printf("%v\n", pages)
	layoutsglob := config.TemplateDir + "/*.html"
	_, err = template.ParseGlob(layoutsglob)
	if err != nil {
		PrintErr("Error Parsing Templates: ", err)
		os.Exit(1)
	}

	for _, page := range pages {
		html := applyTemplates(page)
		fmt.Println(page.Url)
		err = WritePage(page, html)
		if err != nil {
			PrintErr("Cant generate site", err)
			os.Exit(1)
		}
	}
	// Generate index listings...
	for _, dir := range site.Directories {
		html, page, _ := getDirectoryListing(dir)
		err = WritePage(page, html)
		if err != nil {
			PrintErr("Error writing index files", err)
			os.Exit(1)
		}
		WriteJson(page, dir)

	}

	return nil

}

func WritePage(page Page, html string) error {

	outfile := config.PublishDir + page.OutFile
	Printvln(" Writing File:", outfile)
	err := ioutil.WriteFile(outfile, []byte(html), 0644)
	if err != nil {
		return err
	}
	return nil
}

func WriteJson(page Page, dir string) error {
	res, _ := json.Marshal(page.Pages)
	outfile := dir + "/index.json"
	err := ioutil.WriteFile(outfile, []byte(res), 0644)
	if err != nil {
		return err
	}
	return nil
}

// WalkFn that fills SiteStruct with data.
func Walker(fn string, fi os.FileInfo, err error) error {
	if err != nil {
		PrintErr("Walker: ", err)
		return nil
	}

	if fi.IsDir() {
		site.Categories = append(site.Categories, fi.Name())
		site.Directories = append(site.Directories, fn)
		return nil
	} else {
		ext := filepath.Ext(fn)
		if ext == ".md" {
			site.Files = append(site.Files, fn)
		}
		return nil
	}
	return nil

}
