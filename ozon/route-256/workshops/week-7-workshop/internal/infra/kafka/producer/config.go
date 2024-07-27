package producer

import (
	"time"

	"github.com/IBM/sarama"
)

func PrepareConfig(opts ...Option) *sarama.Config {
	c := sarama.NewConfig()

	// алгоритм выбора партиции
	{
		/*
			Кейсы:
				- одинаковые ключи в одной партиции
				- при cleanup.policy = compact останется только последнее сообщение по этому ключу
		*/
		// ручной
		//c.Producer.Partitioner = sarama.NewManualPartitioner

		// случайная партиция
		//c.Producer.Partitioner = sarama.NewRandomPartitioner

		// по кругу
		//c.Producer.Partitioner = sarama.NewRoundRobinPartitioner

		// по ключу
		c.Producer.Partitioner = sarama.NewHashPartitioner
	}

	// acks параметр
	{
		// acks = 0 (none) - ничего не ждем
		//c.Producer.RequiredAcks = sarama.NoResponse
		// acks = 1 (one) - ждем успешной записи ТОЛЬКО на лидер партиции
		//c.Producer.RequiredAcks = sarama.WaitForLocal
		// acks = -1 (all) - ждем успешной записи на лидер партиции и всех in-sync реплик (настроено в кафка кластере)
		c.Producer.RequiredAcks = sarama.WaitForAll
	}

	// семантика exactly once
	{
		/*
			Если хотим добиться семантики exactly once, то выставляем в true

			У продюсера есть счетчик (count).
			Каждое успешно отправленное сообщение учеличивает счетчик (count++).
			Если продюсер не смог отправить сообщение, то счетчик не меняется и отправляется в таком виде в другом сообщение.
			Кафка это видит и начинает сравнивать (в том числе Key) сообщения с одниковыми счетчиками.
			Далее не дает отправить дубль, если Idempotent = true.
		*/
		c.Producer.Idempotent = false
	}

	// повторы ошибочных отправлений
	{
		// число попыток отправить сообщение
		c.Producer.Retry.Max = 100
		// интервалы между попытками отправить сообщение
		c.Producer.Retry.Backoff = 5 * time.Millisecond
	}

	{
		// Уменьшаем пропускную способность, тем самым гарантируем строгий порядок отправки сообщений/батчей
		c.Net.MaxOpenRequests = 1
	}

	// сжатие на клиенте
	{
		// Если хотим сжимать, то задаем нужный уровень кодировщику
		c.Producer.CompressionLevel = sarama.CompressionLevelDefault
		// И сам кодировщик
		c.Producer.Compression = sarama.CompressionGZIP
	}

	{
		/*
			Если эта конфигурация используется для создания `SyncProducer`, оба параметра должны быть установлены
			в значение true, и вы не не должны читать данные из каналов, поскольку это уже делает продьюсер под капотом.
		*/
		c.Producer.Return.Successes = true
		c.Producer.Return.Errors = true
	}

	for _, opt := range opts {
		_ = opt.Apply(c)
	}

	return c
}
