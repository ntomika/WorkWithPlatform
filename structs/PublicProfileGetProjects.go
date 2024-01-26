package structs

type StudentItem struct {
	GroupName             string `json:"groupName"`
	Name                  string `json:"name"`
	Experience            int    `json:"experience"`
	FinalPercentage       *int   `json:"finalPercentage"`
	GoalID                int    `json:"goalId"`
	GoalStatus            string `json:"goalStatus"`
	AmountAnswers         *int   `json:"amountAnswers"`
	AmountReviewedAnswers *int   `json:"amountReviewedAnswers"`
	TypeName              string `json:"__typename"`
}

type School21Queries struct {
	Data struct {
		School21 struct {
			GetStudentProjectsForPublicProfileByStageGroup []StudentItem `json:"getStudentProjectsForPublicProfileByStageGroup"`
		} `json:"school21"`
	} `json:"data"`
}