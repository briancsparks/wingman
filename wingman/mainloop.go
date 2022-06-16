package wingman

import (
  _ "embed"
  "fmt"
  "github.com/getlantern/systray"
  "github.com/hajimehoshi/ebiten/v2"
  "github.com/moutend/go-hook/pkg/keyboard"
  "github.com/moutend/go-hook/pkg/types"
  "github.com/moutend/go-hook/pkg/win32"
  "log"
  "os"
  "os/signal"
  "time"
  "unsafe"
)

/* Copyright © 2022 Brian C Sparks <briancsparks@gmail.com> -- MIT (see LICENSE file) */

//go:embed assets/kb-arrows-white2.ico
var whiteArrows []byte

//go:embed assets/kb-arrows-black.ico
var blackArrows []byte

//go:embed assets/kb-arrows-teal.ico
var tealArrows []byte

//go:embed assets/kb-arrows-red.ico
var redArrows []byte

//go:embed assets/kb-arrows-blue1.ico
var blue1Arrows []byte

//go:embed assets/kb-arrows-blue2.ico
var blue2Arrows []byte

type kbmode int

const (
  normal kbmode = iota
  engaged
)

func (km kbmode) String() string {
  return [...]string{"normal", "engaged"}[km]
}

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

var finish chan struct{}
var finishedIt func()

func RunKbMove0() error {
  // Setup system tray
  systray.Run(RunKbMove, onExit)
  return nil
}

func RunKbMove() {
  // Finish signal
  finish := make(chan struct{}, 1)
  finished := false
  finishedIt = func() {
    if !finished {
      finished = true
      close(finish)
      fmt.Printf("finishedIt - close\n")
    }
    fmt.Printf("finishedIt\n")
  }
  defer finishedIt()

  // Setup keyboard hook
  keyboardChan := make(chan types.KeyboardEvent, 100)
  if err := keyboard.Install(handler, keyboardChan); err != nil {
    fmt.Printf("kb.install\n")
    //return err
    return
  }
  defer keyboard.Uninstall()

  // Handle Ctrl+C
  signalChan := make(chan os.Signal, 1)
  signal.Notify(signalChan, os.Interrupt)

  //fmt.Println("start capturing keyboard input")

  systray.SetIcon(tealArrows)
  fmt.Printf("icon\n")
  quit := systray.AddMenuItem("Quit", "Quit Wingman")

  //// Main loop
  //wg := sync.WaitGroup{}
  //wg.Add(1)

  //go func() {
  //  defer wg.Done()

  for {
    select {
    case <-time.After(1 * time.Minute):
      fmt.Println("Timeout")
      systray.Quit()
      finishedIt()
      return
      //return nil

    case <-signalChan:
      fmt.Println("Shutdown")
      systray.Quit()
      finishedIt()
      //time.Sleep(5 * time.Second)
      return
      //return nil

      //case k := <-keyboardChan:
      //  //fmt.Printf("Keyb: %-12v vkcode: %-12v sccode: %4v t: %10v\n", k.Message, k.VKCode, k.ScanCode, k.Time)
      //  fmt.Printf("%-12v %-12v (%v)\n", k.Message, k.VKCode, k.ScanCode)
      //  continue

    case <-finish:
      fmt.Printf("finish in onReady\n")
      systray.Quit()
      finishedIt()
      return

    case <-quit.ClickedCh:
      fmt.Printf("quit menu clicked\n")
      finishedIt()
      return

    }
  }
  //}()

  //wg.Wait()
}

// --------------------------------------------------------------------------------------------------------------------

func onReady() {
  systray.SetIcon(tealArrows)
  fmt.Printf("icon\n")
  //quit := systray.AddMenuItem("Quit", "Quit Wingman")
  _ = systray.AddMenuItem("Quit", "Quit Wingman")

  //go func() {
  //  fmt.Printf("one\n")
  //  <-finish
  //  fmt.Printf("two\n")
  //}()

  go func() {
    for {
      fmt.Printf("onReady for\n")
      select {
      //case <-quit.ClickedCh:
      //  fmt.Printf("quit menu clicked\n")
      //  finishedIt()
      //  goto ALLALLA
      //  //break
      //  //return

      case <-finish:
        fmt.Printf("finish in onReady\n")
        systray.Quit()
        goto ALLALLA
        //break
        //return
      }
    }
  ALLALLA:
    fmt.Printf("onReady for func done\n")
  }()

  //go func() {
  //  for {
  //    systray.SetTooltip("Ok - No Proxy")
  //    systray.SetIcon(blackArrows)
  //    systray.SetTooltip("Engaged")
  //
  //    time.Sleep(2 * time.Second)
  //  }
  //}()
}

