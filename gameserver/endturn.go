// endturn.go
package main

import (
	"database/sql"
	"fmt"
	"io/ioutil" //required for io with filesystem
	"log"       //required for log
	"strings"   //required for Split database/sql
	"net/http"  //required for http
	"errors"    //required for errors handling
	_ "github.com/go-sql-driver/mysql"
)

type Account struct {
	id    int
	login string
	pass  string
}

type NewChar struct {
	Charname string
	Race      int
	Gender    int
	Avatar    int
}

func create_new_char(new_char *NewChar) {
	fmt.Println("try to create new char")
	fmt.Println(new_char)
}

func char_create(rw http.ResponseWriter, req *http.Request) {
	var new_char NewChar
	err := decodeJSONBody(rw, req, &new_char)
    if err != nil {
        var mr *malformedRequest
        log.Println(req)
        if errors.As(err, &mr) {
            http.Error(rw, mr.msg, mr.status)
        } else {
            log.Println(err.Error())
            http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
        }
        return
    }
    create_new_char(&new_char)
	rw.Write([]byte("char created"))
	fmt.Fprintf(rw, "char: %+v", new_char)
}


func game_heartbeat(rw http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(rw, "OK")
}


func db_check() {
	config_auth_b, err := ioutil.ReadFile("config_auth.txt")
	if err != nil {
		log.Fatal(err)
	}
	config_auth_s := string(config_auth_b)
	var config_auth_a [5]string
	iterator := 0
	for _, line := range strings.Split(strings.TrimSuffix(config_auth_s, "\n"), "\n") {
		fmt.Println(line)
		config_auth_a[iterator] = line
		iterator += 1
	}
	connect_string := "%user:%pass@tcp(%host:%port)/%db"
	connect_string = strings.Replace(connect_string, "%user", config_auth_a[3], 1)
	connect_string = strings.Replace(connect_string, "%pass", config_auth_a[4], 1)
	connect_string = strings.Replace(connect_string, "%host", config_auth_a[0], 1)
	connect_string = strings.Replace(connect_string, "%port", config_auth_a[1], 1)
	connect_string = strings.Replace(connect_string, "%db", config_auth_a[2], 1)
	fmt.Println(connect_string)
	db, err := sql.Open("mysql", connect_string)

	if err != nil {
		fmt.Println("error connecting to db")
		panic(err.Error())
	}
	// defer the close till after the main function has finished
	// executing
	defer db.Close()

	results, err := db.Query("SELECT * FROM dbstat;")
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	for results.Next() {
		var account Account
		// for each row, scan the result into our tag composite object
		err = results.Scan(&account.id, &account.login)
		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}
		//and then print out the tag's Name attribute
		log.Printf(account.login)
	}
}
func main() {
	//http.Handle("/char_create", http.HandlerFunc(char_create))
	http.HandleFunc("/char_create", char_create)
	http.HandleFunc("/game_heartbeat", game_heartbeat)
	log.Println("Starting server on port 6199")
	log.Fatal(http.ListenAndServe(":6199", nil))
	//fmt.Println(account.login)
}
