## Ответы на вопросы из вопросников по вакансиям

### Чем горутины отличаются от потоков ОС?
Горутины - это аналог потоков ОС, но на уровне языка go, соответственно и управляются не ОС, а рантаймом go. В отличие от потоков ОС они имеют стек 2 кб, а потоки ОС 2 мб. Также они добавляют оптимизацию, планировщик может переключать го рутины на запущенном потоке ОС без его закрытия.

### Когда стоит использовать sync.Mutex, а когда — каналы?
Мьютекс стоит использовать тогда когда нужно защитить критическую секцию, то есть некий ресурс от многопоточного использования. Каналы стоит использовать когда нужно передавать некие данные между несколькими потоками или же распределить некую работу по нескольким потокам

### Что такое work stealing в планировщике Go?
Это механизм планировщика при котором процессоры забирают горутины готовые к исполнению у других процессоров, это сделано в качестве оптимизации на случай, если процессор не будет иметь работы

### Как планировщик Go обрабатывает блокирующие операции?
Если речь про сетевые вызовы, то планировщик помещает такие го рутины в network poller, а тот в свою очередь отправляет в linux процесс epoll, а он в свою очередь позволяет обрабатывать тысячи сетевых вызовов. Если же речь про системные вызовы, то есть например работа с файлом, то здесь применяется механизм handoff, планировщик отвязывает поток ОС от процессора, создаётся новый поток ОС или переиспользуется другой и процессор запускает свои го рутины уже на другом потоке

### Какие алгоритмы лежат в основе сборщика мусора в Go?
Алгоритм трёхцветной раскраски, concurrent mark and sweep. Сначала обходятся достижимые объекты, далее этап завершения маркировки и stop the world, далее происходит очистка объектов, на которые никто не ссылается

### Как избежать утечек памяти в Go, особенно при работе с горутинами и каналами?
Всегда закрывать каналы после окончания записи, всегда завершать го рутину после завершения её задачи, например по select <-ctx.Done(), очищать неиспользуемые буферы или пулы

### Зачем в Go используются пустые структуры?
Можно использовать для экономии на копейках в случаях когда требуется проверять некий флаг для множества сущностей, например в каналах или мапах

### Чем интерфейсы в Go отличаются от интерфейсов в других языках?
За счёт утиной типизации, то есть если объект имеет необходимые методы, то он удовлетворяет интерфейсу, не нужно явно указывать, что объект удовлетворяет ему

### Как нарушение Single Responsibility усложняет поддержку?
Усложняет поддержку на счёт повышения связности, изменения в одной части, могут потребовать правок в других местах, то же касается и работоспособности и тестирования, изменение 1 части сломает другую и тестировать будет сложнее

### Как реализовать Dependency Inversion в Go?
Написать интерфейс абстракцию, использовать его как тип в конструкторе создания высокоуровневого объекта, передать в месте создания высокуровневого объекта низкоуровневый, который удовлетворит необходимым методам и присвоить его полю структуры. Теперь в методах высокоуровневого объекта можно использовать методы низкоуровневого. Dependency Inversion готов

### В чем основное отличие Kafka от RabbitMQ?
RabbitMQ это классическая fifo очередь. Kafka это распределённый лог событий, 1 сообщение n потребителей. Kafka хранит данные на диске, хотя некоторые может кешировать, RabbitMQ в оперативной памяти, хотя может сделать снэпшот и выгрузить данные на диск.

### Как гарантировать доставку сообщений в системе с Message Queue?
В RabbitMQ брокер подтверждает запись на диск. В Kafka можно использовать ack: all, то есть подтверждение от всех реплик о получении от продьюсера.
При чтении же, в Rabbit потребитель должен подтвердить получение через ручной ack, а вот в Kafka нужно коммитить оффсет после обработки сообщения, тогда будет at least once

### Какие преимущества у PostgreSQL перед MySQL?
У postgres есть поддержка jsonb с индексацией, а у mysql только json и более медленный. В postgres есть частичные индексы, в mysql нет. В postgres можно ставить расширения вроде postgis и множества других, в mysql только движки таблиц. В postgres полноценная многопоточность, в mysql многопоточность только для отдельных операций вроде сканирования таблиц

### Как работает оператор HAVING и чем он отличается от WHERE?
Оператор having предназначен для фильтрации, работает после группировки и используется для аггрегирующих функций, where же тоже для фильтрации, но для значений столбцов
