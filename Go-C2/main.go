package main 

import (
  "fmt"
  "os"

  "github.com/DUDLEYDANIEL/GolemC2/cmd"
)

func main(){
  if len(os.Args) < 2 {
    fmt.Println("Usage: GolemC2 [server|agent]")
    return
  }

  switch os.Args[1]{
  case "server":
    cmd.RunServer()
  case "agent":
    cmd.RunAgent("https://localhost:8443")
  default:
    fmt.Println("Unknown Command")
  }
}
