# Конфиг файл для сервиса service_ config.yaml
Структура:
``` yaml
serverPort: <int: Порт, на котором работает сервис>
consensus: <int: Необходимое количество проверок специалистов для оценки случая> 
reportPeriod: <duration: Время отчетного периода>

passSalt: <string: Соль для хеширования паролей>
signingKey: <string: Ключ подписи JWT токенов>

postgres: <Информация о БД>
  user: <string: Имя пользователя БД>
  password: <string: Пароль пользователя БД>
  host: <string: Хост БД>
  port: <int: Порт БД>
  database: <string: Наименование БД>

rabbitmq: <Информация о RabbitMQ>
  user: <string: Имя пользователя RabbitMQ>
  password: <string: Пароль пользователя RabbitMQ>
  host: <string: Хост RabbitMQ>
  port: <int: Порт RabbitMQ>

directors: <array: Массив директоров>
  - username: <string: Имя директора>
    password: <string: Пароль директора>
```

Пример:
``` yaml
serverPort: 8080
consensus: 2
reportPeriod: 3m

passSalt: "salt"
signingKey: "sign"

postgres:
  user: "user"
  password: "user"
  host: "postgres"
  port: 5432
  database: "traffic_police_db"

rabbitmq:
  user: "guest"
  password: "guest"
  host: "rabbitmq"
  port: 5672

directors:
  - username: "director1"
    password: "director1"
  - username: "director2"
    password: "director2"

```
# Конфиг файл для сервиса уведомлений notification_config.yaml
Структура:
``` yaml
emailSender: <Информация об отправителе сообщений по почте>
  host: <string: Хост отправителя сообщений>
  port: <string: Порт отправителя сообщений>
  username: <string: Имя пользователя отправителя сообщений>
  password: <string: Пароль пользователя>
  subject: <string: Заголовок сообщения о правонарушении>


rabbitmq: <Информация о RabbitMQ>
  user: <string: Имя пользователя RabbitMQ>
  password: <string: Пароль пользователя RabbitMQ>
  host: <string: Хост RabbitMQ>
  port: <string: Порт RabbitMQ>
```

Пример:
``` yaml
emailSender:
  host: "smtp.gmail.com"
  port: 587
  username: "emailsender@gmail.com"
  password: "secret"
  subject: "Информация о правонарушении"


rabbitmq:
  user: "guest"
  password: "guest"
  host: "rabbitmq"
  port: 5672
```
