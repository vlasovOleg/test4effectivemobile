# Тестовое задание для Effective Mobile на позицию "GO разработчик"

## Задание:

Реализовать сервис, который будет получать по апи ФИО, из открытых апи обогащать ответ наиболее вероятными возрастом, полом и национальностью и сохранять данные в БД. По запросу выдавать инфу о найденных людях. Необходимо реализовать следующее

1. Выставить rest методы

   - Для получения данных с различными фильтрами и пагинацией
   - Для удаления по идентификатору
   - Для изменения сущности
   - 4 Для добавления новых людей в формате

     ```json
     {
     "name": "Dmitriy",
     "surname": "Ushakov",
     "patronymic": "Vasilevich" // необязательно
     }
     ```
2. Корректное сообщение обогатить

   * Возрастом - https://api.agify.io/?name=Dmitriy
   * Полом - https://api.genderize.io/?name=Dmitriy
   * Национальностью - https://api.nationalize.io/?name=Dmitriy
3. Обогащенное сообщение положить в БД postgres (структура БД должна быть создана путем миграций)
4. Покрыть код debug- и info-логами
5. Вынести конфигурационные данные в .env

## Реализованное API

#### Создание person

```http
  Запрос
  POST /person
```

```json
{
    "name": "Dmitriy",
    "surname": "Ushakov",
    "patronymic": "Vasilevich"
}
```

```json
Ответ
Status сode: 201
{
    "status": "OK",
    "id": 39
}
```

#### Получение person

```http
  Запрос
  GET /person/{ID}
```

```json
Ответ
Status сode: 200
{
    "status": "OK",
    "person": {
        "id": 39,
        "name": "Dmitriy",
        "surname": "Ushakov",
        "patronymic": "Vasilevich",
        "age": 43,
        "gender": "male",
        "country": "RU"
    }
}
```

```http
  Запрос
  GET /person/?
```

| Параметр | Описание                        |
| :--------------- | :-------------------------------------- |
| agemin           | Минимальный возраст   |
| agemax           | Максимальный возраст |
| gender           | Пол                                  |
| country          | Страна                            |
| offset           | Смещение                        |
| limit            | Лимит                              |

```json
Ответ
Status сode: 200
{
    "status": "OK",
    "offset": 2,
    "count": 2,
    "persons": [
        {
            "id": 25,
            "name": "Nikita",
            "surname": "Galkin",
            "patronymic": "Matveyevich",
            "age": 44,
            "gender": "male",
            "country": "RU"
        },
        {
            "id": 39,
            "name": "Dmitriy",
            "surname": "Ushakov",
            "patronymic": "Vasilevich",
            "age": 43,
            "gender": "male",
            "country": "RU"
        }
    ]
}
```

#### Обновление person

```http
  Запрос
  PATCH /person/{ID}
```

```json
{
    "name": "Dmitriy",
    "surname": "Ushakov",
    "patronymic": "Vasilevich",
    "age": 1,
    "Gender": "male",
    "Country": "US"
}
```

```json
Ответ
Status сode: 200
{
    "status": "OK"
}
```

#### Удаление person

```http
  Запрос
  DELETE /person/{ID}
```

```json
Ответ
Status сode: 200
{
    "status": "OK",
}
```

#### Возможные ошибки

```json
Status сode: 404
{
    "status": "Error",
    "error": "not found"
}
```

```json
Status сode: 400
{
    "status": "Error",
    "error": "Bad request",
    "message": "*"
}
```

```json
Status сode: 500
{
    "status": "Error",
    "error": "Internal server error"
}
```
