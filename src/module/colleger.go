package module

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"tira.com/src/db"
	"tira.com/src/helper"
	"tira.com/src/model"
)

// Get all collegers
func GetCollegers(w http.ResponseWriter, r *http.Request) {
	helper.CheckMethod(w, r, "GET")
	db := db.SetupDBMssql()

	helper.PrintMessage("Getting college...")

	rows, err := db.Query("SELECT * FROM mahasiswa")

	helper.CheckErr(err)
	var collegers []model.Colleger
	// Foreach collegers
	for rows.Next() {
		var id int
		var nim string
		var name string
		var gendre string
		var className string
		var classType string
		var subject string

		err = rows.Scan(&id, &nim, &name, &gendre, &className, &classType, &subject)

		helper.CheckErr(err)

		collegers = append(collegers, model.Colleger{
			NIM:       nim,
			Name:      name,
			ClassName: className,
			ClassType: classType,
			Subject:   subject})
	}

	var response = model.JsonResColleger{Type: "success", Data: collegers}

	json.NewEncoder(w).Encode(response)
}

// Create a colleger
func CreateColleger(w http.ResponseWriter, r *http.Request) {
	helper.CheckMethod(w, r, "POST")
	var p model.Colleger

	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	nim := p.NIM
	name := p.Name
	gendre := p.Gender
	className := p.ClassName
	classType := p.ClassType
	subject := p.Subject

	var response = model.JsonResColleger{}

	if nim == "" || name == "" || gendre == "" || className == "" || classType == "" || subject == "" {
		response = model.JsonResColleger{Type: "error", Message: "You are missing Field, check your request"}
	} else {
		db := db.SetupDBMssql()

		helper.PrintMessage("Inserting book into DB")

		fmt.Println("Inserting new collegers")

		var lastInsertID string
		err := db.QueryRow(
			"INSERT INTO mahasiswa(nim, name, gender, class_name, class_type, subject) VALUES(@p1, @p2, @p3, @p4, @p5, @p6) select SCOPE_IDENTITY();",
			nim,
			name,
			gendre,
			className,
			className,
			subject).Scan(&lastInsertID)

		helper.CheckErr(err)

		message := "The book has been inserted successfully! with id: " + lastInsertID
		response = model.JsonResColleger{Type: "success", Message: message}
	}

	json.NewEncoder(w).Encode(response)
}

// Delete a colleger
func DeleteColleger(w http.ResponseWriter, r *http.Request) {
	helper.CheckMethod(w, r, "DELETE")
	params := mux.Vars(r)

	nim := params["nim"]

	var response = model.JsonResColleger{}

	if nim == "" {
		response = model.JsonResColleger{Type: "error", Message: "You are missing nim parameter."}
	} else {
		db := db.SetupDBMssql()

		helper.PrintMessage("Deleting book from DB")

		_, err := db.Exec("DELETE FROM mahasiswa where nim = @p1", nim)
		helper.CheckErr(err)

		message := "Colleger has been deleted successfully! NIM: " + nim
		response = model.JsonResColleger{Type: "success", Message: message}
	}

	json.NewEncoder(w).Encode(response)
}

// Delete all collegers
func DeleteCollegers(w http.ResponseWriter, r *http.Request) {
	helper.CheckMethod(w, r, "DELETE")
	db := db.SetupDBMssql()

	helper.PrintMessage("Deleting all colleger...")

	_, err := db.Exec("DELETE FROM mahasiswa")
	helper.CheckErr(err)

	helper.PrintMessage("All collegers have been deleted successfully!")

	var response = model.JsonResColleger{Type: "success", Message: "All collegers have been deleted successfully!"}

	json.NewEncoder(w).Encode(response)
}
