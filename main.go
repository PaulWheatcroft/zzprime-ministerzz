package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"
	"reflect"
	"os"
	"github.com/joho/godotenv"

	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
	"net/http"

	"github.com/go-sql-driver/mysql"
)

type PrimeMinister struct {
	ID			int		`json:"id"`
	FirstName	string	`json:"firstName"`
	LastName	string	`json:"lastName"`
	FromDate	string	`json:"fromDate"`
	ToDate		string	`json:"toDate"`
	Terms		int		`json:"terms"`
	Office		string	`json:"office"`
	Party		string	`json:"party"`
}

var db *sql.DB


func setDataBaseConnection() {
	dsn := mysql.Config{
			User:                 os.Getenv("DB_USER"),
			Passwd:               os.Getenv("DB_PASSWORD"),
			Net:                  "tcp",
			Addr:                 os.Getenv("DB_ADDRESS"),
			DBName:               os.Getenv("DB_NAME"),
			AllowNativePasswords: true,
	}

	// Get a database handle
	var err error
	db, err = sql.Open("mysql", dsn.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}
	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}
}

func getAllPrimeMinisters(context *gin.Context) {
	var primeMinisters []PrimeMinister
	// Data Source Name Properties
	
	setDataBaseConnection()
	//Get all prime ministers
	queryResults, err := db.Query("select * from prime_ministers_tbl")
    if err != nil {
        fmt.Printf("The error is: %v", err)
		return
    }
	for queryResults.Next() {
        var pm PrimeMinister
        if err := queryResults.Scan(&pm.ID, &pm.FirstName, &pm.LastName, &pm.FromDate, &pm.ToDate, &pm.Terms, &pm.Office, &pm.Party); err != nil {
            fmt.Printf("Error: %v", err)
			return
        }
        primeMinisters = append(primeMinisters, pm)
    }

	//Put the results in to JSON and pass to the context
	fmt.Println("primeMinisters = ", reflect.TypeOf(primeMinisters))
	fmt.Println(primeMinisters)
	context.IndentedJSON(http.StatusOK, primeMinisters)
	fmt.Println("The last line is", context.Request)

}

func addPrimeMinister(context *gin.Context) {
	var newPriminister PrimeMinister
	setDataBaseConnection()
	err := context.BindJSON(&newPriminister)
	if err != nil {
		return
	}

	context.IndentedJSON(201, newPriminister)

	fmt.Println("I ran add: ", newPriminister.FirstName)

	db.Exec("INSERT INTO prime_ministers_tbl (first_name, last_name, from_date, to_date, terms, office, party) VALUES (?, ?, ?, ?, ?, ?, ?)",
	newPriminister.FirstName, newPriminister.LastName, newPriminister.FromDate, newPriminister.ToDate, newPriminister.Terms, newPriminister.Office, newPriminister.Party)
	
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	fmt.Printf("I am running")
	router := gin.Default()
	router.Use(cors.New(cors.Config{
        AllowOrigins:     []string{"http://localhost:3000"},
        AllowMethods:     []string{"GET", "POST", "PUT", "PATCH"},
        AllowHeaders:     []string{"Origin"},
        ExposeHeaders:    []string{"Content-Length"},
        AllowCredentials: true,
        MaxAge: 12 * time.Hour,
    }))
	router.GET("/", getAllPrimeMinisters)
	router.POST("/add", addPrimeMinister)
	router.Run("localhost:9099")
}
