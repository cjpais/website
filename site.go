package main

import (
	"log"
	//"fmt"
	"time"

	"net/http"
	//"io/ioutil"
	"html/template"
)

// TODO what are these files going to look like and how can they be written
type Day struct {
	// note that we should be able to load a day from a raw file.
	// basically think of the use cases...
	// should I always be posting from the site?
	// if the site is unavailable can I download the data and upload it?
	Emoji byte
	Date time.Time
	Summary []byte
	Timezone time.Time
	Posts []Post
}

type Post struct {
	// if no emoji then what?
	Emoji byte
	Time time.Time
	Summary []byte
	Content []byte
}

/*
func (p *Post) save() error {
	// TODO add date
	filename := STORE + p.Title + ".site"
	return ioutil.WriteFile(filename)
}*/

var templates = template.Must(template.ParseFiles("templates/index.html"))

func indexHandler(w http.ResponseWriter, r *http.Request) {
	// handle actually loading the page.
	err := templates.ExecuteTemplate(w, "index.html", Post{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// TODO decide to use VueJS to dynamically add components?

func main() {
	http.HandleFunc("/", indexHandler)
	/*
	http.HandleFunc("/post", postHandler)
	http.HandleFunc("/new", newPostHandler)
	http.HandleFunc("/todo", todoHandler)
	*/
	log.Println("Serving...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
