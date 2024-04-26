// +build darwin

/*
**
** activation_darwin
** Disables Dock Icon on darwin systems
**
** https://github.com/fyne-io/fyne/issues/3156#issuecomment-1295732800
*/

package tray

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Cocoa
#import <Cocoa/Cocoa.h>

int SetActivationPolicy(void) {
  [NSApp setActivationPolicy:NSApplicationActivationPolicyAccessory];
  return 0;
}
*/
import "C"
import (
  "fmt"
  "time"
  "fyne.io/fyne/v2"
)

func setActivationPolicy() {
  fmt.Println("Setting ActivationPolicy")
  C.SetActivationPolicy()
}

func InitMacSpecific(myApp fyne.App) {
  myApp.Lifecycle().SetOnStarted(func() {
    go func() {
      time.Sleep(200 * time.Millisecond)
      setActivationPolicy()
    }()
  })
}
