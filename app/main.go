package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
)

type PageVariables struct {
	Color       string
	Title       string
	Description string
}

func main() {
	http.HandleFunc("/", HomePage)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func HomePage(w http.ResponseWriter, r *http.Request) {
	HomePageVars := PageVariables{
		Color:       os.Getenv("COLOR"),
		Title:       os.Getenv("TITLE"),
		Description: os.Getenv("DESCRIPTION"),
	}

	t, err := template.ParseFiles("templates/homepage.html")
	if err != nil {
		log.Print("Template parsing error: ", err)
	}
	err = t.Execute(w, HomePageVars)
	if err != nil {
		log.Print("Template executing error: ", err)
	}
}
