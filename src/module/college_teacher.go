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

// Get all colleger teacher
func GetCollegerTeacher(w http.ResponseWriter, r *http.Request) {
	helper.CheckMethod(w, r, "GET")
	db := db.SetupDBMssql()

	helper.PrintMessage("Getting college teacher...")

	rows, err := db.Query("SELECT * FROM dosen")

	helper.CheckErr(err)
	var collegers []model.CollegeTeacher
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

		collegers = append(collegers, model.CollegeTeacher{
			NIM:       nim,
			Name:      name,
			ClassName: className,
			ClassType: classType,
			Subject:   subject})
	}

	var response = model.JsonResCollegeTeacher{Type: "success", Data: collegers}

	json.NewEncoder(w).Encode(response)
}

// Create a colleger teacher
func CreateCollegeTeacher(w http.ResponseWriter, r *http.Request) {
	helper.CheckMethod(w, r, "POST")
	var p model.CollegeTeacher

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

	var response = model.JsonResCollegeTeacher{}

	if nim == "" || name == "" || gendre == "" || className == "" || classType == "" || subject == "" {
		response = model.JsonResCollegeTeacher{Type: "error", Message: "You are missing Field, check your request"}
	} else {
		db := db.SetupDBMssql()

		helper.PrintMessage("Inserting coleger teacher into DB")

		fmt.Println("Inserting new colleger teacher")

		var lastInsertID string
		err := db.QueryRow(
			"INSERT INTO dosen(nim, name, gender, class_name, class_type, subject) VALUES(@p1, @p2, @p3, @p4, @p5, @p6) select SCOPE_IDENTITY();",
			nim,
			name,
			gendre,
			className,
			className,
			subject).Scan(&lastInsertID)

		helper.CheckErr(err)

		message := "The book has been inserted successfully! with id: " + lastInsertID
		response = model.JsonResCollegeTeacher{Type: "success", Message: message}
	}

	json.NewEncoder(w).Encode(response)
}

// Delete a colleger
func DeleteCollegeTeacher(w http.ResponseWriter, r *http.Request) {
	helper.CheckMethod(w, r, "DELETE")
	params := mux.Vars(r)

	nim := params["nim"]

	var response = model.JsonResCollegeTeacher{}

	if nim == "" {
		response = model.JsonResCollegeTeacher{Type: "error", Message: "You are missing nim parameter."}
	} else {
		db := db.SetupDBMssql()

		helper.PrintMessage("Deleting coleger teacher from DB")

		_, err := db.Exec("DELETE FROM dosen where nim = @p1", nim)
		helper.CheckErr(err)

		message := "Colleger has been deleted successfully! NIM: " + nim
		response = model.JsonResCollegeTeacher{Type: "success", Message: message}
	}

	json.NewEncoder(w).Encode(response)
}

// Delete all collegers
func DeleteCollegeTeachers(w http.ResponseWriter, r *http.Request) {
	helper.CheckMethod(w, r, "DELETE")
	db := db.SetupDBMssql()

	helper.PrintMessage("Deleting all colleger teacher...")

	_, err := db.Exec("DELETE FROM dosen")
	helper.CheckErr(err)

	helper.PrintMessage("All colleger teacher have been deleted successfully!")

	var response = model.JsonResCollegeTeacher{Type: "success", Message: "All colleger teacher have been deleted successfully!"}

	json.NewEncoder(w).Encode(response)
}
