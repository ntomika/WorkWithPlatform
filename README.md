# Work with platform users

/docs/logins  - список пользователей формата `login@student.21-school.ru`
/docs/classList - список классов

.env          - данные для авторизации

Опции:

* check_coins - работает со списком пользователей /docs/logins и выводит в результатате список переданых логинов с количеством Coins
* new_exam_events - работает со списком классов /docs/classList и заводит мероприяти типа TEST или EXAM

Запуск программы:
go run main.go option_name
