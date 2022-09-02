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
  DoTtt5()
  //fmt.Println("mvw called")
  //wingman.RunKbMove0()
  //fmt.Printf("Run done\n")
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

  tttCmd.AddCommand(ttt1Cmd)
  tttCmd.AddCommand(ttt2Cmd)
  tttCmd.AddCommand(ttt3Cmd)
  tttCmd.AddCommand(ttt4Cmd)
  tttCmd.AddCommand(ttt5Cmd)
}

// --------------------------------------------------------------------------------------------------------------------

// ttt5Cmd represents the ttt5 command
var ttt5Cmd = &cobra.Command{
  Use:   "5",
  Short: "wingman.RunKbMove0 [Run main loop]",

  Run: func(cmd *cobra.Command, args []string) {
    fmt.Println("ttt5 called")
    DoTtt5()
  },
}

func DoTtt5() {
  wingman.RunKbMove0()
}

// --------------------------------------------------------------------------------------------------------------------

// ttt1Cmd represents the ttt1 command
var ttt1Cmd = &cobra.Command{
  Use:   "1",
  Short: "wingman.MoveWindowTo (notepad, main)",

  Run: func(cmd *cobra.Command, args []string) {
    fmt.Println("ttt1 called")
    DoTtt1()
  },
}

func DoTtt1() {
  wr := win.RECT{Left: 288, Top: 132, Right: 1828, Bottom: 1019} // Main
  hwnd, _ := wingman.FindWindowByClassName("Notepad")

  rect := wingman.ImageRectangle(wr)
  //fmt.Printf("ttt1: hwnd: %v, rect: %v\n", hwnd, rect)
  _ = wingman.MoveWindowTo(hwnd, rect)
}

// --------------------------------------------------------------------------------------------------------------------

// ttt2Cmd represents the ttt2 command
var ttt2Cmd = &cobra.Command{
  Use:   "2",
  Short: "wingman.MoveActiveWindowTo (main)",

  Run: func(cmd *cobra.Command, args []string) {
    fmt.Println("ttt2 called")
    DoTtt2()
  },
}

func DoTtt2() {
  wr := win.RECT{Left: 288, Top: 132, Right: 1828, Bottom: 1019} // Main

  rect := wingman.ImageRectangle(wr)
  wingman.ThreeTwoOneGo()
  _ = wingman.MoveActiveWindowTo(rect)
}

// --------------------------------------------------------------------------------------------------------------------

// ttt3Cmd represents the ttt3 command
var ttt3Cmd = &cobra.Command{
  Use:   "3",
  Short: "wingman.MoveActiveWindowTo (main)",

  Run: func(cmd *cobra.Command, args []string) {
    fmt.Println("ttt3 called")
    DoTtt3()
  },
}

func DoTtt3() {
  wr := win.RECT{Left: 288, Top: 132, Right: 1828, Bottom: 1019} // Main

  rect := wingman.ImageRectangle(wr)
  wingman.ThreeTwoOneGo()
  _ = wingman.MoveActiveWindowTo(rect)
}

// --------------------------------------------------------------------------------------------------------------------

// ttt4Cmd represents the ttt4 command
var ttt4Cmd = &cobra.Command{
  Use:   "4",
  Short: "wingman.MoveActiveWindowDir (left)",

  Run: func(cmd *cobra.Command, args []string) {
    fmt.Println("ttt4 called")
    DoTtt4()
  },
}

func DoTtt4() {
  wingman.ThreeTwoOneGo()
  _ = wingman.MoveActiveWindowDir(wingman.ScreenLeft, 20)
}