// --------------------------------------------------------------------------------------------------------------------

func onExit() {
  fmt.Printf("onExit\n")
}

// --------------------------------------------------------------------------------------------------------------------
// See: https://github.com/moutend/go-hook/blob/develop/examples/swapkeys/main.go

func handler(c chan<- types.KeyboardEvent) types.HOOKPROC {
  trigger1 := false
  trigger2 := false
  mode := normal
  ctrlDown := false
  shiftDown := false
  winkeyDown := false
  altDown := false
  menuKeyDown := false

  return func(code int32, wParam, lParam uintptr) uintptr {
    var vkcode types.VKCode
    var scancode uint32
    var message types.Message
    callNext := true
    debugMsg := ""

    if lParam == 0 {
      goto NEXT
    }

    vkcode = (*types.KBDLLHOOKSTRUCT)(unsafe.Pointer(lParam)).VKCode
    scancode = (*types.KBDLLHOOKSTRUCT)(unsafe.Pointer(lParam)).ScanCode
    message = types.Message(wParam)

    c <- types.KeyboardEvent{
      Message:         types.Message(wParam),
      KBDLLHOOKSTRUCT: *(*types.KBDLLHOOKSTRUCT)(unsafe.Pointer(lParam)),
    }

    // ----- Engeging code
    if mode == normal {
      if !trigger1 {
        if vkcode == types.VK_RCONTROL && isDown(message) {
          trigger1 = true
          systray.SetIcon(blue1Arrows)
        }
      }

      if trigger1 {
        if vkcode == types.VK_RCONTROL && isUp(message) {
          trigger1 = false /* cancel */
          systray.SetIcon(redArrows)
        } else {
          if vkcode == types.VK_LCONTROL && isUp(message) {
            trigger2 = true
            systray.SetIcon(blue2Arrows)
          }
        }
      }

      if trigger2 {
        if vkcode == types.VK_RCONTROL && isUp(message) {
          mode = engaged
          trigger1 = false
          trigger2 = false
          systray.SetIcon(whiteArrows)

          callNext = true
          goto NEXT
        }
        // Could put 'else' here and fizzle the trigger if it is not right-control key
      }

    } else if mode == engaged {
      if vkcode == types.VK_ESCAPE && isUp(message) {
        mode = normal
        systray.SetIcon(tealArrows)
      }
    }

    // ----- Engaged
    if mode == engaged {
      callNext = false

      if vkcode == types.VK_LMENU || vkcode == types.VK_TAB {
        callNext = true

      } else {
        // Use an immediately-called function to avoid the 'goto NEXT' jumping over declarations
        func() {
          ctrlDown = isCtrlDown(vkcode, message, ctrlDown)
          shiftDown = isShiftDown(vkcode, message, shiftDown)
          winkeyDown = isWinKeyDown(vkcode, message, winkeyDown)
          altDown = isAltDown(vkcode, message, altDown)
          menuKeyDown = isMenuKeyDown(vkcode, message, menuKeyDown)

          if ctrlDown {
            debugMsg += "^"
          }

          if shiftDown {
            debugMsg += "$"
          }

          if winkeyDown {
            debugMsg += "#"
          }

          if altDown {
            debugMsg += "@"
          }

          if menuKeyDown {
            debugMsg += "="
          }
          debugMsg += " "

          if vkcode == types.VK_LEFT && ctrlDown {
            _ = MoveActiveWindowDir(ScreenLeft, 20)
          } else if vkcode == types.VK_LEFT {
            _ = MoveActiveWindowDir(ScreenLeft, 1)

          } else if vkcode == types.VK_UP && ctrlDown {
            _ = MoveActiveWindowDir(ScreenTop, 20)
          } else if vkcode == types.VK_UP {
            _ = MoveActiveWindowDir(ScreenTop, 1)

          } else if vkcode == types.VK_RIGHT && ctrlDown {
            _ = MoveActiveWindowDir(ScreenRight, 20)
          } else if vkcode == types.VK_RIGHT {
            _ = MoveActiveWindowDir(ScreenRight, 1)

          } else if vkcode == types.VK_DOWN && ctrlDown {
            _ = MoveActiveWindowDir(ScreenBottom, 20)
          } else if vkcode == types.VK_DOWN {
            _ = MoveActiveWindowDir(ScreenBottom, 1)
          }

        }()
      }
    }

  NEXT:
    fmt.Printf("%-13v %-13v (%v); wproc: %v, w: %v, t1: %v, t1: %v, m: %v -- %v\n", message, vkcode, scancode, code, wParam, trigger1, trigger2, mode, debugMsg)

    // Do not call this, if you want to eat the key, otherwise the active window will get it.
    if callNext {
      return win32.CallNextHookEx(0, code, wParam, lParam)
    }
    return 1
  }
}

