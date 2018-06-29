package main

import (
	"os"
	"log"
	"fmt"
	"time"
	"bytes"
	"strings"

	"net/http"
	"io/ioutil"
	"image/jpeg"
	"encoding/json"

	"github.com/nfnt/resize"
)

type Moment interface {
	save() error
	saveReq(*http.Request) error
	load(string) error
	time() time.Time
}

const (
	POST = "website.post"
	PHOTO = "website.photo"
)

func getTime(d, t, format string) (time.Time, error) {
	if format == "rfc3339" {
		datetime := d + "T" + t + ":00Z"
		return time.Parse(time.RFC3339, datetime)
	} else if format == "cjcustom" {
		datetime := d + ", " + t
		return time.Parse("Monday, January 02, 2006, 3:04:05 PM:", datetime)
	}
	return time.Time{}, fmt.Errorf("time format incorrect")
}

type Post struct {
	Time time.Time
	Timezone string
	Summary string
	Content string
}

func (p *Post) time() time.Time {
	return p.Time
}

func (p *Post) save() error {
	path := getPath(p.Time)
	filename := path + "/website.post"
	content, err := json.Marshal(p)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filename, content, 0755)
}

func (p *Post) saveReq(r *http.Request) error {
	log.Println("saving post")
	date := strings.TrimSpace(r.Form["date"][0])
	time := strings.TrimSpace(r.Form["time"][0])
	format := r.Form["format"][0]
	p.Time, _ = getTime(date, time, format)
	p.Timezone = r.Form["tz"][0]
	p.Summary = r.Form["summary"][0]
	p.Content = r.Form["content"][0]
	return p.save()
}

func (p *Post) load(path string) error {
	data, err := ioutil.ReadFile(path + "/website.post")
	err = json.Unmarshal(data, p)
	return err
}

type Photo struct {
	Time time.Time
	Timezone string
	Summary string
	Path string
	Fullpath string
}

func (p *Photo) time() time.Time {
	return p.Time
}

func (p *Photo) save() error {
	path := getPath(p.Time)
	filename := path + "/website.photo"
	content, err := json.Marshal(p)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filename, content, 0755)
}

func (p *Photo) saveReq(r *http.Request) error {
	log.Println("saving photo")
	date := strings.TrimSpace(r.Form["date"][0])
	time := strings.TrimSpace(r.Form["time"][0])
	format := r.Form["format"][0]
	p.Time, _ = getTime(date, time, format)
	p.Timezone = r.Form["tz"][0]
	p.Summary = r.Form["summary"][0]

	// handle file
	// TODO multiple photo uploads
	//files := r.MultipartForm.File["images"]

	file, header, err := r.FormFile("images")
	if err != nil {
		return err
	}
	defer file.Close()
	filedata, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}

	p.Fullpath = getPath(p.Time) + "/fullsize_" + header.Filename
	p.Path = getPath(p.Time) + "/" + header.Filename

	// write out the fullsize image
	err = ioutil.WriteFile(p.Fullpath, filedata, 0600)
	if err != nil {
		return err
	}

	// resize the image 
	image, err := jpeg.Decode(bytes.NewReader(filedata))
	if err != nil {
		return err
	}

	m := resize.Resize(620, 0, image, resize.Lanczos3)

	ri, err := os.Create(p.Path)
	if err != nil {
		return err
	}
	defer ri.Close()
	err = jpeg.Encode(ri, m, nil)

	return p.save()
}

func (p *Photo) load(path string) error {
	data, err := ioutil.ReadFile(path + "/website.photo")
	err = json.Unmarshal(data, p)
	return err
}

func getPath(t time.Time) string {
	path := fmt.Sprintf("timeline/%v/%02d/%v/%v", t.Year(), t.Month(), t.Day(), t.Unix())
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.MkdirAll(path, os.ModePerm)
	}
	return path
}
