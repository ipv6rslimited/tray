// +build windows

/*
**
** nowindow_windows
** Disables window from opening on windows
**
** https://github.com/fyne-io/fyne/issues/3156#issuecomment-1295732800
*/

package tray

import (
  "os/exec"
  "golang.org/x/sys/windows"
)

func SetCommandNoWindow(cmd *exec.Cmd) {
  cmd.SysProcAttr = &windows.SysProcAttr{
    CreationFlags: windows.CREATE_NO_WINDOW,
  }
}
