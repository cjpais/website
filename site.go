package main

import (
	"os"
	"log"
	//"fmt"
	//"time"
	"sort"
	"bytes"
	"strconv"

	"net/http"
	//"math/rand"
	"io/ioutil"
	"encoding/json"
	"path/filepath"

	"golang.org/x/crypto/scrypt"
	//"github.com/gorilla/sessions"
)


type Year struct {
	String string
	Int int
	Months []*Month
}

type Month struct {
	//Year *Year
	String string
	Int int
	Days []*Day
}

type Day struct {
	//Month *Month
	String string
	Int int
	Moments []map[string]Moment
}

func loadMoment(path string) (map[string]Moment, error) {
	// go read the dir and check the filename and load proper moment
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}

	moment := map[string]Moment{}
	for _, file := range files {
		switch file.Name() {
			case POST:
				moment[POST] = &Post{}
				moment[POST].load(path)
				return moment, nil
			case PHOTO:
				moment[PHOTO] = &Photo{}
				moment[PHOTO].load(path)
				return moment, nil
		}
	}
	return nil, nil
}

func loadDay(path string) (*Day, error) {
	var moments []map[string]Moment

	files, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}
	for _, file := range files {
		if file.IsDir() {
			mpath := path + "/" + file.Name()
			moment, err := loadMoment(mpath)
			if err != nil {
				return nil, err
			}
			moments = append(moments, moment)
		}
	}
	d, _ := strconv.Atoi(filepath.Base(path))
	day := Day{Int:d, Moments:moments}
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
		days = append([]*Day{day}, days...)
	}
	sort.Slice(days, func(i, j int) bool {
		return days[i].Int > days[j].Int
	})
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
		mo := Month{Int: m, Days: days}
		months = append([]*Month{&mo}, months...)
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
		yr := Year{Int: y, Months: months}
		years = append([]*Year{&yr}, years...)
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

/*
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
		post := createMoment(t, summary[rands], content[randc])
		post.save()
	}
}
*/

func newPage(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "cjpais.com")
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}
	http.ServeFile(w, r, "static/html/new.html")
}

func loginPage(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/html/login.html")
}

func index(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/html/index.html")
}

func login(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "cjpais.com")
	r.ParseForm()

	// TODO replace with constant time compare
	if pass, ok := r.Form["password"]; ok {
		k, _ := scrypt.Key([]byte(pass[0]), salt, 32768, 8, 1, 32)
		if user, ok := r.Form["username"]; ok {
			// use randomly generated username to login
			if bytes.Compare(k, genpass) == 0 && user[0] == genusername {
				log.Println("CJ authenticated from:", r.RemoteAddr)
				session.Values["authenticated"] = true
				session.Save(r, w)
			}
		} else {
			session.Values["authenticated"] = false
		}
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func newMoment(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "cjpais.com")
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	r.ParseMultipartForm(0)
	// TODO sanitize
	switch r.Form["type"][0] {
		case POST:
			post := Post{}
			post.saveReq(r)
		case PHOTO:
			photo := Photo{}
			photo.saveReq(r)
	}
}

func auth(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "cjpais.com")
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		w.Write([]byte("false"))
	}
	w.Write([]byte("true"))
}

func main() {
	// if any command line argument write new days
	if len(os.Args) > 1 {
		os.RemoveAll("timeline/")
		os.MkdirAll("timeline/", os.ModePerm)
		//writeTestPosts()
	}

	// serve filesystem parts
	sfs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", sfs))
	tfs := http.FileServer(http.Dir("timeline"))
	http.Handle("/timeline/", http.StripPrefix("/timeline/", tfs))

	// serve home/login/new
	http.HandleFunc("/", index)
	http.HandleFunc("/login/", loginPage)
	http.HandleFunc("/new/", newPage)

	// api
	http.HandleFunc("/api/days", daysHandler)
	http.HandleFunc("/api/login", login)
	http.HandleFunc("/api/auth", auth)
	http.HandleFunc("/api/new", newMoment)

	log.Println("Serving cjpais.com...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
