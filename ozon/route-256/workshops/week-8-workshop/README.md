# План workshop-а по system-design

## План занятия

- Изучаем теорию по презентации
  - повторение прошлого материала
  - event sourcing
  - saga
  - порождение идентификаторов
  - решардинг
- Разбираемся с нюансами реализации кэша в go
  - window slice
  - ring slice
  - error group and concurrent slice access
  - map + mutex
  - map + rw mutex
  - sharded map
  - sync map
  - read only map (осторожно)
  - error group with slice
- Разбираемся с шардированием
  - рассматриваем dev окружение
  - реализуем http-api
  - реализуем shard-manager
  - реализуем репозиторий
  - рассматриваем цепочку для создания заказа
  - рассматриваем цепочку получения заказа по order_id
  - рассматриваем цепочку получения списка заказов по user_id (ключу шардирования)
  - рассматриваем цепочку получения списка заказов по order_id (не ключу шардирования)
- Практика проектирования ленты новостей
  - event storming
  - проектирование архитектуры
- Свободные вопросы

Ссылки для самостоятельного изучения
- system design
  - [Владимир Маслов — System Design. Как построить распределенную систему и пройти собеседование](https://www.youtube.com/watch?v=popkBBjbAv8)
  - Инфографика [System design resources](https://imgur.com/dZEbZaT)
  - [Как проходят архитектурные секции собеседования в Яндексе: практика дизайна распределённых систем](https://habr.com/ru/companies/yandex/articles/564132/)
  - [Системный дизайн Twitter](https://proselyte.net/tutorials/system-design/twitter/)
  - [System Design](https://github.com/CrazySquirrel/Outtalent/tree/master/System%20Design)
  - [system-design-interview](https://github.com/checkcheckzz/system-design-interview)
  - [Lamport timestamp](https://en.wikipedia.org/wiki/Lamport_timestamp)
  - [Paxos](https://en.wikipedia.org/wiki/Paxos_(computer_science))
  - [The C10K problem](http://www.kegel.com/c10k.html)
  - [Теория отказоустойчивых распределенных систем. Липовский](https://www.youtube.com/watch?v=3yjy19cUWIM)
- работа со временем
  - [Тяжёлое бремя времени. Доклад Яндекса о типичных ошибках в работе со временем](https://habr.com/ru/companies/yandex/articles/463203/)
  - [Заблуждения программистов относительно времени](https://habr.com/ru/articles/146109/)
- golang
  - [50 оттенков Go: ловушки, подводные камни и распространённые ошибки новичков](https://habr.com/ru/companies/vk/articles/314804/)
  - [Разбираемся с новым sync.Map в Go 1.9](https://habr.com/ru/articles/338718/)
  - [Сценарии использования sync.Map](https://pkg.go.dev/sync#Map)
- sql
  - [Больше разработчиков должны знать это о базах данных](https://habr.com/ru/companies/flant/articles/500850/)
  - [use-the-index-luke](https://use-the-index-luke.com/)