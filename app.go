package main

import (
	"fmt"
	"http"
	//"html"
	"template"
	"os"
	"bufio"
	"bytes"
	"strings"
)

type Dataobject struct {
	HtmlFragment string
}

func newFragmentHandler(fragpath string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		index := template.MustParseFile("design.gotemp", nil)

		fragment, err := os.Open(fragpath)
		if err != nil {
			fmt.Print("Could not open fragment")
		}
		text, err := bufio.NewReader(fragment).ReadString(0)
		if err != os.EOF {
			fmt.Print("Reading failed")
		}
		dataobject := Dataobject{HtmlFragment: text}

		if err := index.Execute(w, dataobject); err != nil {
			fmt.Fprint(w, "<h1>Šablona zfailovala</h1>")
		}
		return
	}
}

func handleScanKod(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.String(), http.StatusInternalServerError)
		return
	}
	
	index := template.MustParseFile("design.gotemp", nil)
	
	buffer := bytes.NewBuffer(nil)
	
	var kod string
	var user string
	
	var instructions bool
	
	if(r.FormValue("submit") == ""){
		fmt.Println(r.URL.Path)
		urlsegs := strings.Split(r.URL.Path, "/", 3)
	
		/* leading "/" */
		if len(urlsegs) == 3 {
			kod = urlsegs[2]
		}

		
		for _, v := range r.Cookie {
			fmt.Printf("%v, %v", v.Name, v.Value)
			if v.Name == "user" {
				user = v.Value
				break
			}
		}

		instructions = true
						
	} else {
		//kodtovalidate := r.FormValue("kod")
		user = r.FormValue("user")
		
		if len(kod) < 32 && len(user) < 32 {
			usercooke := http.Cookie {Name: "user", Value: user}
			http.SetCookie(w, &usercooke)		
		}
		
		//validateEntry(user, kod)
	
		fragment := template.MustParseFile("validkod.templatefragment", nil)

		fragment.Execute(buffer, struct{ Status string; Error string; Body_plus int; Body_total int }{"udf", "Errosddfe", 5 , 8})
	
	}
	
	fragment := template.MustParseFile("kod.templatefragment", nil)
		
	fragment.Execute(buffer, struct{ Kod string; User string; Instructions bool}{kod, user, instructions})
	
	dataobject := Dataobject{HtmlFragment: buffer.String()}

	if err := index.Execute(w, dataobject); err != nil {
		fmt.Fprint(w, "<h1>Šablona zfailovala</h1>")
	}
	

	return
}

func setCookie(w http.ResponseWriter, r *http.Request) {
	for _, v := range r.Cookie {
		fmt.Fprintln(w, v.Name)
	}
	cookie := http.Cookie{Name : "name", Value: "value", Path : "/cookie"}
	http.SetCookie(w, &cookie)
	fmt.Fprint(w, "AAA")
	
}

func main() {

	http.HandleFunc("/", newFragmentHandler("uvod.htmlfragment"))
	http.HandleFunc("/uvod/", newFragmentHandler("uvod.htmlfragment"))
	http.HandleFunc("/pravidla/", newFragmentHandler("pravidla.htmlfragment"))
	http.HandleFunc("/kod/", handleScanKod)
	http.HandleFunc("/cookie/", setCookie)
	http.Handle("/static/", http.FileServer(".", "/static/"))
	if err := http.ListenAndServe("192.168.1.102:8080", nil); err != nil {
		fmt.Println("Error: " + err.String())
	}
}

type QRCode struct {
	no int
	code string
	points int
}

type PointAccount struct {
	identifier string
	balance int
}
