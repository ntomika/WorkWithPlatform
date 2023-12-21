package conf

type GetStudentByLogin struct {
	StudentId  string `json:"studentId"`
	UserId     string `json:"userId"`
	SchoolId   string `json:"schoolId"`
	IsActive   bool   `json:"isActive"`
	IsGraduate bool   `json:"isGraduate"`
}

type School21 struct {
	GetStudentByLogin GetStudentByLogin `json:"getStudentByLogin"`
}

type Data struct {
	School21 School21 `json:"school21"`
}

type GetCredentialsByLogin struct {
	Data Data `json:"data"`
}
