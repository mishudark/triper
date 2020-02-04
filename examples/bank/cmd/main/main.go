package main

import (
	"flag"
	"os"
	"time"

	"github.com/golang/glog"
	"github.com/mishudark/triper"
	"github.com/mishudark/triper/examples/bank"
)

func main() {
	flag.Parse()

	commandBus, err := getConfig()
	if err != nil {
		glog.Infoln(err)
		os.Exit(1)
	}

	end := make(chan bool)

	//Create Account
	for i := 0; i < 3; i++ {
		go func() {
			uuid := triper.GenerateUUID()

			//1) Create an account
			account := bank.CreateAccount{
				Owner: "mishudark",
			}

			account.Type = "create_account"

			commandBus.HandleCommand(&account)
			glog.Infof("account %s - account created", uuid)

			//2) Perform a deposit
			time.Sleep(time.Millisecond * 100)
			deposit := bank.PerformDeposit{
				Amount: 300,
			}

			deposit.AggregateID = uuid
			deposit.Version = 1

			commandBus.HandleCommand(&deposit)
			glog.Infof("account %s - deposit performed", uuid)

			//3) Perform a withdrawl
			time.Sleep(time.Millisecond * 100)
			withdrawl := bank.PerformWithdrawal{
				Amount: 249,
			}

			withdrawl.AggregateID = uuid
			withdrawl.Version = 2

			commandBus.HandleCommand(&withdrawl)
			glog.Infof("account %s - withdrawl performed", uuid)
		}()
	}
	<-end
}
