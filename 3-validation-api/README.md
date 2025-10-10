API подтверждения Email
Начинаем собирать API подтверждения Email

Сделать конфигурацию
Email
Password
Address
Добавить модуль verify с одним обработчиком и методами:
POST /send
GET /verify/{hash}
Устанавливаем библиотеку <https://github.com/jordan-wright/email>
Завершаем сбор API и реализуем:

POST /send принимает email, формирует случайный hash и отправляет письмо на указанный email со ссылкой http://localhost:8081/verify/{hash}
При этом локально сохраняет json с email и hash для валидации
После перехода если hash совпадает с сохранённым вадаёт true, если нет false и удаляет запись.