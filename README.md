# Work with platform (GraphQL)

/docs/logins  - список пользователей формата `login@student.21-school.ru`   
/docs/classList - список классов   

.env          - данные для авторизации   
* access_token - токен авторизации формата 'Bearer access_token'
* schoolid - id кампуса (у каждого кампуса свой)
* ContentType - указывает тип и кодировку символов документа (в случае с платформой это 'application/json')


Опции:   
* check_coins - работает со списком пользователей /docs/logins и выводит в результатате список переданых логинов с количеством Coins
* new_exam_events - работает со списком классов /docs/classList и заводит мероприятия типа TEST или EXAM

Запуск программы:   
`go run main.go option_name`
