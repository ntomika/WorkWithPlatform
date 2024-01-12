package conf

type Range struct {
	LeftBorder  int `json:"leftBorder"`
	RightBorder int `json:"rightBorder"`
}
type Level struct {
	LevelCode int   `json:"levelCode"`
	Range     Range `json:"range"`
}
type GetExperiencePublicProfile struct {
	Value            int   `json:"value"`
	Level            Level `json:"level"`
	CookiesCount     int   `json:"cookiesCount"`
	CoinsCount       int   `json:"coinsCount"`
	CodeReviewPoints int   `json:"codeReviewPoints"`
}
type GetStageGroupS21PublicProfile struct {
	WaveId   int    `json:"waveId"`
	WaveName string `json:"waveName"`
	EduForm  string `json:"eduForm"`
}
type GetFeedbackStatisticsAverageScore struct {
	CountFeedback                     int      `json:"countFeedback"`
	GetFeedbackStatisticsAverageScore []string `json:"getFeedbackStatisticsAverageScore"`
}
type GetSchool struct {
	Id        string `json:"id"`
	FullName  string `json:"fullName"`
	ShortName string `json:"shortName"`
	Address   string `json:"address"`
}

type User struct {
	GetSchool GetSchool `json:"getSchool"`
}
type Student struct {
	GetWorkstationByLogin             string                            `json:"getWorkstationByLogin"`
	GetFeedbackStatisticsAverageScore GetFeedbackStatisticsAverageScore `json:"getFeedbackStatisticsAverageScore"`
}
type School21 struct {
	GetAvatarByUserId             string                        `json:"getAvatarByUserId"`
	GetEmailbyUserId              string                        `json:"getEmailbyUserId"`
	GetClassRoomByLogin           string                        `json:"getClassRoomByLogin"`
	GetStageGroupS21PublicProfile GetStageGroupS21PublicProfile `json:"getStageGroupS21PublicProfile"`
	GetExperiencePublicProfile    GetExperiencePublicProfile    `json:"getExperiencePublicProfile"`
}

type Data struct {
	School21 School21 `json:"school21"`
	Student  Student  `json:"Student"`
	User     User     `json:"User"`
}

type PublicProfileGetPersonalInfo struct {
	Data Data `json:"data"`
}
