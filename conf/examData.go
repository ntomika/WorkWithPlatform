package conf

type ExamData struct {
	GoalId             string
	Name               string
	ExamType           string
	Location           string
	MaxStudentCount    int
	BeginDate          string
	EndDate            string
	IsVisible          bool
	IsWaitListActive   bool
	StopRegisterDate   string
	StartRegisterDate  string
	StageSubjectGroups string
}
