package main

import (
	"crypto/rand"
	"fmt"
	"log"
	"math/big"
)

type Transaction struct {
	ID     int64
	Amount float64
}

type Pipeline struct {
	convertIndex float64
}

func NewPipeline(convertIndex float64) *Pipeline {
	return &Pipeline{
		convertIndex: convertIndex,
	}
}

func GetTransactions(amount int) chan *Transaction {
	transactions := make(chan *Transaction)

	go func() {
		for i := 0; i < amount; i++ {
			randNegative, err := rand.Int(rand.Reader, big.NewInt(2))
			if err != nil {
				log.Fatal(err)
			}

			randSum, err := rand.Int(rand.Reader, big.NewInt(1000))
			if err != nil {
				log.Fatal(err)
			}
			resultRandNum, _ := randSum.Float64()

			if randNegative.Int64() == 1 {
				resultRandNum *= -1
			}

			transactions <- &Transaction{
				ID:     int64(i + 1),
				Amount: resultRandNum,
			}
		}
		close(transactions)
	}()

	return transactions
}

func (p *Pipeline) Filter(in <-chan *Transaction) <-chan *Transaction {
	out := make(chan *Transaction)

	go func() {
		for transaction := range in {
			if transaction.Amount >= 0 {
				out <- transaction
				continue
			}
			fmt.Printf("filtered transaction ID %d with negative Amount %f.2\n",
				transaction.ID, transaction.Amount)
		}
		close(out)
	}()

	return out
}

func (p *Pipeline) Convert(in <-chan *Transaction) <-chan *Transaction {
	out := make(chan *Transaction)

	go func() {
		for transaction := range in {
			transaction.Amount *= p.convertIndex
			out <- transaction
		}
		close(out)
	}()

	return out
}

func (p *Pipeline) Result(in <-chan *Transaction) {
	for transaction := range in {
		fmt.Printf("result: transaction ID %d, Amount %f\n", transaction.ID, transaction.Amount)
	}
}

func (p *Pipeline) Run(transactions chan *Transaction) {
	filtered := p.Filter(transactions)
	converted := p.Convert(filtered)
	p.Result(converted)
}

func main() {
	transactions := GetTransactions(5)

	pipeline := NewPipeline(1.5)
	pipeline.Run(transactions)
}
