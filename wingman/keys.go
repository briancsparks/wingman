package wingman

/* Copyright © 2022 Brian C Sparks <briancsparks@gmail.com> -- MIT (see LICENSE file) */

import (
  "fmt"
  "os"
  "os/signal"
  "time"

  "github.com/moutend/go-hook/pkg/keyboard"
  "github.com/moutend/go-hook/pkg/types"
)

func dokeys() error {
  // Buffer size is depends on your need. The 100 is placeholder value.
  keyboardChan := make(chan types.KeyboardEvent, 100)

  if err := keyboard.Install(nil, keyboardChan); err != nil {
    return err
  }

  defer keyboard.Uninstall()

  signalChan := make(chan os.Signal, 1)
  signal.Notify(signalChan, os.Interrupt)

  fmt.Println("start capturing keyboard input")

  for {
    select {
    case <-time.After(1 * time.Minute):
      fmt.Println("Received timeout signal")
      return nil
    case <-signalChan:
      fmt.Println("Received shutdown signal")
      return nil
    case k := <-keyboardChan:
      fmt.Printf("Received %v %v\n", k.Message, k.VKCode)
      continue
    }
  }

  // not reached
  return nil
}
