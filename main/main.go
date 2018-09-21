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

//Project object represent each project
type Project struct {
	ID          int
	Title       string
	Description string
	Duration    float64
	Cost        float64
	Sector      string
	Time        []uint8
	Image       string
}

//Comments are either "Support" or "Against"
type Comments struct {
	Comments string
	Time     []uint8
}

func main() {
	http.Handle("/public/", http.StripPrefix(("/public/"), http.FileServer(http.Dir("src/sca1/public"))))
	http.HandleFunc("/", myHomePage)
	http.HandleFunc("/signup", signup)
	http.HandleFunc("/login", login)
	http.HandleFunc("/new-project", newProject)
	http.HandleFunc("/comment", comment)
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

	t, _ := template.ParseFiles("src/sca1/templates/home.html", "src/sca1/templates/header.html", "src/sca1/templates/footer.html", "src/sca1/templates/comment.html")

	err = t.ExecuteTemplate(w, "header.html", nil)
	checkerr(err)

	//iterating over the database to pull out Projects data and supplying it to the front-end...
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
		myProject := Project{id, title, description, duration, cost, sector, time, imageName}
		// myComments := Comments{comments}
		commentStmt, err := db.Prepare("SELECT comment, action, time FROM comment where project_id=?")
		checkerr(err)

		commentRes, err := commentStmt.Query(id)
		checkerr(err)
		err = t.ExecuteTemplate(w, "home.html", myProject)
		checkerr(err)

		for commentRes.Next() {
			var comment string
			var action int
			var time []uint8

			err := commentRes.Scan(&comment, &action, &time)
			checkerr(err)

			myComment := Comments{comment, time}
			if myComment.Comments != "" {
				t.ExecuteTemplate(w, "comment", myComment)
			}
		}

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

		myHomePage(w, req)

	} else {
		fmt.Println("Method is get...")
	}
}

//A func that handls comment
func comment(w http.ResponseWriter, req *http.Request) {
	if req.Method == "POST" {
		err := req.ParseForm()
		checkerr(err)

		for k, v := range req.Form {
			fmt.Printf("%s: %s\n", k, v)
		}
		comment := req.FormValue("comment")
		projectID := req.FormValue("project_id")
		action := req.FormValue("action")
		// actio

		fmt.Println("The id is: " + projectID)
		fmt.Println("action is :" + action)
		// checkerr(err)

		db, err := sql.Open("mysql", "root:atinuke22@/test")
		checkerr(err)
		defer db.Close()

		stmt, err := db.Prepare("INSERT comment SET project_id=?, comment=?, action=?")
		checkerr(err)

		res, err := stmt.Exec(projectID, comment, action)
		checkerr(err)
		fmt.Println(res.LastInsertId())

		myHomePage(w, req)

	}
}

//a simple func to checkerr
func checkerr(err error) {
	if err != nil {
		panic(err)
	}
}
