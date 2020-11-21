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
	Portrait  int
}

func create_new_char(new_char *NewChar) bool {
	fmt.Println("try to create new char")
	//checks
	
	//check race
	var db *sql.DB
	db = create_db_pool()
	defer db.Close()
	results, err := db.Query("SELECT id_race FROM races;")
	var races []int
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	for results.Next() {
		var race_id int
		// for each row, scan the result into our tag composite object
		//err = results.Scan(&account.id, &account.login)
		err = results.Scan(&race_id)
		races = append(races, race_id)
		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}
	}
	var race_check bool
	printSlice(races)
	race_check = int_slice_contains(races, new_char.Race)
	fmt.Println(race_check)
	if race_check == false {
		return false
	}
	
	//check gender
	var gender_check bool
	var genders []int
	genders = append(genders, 0)
	genders = append(genders, 1)
	gender_check = int_slice_contains(genders, new_char.Gender)
	if gender_check == false {
		return false
	}
	fmt.Println(gender_check)
	fmt.Println(new_char)
	
	//check portrait
	db = create_db_pool()
	defer db.Close()
	port_query := `SELECT id_portrait FROM portraits 
                   WHERE id_race = ? AND 
                   gender = ? AND
                   allowed_to_use = 1`
	results, err = db.Query(port_query, new_char.Race, new_char.Gender)
	var portraits []int
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	for results.Next() {
		var portrait_id int
		err = results.Scan(&portrait_id)
		portraits = append(portraits, portrait_id)
		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}
	}
	var portrait_check bool
	printSlice(portraits)
	portrait_check = int_slice_contains(portraits, new_char.Portrait)
	fmt.Println(portrait_check)
	if portrait_check == false {
		return false
	}
	//check settings, whether it is allow to use a custom nickname or not
	//check uniq of the nickname
	return true
}

//HELPERS
func printSlice(s []int) {
	fmt.Printf("len=%d cap=%d %v\n", len(s), cap(s), s)
}
func int_slice_contains(s []int, e int) bool {
    for _, a := range s {
        if a == e {
            return true
        }
    }
    return false
}

func char_create(rw http.ResponseWriter, req *http.Request) {
	var new_char NewChar
	var char_create_res bool
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
	char_create_res = create_new_char(&new_char)
	//rw.Write([]byte("char creation res %v"))
	fmt.Fprintf(rw, "char creation res %v\n", char_create_res)
	fmt.Fprintf(rw, "char: %+v", new_char)
}

func game_heartbeat(rw http.ResponseWriter, req *http.Request) {
	db_check()
	fmt.Fprintf(rw, "OK")
}

func create_db_pool() *sql.DB {
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
	return db
}

func db_check() {
	var db *sql.DB
	db = create_db_pool()
	defer db.Close()
	results, err := db.Query("SELECT * FROM dbstat;")
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	for results.Next() {
		//var account Account
		var db_con string
		// for each row, scan the result into our tag composite object
		//err = results.Scan(&account.id, &account.login)
		err = results.Scan(&db_con)
		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}
		//and then print out the tag's Name attribute
		//log.Printf(account.login)
		log.Printf(db_con)
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
