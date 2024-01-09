package main

import (
	"bufio"
	"conf"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

const URL = "https://edu.21-school.ru/services/graphql"

func openFiles(path string) []string {
	logins := make([]string, 8, 16)

	// Открываем файл для чтения
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// Создаем сканер для чтения файла
	scanner := bufio.NewScanner(file)

	// Читаем файл построчно и записываем в масив логинов
	i := 0
	for scanner.Scan() {
		logins[i] = scanner.Text()
		i++
	}

	// Проверяем наличие ошибок после завершения сканирования
	if err := scanner.Err(); err != nil {
		fmt.Println("Error scanning file:", err)
	}
	return logins
}

func parsSchemaWithIDs(datasByLogin conf.GetCredentialsByLogin, login string) *strings.Reader {

	payloadExample := strings.NewReader("{\"query\":\"query PublicProfileGetPersonalInfo($userId: UUID!, $studentId: UUID!, $login: String!, $schoolId: UUID!) {\\n  school21 {\\n    getAvatarByUserId(userId: $userId)\\n    getStageGroupS21PublicProfile(studentId: $studentId) {\\n      waveId\\n      waveName\\n      eduForm\\n      \\n    }\\n    getExperiencePublicProfile(userId: $userId) {\\n      value\\n      level {\\n        levelCode\\n        range {\\n          leftBorder\\n          rightBorder\\n          \\n        }\\n        \\n      }\\n      cookiesCount\\n      coinsCount\\n      codeReviewPoints\\n      \\n    }\\n    getEmailbyUserId(userId: $userId)\\n    getClassRoomByLogin(login: $login) {\\n      id\\n      number\\n      floor\\n      \\n    }\\n    \\n  }\\n  student {\\n    getWorkstationByLogin(login: $login) {\\n      workstationId\\n      hostName\\n      row\\n      number\\n      \\n    }\\n    getFeedbackStatisticsAverageScore(studentId: $studentId) {\\n      countFeedback\\n      feedbackAverageScore {\\n        categoryCode\\n        categoryName\\n        value\\n        \\n      }\\n      \\n    }\\n    \\n  }\\n  user {\\n    getSchool(schoolId: $schoolId) {\\n      id\\n      fullName\\n      shortName\\n      address\\n      \\n    }\\n    \\n  }\\n}\",\"variables\":{\"userId\":\"25ecd411-4c49-45ed-a422-ac19082bb943\",\"studentId\":\"8311ac3b-1553-47ff-a92f-190c905f8689\",\"schoolId\":\"6bfe3c56-0211-4fe1-9e59-51616caac4dd\",\"login\":\"abathurm@student.21-school.ru\"}}")

	var payloadTmp map[string]interface{}
	if err := json.NewDecoder(payloadExample).Decode(&payloadTmp); err != nil {
		panic(err)
	}
	variables := payloadTmp["variables"].(map[string]interface{})
	variables["login"] = login
	variables["userId"] = datasByLogin.DataUser.School21User.GetStudentByLogin.UserId
	variables["studentId"] = datasByLogin.DataUser.School21User.GetStudentByLogin.StudentId

	tmpP, err := json.Marshal(payloadTmp)
	if err != nil {
		panic(err)
	}

	payloadNew := strings.NewReader(string(tmpP))

	return payloadNew
}

func parsSchemaWithLogin(login string) *strings.Reader {

	payloadExample := strings.NewReader("{\"query\":\"query GetCredentialsByLogin($login: String!) {\\n  school21 {\\n    getStudentByLogin(login: $login) {\\n      studentId\\n      userId\\n      schoolId\\n      isActive\\n      isGraduate\\n      __typename\\n    }\\n    __typename\\n  }\\n}\",\"variables\":{\"login\":\"abathurm@student.21-school.ru\"}}")

	var payloadTmp map[string]interface{}
	if err := json.NewDecoder(payloadExample).Decode(&payloadTmp); err != nil {
		panic(err)
	}
	variables := payloadTmp["variables"].(map[string]interface{})
	variables["login"] = login

	tmpP, err := json.Marshal(payloadTmp)
	if err != nil {
		panic(err)
	}

	payloadNew := strings.NewReader(string(tmpP))

	return payloadNew
}

func handlerRequest(method string, payload *strings.Reader) []byte {
	client := &http.Client{}
	req, err := http.NewRequest(method, URL, payload)
	if err != nil {
		panic(err)
	}

	// Загрузка переменных окружения из файла .env
	if godotenv.Load(".env") != nil {
		panic(err)
	}

	req.Header.Add("userrole", os.Getenv("userrole"))
	req.Header.Add("schoolid", os.Getenv("schoolid"))
	req.Header.Add("Content-Type", os.Getenv("ContentType"))
	req.Header.Add("Authorization", os.Getenv("access_token"))
	req.Header.Add("Cookie", os.Getenv("Cookie"))

	res, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	return body
}

func checkCoins() {
	loginsList := openFiles("docs/logins")

	fmt.Println("\tLogin\t\t\t", "Coins\n", "-------------------------------------")

	for _, login := range loginsList {
		if login == "" {
			break
		}

		// Шаг 1: получить необходимые IDs (userId и studentId) пользователя
		// Разобрать запрос на query и variables, подставить нужный логин в поле login в variables
		payloadForLoginData := parsSchemaWithLogin(login)
		// Отправка запроса с нужным логином
		body := handlerRequest("POST", payloadForLoginData)
		var datasByLogin conf.GetCredentialsByLogin
		// Переложить ответ (body) в структуру datasByLogin
		jsonErr := json.Unmarshal(body, &datasByLogin)
		if jsonErr != nil {
			fmt.Println("Step 1: JSON encoding error:", jsonErr)
		}

		// Шаг 2: подставить полученные ранее userId и studentId в новый запрос на получение информации о пользователе
		// Разобрать запрос на query и variables, подставить нужные IDs (userId, studentId) в variables
		payloadForIDsData := parsSchemaWithIDs(datasByLogin, login)
		// Отправка запроса c нужными userId, studentId и login
		body = handlerRequest("POST", payloadForIDsData)
		var personalInfo conf.PublicProfileGetPersonalInfo
		// Переложить ответ (body) в структуру personalInfo
		jsonErr = json.Unmarshal(body, &personalInfo)
		if jsonErr != nil {
			fmt.Println("Step 2: JSON encoding error:", jsonErr)
		}

		fmt.Println(personalInfo.Data.School21.GetEmailbyUserId, " | ",
			personalInfo.Data.School21.GetExperiencePublicProfile.CoinsCount)
	}
}

func createMapClassWithID() (ClassWithID map[string]string) {
	ClassWithID = map[string]string{
		"22_04_MSK":    "1137037",
		"22_10_MSK":    "2329858",
		"22_11_MSK_11": "2445422",
		"22_11_MSK_8":  "2445424",
		"22_11_MSK_9":  "2445425",
		"22_11_MSK_3":  "2445429",
		"22_11_MSK_4":  "2445430",
		"22_11_MSK_5":  "2445564",
		"22_11_MSK_6":  "2445565",
		"22_11_MSK_7":  "2445566",
		"23_04_MSK_1":  "2885181",
		"23_04_MSK_2":  "2885182",
		"23_04_MSK_3":  "2885183",
		"23_04_MSK_4":  "2885184",
		"23_04_MSK_5":  "2885185",
		"23_04_MSK_6":  "2885186",
		"23_04_MSK_7":  "2885187",
		"23_04_MSK":    "2885282",
		"23_10_MSK":    "2897827",
		"23_12_MSK":    "2899307",
		"Енот":         "1111575",
	}
	return
}

// func newExamEvents() {
// 	classList := openFiles("docs/classList")

// 	allClassWithID := createMapClassWithID()
// 	var examStruct conf.ExamData

// 	fmt.Scanf("Название мероприятия: %s", &examStruct.Name)
// 	if examStruct.Name == "DevOps Exam" {
// 		examStruct.GoalId = "66062"
// 		examStruct.ExamType = "TEST"
// 	}

// 	fmt.Scanf("Место проведения (город кампуса или названия кластера/кластеров): %s", &examStruct.Location)
// 	fmt.Scanf("Количество участников: %d", &examStruct.MaxStudentCount)
// 	fmt.Scanf("beginDate format yyyy-mm-ddThh:mm:00.000Z (Москва минус 3 часа): %s", &examStruct.BeginDate)
// 	fmt.Scanf("endDate yyyy-mm-ddThh:mm:00.000Z (Москва минус 3 часа): %s", &examStruct.EndDate)
// 	fmt.Scanf("Видимость для участников (true - видимое, false - скрытое) : %t", &examStruct.IsVisible)
// 	fmt.Scanf("Лист ожидания (true - включено, false - выключено) : %t", &examStruct.IsWaitListActive)
// 	fmt.Scanf("stopRegisterDate yyyy-mm-ddThh:mm:00.000Z (Москва -3ч) %s", &examStruct.StopRegisterDate)
// 	fmt.Scanf("startRegisterDate yyyy-mm-ddThh:mm:00.000Z (Москва -3ч) %s", &examStruct.StartRegisterDate)


// 	for _, class := range classList {
// 		classID, ok := allClassWithID[class]
// 		if ok {
// 			examStruct.StageSubjectGroups = classID
// 			fmt.Printf(classID)
// 		}
// 	}

// }

func main() {

	if os.Args[1] == "check_coins" {
		checkCoins()
	} else if os.Args[1] == "new_exam_events" {
		// newExamEvents()
		fmt.Println("new_exam_events")
	} else {
		fmt.Println("Incorrect argument.\nSelect one of arguments: check_coins, new_events.\nUse: go run main.go argument")
	}

}
