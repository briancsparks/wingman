package wingman

import (
  "fmt"
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

func RunKbMove() error {
  // Buffer size is depends on your need. The 100 is placeholder value.
  keyboardChan := make(chan types.KeyboardEvent, 100)

  if err := keyboard.Install(handler, keyboardChan); err != nil {
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

      //case k := <-keyboardChan:
      //  //fmt.Printf("Keyb: %-12v vkcode: %-12v sccode: %4v t: %10v\n", k.Message, k.VKCode, k.ScanCode, k.Time)
      //  fmt.Printf("%-12v %-12v (%v)\n", k.Message, k.VKCode, k.ScanCode)
      //  continue

    }
  }

}

// --------------------------------------------------------------------------------------------------------------------
// See: https://github.com/moutend/go-hook/blob/develop/examples/swapkeys/main.go

func handler(c chan<- types.KeyboardEvent) types.HOOKPROC {
  trigger1 := false
  trigger2 := false
  mode := normal

  return func(code int32, wParam, lParam uintptr) uintptr {
    var vkcode types.VKCode
    var scancode uint32
    var message types.Message
    callNext := true

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
        }
      }

      if trigger1 {
        if vkcode == types.VK_RCONTROL && isUp(message) {
          trigger1 = false /* cancel */
        } else {
          if vkcode == types.VK_LCONTROL && isUp(message) {
            trigger2 = true
          }
        }
      }

      if trigger2 {
        if vkcode == types.VK_RCONTROL && isUp(message) {
          mode = engaged
          trigger1 = false
          trigger2 = false
          callNext = true
          goto NEXT
        }
        // Could put 'else' here and fizzle the trigger if it is not right-control key
      }

    } else if mode == engaged {
      if vkcode == types.VK_ESCAPE && isUp(message) {
        mode = normal
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
          //ctrl := isCtrlDown(vkcode, message)
          _ = isCtrlDown(vkcode, message)
          _ = isShiftDown(vkcode, message)
          _ = isWinKeyDown(vkcode, message)
          _ = isAltDown(vkcode, message)
          _ = isMenuKeyDown(vkcode, message)
        }()
      }
    }

  NEXT:
    fmt.Printf("%-13v %-13v (%v); wproc: %v, w: %v, t1: %v, t1: %v, m: %v\n", message, vkcode, scancode, code, wParam, trigger1, trigger2, mode)

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

func isCtrlDown(k types.VKCode, m types.Message) bool {
  return isCtrl(k) && isDown(m)
}

func isShiftDown(k types.VKCode, m types.Message) bool {
  return isShift(k) && isDown(m)
}

func isWinKeyDown(k types.VKCode, m types.Message) bool {
  return isWinKey(k) && isDown(m)
}

func isAltDown(k types.VKCode, m types.Message) bool {
  return isAlt(k) && isDown(m)
}

func isMenuKeyDown(k types.VKCode, m types.Message) bool {
  return isMenuKey(k) && isDown(m)
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
