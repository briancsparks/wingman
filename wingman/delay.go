package wingman

import (
  "fmt"
  "time"
)

func ThreeTwoOneGo() {
  LoudCountdown(3)
}

func LoudCountdown(n int) {
  for ; n > 0; n-- {
    fmt.Printf("%d", n)
    if n > 1 {
      fmt.Printf(", ")
    } else if n == 1 {
      fmt.Printf("... ")
    }

    time.Sleep(1 * time.Second)
  }
  fmt.Printf("GO!\n")
}

