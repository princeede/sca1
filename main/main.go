package main

import (
	"html/template"
	"log"
	"net/http"
	// "text/template"
)

//Post is good
type Post struct {
	Title       string
	Description string
	Time        string
}

func main() {
	http.Handle("/public/", http.StripPrefix(("/public/"), http.FileServer(http.Dir("src/sca1/public"))))
	http.HandleFunc("/", myHomePage)
	http.HandleFunc("/signup", signup)
	http.HandleFunc("/login", login)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServer: ", err)
	}
}

func myHomePage(w http.ResponseWriter, req *http.Request) {
	t, _ := template.ParseFiles("src/sca1/templates/home.html", "src/sca1/templates/header.html", "src/sca1/templates/footer.html")
	Description := "Etiam mollis lectus et laoreet faucibus. Suspendisse facilisis tellus velit, varius scelerisque leo sollicitudin at. Aliquam erat volutpat. Sed et varius quam. Cras aliquam odio a odio ornare, ut dapibus tortor cursus. Vestibulum eget erat libero. Duis ornare ligula ac elit gravida, at hendrerit dolor cursus. Proin sit amet magna eros. Proin posuere nibh a consectetur aliquam. Aenean non ligula faucibus urna finibus blandit a ut ligula."
	Title := "Trans Mountain's delay costs Canadians $40 million every day."
	myPost := Post{Title, Description, "2 days"}
	t.ExecuteTemplate(w, "header.html", nil)
	t.ExecuteTemplate(w, "home.html", myPost)
	t.ExecuteTemplate(w, "footer.html", nil)

}

func signup(w http.ResponseWriter, req *http.Request) {
	t, _ := template.ParseFiles("src/sca1/templates/signup.html", "src/sca1/templates/header.html", "src/sca1/templates/footer.html")
	t.ExecuteTemplate(w, "header.html", nil)
	t.ExecuteTemplate(w, "signup.html", nil)
	t.ExecuteTemplate(w, "footer.html", nil)

}

func login(w http.ResponseWriter, req *http.Request) {
	t, _ := template.ParseFiles("src/sca1/templates/login.html", "src/sca1/templates/header.html", "src/sca1/templates/footer.html")
	t.ExecuteTemplate(w, "header.html", nil)
	t.ExecuteTemplate(w, "login.html", nil)
	t.ExecuteTemplate(w, "footer.html", nil)

}
