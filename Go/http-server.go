package main

import (
	"crypto/subtle"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"

	_ "github.com/go-sql-driver/mysql"

	"github.com/gorilla/handlers"
)

const (
	connPort       = "8080"
	adminUser      = "admin"
	adminPassword  = "admin"
	authRealm      = "Пжлст введите свой логин и пароль"
	driverName     = "mysql"
	dataSourceName = "admin:admin@tcp(db:3306)/mydb"
)

var db *sql.DB
var connectionError error

func init() {
	db, connectionError = sql.Open(driverName, dataSourceName)
	if connectionError != nil {
		log.Fatal("Ошибка подключения к базе данных :: ", connectionError)
	}
}

type Emloyee struct {
	ID   int    `json:"uid"`
	Name string `json:"name"`
}

var getHome = basicAuthMid(
	func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("REST-сервис базы данных Employee"))
	}, authRealm)

var getCurrentDb = basicAuthMid(
	func(w http.ResponseWriter, r *http.Request) {

		rows, err := db.Query("SELECT DATABASE() as db")
		if err != nil {
			log.Print("Ошибка выполнения запроса :: ", err)
			return
		}

		var db string
		for rows.Next() {
			rows.Scan(&db)
		}
		fmt.Fprintf(w, "Текущая база данных :: %s", db)
	}, authRealm)

var createEmployee = basicAuthMid(
	func(w http.ResponseWriter, r *http.Request) {

		vals := r.URL.Query()
		name, ok := vals["name"]

		if ok {
			log.Print("Добавление записи в базу данных для : ", name[0])
			stmt, err := db.Prepare("INSERT employee SET name=?")
			if err != nil {
				log.Print("Ошибка при подготовке запроса :: ", err)
				return
			}
			result, err := stmt.Exec(name[0])
			if err != nil {
				log.Print("Ошибка выполнения запроса :: ", err)
				return
			}

			id, err := result.LastInsertId()
			fmt.Fprintf(w, "id последней записи :: %s", strconv.FormatInt(id, 10))
		} else {
			fmt.Fprintf(w, "Ошибка query-параметра")
		}

	}, authRealm)

var readEmployees = basicAuthMid(
	func(w http.ResponseWriter, r *http.Request) {

		log.Print("Чтение записи из базы данных")
		rows, err := db.Query("SELECT * FROM employee")
		if err != nil {
			log.Print("Ошибка выполнения запроса :: ", err)
			return
		}
		employees := []Emloyee{}
		for rows.Next() {
			var uid int
			var name string
			err = rows.Scan(&uid, &name)
			employee := Emloyee{ID: uid, Name: name}
			employees = append(employees, employee)
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(employees)

	}, authRealm)

var updateEmployee = basicAuthMid(
	func(w http.ResponseWriter, r *http.Request) {

		vars := mux.Vars(r)
		id := vars["id"]
		vals := r.URL.Query()
		name, ok := vals["name"]
		if ok {
			log.Print("Обновление записи базы данных для id :: ", id)
			stmt, err := db.Prepare("UPDATE employee SET name=? WHERE uid=?")
			if err != nil {
				log.Print("Ошибка при подготовке запроса :: ", err)
				return
			}
			result, err := stmt.Exec(name[0], id)
			if err != nil {
				log.Print("Ошибка выполнения запроса :: ", err)
				return
			}
			rowsAffected, err := result.RowsAffected()
			fmt.Fprintf(w, "Число обновленных строк в базе данных :: %d", rowsAffected)
		} else {
			fmt.Fprintf(w, "Ошибка query-параметра")
		}

	}, authRealm)

var deleteEmployee = basicAuthMid(
	func(w http.ResponseWriter, r *http.Request) {
		vals := r.URL.Query()
		name, ok := vals["name"]
		if ok {
			log.Print("Удаление записи в базе данных для name :: ", name[0])
			stmt, err := db.Prepare("DELETE FROM employee WHERE name=?")
			if err != nil {
				log.Print("Ошибка при подготовке запроса :: ", err)
				return
			}
			result, err := stmt.Exec(name[0])
			if err != nil {
				log.Print("Ошибка выполнения запроса :: ", err)
				return
			}
			rowsAffected, err := result.RowsAffected()
			fmt.Fprintf(w, "Число удаленных записей в базе данных :: %d", rowsAffected)

		} else {
			fmt.Fprintf(w, "Ошибка query-параметра")
		}

	}, authRealm)

func basicAuthMid(handler http.HandlerFunc, realm string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, pass, ok := r.BasicAuth()
		if !ok || subtle.ConstantTimeCompare([]byte(user),
			[]byte(adminUser)) != 1 || subtle.ConstantTimeCompare([]byte(pass),
			[]byte(adminPassword)) != 1 {
			w.Header().Set("WWW-Authenticate", `Basic realm="`+realm+`"`)
			w.WriteHeader(401)
			w.Write([]byte("Неаутентифицированный доступ к приложению.\n"))
			return
		}
		handler(w, r)
	}
}

func main() {
	router := mux.NewRouter()
	logFile, errfile := os.OpenFile("http-server.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)

	router.Handle("/", handlers.LoggingHandler(logFile, getHome)).Methods("GET")
	router.Handle("/database", handlers.LoggingHandler(logFile, getCurrentDb)).Methods("GET")
	router.Handle("/employee", handlers.LoggingHandler(logFile, createEmployee)).Methods("POST")
	router.Handle("/employee", handlers.LoggingHandler(logFile, readEmployees)).Methods("GET")
	router.Handle("/employee", handlers.LoggingHandler(logFile, deleteEmployee)).Methods("DELETE")
	router.Handle("/employee/{id}", handlers.CombinedLoggingHandler(logFile, updateEmployee)).Methods("PUT")

	defer db.Close()

	errserve := http.ListenAndServe(":"+connPort, handlers.CompressHandler(router))

	if errserve != nil {
		log.Fatal("Ошибка запуска http-сервера : ", errserve)
		return
	}

	if errfile != nil {
		log.Fatal("Ошибка создания открытия лог-файла : ", errfile)
		return
	}
}
