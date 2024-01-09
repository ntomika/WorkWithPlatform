package conf

type ExamData struct {
	goalId             string
	name               string
	examType           string
	location           string
	maxStudentCount    int
	beginDate          string
	endDate            string
	isVisible          bool
	isWaitListActive   bool
	stopRegisterDate   string
	startRegisterDate  string
	stageSubjectGroups string
}