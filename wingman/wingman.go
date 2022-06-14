package wingman

/* Copyright © 2022 Brian C Sparks <briancsparks@gmail.com> -- MIT (see LICENSE file) */

type ScreenDir int32

const (
  ScreenLeft ScreenDir = iota
  ScreenTop
  ScreenRight
  ScreenBottom
)
