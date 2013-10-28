//
//  Task worker.
//  Connects PULL socket to tcp://localhost:5557
//  Collects workloads from ventilator via that socket
//  Connects PUSH socket to tcp://localhost:5558
//  Sends results to sink via that socket
//

package main

import (
	zmq "github.com/pebbe/zmq4"

	"fmt"
	"encoding/json"
	"syscall"
	"strconv"
)

func main() {
	//  Socket to receive messages on
	receiver, _ := zmq.NewSocket(zmq.PULL)
	defer receiver.Close()
	receiver.Connect("tcp://localhost:5557")

	//  Socket to send messages to
	sender, _ := zmq.NewSocket(zmq.PUSH)
	defer sender.Close()
	sender.Connect("tcp://localhost:5558")

	//  Process tasks forever
	for {
		workload := new(PrimeWorkload);
		bytes, _ := receiver.RecvBytes(0)
		json.Unmarshal(bytes, workload);

		//  Simple progress indicator for the viewer
		fmt.Println(workload.Number, strconv.Itoa(syscall.Getpid())+".")

		workload.IsPrime = workload.BigInt().ProbablyPrime(10)

		result,_ := json.Marshal(workload)

		//  Send results to sink
		sender.SendBytes(result, 0)
	}
}
