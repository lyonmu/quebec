# Kafka

```go
	type Payload struct {
		Test string `json:"test"`
	}
	func (p *Payload) Reset() {
		p.Test = ""
	}

	// producer
	kafkaProducer, err := kafka.NewJSONProducer[string, Payload]([]string{"localhost:9092"}, "test-topic")
	if err != nil {
		panic(err)
	}
	ctx := context.Background()
	err = kafkaProducer.Produce(ctx, nil, &Payload{Test: "test"})
	if err != nil {
		panic(err)
	}

	// raw producer
	kafkaProducer, err := kafka.NewRawProducer([]string{"localhost:9092}", "test-topic")
	if err != nil {
		panic(err)
	}
	ctx := context.Background()
	bytes := []byte("xxx")
	err = kafkaProducer.Produce(ctx, nil, &bytes)
	if err != nil {
		panic(err)
	}


	// subscriber
	// normal reader
	reader := kafka.NewReader[string, Payload](func() (*string, *Payload) {
		return nil, &Payload{}
	})

	// pool reader
	reader := kafka.NewPoolReader[string, Payload]()
	// json subscriber
	kafakSubscriber, err := kafka.NewJSONSubscriber[string, Payload]([]string{"localhost:9092"}, "test-topic", "test-group", reader)
	if err != nil {
		panic(err)
	}

	// raw reader
	reader := kafka.NewRawReader[bytes.Buffer, bytes.Buffer]()
	kafakSubscriber, err := kafka.NewSubscriber[bytes.Buffer, bytes.Buffer]([]string{"localhost:9092"}, "test-topic", "test-group", reader)
	if err != nil {
		panic(err)
	}

	err = kafakSubscriber.Subscribe(ctx, func(key *bytes.Buffer, payload *bytes.Buffer, err error) {
	})

```
