package main

import (
	"database/sql"
	"html/template"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type PageVariables struct {
	Color       string
	Title       string
	Description string
}

type APIPageVariables struct {
	Color         string
	Title         string
	Endpoint      string
	StatusMessage string
}

type DBConnectionPageVars struct {
	Color       string
	Title       string
	DB_Host     string
	DB_Database string
	DB_User     string
	Status      bool
	Message     string
}

func main() {
	http.HandleFunc("/", HomePage)
	http.HandleFunc("/mysql", MySQL)
	http.HandleFunc("/api", APIStatus)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func APIStatus(w http.ResponseWriter, r *http.Request) {
	endpoint, isSet := os.LookupEnv("API_ENDPOINT")
	if !isSet {
		endpoint = "https://httpbin.org/get"
	}
	client := &http.Client{Timeout: 10 * time.Second}
	response, err := client.Get(endpoint)
	var message string

	if err, ok := err.(net.Error); ok && err.Timeout() {
		// A timeout error occurred
		message = "Timeout:" + err.Error()
	} else if err != nil {
		log.Fatal(err)
		// This was an error, but not a timeout
	} else {
		message = "Success- " + strconv.Itoa(response.StatusCode)

	}

	log.Print(response)

	APIPageVars := APIPageVariables{
		Endpoint:      endpoint,
		StatusMessage: message,
	}

	t, err := template.ParseFiles("templates/api.html")
	if err != nil {
		log.Print("Template parsing error: ", err)
	}
	err = t.Execute(w, APIPageVars)
	if err != nil {
		log.Print("Template executing error: ", err)
	}
}

func MySQL(w http.ResponseWriter, r *http.Request) {
	DB_Host := os.Getenv("DB_Host")
	DB_Database := os.Getenv("DB_Database")
	DB_User := os.Getenv("DB_User")
	DB_Password := os.Getenv("DB_Password")
	var status bool
	var statusMessage string

	mysql_string := DB_User + ":" + DB_Password + "@tcp(" + DB_Host + ":3306)/" + DB_Database + "?timeout=5s"
	db, err := sql.Open("mysql", mysql_string)
	if err != nil {
		log.Printf("Error %s when opening DB\n", err)
		status = false
	}
	defer db.Close()

	err = db.Ping()

	if err != nil {
		statusMessage = err.Error()
		status = false
	} else {
		statusMessage = "Connected"
		status = true
	}

	MySQLPageVars := DBConnectionPageVars{
		Color:       os.Getenv("COLOR"),
		Title:       "MySQL DB Connection",
		DB_Host:     DB_Host,
		DB_Database: DB_Database,
		DB_User:     DB_User,
		Status:      status,
		Message:     statusMessage,
	}

	t, err := template.ParseFiles("templates/mysql.html")
	if err != nil {
		log.Print("Template parsing error: ", err)
	}
	err = t.Execute(w, MySQLPageVars)
	if err != nil {
		log.Print("Template executing error: ", err)
	}
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
