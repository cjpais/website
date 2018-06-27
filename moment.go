package main

import (
	"os"
	"io"
	"log"
	"fmt"
	"time"

	"net/http"
	"io/ioutil"
	"encoding/json"
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

func getTime(d, t string) (time.Time, error) {
	datetime := d + "T" + t + ":00Z"
	return time.Parse(time.RFC3339, datetime)
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
	date := r.Form["date"][0]
	time := r.Form["time"][0]
	p.Time, _ = getTime(date, time)
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
	date := r.Form["date"][0]
	time := r.Form["time"][0]
	p.Time, _ = getTime(date, time)
	p.Timezone = r.Form["tz"][0]
	p.Summary = r.Form["summary"][0]

	// handle file
	file, header, err := r.FormFile("image")
	p.Path = getPath(p.Time) + "/" + header.Filename
	f, err := os.OpenFile(p.Path, os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		return err
	}
	defer file.Close()
	defer f.Close()
	io.Copy(f, file)
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
