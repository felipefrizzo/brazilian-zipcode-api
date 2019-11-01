package main 

import (
	"github.com/felipefrizzo/brazilian-zipcode-api/cmd"
)

func main()  {
	server := &server.Server{}
	server.Initialize()
	server.Run(":8000")
}