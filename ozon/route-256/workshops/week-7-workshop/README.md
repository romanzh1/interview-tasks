# Workshop по Apache Kafka

## План:
- знакомимся с локальным окружением
  - поднимаем apache kafka в docker окружении
  - изучаем физическое хранение
  - изучаем как через консоль управлять кластером
- изучаем консольные утилиты
  - kafka-topics
  - kafka-console-producer
  - kafka-console-consumer
  - kafka-consumer-groups
  - kafka-get-offsets
- знакомимся с ui
  - изучаем kafka-ui
  - топики
  - партиции
  - сообщения
  - публикация сообщений
- знакомимся с библиотекой sarama
  - подключаем kafka в go проекте
  - обсуждаем структуру проекта
  - изучаем конфигурационные параметры
- изучаем код синхронного продюсера
- изучаем запуск приложения
- пушим сообщения с hashPartitioner-ом но без ключа (несколько раз)
    - убеждаемся что сообщения случайно распределяются по партициям
- пушим сообщения с hashPartitioner-ом с ключом по orderID
    - убеждаемся что соблюдается упорядочивание в рамках ключа
- устанавливаем номер партиции в сообщении, видиим что для hashPartitioner-а это ни на что не влияет
- устанавливаем ManualPartitioner, запускаем без указания партиции
    - партиция везде = zero value
- не нулевую партицию, проверяем что все сообщения в ней
- выбираем RoundRobinPartitioner
    - убеждаемся что сообщения раскладываются по порядку
- выбираем RandomPartitioner
    - несмотря на указание ключа партиции выбираются случайно
- изучаем асинхронный продьюсер
- изучаем консьюмер
- изучаем консьюмер группу
- проверяем как работает консьюмер группа если запустить несколько инстансов
- докеризируем консьюмер группу
- обсуждаем kcat https://github.com/edenhill/kcat

## Знакомимся с локальным окружением

- [docker-compose.yml](build/dev/docker-compose.yml)
- [Makefile](Makefile)

## Запускаем локальное окружение

```bash
make compose-up
```

```bash
# Тут хранится конфигурация кафки
ls -al /etc/kafka/

# Тут хранятся партиции топиков
ls -al /tmp/kraft-combined-logs

# Тут хранятся сегменты партиции
ls -al /tmp/kraft-combined-logs/route256-example-0/

# Набор программ для управление kafka
ls -al /bin/kafka-*
```

kafka

## Программа kafka-topics
```bash
# Создаем топик
# replication-factor нельзя поставить больше 1 т.к. только один брокер настроен
kafka-topics --create --topic route256-yet-another --partitions 3 --replication-factor 1 --if-not-exists --bootstrap-server kafka0:29092

# Список топиков
kafka-topics --list --bootstrap-server kafka0:29092

# Просмотр информации о топике
kafka-topics --describe --bootstrap-server kafka0:29092 --topic route256-yet-another

# Добавляем партиции в топик
kafka-topics --bootstrap-server kafka0:29092 --alter --topic route256-yet-another --partitions 10

# Уменьшать количество партиций нельзя
kafka-topics --bootstrap-server kafka0:29092 --alter --topic route256-yet-another --partitions 5

# Удаляем топик
kafka-topics --bootstrap-server kafka0:29092 --delete --topic route256-yet-another
```

## kafka-console-producer

```bash
# Пишем событие в кафку
kafka-console-producer --topic route256-example --bootstrap-server kafka0:29092
# parse.key=true
kafka-console-producer --topic route256-example --property parse.key=true --bootstrap-server kafka0:29092
# key.separator
kafka-console-producer --topic route256-example --property parse.key=true --property key.separator=: --property parse.key=true --bootstrap-server kafka0:29092
#kafka-console-producer --producer-property partitioner.class=org.apache.kafka.clients.producer.RoundRobinPartitioner --topic route256-example --bootstrap-server kafka0:29092
#kafka-console-producer --producer-property partitioner.class=org.apache.kafka.clients.producer.UniformStickyPartitioner --topic route256-example --bootstrap-server kafka0:29092
```

## kafka-console-consumer
```bash
# Читаем сообщения с конца
kafka-console-consumer --topic route256-example --bootstrap-server kafka0:29092

# Читаем сообщения с начала
kafka-console-consumer --topic route256-example --bootstrap-server kafka0:29092 --from-beginning

# Читаем сообщения из определенной партиции
kafka-console-consumer --topic route256-example --bootstrap-server kafka0:29092 --partition 1 --from-beginning

# Читаем с определенного офсета отсчитывая от начала
kafka-console-consumer --topic route256-example --bootstrap-server kafka0:29092 --partition 1 --offset 2

# Читаем с помощью консьюмер группы
## Сначала показываем что мы переподключаемся к топики и дочитываем с того же места
kafka-console-consumer --topic route256-example --bootstrap-server kafka0:29092 --group "kafka-cli"

## Делаем 3 консьюмера к одной группе и показываем что сообщения распределяются между ними
## docker exec -ti 2b3f2244c594 sh
## kafka-console-consumer --topic route256-yet-another --bootstrap-server kafka0:29092 --group "kafka-cli"

# Можем читать сразу несколько топиков
kafka-console-consumer --whitelist 'route256-yet-another|route256-example' --bootstrap-server kafka0:29092 --group "kafka-cli"
```

## kafka-consumer-groups

