//
//  Task ventilator.
//  Binds PUSH socket to tcp://localhost:5557
//  Sends batch of tasks to workers via that socket
//

package main

import (
	zmq "github.com/pebbe/zmq4"
	rest "github.com/ungerik/go-rest"

	"fmt"
	"encoding/json"
	"strconv"
)

var sender, receiver *zmq.Socket;

func main(){
	//  Prepare our socket
	receiver, _ = zmq.NewSocket(zmq.PULL)
	defer receiver.Close()
	receiver.Bind("tcp://*:5558")
	//  Socket to send messages on
	sender, _ = zmq.NewSocket(zmq.PUSH)
	defer sender.Close()
	sender.Bind("tcp://*:5557")

	rest.DontCheckRequestMethod = true
	rest.HandlePost("/getPrimeFactors", getPrimeFactors)
	stopServerChan := make(chan bool)
	rest.HandleGet("/close", func() string {
			stopServerChan <- true
			return "stoping server..."
		})
	rest.RunServer("0.0.0.0:88", stopServerChan)
}

func getPrimeFactors(work *PrimeWorkload) string {
	resultChan := make(chan string);
	resultCountChan := make(chan int);

	go receive(work, resultChan, resultCountChan);
	go sendOut(work, resultCountChan)

	result := <- resultChan

	return result
}

func sendOut(work *PrimeWorkload, resultCount chan int) {

	fmt.Println("Sending tasks to workers...")

	num := work.BigInt()

	max := num.Int64()/2
	var possFactor int64 = 2
	requestIsPrime(possFactor)
	count := 1
	for possFactor = 3; possFactor < max; possFactor+=2 {
		requestIsPrime(possFactor)
		count++
	}
	resultCount <- count
}

func requestIsPrime(possFactor int64){
	fmt.Println("sending request for ", possFactor)
	var workload PrimeWorkload;
	workload.Number = strconv.FormatInt(possFactor,10);
	bytes,_ := json.Marshal(workload);
	sender.SendBytes(bytes, 0)
}

func receive(work *PrimeWorkload, resultChan chan string, resultCountChan chan int) {

	resultCount := <- resultCountChan
	stringResult := "Factors: "
	num := work.BigInt().Int64()

	//  Process 100 confirmations
	for i := 0; i < resultCount; i++ {
		bytes,_ := receiver.RecvBytes(0)
		var workResult PrimeWorkload;
		json.Unmarshal(bytes, &workResult)
		factor := workResult.BigInt().Int64()

		if (workResult.IsPrime){
			for ;num%factor == 0; {
				num=num/factor
				stringResult += strconv.FormatInt(factor,10)+"*"
			}
		}
	}

	//  Calculate and report duration of batch
	resultChan <- stringResult+"1"
}
