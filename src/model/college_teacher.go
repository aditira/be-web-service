package model

type CollegeTeacher struct {
	NIM       string `json:"nim"`
	Name      string `json:"name"`
	Gender    string `json:"gender"`
	ClassName string `json:"class_name"`
	ClassType string `json:"class_type"`
	Subject   string `json:"subject"`
}

type JsonResCollegeTeacher struct {
	Type    string           `json:"type"`
	Data    []CollegeTeacher `json:"data"`
	Message string           `json:"message"`
}
