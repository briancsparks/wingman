package wingman

/* Copyright © 2022 Brian C Sparks <briancsparks@gmail.com> -- MIT (see LICENSE file) */

import (
  "fmt"
  "github.com/lxn/win"
  "image"
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

// Left     Min.X
// Top      Min.Y
// Right    Max.X
// Bottom   Max.Y

// --------------------------------------------------------------------------------------------------------------------

func MoveWindowDir(hwnd win.HWND, dir ScreenDir, count int) error {
  fmt.Printf("MoveWindowDir %v, %v\n", dir, count)

  rect, err := GetWindowRect(hwnd)
  if err != nil {
    return err
  }

  switch dir {
  case ScreenLeft:
    rect.Min.X -= count
    rect.Max.X -= count
  case ScreenTop:
    rect.Min.Y -= count
    rect.Max.Y -= count
  case ScreenRight:
    rect.Max.X += count
    rect.Min.X += count
  case ScreenBottom:
    rect.Min.Y += count
    rect.Max.Y += count
  }

  return MoveWindowTo(hwnd, rect)
}

// -----------------------------------------------------------------------------------------

// MoveActiveWindowDir moves the active window in a certain direction, the given count.
func MoveActiveWindowDir(dir ScreenDir, count int) error {
  hwnd := win.GetForegroundWindow()
  return MoveWindowDir(hwnd, dir, count)
}

// --------------------------------------------------------------------------------------------------------------------

func ExpandWindowDir(hwnd win.HWND, dir ScreenDir, count int) error {
  fmt.Printf("ExpandWindowDir %v, %v\n", dir, count)

  rect, err := GetWindowRect(hwnd)
  if err != nil {
    return err
  }

  switch dir {
  case ScreenLeft:
    rect.Min.X -= count
  case ScreenTop:
    rect.Min.Y -= count
  case ScreenRight:
    rect.Max.X += count
  case ScreenBottom:
    rect.Max.Y += count
  }

  return MoveWindowTo(hwnd, rect)
}

// -----------------------------------------------------------------------------------------

func ExpandActiveWindowDir(dir ScreenDir, count int) error {
  hwnd := win.GetForegroundWindow()
  return ExpandWindowDir(hwnd, dir, count)
}

// --------------------------------------------------------------------------------------------------------------------

func ShrinkWindowDir(hwnd win.HWND, dir ScreenDir, count int) error {
  fmt.Printf("ShrinkWindowDir %v, %v\n", dir, count)

  rect, err := GetWindowRect(hwnd)
  if err != nil {
    return err
  }

  switch dir {
  case ScreenLeft:
    rect.Max.X -= count
  case ScreenTop:
    rect.Max.Y -= count
  case ScreenRight:
    rect.Min.X += count
  case ScreenBottom:
    rect.Min.Y += count
  }

  return MoveWindowTo(hwnd, rect)
}

// -----------------------------------------------------------------------------------------

func ShrinkActiveWindowDir(dir ScreenDir, count int) error {
  hwnd := win.GetForegroundWindow()
  return ShrinkWindowDir(hwnd, dir, count)
}

// --------------------------------------------------------------------------------------------------------------------

func MoveWindowToLoc(hwnd win.HWND, loc WindowLocation) error {
  success := SetWindowRect(hwnd, location(loc))
  if !success {
    return fmt.Errorf("could not move window to %v", loc)
  }
  return nil
}

// -----------------------------------------------------------------------------------------

func MoveActiveWindowToLoc(loc WindowLocation) error {
  hwnd := win.GetForegroundWindow()
  return MoveWindowToLoc(hwnd, loc)
}

// --------------------------------------------------------------------------------------------------------------------

func MoveWindowTo(hwnd win.HWND, ir image.Rectangle) error {
  success := SetWindowRect(hwnd, ir)
  if !success {
    return fmt.Errorf("could not move window to %v", ir)
  }
  return nil
}

// -----------------------------------------------------------------------------------------

func MoveActiveWindowTo(rect image.Rectangle) error {
  hwnd := win.GetForegroundWindow()
  return MoveWindowTo(hwnd, rect)
}

// ====================================================================================================================
// Other high-level interfaces

// --------------------------------------------------------------------------------------------------------------------

func MoveWindowClassTo(nameToFind, loc string) {
  fmt.Printf("MWCT %s, %s\n", nameToFind, loc)
  MoveWindowClassToLoc(nameToFind, LocationByInitials(loc))
}

// --------------------------------------------------------------------------------------------------------------------

func MoveWindowClassToLoc(nameToFind string, loc WindowLocation) {
  hwnd, err := FindWindowByClassName(nameToFind)
  if err != nil {
    log.Panic(err)
  }

  fmt.Printf("Move window to 0x%x: %v\n", hwnd, loc)
  MoveWindowToLocW(hwnd, loc)
}

// ====================================================================================================================
// Might be useful, so exported

// TODO: When moving window to the other monitor, the first move is always wrong.

// --------------------------------------------------------------------------------------------------------------------

func SetWindowPos(hwnd win.HWND, x0, y0, x1, y1 int) bool {
  // TODO: Normalize rect coordinates (x0 < x1, y0 < y1)
  return SetWindowPosW(hwnd, int32(x0), int32(y0), int32(x1 - x0), int32(y1 - y0))
}

// --------------------------------------------------------------------------------------------------------------------

func SetWindowRect(hwnd win.HWND, ir image.Rectangle) bool {
  rect := WinRECT(ir)
  return SetWindowRectW(hwnd, rect)
}

// --------------------------------------------------------------------------------------------------------------------

func GetWindowRect(hwnd win.HWND) (image.Rectangle, error) {
  var rect image.Rectangle
  wrect, err := GetWindowRectW(hwnd)
  if err != nil {
    return rect, err
  }
  rect = image.Rect(int(wrect.Left), int(wrect.Top), int(wrect.Right), int(wrect.Bottom))
  return rect, nil
}

// ====================================================================================================================

// --------------------------------------------------------------------------------------------------------------------

func WinRECT(ir image.Rectangle) win.RECT {
  rect := win.RECT{Left: int32(ir.Min.X), Top: int32(ir.Min.Y), Right: int32(ir.Max.X), Bottom: int32(ir.Max.Y)}
  return rect
}

// --------------------------------------------------------------------------------------------------------------------

func ImageRectangle(wr win.RECT) image.Rectangle {
  rect := image.Rect(int(wr.Left), int(wr.Top), int(wr.Right), int(wr.Bottom))
  return rect
}

// --------------------------------------------------------------------------------------------------------------------

func Width(r win.RECT) int32 {
  return r.Right - r.Left
}

func Height(r win.RECT) int32 {
  return r.Bottom - r.Top
}

// --------------------------------------------------------------------------------------------------------------------

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

// --------------------------------------------------------------------------------------------------------------------

func locationW(loc WindowLocation) win.RECT {
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
  return locationW(FullSharedM1)
}

// --------------------------------------------------------------------------------------------------------------------

func location(loc WindowLocation) image.Rectangle {
  switch loc {

  // "m"
  case MainM1:
    return image.Rect(288, 132, 1828, 1019)

  // "f1"
  case FullishM1:
    return image.Rect(223, 11, 1910, 1039)

  // "s1"
  case FullSharedM1:
    return image.Rect(223, 11, 1910, 1039)

  // "l1"
  case LeftMarginishM1:
    return image.Rect(28, 147, 1156, 1050)

  // "f2"
  case FullishM2:
    return image.Rect(-3486, -1058, -777, 607)

  // "s2"
  case FullSharedM2:
    return image.Rect(-3486, -1058, -777, 607)

  // "l2"
  case LeftMarginishM2:
    return image.Rect(-3798, -513, -2740, 616)
  }

  // FullSharedM1
  return location(FullSharedM1)
}

// --------------------------------------------------------------------------------------------------------------------

func MoveWindowDirW(hwnd win.HWND, dir ScreenDir, count int32) error {
  fmt.Printf("MoveWindowDirW %v, %v\n", dir, count)

  rect, err := GetWindowRectW(hwnd)
  if err != nil {
    return err
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

  return MoveWindowToW(hwnd, rect)
}

// --------------------------------------------------------------------------------------------------------------------

func ExpandWindowDirW(hwnd win.HWND, dir ScreenDir, count int32) error {
  fmt.Printf("ExpandWindowDirW %v, %v\n", dir, count)
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

  return MoveWindowToW(hwnd, rect)
}

// --------------------------------------------------------------------------------------------------------------------

func ShrinkWindowDirW(hwnd win.HWND, dir ScreenDir, count int32) error {
  fmt.Printf("ShrinkWindowDirW %v, %v\n", dir, count)
  var rect win.RECT
  success := win.GetWindowRect(hwnd, &rect)
  if !success {
    return fmt.Errorf("could not move window to %v", rect)
  }

  switch dir {
  case ScreenLeft:
    rect.Right -= count
  case ScreenTop:
    rect.Bottom -= count
  case ScreenRight:
    rect.Left += count
  case ScreenBottom:
    rect.Top += count
  }

  return MoveWindowToW(hwnd, rect)
}

// --------------------------------------------------------------------------------------------------------------------

func MoveWindowToLocW(hwnd win.HWND, loc WindowLocation) error {
  success := SetWindowRectW(hwnd, locationW(loc))
  if !success {
    return fmt.Errorf("could not move window to %v", loc)
  }
  return nil
}

// --------------------------------------------------------------------------------------------------------------------

func MoveWindowToW(hwnd win.HWND, rect win.RECT) error {
  success := SetWindowRectW(hwnd, rect)
  if !success {
    return fmt.Errorf("could not move window to %v", rect)
  }
  return nil
}

// -----------------------------------------------------------------------------------------

func MoveActiveWindowToW(rect win.RECT) error {
  hwnd := win.GetForegroundWindow()
  return MoveWindowToW(hwnd, rect)
}

// TODO: When moving window to the other monitor, the first move is always wrong.

// --------------------------------------------------------------------------------------------------------------------

func SetWindowPosW(hwnd win.HWND, x, y, dx, dy int32) bool {
  if x == 0 || y == 0 || dx == 0 || dy == 0 {
    fmt.Printf("NOT setting windows pos for %x: Rect: x: %v, y: %v, dx: %v, dy: %v\n", hwnd, x, y, dx, dy)
    return false
  }

  fmt.Printf("setting windows pos for %v (0x%x): Rect: x: %v, y: %v, dx: %v, dy: %v\n", hwnd, hwnd, x, y, dx, dy)
  return win.SetWindowPos(hwnd, 0, x, y, dx, dy, 0)
}

// --------------------------------------------------------------------------------------------------------------------

func SetWindowRectW(hwnd win.HWND, rect win.RECT) bool {
  return SetWindowPosW(hwnd, rect.Left, rect.Top, Width(rect), Height(rect))
}

// --------------------------------------------------------------------------------------------------------------------

func GetWindowRectW(hwnd win.HWND) (win.RECT, error) {
  var rect win.RECT
  success := win.GetWindowRect(hwnd, &rect)
  if !success {
    return rect, fmt.Errorf("could not move window to %v", rect)
  }
  return rect, nil
}

