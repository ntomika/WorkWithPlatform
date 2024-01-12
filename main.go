package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"structs"
	"time"

	"github.com/joho/godotenv"
)

const URL = "https://edu.21-school.ru/services/graphql"

func readFiles(path string) []string {
	datas := make([]string, 8, 16)
	// Открываем файл для чтения
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// Создаем сканер для чтения файла
	scanner := bufio.NewScanner(file)

	// Читаем файл построчно и записываем в масив данных
	for i := 0; scanner.Scan(); i++ {
		datas[i] = scanner.Text()
	}

	// Проверяем наличие ошибок после завершения сканирования
	if err := scanner.Err(); err != nil {
		fmt.Println("Error scanning file:", err)
	}
	return datas
}

func requestNewExam(examStruct structs.ExamData) *strings.Reader {
	payloadExample := strings.NewReader("{\"query\":\"mutation EventsCreateExam($exam: ExamInput!) {\\n  businessAdmin {\\n    createExam(exam: $exam) {\\n      examId\\n      __typename\\n    }\\n    __typename\\n  }\\n}\",\"variables\":{\"exam\":{\"goalId\":\"66062\",\"name\":\"DevOps Exam\",\"examType\":\"TEST\",\"location\":\"Moscow\",\"maxStudentCount\":20,\"beginDate\":\"2024-01-10T17:00:00.000Z\",\"endDate\":\"2024-01-10T18:00:00.000Z\",\"isVisible\":false,\"isWaitListActive\":false,\"stopRegisterDate\":\"2024-01-10T16:59:00.000Z\",\"startRegisterDate\":\"2024-01-10T16:40:00.000Z\",\"stageSubjectGroups\":[\"0\"]}}}")

	var payloadTmp map[string]interface{}
	if err := json.NewDecoder(payloadExample).Decode(&payloadTmp); err != nil {
		panic(err)
	}
	variables := payloadTmp["variables"].(map[string]interface{})
	exam := variables["exam"].(map[string]interface{})
	exam["goalId"] = examStruct.GoalId
	exam["name"] = examStruct.Name
	exam["examType"] = examStruct.ExamType
	exam["location"] = examStruct.Location
	exam["maxStudentCount"] = examStruct.MaxStudentCount
	exam["beginDate"] = examStruct.BeginDate
	exam["endDate"] = examStruct.EndDate
	exam["stopRegisterDate"] = examStruct.StopRegisterDate
	exam["startRegisterDate"] = examStruct.StartRegisterDate
	exam["stageSubjectGroups"] = examStruct.StageSubjectGroups[0]

	tmpP, err := json.Marshal(payloadTmp)
	if err != nil {
		panic(err)
	}

	payloadNew := strings.NewReader(string(tmpP))

	return payloadNew
}

