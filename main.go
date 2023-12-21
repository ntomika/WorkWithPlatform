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

func openFiles() []string {
	logins := make([]string, 8, 16)

	// Открываем файл для чтения
	file, err := os.Open("docs/logins")
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

func main() {

	loginsList := openFiles()

	fmt.Println("\tLogin\t\t\t", "Coins\n", "-------------------------------------")

	for _, login := range loginsList {
		if login == "" {
			break
		}

		payloadForLoginData := parsSchemaWithLogin(login)
		body := handlerRequest("POST", payloadForLoginData)
		var datasByLogin conf.GetCredentialsByLogin
		jsonErr := json.Unmarshal(body, &datasByLogin)
		if jsonErr != nil {
			fmt.Println("Step 1: JSON encoding error:", jsonErr)
		}

		payloadForIDsData := parsSchemaWithIDs(datasByLogin, login)
		body = handlerRequest("POST", payloadForIDsData)
		var personalInfo conf.PublicProfileGetPersonalInfo
		jsonErr = json.Unmarshal(body, &personalInfo)
		if jsonErr != nil {
			fmt.Println("Step 2: JSON encoding error:", jsonErr)
		}

		fmt.Println(personalInfo.Data.School21.GetEmailbyUserId, " | ",
			personalInfo.Data.School21.GetExperiencePublicProfile.CoinsCount)
	}
}
