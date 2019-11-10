// Construido como parte da disciplina de Sistemas Distribuidos
// Semestre 2018/2  -  PUCRS - Escola Politecnica
// Estudantes:  Andre Antonitsch e Rafael Copstein
// Professor: Fernando Dotti  (www.inf.pucrs.br/~fldotti)
// Algoritmo baseado no livro:
// Introduction to Reliable and Secure Distributed Programming
// Christian Cachin, Rachid Gerraoui, Luis Rodrigues

package main

import (
	"bufio"
	"fmt"
	"os"

	BEB "./BEB"
)

func main() {

	if len(os.Args) < 2 {
		fmt.Println("Please specify at least one address:port!")
		return
	}

	addresses := os.Args[1:]
	fmt.Println(addresses)

	beb := BEB.BestEffortBroadcast_Module{
		Req: make(chan BEB.BestEffortBroadcast_Req_Message),
		Ind: make(chan BEB.BestEffortBroadcast_Ind_Message)}

	beb.Init(addresses[0])

	// enviador de broadcasts
	go func() {

		scanner := bufio.NewScanner(os.Stdin)
		var msg string
		var op string

		for {
			fmt.Println("----------------------")
			fmt.Println("Digite 1 para enviar msg")
			fmt.Println("Digite 2 para rever as msgs")
			fmt.Print("-> ")
			if scanner.Scan() {
				op = scanner.Text()
				switch op {
				case "1":
					fmt.Print("Enviar msg: ")
					if scanner.Scan() {
						msg = scanner.Text()
					}
					req := BEB.BestEffortBroadcast_Req_Message{
						Addresses: addresses[1:],
						Message:   msg}
					beb.Req <- req
				case "2":

				}
			}

		}
	}()

	// receptor de broadcasts
	go func() {
		for {
			in := <-beb.Ind
			fmt.Printf("Message from %v: %v\n", in.From, in.Message)
		}
	}()

	blq := make(chan int)
	<-blq
}

/*
go run chat.go 127.0.0.1:5001  127.0.0.1:6001
go run chat.go 127.0.0.1:6001  127.0.0.1:5001
*/

/*
go func() {

		scanner := bufio.NewScanner(os.Stdin)
		var msg string

		for {
			if scanner.Scan() {
				msg = scanner.Text()
			}
			req := BEB.BestEffortBroadcast_Req_Message{
				Addresses: addresses[1:],
				Message:   msg}
			beb.Req <- req
		}
	}()

*/