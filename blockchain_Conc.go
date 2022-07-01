package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Frame struct {
	Cmd    string   `json:"cmd"`
	Sender string   `json:"sender"`
	Data   []string `json:"data"`
}

type Info struct {
	nextNode string
	nextNum  int
	imFirst  bool
	cont     int
}

type InfoCons struct {
	contA, contB int
}

type Block struct {
	timestamp time.Time
	operacion []string
	prev_hash []byte
	hash      []byte
}

var (
	host         string
	myNum        int
	chRemotes    chan []string
	chInfo       chan Info
	chCons       chan InfoCons
	readyToStart chan bool
	participants int
)

func main() {
	HandleRequest()
	rand.Seed(time.Now().UnixNano())
	if len(os.Args) == 1 {
		log.Println("Hostname not given")
	} else {
		host = os.Args[1]
		chRemotes = make(chan []string, 1)
		chInfo = make(chan Info, 1)
		chCons = make(chan InfoCons, 1)
		readyToStart = make(chan bool, 1)

		/*abc := []string{" se deposito 50"}*/

		abc := []string{"id="}
		bloque_after := Blocks(abc, []byte{})
		fmt.Println(" bloque numero:")
		Print(bloque_after)

		for i := 0; i < 3; i++ {

			cda := []string{"id:"}
			bloque_before := Blocks(cda, bloque_after.hash)
			fmt.Println(" bloque") //bloque numero: , i
			Print(bloque_before)

			var contenido int
			for f := 0; f < 200; f++ {
				contenido = rand.Intn(101)
			}

			fmt.Println("\t\t\t", contenido)

		}

		chRemotes <- []string{}
		if len(os.Args) >= 3 {
			connectToNode(os.Args[2])
		}
		if len(os.Args) == 4 {
			switch os.Args[3] {
			case "agrawalla":
				go startAgrawalla()
			case "consensus":
				go startConsensus()

			}
		}
		server()
	}
	//datos de prueba de los bloques 1 2 y 3

}

func startAgrawalla() {
	time.Sleep(3 * time.Second)
	remotes := <-chRemotes
	chRemotes <- remotes
	for _, remote := range remotes {
		send(remote, Frame{"agrawalla", host, []string{}}, nil)
	}
	handleAgrawalla()
}

func startConsensus() {
	remotes := <-chRemotes
	for _, remote := range remotes {
		log.Printf("%s: notifying %s\n", host, remote)
		send(remote, Frame{"consensus", host, []string{}}, nil)
	}
	chRemotes <- remotes
	handleConsensus()
}

func connectToNode(remote string) {
	remotes := <-chRemotes
	remotes = append(remotes, remote)
	chRemotes <- remotes
	if !send(remote, Frame{"hello", host, []string{}}, func(cn net.Conn) {
		dec := json.NewDecoder(cn)
		var frame Frame
		dec.Decode(&frame)
		remotes := <-chRemotes
		remotes = append(remotes, frame.Data...)
		chRemotes <- remotes
		log.Printf("%s: friends0: %s\n", host, remotes)
	}) {
		log.Printf("%s: unable to connect to %s\n", host, remote)
	}
}

func send(remote string, frame Frame, callback func(net.Conn)) bool {
	if cn, err := net.Dial("tcp", remote); err == nil {
		defer cn.Close()
		enc := json.NewEncoder(cn)
		enc.Encode(frame)

		createCSV("", "", "")
		if callback != nil {
			callback(cn)
		}
		return true
	} else {
		log.Printf("%s: can't connect to %s\n", host, remote)
		idx := -1
		remotes := <-chRemotes
		for i, rem := range remotes {
			if remote == rem {
				idx = i
				break
			}
		}
		if idx >= 0 {
			remotes[idx] = remotes[len(remotes)-1]
			remotes = remotes[:len(remotes)-1]
		}
		chRemotes <- remotes
		return false
	}
}

