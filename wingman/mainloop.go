package wingman

import (
  "fmt"
  "github.com/hajimehoshi/ebiten/v2"
  "github.com/moutend/go-hook/pkg/keyboard"
  "github.com/moutend/go-hook/pkg/types"
  "log"
  "os"
  "os/signal"
  "time"
)

/* Copyright © 2022 Brian C Sparks <briancsparks@gmail.com> -- MIT (see LICENSE file) */

type Game struct{}

func (g *Game) Update() error {
  return nil
}

func (g *Game) Draw(screen *ebiten.Image) {

}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
  return 320, 240
}

// --------------------------------------------------------------------------------------------------------------------

func RunKbMove() error {
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
      fmt.Println("Timeout")
      return nil

    case <-signalChan:
      fmt.Println("Shutdown")
      return nil

    case k := <-keyboardChan:
      //fmt.Printf("Keyb: %-12v vkcode: %-12v sccode: %4v t: %10v\n", k.Message, k.VKCode, k.ScanCode, k.Time)
      fmt.Printf("%-12v %-12v (%v)\n", k.Message, k.VKCode, k.ScanCode)
      continue

    }
  }

}

// --------------------------------------------------------------------------------------------------------------------

func RunOverlay() {
  game := &Game{}

  ebiten.SetWindowSize(640, 480)
  ebiten.SetWindowTitle("Wingman")

  if err := ebiten.RunGame(game); err != nil {
    log.Fatal(err)
  }
}
