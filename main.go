package main 

import ( 
  "log" 
  "net/http" 
  "text/template"
  "path/filepath"
  "sync"
  "flag"

) 


type templateHandler struct { 
	once     sync.Once 
	filename string 
	templ    *template.Template 
} 
// ServeHTTP handles the HTTP request. 
func (this *templateHandler) ServeHTTP(w http.ResponseWriter, r  *http.Request) { 
	this.once.Do(func() { 
		this.templ = template.Must(
						template.ParseFiles(
							filepath.Join("templates", this.filename),
						),
					) 
	}) 
	this.templ.Execute(w, r) 
}


func main() {   
	var addr = flag.String("addr", ":8080", "The addr of the  application.") 
	flag.Parse() // parse the flags 
	r := newRoom() 
	http.Handle("/", &templateHandler{filename: "chat.html"}) 
	http.Handle("/room", r) 
	// get the room going 
	go r.run() 
	// start the web server 
	log.Println("Starting web server on", *addr) 
	if err := http.ListenAndServe(*addr, nil); err != nil { 
	  log.Fatal("ListenAndServe:", err) 
	} 
} 