// --------------------------------------------------------------------------------------------------------------------

func isDown(m types.Message) bool {
  return m == types.WM_KEYDOWN || m == types.WM_SYSKEYDOWN
}

func isUp(m types.Message) bool {
  return m == types.WM_KEYUP || m == types.WM_SYSKEYUP
}

func isCtrlDown(k types.VKCode, m types.Message, current bool) bool {
  if isCtrl(k) {
    return isDown(m)
  }
  return current
}

func isShiftDown(k types.VKCode, m types.Message, current bool) bool {
  if isShift(k) {
    return isDown(m)
  }
  return current
}

func isWinKeyDown(k types.VKCode, m types.Message, current bool) bool {
  if isWinKey(k) {
    return isDown(m)
  }
  return current
}

func isAltDown(k types.VKCode, m types.Message, current bool) bool {
  if isAlt(k) {
    return isDown(m)
  }
  return current
}

func isMenuKeyDown(k types.VKCode, m types.Message, current bool) bool {
  if isMenuKey(k) {
    return isDown(m)
  }
  return current
}

func isCtrl(k types.VKCode) bool {
  return k == types.VK_LCONTROL || k == types.VK_RCONTROL
}

func isShift(k types.VKCode) bool {
  return k == types.VK_LSHIFT || k == types.VK_RSHIFT
}

func isWinKey(k types.VKCode) bool {
  return k == types.VK_LWIN || k == types.VK_RWIN
}

func isAlt(k types.VKCode) bool {
  return k == types.VK_LMENU || k == types.VK_RMENU
}

func isMenuKey(k types.VKCode) bool {
  return k == types.VK_APPS
}

//VK_LMENU               VKCode = 0xA4 // Left MENU key
//VK_RMENU               VKCode = 0xA5 // Right MENU key
//VK_LWIN                VKCode = 0x5B // Left Windows key (Natural keyboard)
//VK_RWIN                VKCode = 0x5C // Right Windows key (Natural keyboard)
//VK_APPS                VKCode = 0x5D // Applications key (Natural keyboard)

//WM_KEYDOWN    VK_LEFT       (75); wproc: 0, w: 256, t1: false, t1: false, m: normal
//WM_KEYUP      VK_LEFT       (75); wproc: 0, w: 257, t1: false, t1: false, m: normal
//WM_KEYDOWN    VK_UP         (72); wproc: 0, w: 256, t1: false, t1: false, m: normal
//WM_KEYUP      VK_UP         (72); wproc: 0, w: 257, t1: false, t1: false, m: normal
//WM_KEYDOWN    VK_RIGHT      (77); wproc: 0, w: 256, t1: false, t1: false, m: normal
//WM_KEYUP      VK_RIGHT      (77); wproc: 0, w: 257, t1: false, t1: false, m: normal
//WM_KEYDOWN    VK_DOWN       (80); wproc: 0, w: 256, t1: false, t1: false, m: normal
//WM_KEYUP      VK_DOWN       (80); wproc: 0, w: 257, t1: false, t1: false, m: normal
//WM_KEYDOWN    VK_LCONTROL   (29); wproc: 0, w: 256, t1: false, t1: false, m: normal
//WM_KEYDOWN    VK_C          (46); wproc: 0, w: 256, t1: false, t1: false, m: normal

