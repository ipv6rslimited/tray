package main

import (
  "fyne.io/fyne/v2/app"
  "github.com/ipv6rslimited/tray"
  "os"
  "fmt"
)

func main() {
  myApp := app.New()

  if len(os.Args) < 2 {
    fmt.Printf("Usage: %s path/to/config.json", os.Args[0])
    os.Exit(1)
  }

  configFilePath := os.Args[1]

  tray.Tray(myApp, configFilePath)

  myApp.Run()
}
