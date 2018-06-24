package main

import (
	"os"
	"log"
	"fmt"
	"time"
	"strconv"

	"net/http"
	"math/rand"
	"io/ioutil"
	"encoding/json"
	"path/filepath"
)

type Year struct {
	Year int
	Months []*Month
}

type Month struct {
	Year *Year
	Month int
	Days []*Day
}

type Day struct {
	Month *Month
	Day int
	Posts []*Post
}

type Post struct {
	Time time.Time
	//Timezone string
	Day *Day
	Summary string
	Content string
}

func (p *Post) save() error {
	t := p.Time
	path := fmt.Sprintf("timeline/%v/%02d/%v/%v", t.Year(), t.Month(), t.Day(), t.Unix())
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.MkdirAll(path, os.ModePerm)
	}
	sfilename := path + "/summary"
	cfilename := path + "/content"
	summary := []byte(p.Summary)
	content := []byte(p.Content)
	e := ioutil.WriteFile(sfilename, summary, 0755)
	if e != nil {
		return e
	}
	return ioutil.WriteFile(cfilename, content, 0755)
}

func createPost(t time.Time, summary, content string) *Post {
	return &Post{Time: t, Summary: summary, Content: content}
}

func loadPost(path string) (*Post, error) {
	summary := ""
	content := ""
	sfilename := path + "/summary"
	cfilename := path + "/content"
	timeint, _ := strconv.ParseInt(filepath.Base(path), 10, 64)
	time := time.Unix(timeint, 0)
	if _, err := os.Stat(sfilename); err == nil {
		// summary exists so load it
		data, err := ioutil.ReadFile(sfilename)
		if err != nil {
			return nil, err
		}
		summary = string(data[:])
	}
	if _, err := os.Stat(cfilename); err == nil {
		data, err := ioutil.ReadFile(cfilename)
		if err != nil {
			return nil, err
		}
		content = string(data[:])
	}
	return &Post{Time: time, Summary: summary, Content: content}, nil
}

func loadDay(path string) (*Day, error) {
	var posts []*Post
	// get posts
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}
	for _, file := range files {
		if file.IsDir() {
			ppath := path + "/" + file.Name()
			//log.Println("loading post", ppath)
			post, err := loadPost(ppath)
			if err != nil {
				return nil, err
			}
			posts = append(posts, post)
		}
	}
	d, _ := strconv.Atoi(filepath.Base(path))
	day := Day{Day: d, Posts: posts}
	return &day, nil
}

func loadDays(path string) ([]*Day, error) {
	var days []*Day
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}
	for _, file := range files {
		dpath := path + "/" + file.Name()
		//log.Println("loading day", dpath)
		day, err := loadDay(dpath)
		if err != nil {
			return nil, err
		}
		days = append(days, day)
	}
	return days, nil
}

func loadMonths(path string) ([]*Month, error) {
	var months []*Month
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}
	for _, file := range files {
		// for each year get months
		mpath := path + "/" + file.Name()
		//log.Println("loading month", mpath)
		days, err := loadDays(mpath)
		if err != nil {
			return nil, err
		}
		m, _ := strconv.Atoi(file.Name())
		mo := Month{Month: m, Days: days}
		months = append(months, &mo)
	}
	return months, nil
}

func loadYears() ([]*Year, error) {
	var years []*Year
	files, err := ioutil.ReadDir("timeline")
	if err != nil {
		return nil, err
	}
	for _, file := range files {
		// for each year get months
		path := "timeline/" + file.Name()
		//log.Println("loading year", path)
		months, err := loadMonths(path)
		if err != nil {
			return nil, err
		}
		y, _ := strconv.Atoi(file.Name())
		yr := Year{Year: y, Months: months}
		years = append(years, &yr)
	}
	return years, nil
}

func daysHandler(w http.ResponseWriter, r *http.Request) {
	// get posts from all time
	years, err := loadYears()
	if err != nil {
		log.Println(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(years)
}

func writeTestPosts() {
	log.Println("Writing New Set of Days")
	// dirty constants
	summary := []string{"", "lorem ipsum"}
	content := []string{"Lorem ipsum dolor sit amet, consectetur adipiscing elit. Donec mattis erat nec vehicula suscipit. Aenean posuere ante justo, a pretium risus aliquet nec. Nullam accumsan augue odio, vitae suscipit ligula lobortis nec. Curabitur semper at mauris id tempor. In vitae ipsum quis libero commodo facilisis. In ac faucibus ipsum. Sed fringilla enim non gravida imperdiet. Curabitur vitae est blandit, viverra sapien vel, rutrum nisl.", "hi", "test", "noooooo", "CJ"}
	numPosts := 1000

	min := time.Date(2014, 5, 30, 0, 0, 0, 0, time.UTC).Unix()
	max := time.Now().Unix()
	dt := max - min

	for i := 0; i < numPosts; i++ {
		t := time.Unix(rand.Int63n(dt) + min, 0)
		randc := rand.Intn(len(content))
		rands := rand.Intn(len(summary))
		post := createPost(t, summary[rands], content[randc])
		post.save()
	}
}

func main() {
	// if any command line argument write new days
	if len(os.Args) > 1 {
		os.RemoveAll("timeline/")
		os.MkdirAll("timeline/", os.ModePerm)
		writeTestPosts()
	}

	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	http.Handle("/", http.FileServer(http.Dir("static/html")))
	http.HandleFunc("/days/", daysHandler)
	log.Println("Serving cjpais.com...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
