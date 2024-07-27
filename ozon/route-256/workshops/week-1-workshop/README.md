# Workshop #1

# Задача

Создать сервис отзывов на товары, который будет иметь два HTTP API:

1. POST products/{id}/reviews - Создать отзыв.
2. GET products/{id}/reviews - Получить все отзывы по товару

## Создание отзыва

Все параметры обязательны.

```http request
POST products/{id}/reviews
{
  "sku": 100,
  "comment": "Отлично",
  "user_id": 20
}
```

## Получение отзывов

```http request
GET products/{id}/reviews
{
  "reviews": [
     {
        "sku": 100,
        "comment": "Отлично",
        "user_id": 20
      }
   ]
}
```

# Полезные ссылки
+ Структура проекта - https://github.com/golang-standards/project-layout/blob/master/README_ru.md  
+ Пример имплементации чистой архитектуры - https://github.com/bxcodec/go-clean-arch?tab=readme-ov-file