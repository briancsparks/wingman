package cmd

/* Copyright © 2022 Brian C Sparks <briancsparks@gmail.com> -- MIT (see LICENSE file) */

import (
  "fmt"
  "github.com/briancsparks/wingman/wingman"
  "github.com/lxn/win"
  "github.com/spf13/cobra"
)

// tttCmd represents the ttt command
var tttCmd = &cobra.Command{
  Use:   "ttt",
  Short: "Run the currently-in-development command",
  Long: `You can

Test what yer working on.`,
  Run: func(cmd *cobra.Command, args []string) {
    fmt.Println("ttt called")
    DoTtt()
  },
}

func DoTtt() {
  wingman.ThreeTwoOneGo()
  _ = wingman.MoveActiveWindowDir(wingman.ScreenLeft, 20)
}

func DoTtt3() {
  rect := win.RECT{Left: 288, Top: 132, Right: 1828, Bottom: 1019} // Main
  //rect := win.RECT{Left: 28, Top: 147, Right: 1156, Bottom: 1050} // Left
  //hwnd, _ := wingman.FindWindowByClassName("Notepad")

  wingman.ThreeTwoOneGo()
  _ = wingman.MoveActiveWindowTo(rect)
}

func DoTtt2() {
  rect := win.RECT{Left: 288, Top: 132, Right: 1828, Bottom: 1019} // Main
  //rect := win.RECT{Left: 28, Top: 147, Right: 1156, Bottom: 1050} // Left
  //hwnd, _ := wingman.FindWindowByClassName("Notepad")

  wingman.ThreeTwoOneGo()
  _ = wingman.MoveActiveWindowTo(rect)
}

func DoTtt1() {
  rect := win.RECT{Left: 288, Top: 132, Right: 1828, Bottom: 1019} // Main
  //rect := win.RECT{Left: 28, Top: 147, Right: 1156, Bottom: 1050} // Left
  hwnd, _ := wingman.FindWindowByClassName("Notepad")
  _ = wingman.MoveWindowTo(hwnd, rect)
}

func init() {
  rootCmd.AddCommand(tttCmd)

  // Here you will define your flags and configuration settings.

  // Cobra supports Persistent Flags which will work for this command
  // and all subcommands, e.g.:
  // tttCmd.PersistentFlags().String("foo", "", "A help for foo")

  // Cobra supports local flags which will only run when this command
  // is called directly, e.g.:
  // tttCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
