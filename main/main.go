package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
)

//Post is good
type Post struct {
	Title       string
	Description string
	Time        string
}

//Project object represent each project
type Project struct {
	Title       string
	Description string
	Duration    float64
	Cost        float64
	Sector      string
	Time        []uint8
	Image       string
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

//Handels requests to the hompage

func myHomePage(w http.ResponseWriter, req *http.Request) {

	db, err := sql.Open("mysql", "root:atinuke22@tcp(127.0.0.1:3306)/test")
	checkerr(err)
	defer db.Close()
	rows, err := db.Query("SELECT * FROM project ORDER BY time DESC")
	checkerr(err)

	t, _ := template.ParseFiles("src/sca1/templates/home.html", "src/sca1/templates/header.html", "src/sca1/templates/footer.html")
	// Description := "Etiam mollis lectus et laoreet faucibus. Suspendisse facilisis tellus velit, varius scelerisque leo sollicitudin at. Aliquam erat volutpat. Sed et varius quam. Cras aliquam odio a odio ornare, ut dapibus tortor cursus. Vestibulum eget erat libero. Duis ornare ligula ac elit gravida, at hendrerit dolor cursus. Proin sit amet magna eros. Proin posuere nibh a consectetur aliquam. Aenean non ligula faucibus urna finibus blandit a ut ligula."
	// Title := "Trans Mountain's delay costs Canadians $40 million every day."

	t.ExecuteTemplate(w, "header.html", nil)
	for rows.Next() {
		var id int
		var uid int
		var title string
		var description string
		var cost float64
		var duration float64
		var sector string
		var time []uint8

		err := rows.Scan(&id, &uid, &title, &description, &cost, &duration, &time, &sector)
		checkerr(err)

		//Get the image from project_image table
		stmt, err := db.Prepare("SELECT image_name FROM project_media WHERE project_id =?")
		checkerr(err)
		var imageName string
		err = stmt.QueryRow(id).Scan(&imageName)
		// fmt.Println(imageName)
		myProject := Project{title, description, duration, cost, sector, time, imageName}
		t.ExecuteTemplate(w, "home.html", myProject)
	}

	t.ExecuteTemplate(w, "footer.html", nil)

}

//Handels requests to the signup page
func signup(w http.ResponseWriter, req *http.Request) {
	t, _ := template.ParseFiles("src/sca1/templates/signup.html", "src/sca1/templates/header.html", "src/sca1/templates/footer.html")
	t.ExecuteTemplate(w, "header.html", nil)
	t.ExecuteTemplate(w, "signup.html", nil)
	t.ExecuteTemplate(w, "footer.html", nil)

}

//Handels requests to the login page
func login(w http.ResponseWriter, req *http.Request) {
	t, _ := template.ParseFiles("src/sca1/templates/login.html", "src/sca1/templates/header.html", "src/sca1/templates/footer.html")
	t.ExecuteTemplate(w, "header.html", nil)
	t.ExecuteTemplate(w, "login.html", nil)
	t.ExecuteTemplate(w, "footer.html", nil)

}

//Handels requests for adding a new project

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
		f, err := os.OpenFile("src/sca1/public/img/test/"+image, os.O_WRONLY|os.O_CREATE, 0666)
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
		// time := make([]uint8, 15)

		// }

		db, err := sql.Open("mysql", "root:atinuke22@tcp(127.0.0.1:3306)/test")

		//Prepare stmt top avoid sql injection
		stmt, err := db.Prepare("INSERT project SET user_id=?, title=?, description=?, cost=?, duration=?, sector=?")
		checkerr(err)

		//execute the statement abd store the returned value in res
		res, err := stmt.Exec(1, title, description, cost, duration, sector)
		checkerr(err)

		//
		id, err := res.LastInsertId()

		stmtImg, err := db.Prepare("INSERT project_media SET project_id=?, image_name=?")
		checkerr(err)

		resImg, err := stmtImg.Exec(id, image)
		checkerr(err)
		fmt.Println(resImg.LastInsertId())

		// t, err := template.ParseFiles("src/sca1/templates/temp.html", "src/sca1/templates/header.html", "src/sca1/templates/footer.html")
		// // if err != nil {
		// t.ExecuteTemplate(w, "header.html", nil)
		// t.ExecuteTemplate(w, "temp.html", Project{title, description, duration, cost, sector, time, image})
		// t.ExecuteTemplate(w, "footer.html", nil)

		myHomePage(w, req)

		// fmt.Println(title)
		// fmt.Println(description)
		// fmt.Println(duration)
		// fmt.Println(cost)
		// fmt.Println(sector)
		// fmt.Println(time)

	} else {
		fmt.Println("Method is get...")
	}
}

// func showTime(sec int64) (time string) {
// 	if sec < 60 {
// 		return strconv.FormatInt(sec, 10) + " sec ago"
// 	}
// 	if sec >= 60 && sec < 60*60 {
// 		return strconv.FormatInt(int64(math.Floor(float64(sec/60))), 10) + " min ago"
// 	}
// 	if sec >= 60*60 && sec < 60*60*60 {
// 		return strconv.FormatInt(int64(math.Floor(float64(sec/60*60))), 10) + " hours ago"
// 	}
// }

func checkerr(err error) {
	if err != nil {
		panic(err)
	}
}
