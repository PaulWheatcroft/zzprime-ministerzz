package main

import (
	"database/sql"
	"fmt"
	"log"
	"reflect"
	"encoding/json"
	"os"
	"github.com/joho/godotenv"

	"github.com/gin-gonic/gin"
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
	primeMinistersToJson, err := json.Marshal(primeMinisters)
	if err != nil {
		fmt.Println("The JSON conversion went wrong")
	}
	fmt.Println("primeMinisters = ", reflect.TypeOf(primeMinisters))
	fmt.Println(primeMinisters)
	fmt.Println(string(primeMinistersToJson))
	context.IndentedJSON(http.StatusOK, primeMinisters)

}

// func parsePrimeMinister(jsonBuffer []byte) ([]PrimeMinister, error) {
// 	// We create an empty array
//     primeMinister := []PrimeMinister{}

//     // Unmarshal the json into it. this will use the struct tag
//     err := json.Unmarshal(jsonBuffer, &primeMinister)
//     if err != nil {
//         return nil, err
//     }

//     // the array is now filled with users
// 	return primeMinister, nil
// }

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
	router.GET("/", getAllPrimeMinisters)
	router.POST("/add", addPrimeMinister)
	router.Run("localhost:9099")
}