```bash
# Список консьюмер групп
kafka-consumer-groups --bootstrap-server kafka0:29092 --list

# Детальная информация о консьюмер группе
kafka-consumer-groups --bootstrap-server kafka0:29092 --describe --group kafka-cli

# Сброс офсета для консьюмер группы для всех партиций
# необходимо отключить консьюмеров
#--to-current                            Reset offsets to current offset.       
#--to-datetime <String: datetime>        Reset offsets to offset from datetime. Format: 'YYYY-MM-DDTHH:mm:SS.sss'  
#--to-earliest                           Reset offsets to earliest offset.      
#--to-latest                             Reset offsets to latest offset.        
#--to-offset <Long: offset>              Reset offsets to a specific offset. 
kafka-consumer-groups --bootstrap-server kafka0:29092 --reset-offsets --group kafka-cli --topic route256-example --to-offset 15 --execute

# Для конкретной партиций 0 и 2
kafka-consumer-groups --bootstrap-server kafka0:29092 --reset-offsets --group kafka-cli --topic route256-example:0,2 --to-offset 12 --execute

# Получение списка консьюмеров группы и распределение партиций по ним
kafka-consumer-groups --bootstrap-server kafka0:29092 --describe --members --verbose --group kafka-cli

# Сброс офестов для всех топиков
kafka-consumer-groups --bootstrap-server kafka0:29092 --reset-offsets --group kafka-cli --all-topics --to-earliest --execute

# Удаление группы
kafka-consumer-groups --bootstrap-server kafka0:29092 --delete --group kafka-cli
```

## kafka-get-offsets
```bash
# Получение офсетов для всех топиков
kafka-get-offsets --bootstrap-server kafka0:29092

# Офсет отдельного топика
kafka-get-offsets --bootstrap-server kafka0:29092 --topic route256-example
# Офсет партиции отдельного топика
kafka-get-offsets --bootstrap-server kafka0:29092 --topic-partitions route256-example:0
```

## Знакомимся с ui

http://localhost:8080

## Знакомимся с библиотекой sarama

https://github.com/IBM/sarama

```bash
go get github.com/IBM/sarama
```

## sync-producer

```bash
# проверяем стратегии распределения сообщений по партициям (настройки по-умолчанию)
go run ./cmd/sync_producer

# для round robin стратегии
go run ./cmd/sync_producer -count=3
```

- изучаем устройство
- рассматриваем различные стратегии партиционирования
- передачу и игнорирования ключа в сообщении
- ручное указание партиции при различных стратегиях
- сравниваем как сообщения с одинаковым ключом распределяются по партициям

## async-producer

```bash
# проверяем стратегии распределения сообщений по партициям (настройки по-умолчанию)
go run ./cmd/async_producer
```

- изучаем устройство
- запускаем с различными комбинациями FlushMessages
- запускаем с различными комбинациями FlushFrequency
- обращаем внимание на gracefully shutdown


## consumer

```bash
go run ./cmd/consumer

go run ./cmd/sync_producer -repeat-count=1 -interval=1s
```

- изучаем устройство
- очистить топик (для наглядности)
- запускаем продюсер или публикуем сообщения вручную (увеличить таймаут)
- повторно запускаем консьюмер, убеждаемся что все сообщения перечтаны сначала

## consumer-group

```bash
go run ./cmd/consumer_group

go run ./cmd/sync_producer -count=1000 -interval=1s

# смотрим распределение партиций по консьюмер группе
kafka-consumer-groups --bootstrap-server kafka0:29092 --describe --members --group route256-consumer-group

docker-compose -p route-256-ws-7 -f "./build/dev/docker-compose.yml" logs -f
```

- изучаем устройство
- очистить топик (для наглядности)
- запускаем продюсер или публикуем сообщения вручную (увеличить таймаут)
- смотрим логи кафки (docker)
- запускаем по одному дополнительные консьюмеры, смотрим за перебалансировкой группы
  - убеждаемся что 3-й консьюмер не получает сообщения
  - убеждаемся что каждый консьюмер читает свою партицию
- выключаем один из консьюмеров, смотрим за перебалансировкой группы
  - убеждаемся что 3-й консьюмер начал получать сообщения
- проверяем работу комита смещений
  - автокомит
    - пушим сообщения,
    - запускаем консьюмер,
    - останавливаем констьюмер,
    - повторно запускаем
    - смотрим смещения в консьюмер группе
    - убеждаемся что нет повторной обработки сообщений
  - автокоммит без MarkMessage()
    - убеждаемся что фиксации смещения не происходит
  - эмулируем панику с автокомитом
    - убеждаемся что при панике смещение не фиксируется
    - обсуждаем что должен быть recover в хендлере
    - обсуждаем почему нельзя делать фиксацию смещения при shutdown приложения
  - с выключенным автокоммитом, но без фиксации
    - убеждаемся что фиксации смещения не происходит
  - с ручной фиксацией смещения без автокомита
    - убеждаемся что фиксация смещения происходит

## Заппускаем несколько контейнеров в группе

```bash
# docker-compose -p route256 -f ./build/dev/docker-compose.yml up go-consumer
# запускаем 3 консьюмера в докере
make compose-up

# запускаем продюсер
go run ./cmd/sync_producer -repeat-count=1 -count=1000

# смотрим логи
docker-compose -p route-256-ws-7 -f "./build/dev/docker-compose.yml" logs -f
```

- изучаем multistage сборку
- правим docker-compose
- запускаем несколько консьюмеров
- читаем логи докера
- запускаем публикацию сообщеий
- изучаем распределение сообщений между консьюмерами
- постепенно отключаем консьюмеры/включаем, смотрим логи
