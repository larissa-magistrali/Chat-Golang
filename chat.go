package main

import (
	"bufio"
	"fmt"
	"os"

	BEB "./BEB"
)

var channels []string // Historico do chat

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Especifique pelo menos um endereço: porta!")
		return
	}

	addresses := os.Args[1:]
	fmt.Println("Enderecos:", addresses)

	beb := BEB.Modulo{
		Req: make(chan BEB.Envia_Mensagem),
		Ind: make(chan BEB.Recebe_Mensagem),
		NovoUsuario: make(chan BEB.Envia_Novo_Usuario), //Canal que envia um novo usuario
		RecebeUsuario: make(chan BEB.Recebe_Usuario), // Canal que recebe um novo usuario
		NovoGrupo: make(chan BEB.Envia_Novo_Grupo), // Canal que recebe dados de um grupo (usuarios e historico)
		RecebeGrupo: make(chan BEB.Recebe_Grupo)} // Canal que recebe dados de um grupo (usuarios e historico)

	beb.Init(addresses[0])

	fmt.Println("-------------COMANDOS---------------")
	fmt.Println("1) Enviar mensagem")
	fmt.Println("2) Visualizar histórico de mensagens")
	fmt.Println("3) Entrar em um chat")
	fmt.Println("4) Visualizar membros do chat")
	fmt.Println("------------------------------------")

	// enviador de broadcasts
	go func() {

		scanner := bufio.NewScanner(os.Stdin)
		var msg string
		var op string
		var ip string

		for {
			if scanner.Scan() {
				op = scanner.Text()
				switch op {
					case "1":
						fmt.Print("Enviar mensagem: ")
						if scanner.Scan() {
							msg = scanner.Text()
							channels = append(channels, addresses[0] + ":" + msg)
						}
						req := BEB.Envia_Mensagem{
							Addresses: addresses[1:],
							IpCorreto: addresses[0],
							Message:   msg}
						beb.Req <- req
					case "2":
						fmt.Printf("%+v\n", channels)
					case "3":
						fmt.Print("Ip: ") // Ip de algum usuario do chat que deseja entrar
						if scanner.Scan() {
							ip = scanner.Text()
						}
						req := BEB.Envia_Novo_Usuario{ // Envia o novo usuario para esse ip
							Address: ip,
							IpCorreto: addresses[0],
							Tag: "0"}
						beb.NovoUsuario <- req
						fmt.Print("Você entrou no chat!")
					case "4":
						fmt.Println(addresses)
				}
			}
		}
	}()

	// Rotina responsavel por receber novas mensagens
	go func() {
		for {
			mensagemRecebida := <-beb.Ind
			fmt.Printf("Mensagem de %v: %v\n", mensagemRecebida.IpCorreto, mensagemRecebida.Message)
			channels = append(channels, mensagemRecebida.IpCorreto + ":" + mensagemRecebida.Message)
		}
	}()

	// Rotina responsavel por receber dados de um grupo
	go func() {
		for {
			grupoRecebido := <-beb.RecebeGrupo
			addresses = append(addresses, grupoRecebido.Addresses...)
			channels = append(channels, grupoRecebido.Historico...)
			
		}
	}()

	// Rotina responsavel por receber usuarios novos
	go func() {
		for {
			usuarioRecebido := <-beb.RecebeUsuario
			//fmt.Printf("Usuario recebido de %v: %v\n", usuarioRecebido.From, usuarioRecebido.IpCorreto)
			
			// Tag responsavel por indicar se o usuario deve espalhar o novo usuario
			// 1 - Usuario novo solicita para um usuario do chat que quer entrar no chat
			// 2 - Usuario do chat manda o ip desse usuario novo para todos os outros usuarios do chat, para que todos consigam o adicionar e conversar
			// 3 - Porem os usuarios que receberam nao devem espalhar o ip de novo, pois isso ja esta sendo feito
			// A tag controla isso

			// Se a tag for 0 o ip deve ser espalhado
			if (usuarioRecebido.Tag == "0"){ 
				// Manda o novo usuario para todos os outros usuarios
				for i := 1; i < len(addresses); i++ {
					req := BEB.Envia_Novo_Usuario{
						Address: addresses[i],
						IpCorreto: usuarioRecebido.IpCorreto,
						Tag: "1"} // Adiciona tag 1 para avisar que nao deve espalhar mais
					beb.NovoUsuario <- req	
				}

				// Manda todos os usuarios e o historico de conversa para o novo membro
				req := BEB.Envia_Novo_Grupo{
					Addresses: addresses,
					IpCorreto: usuarioRecebido.IpCorreto,
					Historico: channels}
				beb.NovoGrupo <- req	
			
			}
			fmt.Printf(usuarioRecebido.IpCorreto + " entrou no chat!")
			// Adiciona o novo usuario a lista de participantes
			addresses = append(addresses, usuarioRecebido.IpCorreto)
		}
	}()

	blq := make(chan int)
	<-blq
}