//WM_KEYDOWN    VK_F1         (59); wproc: 0, w: 256, t1: false, t1: false, m: normal
//WM_KEYUP      VK_F1         (59); wproc: 0, w: 257, t1: false, t1: false, m: normal
//WM_KEYDOWN    VK_F10        (68); wproc: 0, w: 256, t1: false, t1: false, m: normal
//WM_KEYUP      VK_F10        (68); wproc: 0, w: 257, t1: false, t1: false, m: normal
//WM_KEYDOWN    VK_F11        (87); wproc: 0, w: 256, t1: false, t1: false, m: normal
//WM_KEYUP      VK_F11        (87); wproc: 0, w: 257, t1: false, t1: false, m: normal
//WM_KEYDOWN    VK_F12        (88); wproc: 0, w: 256, t1: false, t1: false, m: normal
//WM_KEYUP      VK_F12        (88); wproc: 0, w: 257, t1: false, t1: false, m: normal
//WM_KEYDOWN    VK_MEDIA_PLAY_PAUSE (0); wproc: 0, w: 256, t1: false, t1: false, m: normal
//WM_KEYUP      VK_MEDIA_PLAY_PAUSE (0); wproc: 0, w: 257, t1: false, t1: false, m: normal
//WM_KEYDOWN    VK_MEDIA_PREV_TRACK (0); wproc: 0, w: 256, t1: false, t1: false, m: normal
//WM_KEYUP      VK_MEDIA_PREV_TRACK (0); wproc: 0, w: 257, t1: false, t1: false, m: normal
//WM_KEYDOWN    VK_MEDIA_NEXT_TRACK (0); wproc: 0, w: 256, t1: false, t1: false, m: normal
//WM_KEYUP      VK_MEDIA_NEXT_TRACK (0); wproc: 0, w: 257, t1: false, t1: false, m: normal
//WM_KEYDOWN    VK_VOLUME_MUTE (0); wproc: 0, w: 256, t1: false, t1: false, m: normal
//WM_KEYUP      VK_VOLUME_MUTE (0); wproc: 0, w: 257, t1: false, t1: false, m: normal
//WM_KEYDOWN    VK_VOLUME_DOWN (0); wproc: 0, w: 256, t1: false, t1: false, m: normal
//WM_KEYUP      VK_VOLUME_DOWN (0); wproc: 0, w: 257, t1: false, t1: false, m: normal
//WM_KEYDOWN    VK_VOLUME_UP  (0); wproc: 0, w: 256, t1: false, t1: false, m: normal
//WM_KEYUP      VK_VOLUME_UP  (0); wproc: 0, w: 257, t1: false, t1: false, m: normal
//WM_KEYDOWN    VK_CAPITAL    (58); wproc: 0, w: 256, t1: false, t1: false, m: normal
//WM_KEYUP      VK_CAPITAL    (58); wproc: 0, w: 257, t1: false, t1: false, m: normal
//WM_KEYDOWN    VK_CAPITAL    (58); wproc: 0, w: 256, t1: false, t1: false, m: normal
//WM_KEYUP      VK_CAPITAL    (58); wproc: 0, w: 257, t1: false, t1: false, m: normal

// --------------------------------------------------------------------------------------------------------------------
// See: https://github.com/moutend/go-hook/blob/develop/examples/swapkeys/main.go

func handlerOrig(c chan<- types.KeyboardEvent) types.HOOKPROC {
  counter := 0

  return func(code int32, wParam, lParam uintptr) uintptr {
    if lParam == 0 {
      goto NEXT
    }

    c <- types.KeyboardEvent{
      Message:         types.Message(wParam),
      KBDLLHOOKSTRUCT: *(*types.KBDLLHOOKSTRUCT)(unsafe.Pointer(lParam)),
    }

    switch (*types.KBDLLHOOKSTRUCT)(unsafe.Pointer(lParam)).VKCode {
    case types.VK_A:
      if counter == 1 {
        counter = 0
        goto NEXT
      }
      switch types.Message(wParam) {
      case types.WM_KEYDOWN:
        //go kbB.Press()
      case types.WM_KEYUP:
        //go kbB.Release()
      }

      counter = 1

      return 1
    case types.VK_B:
      if counter == 1 {
        counter = 0
        goto NEXT
      }
      switch types.Message(wParam) {
      case types.WM_KEYDOWN:
        //go kbA.Press()
      case types.WM_KEYUP:
        //go kbA.Release()
      }

      counter = 1

      return 1
    default:
    }

  NEXT:

    return win32.CallNextHookEx(0, code, wParam, lParam)
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
