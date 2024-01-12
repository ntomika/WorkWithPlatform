package structs

type GetStudentByLogin struct {
	StudentId  string `json:"studentId"`
	UserId     string `json:"userId"`
	SchoolId   string `json:"schoolId"`
	IsActive   bool   `json:"isActive"`
	IsGraduate bool   `json:"isGraduate"`
}

type School21User struct {
	GetStudentByLogin GetStudentByLogin `json:"getStudentByLogin"`
}

type DataUser struct {
	School21User School21User `json:"school21"`
}

type GetCredentialsByLogin struct {
	DataUser DataUser `json:"data"`
}
