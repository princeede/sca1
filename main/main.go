package main

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
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
	http.HandleFunc("/new-project", newProject)
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

//Project object represent each project
type Project struct {
	Title       string
	Description string
	Duration    float64
	Cost        float64
	Sector      string
	Time        int64
	Image       string
}

func newProject(w http.ResponseWriter, req *http.Request) {
	if req.Method == "POST" {
		err := req.ParseMultipartForm(32 << 20)
		if err != nil {
			fmt.Print("We are not good")
		}

		file, handler, err := req.FormFile("media")
		if err != nil {
			fmt.Println(err)
			return
		}

		defer file.Close()
		image := handler.Filename
		f, err := os.OpenFile("src/sca1/public/img/test/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer f.Close()
		io.Copy(f, file)
		title := req.FormValue("title")
		description := req.FormValue("description")
		duration, _ := strconv.ParseFloat(req.FormValue("duration"), 64)
		cost, _ := strconv.ParseFloat(req.FormValue("cost"), 64)
		sector := req.FormValue("sector")
		time := time.Now().Unix()

		t, err := template.ParseFiles("src/sca1/templates/temp.html", "src/sca1/templates/header.html", "src/sca1/templates/footer.html")
		// if err != nil {
		t.ExecuteTemplate(w, "header.html", nil)
		t.ExecuteTemplate(w, "temp.html", Project{title, description, duration, cost, sector, time, image})
		t.ExecuteTemplate(w, "footer.html", nil)

		fmt.Println(title)
		fmt.Println(description)
		fmt.Println(duration)
		fmt.Println(cost)
		fmt.Println(sector)
		fmt.Println(time)
		// }

	} else {
		fmt.Println("Method is get...")
	}
}

func myTemplate(templateName string) (temp *template.Template, err error) {
	t, err := template.ParseFiles("src/sca1/templates/"+templateName+".html", "src/sca1/templates/header.html", "src/sca1/templates/footer.html")
	return t, err
}
