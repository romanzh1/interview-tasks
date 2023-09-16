 // Дан массив meetings, в котором каждый элемент meeting[i] - это пара двух чисел [startTime, endTime].
// Необходимо объединить все накладывающиеся друг на друга встречи и вернуть массив с объединенными встречами, покрывающих те же временные слоты.
//
// Input: [[1,3], [2,6], [8,10], [15,18]]
// Output: [[1,6], [8,10], [15,18]]
// Explanation: Интервалы [1,3] и [2,6] пересекаются => объединяем в [1,6].


func join(meet [][]int) [][]int{
    intervals := make([][]int, 0)

    for i := 0; i < len(meet) - 1; i++{
        if meet[i][1] > meet[i+1][0]{
            intervals = append(intervals, []int{meet[i][0], meet[i+1][1]})
        } else {
            intervals = append(intervals, meet[i])
        }
    }

    return intervals
}

//////////////////

// Поочередно выполнит http запросы по предложенному списку ссылок

// в случае получения http-кода ответа на запрос "200 OK" печатаем на экране "адрес url - ok"
// в случае получения http-кода ответа на запрос отличного от "200 OK" либо в случае ошибки печатаем на экране "адрес url - not ok"

// Модифицируйте программу таким образом, чтобы использовались каналы для коммуникации основного потока с горутинами. Пример:

// Запросы по списку выполняются в горутинах.
// Печать результатов на экран происходит в основном потоке

package main

func main() {
  var urls = []string{
    "http://ozon.ru",
    "https://ozon.ru",
    "http://google.com",
    "http://somesite.com",
    "http://non-existent.domain.tld",
    "https://ya.ru",
    "http://ya.ru",
    "http://ёёёё",
  }

  wg := WaitGroup{}

  result := make(chan, len(url))
  wg.Add(len(urls))
  for i := 0; i < len(urls); i++ {
        go func(int i){
            defer wg.Done()
            result <- worker(urls[i])
        }(i)
  }

  go func(){
      wg.Add(1)
      fmt.Println(<-result)
      wg.Done()
  }
  wg.Wait()
}



func worker(string url) string{
    get, err := http.Get(url)
    if err != nil {
        return fmt.Sprintf("url %d is not ok", get.code)
    }
}

//////////////////////////////////

func Get(url string) (resp, error)

type resp struct {
    code int
}

//////////////////////////////////

// Нужно описать модель библиотеки. Есть 3 сущности: "Автор", "Книга", "Читатель"
// Физически книга только одна и может быть только у одного читателя.
// Нужно составить таблицы для библиотеки так что бы это учесть.


Автор
    ID
    Name

Автор-книга
    ID_book
    ID_author

Книга
    ID
    ID_reader
    Name

Читатель
    ID
    Name

/////////////

// Написать запрос - выбрать названия всех книг в библиотеке у которых больше 3 авторов

SELECT name, COUNT(ID_author) as auth
FROM Книга AS k
INNER JOIN Автор-книга AS ab ON ab.ID_author = k.ID
GROUP BY name
HAVING auth > 3


