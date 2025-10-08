Приложение подтверждения email
Начинаем собирать API подтверждения Email
- Сделать конфигурацию
- Email
- Password
- Address
- Добавить модуль verify с одним обработчиком и методами:
- POST /send
- GET /verify/{hash}
- Устанавливаем библиотеку
https://github.com/jordan-wright/email