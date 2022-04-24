package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"io/ioutil"
	"strings"
	"encoding/json"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	sm "backend/stringMatching"
)

type Penyakit struct {
	NamaPenyakit string
	DNA          string
}

type HasilPrediksi struct {
	TanggalPrediksi string
	NamaPasien string
	PenyakitPrediksi string
	TingkatKemiripan int
	Status int
}

func getEnv(key string) string {
	// load .env file
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	return os.Getenv(key)
}

func openDatabase() *sql.DB {
	// Open database connection.
	db, err := sql.Open("mysql", getEnv("DATABASE_USERNAME")+":"+getEnv("DATABASE_PASSWORD")+"@tcp("+getEnv("DATABASE_PORT")+")/"+getEnv("DATABASE_NAME"))

	// if there is an error opening the connection, handle it
	if err != nil {
		panic(err.Error())
	}
	return db
}

func setupResponse(w *http.ResponseWriter, req *http.Request) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
    (*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
    (*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}


func getDetailPrediction(res http.ResponseWriter, req *http.Request) {
	setupResponse(&res, req)
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Fatalln(err)
	}
	
	//Convert the body to type string
	string_body := string(body)
	if strings.Split(string_body, ":")[0] != "" {
		data := strings.Split(string_body, ":")[1] 
		input := ""
		two_arguments := false
		for i := 1; i < len(data) - 2; i++ {
			input += string(data[i])
			if string(data[i]) == " " {
				two_arguments = true
			}
		}

		db := openDatabase()
		result := []HasilPrediksi{}
		if two_arguments {
			date := strings.Split(input, " ")[0]
			disease := strings.Split(input, " ")[1]

			// Db query for hasilprediksi table
			db_result, err := db.Query("SELECT * FROM hasilprediksi WHERE TanggalPrediksi = '" + date + "' AND PenyakitPrediksi = '" + disease + "'")
			if err != nil {
				panic(err.Error())
			}

			for db_result.Next() {
				var hasil HasilPrediksi

				// Get hasil for each row
				err = db_result.Scan(&hasil.TanggalPrediksi, &hasil.NamaPasien, &hasil.PenyakitPrediksi, &hasil.TingkatKemiripan, &hasil.Status)
				if err != nil {
					panic(err.Error()) // proper error handling instead of panic in your app
				}
				// Append hasil to result
				result = append(result, hasil)
			}
		} else {

		}
		// Convert result to []byte
		marshal, err := json.Marshal(result)
		if err != nil {
			fmt.Println(err)
		}
		// Send to frontend
		res.Write(marshal)
	}
}

func main() {
	sm.BoyerMoore("a pattern matching algorithm", "rithm")
	sm.BoyerMoore("abacaabadcabacabaabb", "abacab")

	var check bool = sm.Regex("AGTC")
	if check {
		fmt.Println("Benar")
	} else {
		fmt.Println("Salah")
	}
	// Server
	http.HandleFunc(getEnv("BASE_PORT")+"/get-detailprediction", getDetailPrediction)

	fmt.Println("Starting server at port " + getEnv("BACKEND_PORT"))
	if err := http.ListenAndServe(":"+getEnv("BACKEND_PORT"), nil); err != nil {
		log.Fatal(err)
	}
}