package main

import (
	"flag"
	"fmt"
	"sync"
	"time"
)

type Cashier struct {
	CashierId int
	Status    bool //occupied or not
}

type Customer struct {
	CustomerId int
}

func Process(cashierChannel chan Cashier, customerChannel chan Customer, timePerCustomer int) {

	fmt.Println(time.Now().Format("2006-01-02 03:04:05")+" --> ", "Bank Simulation Started")
	wg := sync.WaitGroup{}
	for {
		select {
		case customer := <-customerChannel:
			cashier := <-cashierChannel
			wg.Add(1)
			go func() {
				defer wg.Done()
				fmt.Println(time.Now().Format("2006-01-02 03:04:05")+" --> ", "Cashier ", cashier.CashierId, ": Customer ", customer.CustomerId, "Started")
				time.Sleep(time.Duration(timePerCustomer) * time.Second)
				fmt.Println(time.Now().Format("2006-01-02 03:04:05")+" --> ", "Cashier ", cashier.CashierId, ": Customer ", customer.CustomerId, "Completed")
				cashier.Status = false
				cashierChannel <- cashier
			}()
		default:
			wg.Wait()
			return
		}
	}
}

func iterateCashiers(cashiers []Cashier, cashierChannel chan Cashier) {
	for _, cashier := range cashiers {
		cashier.Status = true
		cashierChannel <- cashier
	}
}

func iterateCustomers(customers []Customer, customerChannel chan Customer) {
	for _, customer := range customers {
		customerChannel <- customer
	}
}

func main() {

	numCashiers := flag.Int("numCashiers", 2, "an int")
	numCustomers := flag.Int("numCustomers", 20, "an int")
	timePerCustomer := flag.Int("timePerCustomer", 3, "an int")

	flag.Parse()

	//onboard cashiers
	cashiers := make([]Cashier, *numCashiers)
	for i := 0; i < *numCashiers; i++ {
		cashiers[i].CashierId = i
		cashiers[i].Status = false //not occupied
	}

	//onboard customers
	customers := make([]Customer, *numCustomers)
	for i := 0; i < *numCustomers; i++ {
		customers[i].CustomerId = i
	}
	cashierChannel := make(chan Cashier, *numCashiers)
	customerChannel := make(chan Customer, *numCustomers)

	iterateCashiers(cashiers, cashierChannel)
	iterateCustomers(customers, customerChannel)

	Process(cashierChannel, customerChannel, *timePerCustomer)
	fmt.Println(time.Now().Format("2006-01-02 03:04:05")+" --> ", "Bank Simulated Ended")
	close(customerChannel)

}