func requestPublicProfileGetPersonalInfo(datasByLogin structs.GetCredentialsByLogin, login string) *strings.Reader {

	payloadExample := strings.NewReader("{\"query\":\"query PublicProfileGetPersonalInfo($userId: UUID!, $studentId: UUID!, $login: String!, $schoolId: UUID!) {\\n  school21 {\\n    getAvatarByUserId(userId: $userId)\\n    getStageGroupS21PublicProfile(studentId: $studentId) {\\n      waveId\\n      waveName\\n      eduForm\\n      \\n    }\\n    getExperiencePublicProfile(userId: $userId) {\\n      value\\n      level {\\n        levelCode\\n        range {\\n          leftBorder\\n          rightBorder\\n          \\n        }\\n        \\n      }\\n      cookiesCount\\n      coinsCount\\n      codeReviewPoints\\n      \\n    }\\n    getEmailbyUserId(userId: $userId)\\n    getClassRoomByLogin(login: $login) {\\n      id\\n      number\\n      floor\\n      \\n    }\\n    \\n  }\\n  student {\\n    getWorkstationByLogin(login: $login) {\\n      workstationId\\n      hostName\\n      row\\n      number\\n      \\n    }\\n    getFeedbackStatisticsAverageScore(studentId: $studentId) {\\n      countFeedback\\n      feedbackAverageScore {\\n        categoryCode\\n        categoryName\\n        value\\n        \\n      }\\n      \\n    }\\n    \\n  }\\n  user {\\n    getSchool(schoolId: $schoolId) {\\n      id\\n      fullName\\n      shortName\\n      address\\n      \\n    }\\n    \\n  }\\n}\",\"variables\":{\"userId\":\"\",\"studentId\":\"\",\"schoolId\":\"6bfe3c56-0211-4fe1-9e59-51616caac4dd\",\"login\":\"\"}}")

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

func requestGetCredentialsByLogin(login string) *strings.Reader {

	payloadExample := strings.NewReader("{\"query\":\"query GetCredentialsByLogin($login: String!) {\\n  school21 {\\n    getStudentByLogin(login: $login) {\\n      studentId\\n      userId\\n      schoolId\\n      isActive\\n      isGraduate\\n      __typename\\n    }\\n    __typename\\n  }\\n}\",\"variables\":{\"login\":\"\"}}")

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

func getCoinsList() {
	loginsList := readFiles("docs/logins")

	fmt.Println("\tLogin\t\t\t", "Coins\n", "-------------------------------------")

	for _, login := range loginsList {
		if login == "" {
			break
		}

		// Шаг 1: получить необходимые IDs (userId и studentId) пользователя
		// Разобрать запрос на query и variables, подставить нужный логин в поле login в variables
		payloadForLoginData := requestGetCredentialsByLogin(login)
		// Отправка запроса с нужным логином
		body := handlerRequest("POST", payloadForLoginData)
		var datasByLogin structs.GetCredentialsByLogin
		// Переложить ответ (body) в структуру datasByLogin
		jsonErr := json.Unmarshal(body, &datasByLogin)
		if jsonErr != nil {
			fmt.Println("Step 1: JSON encoding error:", jsonErr)
		}

		// Шаг 2: подставить полученные ранее userId и studentId в новый запрос на получение информации о пользователе
		// Разобрать запрос на query и variables, подставить нужные IDs (userId, studentId) в variables
		payloadForIDsData := requestPublicProfileGetPersonalInfo(datasByLogin, login)
		// Отправка запроса c нужными userId, studentId и login
		body = handlerRequest("POST", payloadForIDsData)
		var personalInfo structs.PublicProfileGetPersonalInfo
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
		"ТестДевопс":   "2896773",
	}
	return
}

func timeToEvent(dataAndTime string) (newTime string) {

	// Парсим строку в time.Time с использованием указанного формата
	parsedTime, err := time.Parse("2006-01-02 15:04:05", dataAndTime)
	if err != nil {
		fmt.Println("Ошибка парсинга времени:", err)
		return
	}

	// Вычитаем 3 часа
	utcTime := parsedTime.Add(-3 * time.Hour)

	// Форматируем time.Time обратно в строку с новым форматом
	newTime = utcTime.Format("2006-01-02T15:04:05.000Z")

	return
}

func timeToRegister(dataAndTimeStartEvent string) (timeToStart, timeToEnd string) {
	tmpStartTime := time.Now()
	utcTime := tmpStartTime.Add(-2 * time.Hour)
	timeToStart = utcTime.Format("2006-01-02T15:04:05.000Z")

	// Парсим строку в time.Time с использованием указанного формата
	tmpEndTime, err := time.Parse("2006-01-02 15:04:05", dataAndTimeStartEvent)
	if err != nil {
		fmt.Println("Ошибка парсинга времени:", err)
		return
	}

	utcTime = tmpEndTime.Add(-3*time.Hour + (-1)*time.Minute)
	timeToEnd = utcTime.Format("2006-01-02T15:04:05.000Z")

	return
}

func checkInfo(examStruct structs.ExamData) (correct int) {
	fmt.Println("\nПроверка введеных данных:\n----------------------------")

	fmt.Println("Название мероприятия:\t\t\t", examStruct.Name)
	fmt.Println("Тип экзамена:\t\t\t\t", examStruct.ExamType)
	fmt.Println("Место:\t\t\t\t\t", examStruct.Location)
	fmt.Println("Максимальное количество участников:\t", examStruct.MaxStudentCount)
	fmt.Println("Дата и время мероприятия (время UTC):\t", examStruct.BeginDate, "-", examStruct.EndDate)
	fmt.Println("Дата и время регистрации (время UTC):\t", examStruct.StartRegisterDate, "-", examStruct.StopRegisterDate)

	var tmpFlag string
	fmt.Println("Все верно? (y/n)")
	fmt.Scan(&tmpFlag)
	if tmpFlag == "y" {
		correct = 1
	}
	return
}

func getFullStruct(examStruct *structs.ExamData) {
	var NumExam int
	var tmpData, tmpTime string

	fmt.Println("Для какого экзамена заводим мероприятия (указать цифру):\n1. DevOps Exam\n2. CPP Exam\n3. CPPE-T\n4. Core final exam\n5. Core final testing")

beginingForSwitch:
	for {
		fmt.Scan(&NumExam)
		switch NumExam {
		case 1:
			examStruct.Name = "DevOps Exam"
			examStruct.GoalId = "66062"
			examStruct.ExamType = "TEST"
		case 2:
			examStruct.Name = "CPP Exam"
			examStruct.GoalId = "57572"
			examStruct.ExamType = "EXAM"
		case 3:
			examStruct.Name = "CPPE-T"
			examStruct.GoalId = "57573"
			examStruct.ExamType = "TEST"
		case 4:
			examStruct.Name = "Core final exam"
			examStruct.GoalId = "59889"
			examStruct.ExamType = "EXAM"
		case 5:
			examStruct.Name = "Core final testing"
			examStruct.GoalId = "59888"
			examStruct.ExamType = "TEST"
		default:
			fmt.Println("Нужно ввести только цифру из предложенного списка")
			continue beginingForSwitch
		}
		break
	}

	fmt.Print("Место проведения (город кампуса или названия кластера/кластеров):\n")
	fmt.Scan(&examStruct.Location)

	fmt.Println("Максимальное количество участников:")
	fmt.Scan(&examStruct.MaxStudentCount)

	fmt.Println("Дата проведения мероприятия (формат: 'yyyy-mm-dd'):")
	fmt.Scan(&tmpData)
	fmt.Println("Время НАЧАЛА мероприятия (формат: 'hh:mm:ss'). Указывать время по МСК:")
	fmt.Scan(&tmpTime)
	examStruct.BeginDate = timeToEvent(tmpData + " " + tmpTime)

	examStruct.StartRegisterDate, examStruct.StopRegisterDate = timeToRegister(tmpData + " " + tmpTime)

	fmt.Println("Время ОКОНЧАНИЯ мероприятия (формат: 'hh:mm:ss'). Указывать время по МСК):")
	fmt.Scan(&tmpTime)
	examStruct.EndDate = timeToEvent(tmpData + " " + tmpTime)
}

func newExamEvents() {
	classList := readFiles("docs/classList")

	fmt.Println("Список классов, для которых будут заведены мероприятия:")
	for _, className := range classList {
		if className == "" {
			break
		}
		fmt.Println(className)
	}
	fmt.Println("Список корректен? (y/n)")
	var tmpFlag string
	fmt.Scan(&tmpFlag)
	if tmpFlag == "n" {
		fmt.Println("Проверьте список классов в docs/classList и запустите программу снова")
		return
	} else {
		allClassWithID := createMapClassWithID()
		examStruct := structs.ExamData{
			IsWaitListActive: true,
			IsVisible:        true,
		}

		getFullStruct(&examStruct)

		if checkInfo(examStruct) == 0 {
			fmt.Println("Запустите программу снова")
			return
		}

		for _, class := range classList {
			if class == "" {
				break
			}
			classID, ok := allClassWithID[class]
			if ok {
				fmt.Println("\nКласс", class)
				examStruct.StageSubjectGroups[0] = classID

				// Разобрать запрос на query и variables, подставить все значения из examStruct в variables
				neededRequest := requestNewExam(examStruct)

				// Отправка запроса с нужными variables
				body := handlerRequest("POST", neededRequest)

				fmt.Println("Ответ:\n", string(body))
			} else {
				fmt.Println("The class named '", class, "' was not found")
			}
		}
	}

}

func main() {
	if len(os.Args) == 1 {
		fmt.Println("Select one of arguments: get_coins_list, create_exam_events.\nUse: go run main.go argument")
		return
	}

	if os.Args[1] == "get_coins_list" {
		getCoinsList()
	} else if os.Args[1] == "create_exam_events" {
		newExamEvents()
	} else {
		fmt.Println("Incorrect argument.\nSelect one of arguments: get_coins_list, create_exam_events.\nUse: go run main.go argument")
	}

}
