// +build !windows

/*
**
** nowindow_others
** Dummy function for non windows systems
**
** Distributed under the COOL License.
**
** Copyright (c) 2024 IPv6.rs <https://ipv6.rs>
** All Rights Reserved
**
*/

package tray

import "os/exec"

func SetCommandNoWindow(cmd *exec.Cmd) {
}
