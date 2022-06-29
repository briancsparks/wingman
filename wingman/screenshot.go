package wingman

import (
  ss "github.com/kbinani/screenshot"
  "github.com/lxn/win"
  "image"
  "image/png"
  "os"
)

/* Copyright © 2022 Brian C Sparks <briancsparks@gmail.com> -- MIT (see LICENSE file) */


// --------------------------------------------------------------------------------------------------------------------

func SaveDisplay(name string) error {
  return SaveDisplayN(0, name)
}

// --------------------------------------------------------------------------------------------------------------------

func SaveDisplayN(n int, name string) error {
  img, err := ss.CaptureDisplay(n)
  if err != nil {
    return err
  }

  return save(img, name)
}

// --------------------------------------------------------------------------------------------------------------------

func SaveWindow(hwnd win.HWND, name string) error {
  img, err := CaptureWindow(hwnd)
  if err != nil {
    return err
  }
  return save(img, name)
}

// --------------------------------------------------------------------------------------------------------------------

func CaptureWindow(hwnd win.HWND) (*image.RGBA, error) {
  rect, err := GetWindowRect(hwnd)
  if err != nil {
    return nil, err
  }

  img, err := ss.Capture(rect.Min.X, rect.Min.Y, rect.Max.X, rect.Max.Y)
  if err != nil {
    return nil, err
  }
  return img, nil
}

// --------------------------------------------------------------------------------------------------------------------

func save(img *image.RGBA, name string) error {
  file, err := os.Create(name)
  if err != nil {
    return err
  }
  defer file.Close()

  png.Encode(file, img)
  return nil
}

