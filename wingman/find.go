package wingman

/* Copyright © 2022 Brian C Sparks <briancsparks@gmail.com> -- MIT (see LICENSE file) */

import (
  "github.com/lxn/win"
  "syscall"
)

// --------------------------------------------------------------------------------------------------------------------

func FindWindowByClassName(nameToFind string) (win.HWND, error) {

  name, err := syscall.UTF16PtrFromString(nameToFind)
  if err != nil {
    //log.Panic(err)
    return 0, err
  }

  hwnd := win.FindWindow(name, nil)
  return hwnd, nil
}
