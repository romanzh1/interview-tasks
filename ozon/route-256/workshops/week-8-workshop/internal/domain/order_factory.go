package domain

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/fossoreslp/go-uuid-v4"
)

type orderIDGenerator interface {
	Generate() int64
}

type userIDGenerator interface {
	Generate() int64
}

type descriptionGenerator interface {
	Generate() string
}

type createdAtGenerator interface {
	Generate() time.Time
}

type OrderFactory struct {
	//OrderIDGenerator     orderIDGenerator
	UserIDGenerator      userIDGenerator
	DescriptionGenerator descriptionGenerator
	CreatedAtGenerator   createdAtGenerator
}

func (f *OrderFactory) Create() Order {
	return Order{
		//ID:          f.OrderIDGenerator.Generate(),
		UserID:      f.UserIDGenerator.Generate(),
		Description: f.DescriptionGenerator.Generate(),
		CreatedAt:   f.CreatedAtGenerator.Generate(),
	}
}

func NewDefaultFactory() *OrderFactory {
	return New(
		//NewSeqGen(0),
		NewRndGen(int64(1000)),
		NewUUIDv4Generator(),
		&Clock{},
	)
}

func New(
	//orderIDGen orderIDGenerator,
	userIDGen userIDGenerator,
	idemKeyGen descriptionGenerator,
	momentGen createdAtGenerator,
) *OrderFactory {
	return &OrderFactory{
		//OrderIDGenerator:     orderIDGen,
		UserIDGenerator:      userIDGen,
		DescriptionGenerator: idemKeyGen,
		CreatedAtGenerator:   momentGen,
	}
}

type SeqGen struct {
	cur int
}

func NewSeqGen(start int) *SeqGen {
	return &SeqGen{
		cur: start,
	}
}

func (g *SeqGen) Generate() int64 {
	val := g.cur
	g.cur++

	return int64(val)
}

type RndGen struct {
	max int64
}

func NewRndGen(max int64) *RndGen {
	return &RndGen{
		max: max,
	}
}

func (g *RndGen) Generate() int64 {
	return rand.Int63n(g.max)
}

type Clock struct{}

func (c *Clock) Generate() time.Time {
	return time.Now()
}

type UUIDv4Generator struct {
	gen uuid.UUID
}

func NewUUIDv4Generator() *UUIDv4Generator {
	gen, err := uuid.New()
	if err != nil {
		panic(fmt.Sprintf("NewUUIDv4Generator failed: %s", err.Error()))
	}

	return &UUIDv4Generator{
		gen: gen,
	}
}

func (g *UUIDv4Generator) Generate() string {
	return g.gen.String()
}
