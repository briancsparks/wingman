package wingman

/* Copyright © 2022 Brian C Sparks <briancsparks@gmail.com> -- MIT (see LICENSE file) */

import (
  "fmt"
  "github.com/lxn/win"
  "log"
)

//GetActiveWindow

type WindowLocation int32

const (
  MainM1 WindowLocation = iota

  FullishM1
  FullSharedM1
  LeftMarginishM1

  FullishM2
  FullSharedM2
  LeftMarginishM2
)

// ====================================================================================================================
// Main External interface

func MoveWindowDir(hwnd win.HWND, dir ScreenDir, count int32) error {
  var rect win.RECT
  success := win.GetWindowRect(hwnd, &rect)
  if !success {
    return fmt.Errorf("could not move window to %v", rect)
  }

  switch dir {
  case ScreenLeft:
    rect.Left -= count
    rect.Right -= count
  case ScreenTop:
    rect.Top -= count
    rect.Bottom -= count
  case ScreenRight:
    rect.Right += count
    rect.Left += count
  case ScreenBottom:
    rect.Top += count
    rect.Bottom += count
  }

  return MoveWindowTo(hwnd, rect)
}

func MoveActiveWindowDir(dir ScreenDir, count int32) error {
  hwnd := win.GetForegroundWindow()
  return MoveWindowDir(hwnd, dir, count)
}

func ExpandWindowDir(hwnd win.HWND, dir ScreenDir, count int32) error {
  var rect win.RECT
  success := win.GetWindowRect(hwnd, &rect)
  if !success {
    return fmt.Errorf("could not move window to %v", rect)
  }

  switch dir {
  case ScreenLeft:
    rect.Left -= count
  case ScreenTop:
    rect.Top -= count
  case ScreenRight:
    rect.Right += count
  case ScreenBottom:
    rect.Bottom += count
  }

  return MoveWindowTo(hwnd, rect)
}

func ExpandActiveWindowDir(dir ScreenDir, count int32) error {
  hwnd := win.GetForegroundWindow()
  return ExpandWindowDir(hwnd, dir, count)
}

func MoveWindowTo(hwnd win.HWND, rect win.RECT) error {
  success := SetWindowRect(hwnd, rect)
  if !success {
    return fmt.Errorf("could not move window to %v", rect)
  }
  return nil
}

func MoveActiveWindowTo(rect win.RECT) error {
  hwnd := win.GetForegroundWindow()
  return MoveWindowTo(hwnd, rect)
}

// --------------------------------------------------------------------------------------------------------------------

func MoveWindowClassTo(nameToFind, loc string) {
  fmt.Printf("MWCT %s, %s\n", nameToFind, loc)
  MoveWindowClassToLoc(nameToFind, LocationByInitials(loc))
}

// ====================================================================================================================
// Other high-level interfaces

// --------------------------------------------------------------------------------------------------------------------

func MoveWindowToLoc(hwnd win.HWND, loc WindowLocation) {
  SetWindowRect(hwnd, location(loc))
}

// --------------------------------------------------------------------------------------------------------------------

func MoveWindowClassToLoc(nameToFind string, loc WindowLocation) {
  hwnd, err := FindWindowByClassName(nameToFind)
  if err != nil {
    log.Panic(err)
  }

  fmt.Printf("Move window to 0x%x: %v\n", hwnd, loc)
  MoveWindowToLoc(hwnd, loc)
}

// ====================================================================================================================
// Might be useful, so exported

// TODO: When moving window to the other monitor, the first move is always wrong.

// --------------------------------------------------------------------------------------------------------------------

func SetWindowPos(hwnd win.HWND, x, y, dx, dy int32) bool {
  if x == 0 || y == 0 || dx == 0 || dy == 0 {
    fmt.Printf("NOT setting windows pos for %x: Rect: x: %v, y: %v, dx: %v, dy: %v\n", hwnd, x, y, dx, dy)
    return false
  }

  fmt.Printf("setting windows pos for %v (0x%x): Rect: x: %v, y: %v, dx: %v, dy: %v\n", hwnd, hwnd, x, y, dx, dy)
  return win.SetWindowPos(hwnd, 0, x, y, dx, dy, 0)
}

// --------------------------------------------------------------------------------------------------------------------

func SetWindowRect(hwnd win.HWND, rect win.RECT) bool {
  return SetWindowPos(hwnd, rect.Left, rect.Top, Width(rect), Height(rect))
}

func Width(r win.RECT) int32 {
  return r.Right - r.Left
}

func Height(r win.RECT) int32 {
  return r.Bottom - r.Top
}

func LocationByInitials(loc string) WindowLocation {
  switch loc {
  case "m":
    return MainM1
  case "f1":
    return FullishM1
  case "s1":
    return FullSharedM1
  case "l1":
    return LeftMarginishM1
  case "f2":
    return FullishM2
  case "s2":
    return FullSharedM2
  case "l2":
    return LeftMarginishM2
  }

  return MainM1
}

func location(loc WindowLocation) win.RECT {
  switch loc {

  // "m"
  case MainM1:
    //(288, 132)-(1828, 1019), 1540x887
    //(8, 51)-(1532, 879), 1524x828
    return win.RECT{Left: 288, Top: 132, Right: 1828, Bottom: 1019}

  // "f1"
  case FullishM1:
    //(223, 11)-(1910, 1039), 1687x1028
    //(8, 0)-(1679, 1020), 1671x1020
    return win.RECT{Left: 223, Top: 11, Right: 1910, Bottom: 1039}

  // "s1"
  case FullSharedM1:
    return win.RECT{Left: 223, Top: 11, Right: 1910, Bottom: 1039}

  // "l1"
  case LeftMarginishM1:
    //(28, 147)-(1156, 1050), 1128x903
    //(8, 31)-(1120, 895), 1112x864
    return win.RECT{Left: 28, Top: 147, Right: 1156, Bottom: 1050}

  // "f2"
  case FullishM2:
    //(-3486, -1058)-(-777, 607), 2709x1665
    //(7, 0)-(2701, 1657), 2694x1657
    return win.RECT{Left: -3486, Top: -1058, Right: -777, Bottom: 607}

  // "s2"
  case FullSharedM2:
    return win.RECT{Left: -3486, Top: -1058, Right: -777, Bottom: 607}

  // "l2"
  case LeftMarginishM2:
    //(-3798, -513)-(-2740, 616), 1058x1129
    //(7, 31)-(1051, 1383), 1044x1352
    return win.RECT{Left: -3798, Top: -513, Right: -2740, Bottom: 616}
  }

  // FullSharedM1
  return location(FullSharedM1)
}
