package main

import (
	"io"
	"os"
	"log"
	//"fmt"
	"time"
	"flag"
	"sort"
	"bytes"
	"strconv"

	"net/http"
	//"math/rand"
	"io/ioutil"
	//"crypto/tls"
	"encoding/json"
	"path/filepath"

	"golang.org/x/crypto/scrypt"
	//"golang.org/x/crypto/acme/autocert"
	//"github.com/gorilla/sessions"
)

var (
	clear = flag.Bool("clear", false, "Clears the timeline")
	prod = flag.Bool("prod", false, "If in production or not")
)

type Year struct {
	//String string
	Int int
	Months []*Month
}

type Month struct {
	//Year *Year
	//String string
	Int int
	Days []*Day
}

type Day struct {
	//Month *Month
	//String string
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
			moments = append([]map[string]Moment{moment}, moments...)
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

func logout(w http.ResponseWriter, r *http.Request) {
	log.Println("logout")
	session, _ := store.Get(r, "cjpais.com")
	session.Values["authenticated"] = false
	session.Save(r, w)
	http.Redirect(w, r, "/", http.StatusSeeOther)
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
	//var data []byte;
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
	daysHandler(w, r)
}

func removeMoment(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "cjpais.com")
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	r.ParseMultipartForm(0)
	// assume all times have independent data (not strictly the case)
	//momentType := r.Form["type"][0]
	t := r.Form["time"][0]
	time, err := time.Parse(time.RFC3339, t)
	if err != nil {
		log.Println("error converting time string", err)
		return
	}
	path := getPath(time)
	os.RemoveAll(path)

	// if nothing left in the dir remove it
	daypath := filepath.Dir(path)
	files, err := ioutil.ReadDir(daypath)
	if err != nil {
		log.Println("error reading directory", daypath, err)
		return
	}
	if len(files) == 0 {
		os.Remove(daypath)
	}

	daysHandler(w, r)
}

func faviconHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "static/img/favicon.ico")
}

func auth(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "cjpais.com")
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		w.Write([]byte("false"))
	} else {
		w.Write([]byte("true"))
	}
}

func main() {
	logfile, err := os.OpenFile("logs/website.log", os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		log.Println("error openinng file", err)
		os.Exit(1)
	}
	defer logfile.Close()
	mw := io.MultiWriter(os.Stdout, logfile)
	log.SetOutput(mw)

	flag.Parse()

	// if any command line argument write new days
	if *clear {
		log.Printf("removing timeline")
		os.RemoveAll("timeline/")
		os.MkdirAll("timeline/", os.ModePerm)
	}

	// serve filesystem parts
	sfs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", sfs))
	tfs := http.FileServer(http.Dir("timeline"))
	http.Handle("/timeline/", http.StripPrefix("/timeline/", tfs))

	// serve favicon
	http.HandleFunc("/favicon.ico", faviconHandler)

	// serve home/login/new
	http.HandleFunc("/", index)
	http.HandleFunc("/login/", loginPage)
	http.HandleFunc("/logout/", logout)
	http.HandleFunc("/new/", newPage)

	// api
	http.HandleFunc("/api/days", daysHandler)
	http.HandleFunc("/api/login", login)
	http.HandleFunc("/api/auth", auth)
	http.HandleFunc("/api/new", newMoment)
	http.HandleFunc("/api/remove", removeMoment)

	// set up server
	if *prod {
		log.Println("production server launched")
		/*
		d := []string{"cjpais.com", "www.cjpais.com", "206.189.227.218"}
		m := &autocert.Manager{
		    Cache:      autocert.DirCache("certs"),
		    Prompt:     autocert.AcceptTOS,
		    HostPolicy: autocert.HostWhitelist(d...),
		}
		s := &http.Server{
		    Addr:      ":https",
		    TLSConfig: &tls.Config{
			ServerName: "www.cjpais.com",
			GetCertificate: m.GetCertificate,
			},
		}
		*/

		//go http.ListenAndServe(":http", m.HTTPHandler(nil))
		log.Println(http.ListenAndServeTLS(":https",
						"/etc/letsencrypt/live/www.cjpais.com/fullchain.pem",
						"/etc/letsencrypt/live/www.cjpais.com/privkey.pem",
						nil))

		//listener := autocert.NewListener("www.cjpais.com")
		//log.Fatal(http.Serve(listener, nil))
	} else {
		log.Println("dev server launched on port :8080")
		log.Fatal(http.ListenAndServe(":8080", nil))
	}
}