func createCSV(port string, msg string, hash string) {
	empData := [][]string{
		{port, msg, hash},
	}

	csvFile, err := os.Create("data.csv")

	if err != nil {
		log.Fatalf("failed creating file: %s", err)
	}

	csvwriter := csv.NewWriter(csvFile)

	for _, empRow := range empData {
		_ = csvwriter.Write(empRow)
	}
	csvwriter.Flush()
	csvFile.Close()
}

func ReadCSVFromHttpRequest(res http.ResponseWriter, req *http.Request) {

	file, err := os.Open("./data.csv")
	if err != nil {
		log.Fatal("Unable to read input", err)
	}
	reader := csv.NewReader(file)
	var results [][]string
	for {
		// read one row from csv
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		// add record to result set
		results = append(results, record)
	}

	res.Header().Set("Content-Type", "application/json; charset=utf-8")
	res.Header().Set("Access-Control-Allow-Origin", "*")
	res.Header().Set("Access-Control-Allow-Headers", "Content-Type,access-control-allow-origin, access-control-allow-headers")
	jsonBytes, _ := json.MarshalIndent(results, "", " ")
	io.WriteString(res, string(jsonBytes))
}

func HandleRequest() {
	var port = "9000"
	http.HandleFunc("/data", ReadCSVFromHttpRequest)
	fmt.Printf("Corriendo desde el puerto :%s\n", port)
	fmt.Printf("Llamar al dataset localhost:%s/data\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func server() {
	if ln, err := net.Listen("tcp", host); err == nil {
		defer ln.Close()
		log.Printf("Listening on %s\n", host)

		for {
			if cn, err := ln.Accept(); err == nil {
				go fauxDispatcher(cn)

			} else {
				log.Printf("%s: cant accept connection.\n", host)
			}
		}
	} else {
		log.Printf("Can't listen on %s\n", host)
	}
}

func fauxDispatcher(cn net.Conn) {
	defer cn.Close()
	dec := json.NewDecoder(cn)
	frame := &Frame{}
	dec.Decode(frame)
	switch frame.Cmd {
	case "hello":
		handleHello(cn, frame)
	case "add":
		handleAdd(frame)
	case "agrawalla":
		handleAgrawalla()
	case "num":
		handleNum(frame)
	case "start":
		handleStart()
	case "consensus":
		handleConsensus()
	case "vote":
		handleVote(frame)
		/*case "potato":
		handlePotato(frame)*/
	}
}

func handleHello(cn net.Conn, frame *Frame) {
	enc := json.NewEncoder(cn)
	remotes := <-chRemotes
	enc.Encode(Frame{"<response>", host, remotes})
	notification := Frame{"add", host, []string{frame.Sender}}
	for _, remote := range remotes {
		send(remote, notification, nil)
	}
	remotes = append(remotes, frame.Sender)
	log.Printf("%s: friends1: %s\n", host, remotes)
	chRemotes <- remotes

}

//agregar block *Block para trabajar con bloques
func handleAdd(frame *Frame) {
	remotes := <-chRemotes
	remotes = append(remotes, frame.Data...)
	log.Printf("%s: friends2: %s\n", host, remotes)
	chRemotes <- remotes
	/*
		go Blocks(block.operacion, block.prev_hash)
		Print(block)
	*/
}
func handleAgrawalla() {
	myNum = rand.Intn(1000000000)
	log.Printf("%s: my number is %d\n", host, myNum)
	msg := Frame{"num", host, []string{strconv.Itoa(myNum)}}
	remotes := <-chRemotes
	chRemotes <- remotes
	for _, remote := range remotes {
		send(remote, msg, nil)
	}
	chInfo <- Info{"", 1000000001, true, 0}
}

func handleNum(frame *Frame) {
	if num, err := strconv.Atoi(frame.Data[0]); err == nil {
		info := <-chInfo
		//log.Printf("from %v\n", frame)
		if num > myNum {
			if num < info.nextNum {
				info.nextNum = num
				info.nextNode = frame.Sender
			}
		} else {
			info.imFirst = false
		}
		info.cont++
		chInfo <- info
		remotes := <-chRemotes
		chRemotes <- remotes
		if info.cont == len(remotes) {
			if info.imFirst {
				log.Printf("%s: I'm first!\n", host)
				criticalSection()
			} else {
				readyToStart <- true
			}
		}
	} else {
		log.Printf("%s: can't convert %v\n", host, frame)
	}
}

func handleStart() {
	<-readyToStart
	criticalSection()
}
func handleConsensus() {
	time.Sleep(3 * time.Second)
	var op string
	/*fmt.Print("A o B, elige: ")
	fmt.Scanf("%s\n", &op)*/

	if rand.Intn(100) > 50 {
		op = "A"
	} else {
		op = "B"
	}
	info := InfoCons{0, 0}
	if op == "A" {
		info.contA++
	} else {
		info.contB++
	}
	chCons <- info

	remotes := <-chRemotes
	participants = len(remotes) + 1
	for _, remote := range remotes {
		log.Printf("%s: sending %s to %s\n", host, op, remote)
		send(remote, Frame{"vote", host, []string{op}}, nil)
	}
	chRemotes <- remotes
}
func handleVote(frame *Frame) {
	vote := frame.Data[0]
	info := <-chCons
	if vote == "A" {
		info.contA++
	} else {
		info.contB++
	}
	chCons <- info
	log.Printf("%s: %s voted %s\n", host, frame.Sender, vote)
	if info.contA+info.contB == participants {
		if info.contA > info.contB {
			log.Printf("%s the A won\n", host)
		} else {
			log.Printf("%s the A won\n", host)
		}
	}
}

/*
func handlePotato(frame *Frame) {
	if num, err := strconv.Atoi(frame.Data[0]); err == nil {
		log.Printf("%s: recibí %d\n", host, num)
		if num == 0 {
			log.Printf("%s: perdí\n", host)
		} else {
			for len(remotes) > 0 {
				remote := remotes[rand.Intn(len(remotes))]
				data := []string{strconv.Itoa(num - 1)}
				time.Sleep(100 * time.Millisecond)
				if send(remote, Frame{"potato", host, data}, nil) {
					break
				}
			}
		}
	} else {
		log.Printf("%s: can't convert %s to number\n", host, frame.Data)
	}
}
*/

func criticalSection() {
	log.Printf("%s: my time has come!\n", host)
	info := <-chInfo
	if info.nextNode != "" {
		log.Printf("%s: letting %s start\n", host, info.nextNode)
		send(info.nextNode, Frame{"start", host, []string{}}, nil)
	} else {
		log.Printf("%s: I was the last one :(\n", host)
	}
}

//crea bloque
func Blocks(operacion []string, prev_hash []byte) *Block {
	currentTime := time.Now()
	return &Block{
		timestamp: currentTime,
		operacion: operacion,
		prev_hash: prev_hash,
		hash:      NewHash(),
	}

}

//funcion de encriptado
//https://golangbyexample.com/generate-uuid-guid-golang/
func NewHash() []byte {
	uuidWithHyphen := uuid.New()
	uuid := strings.Replace(uuidWithHyphen.String(), "-", "", -1)
	return []byte(uuid)
}

//pinta los bloques
func Print(block *Block) {

	fmt.Printf("\ttime: %s\n", block.timestamp.String())
	fmt.Printf("\tprev_hash: %s\n", block.prev_hash)
	fmt.Printf("\thash: %s\n", block.hash)
	Operacion(block)

}

//pinta las operaciones
func Operacion(block *Block) {
	fmt.Println("\toperacion: ")
	for i, operacion := range block.operacion {
		fmt.Printf("\t\t%v: %q\n", i, operacion)
	}
}

/*
func potatoGenerator() {
	for {
		time.Sleep(5 * time.Second)
		for len(remotes) > 0 {
			remote := remotes[rand.Intn(len(remotes))]
			data := []string{strconv.Itoa(rand.Intn(20) + 10)}
			if send(remote, Frame{"potato", host, data}, nil) {
				break
			}
		}
	}
}
*/
