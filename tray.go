/*
**
** tray
** Provides a system tray icon generated from a json
**
** Distributed under the COOL License.
**
** Copyright (c) 2024 IPv6.rs <https://ipv6.rs>
** All Rights Reserved
**
*/

package tray

import (
  "runtime"
  "strings"
  "encoding/json"
  "io/ioutil"
  "path/filepath"
  "log"
  "fmt"
  "github.com/fsnotify/fsnotify"
  "os"
  "os/exec"
  "fyne.io/fyne/v2"
  "fyne.io/fyne/v2/driver/desktop"
  "fyne.io/fyne/v2/theme"
)


type MenuItem struct {
  Title string      `json:"title"`
  Exec  string      `json:"exec"`
  Icon  string      `json:"icon"`
  Items []MenuItem  `json:"items"`
}

type TrayConfig struct {
  Icon  string      `json:"icon"`
  Title string      `json:"title"`
  Items []MenuItem  `json:"items"`
}


func Tray(app fyne.App, configPath string) {
  if err := loadAndSetupTray(app, configPath); err != nil {
    log.Fatalf("Failed to set up tray: %v", err)
  }
  initMacSpecific(app)
  watchConfigFile(app, configPath)
}

func loadConfig(filename string) (*TrayConfig, error) {
  bytes, err := ioutil.ReadFile(filename)
  if err != nil {
    return nil, err
  }
  var config TrayConfig
  if err := json.Unmarshal(bytes, &config); err != nil {
    return nil, err
  }
  return &config, nil
}

func executeCommand(command string) {
  if command == "" {
    return
  } else if command == "EXIT" {
    os.Exit(0)
  }

  execDir, err := os.Executable()
  if err != nil {
    log.Printf("Failed to get executable path: %v", err)
  }
  execDir = filepath.Dir(execDir)

  parentDir := filepath.Dir(execDir)

  homeDir, err := os.UserHomeDir()
  if err != nil {
    fmt.Println("Failed to get home directory:", err)
    os.Exit(1)
  }
  var appPath string
  if runtime.GOOS == "windows" {
    appData := os.Getenv("LOCALAPPDATA")
    if appData == "" {
      appData = filepath.Join(homeDir, "AppData", "Local")
    }
    appPath = appData
  } else {
    appPath = homeDir
  }

  command = strings.Replace(command, "_CURRENTPATH", "\"" + execDir + "\"", -1)
  command = strings.Replace(command, "_PARENTPATH", "\"" + parentDir + "\"", -1)
  command = strings.Replace(command, "_HOMEPATH", "\"" + homeDir + "\"", -1)
  command = strings.Replace(command, "_LOCALDATA", "\"" + appPath + "\"", -1)

  var cmd *exec.Cmd

  setCommandNoWindow(cmd)

  if runtime.GOOS == "windows" {
    if strings.HasSuffix(command, ".sh") {
      command = strings.Replace(command, ".sh", ".ps1", -1)
    }
    if strings.HasSuffix(command, ".ps1") {
      cmd = exec.Command("powershell", "-ExecutionPolicy", "Bypass", "-File", command)
    } else {
      cmd = exec.Command("powershell", "-ExecutionPolicy", "Bypass", "-Command", command)
    }
  } else {
    cmd = exec.Command("bash", "-c", command)
  }
  if err := cmd.Start(); err != nil {
    log.Printf("Failed to execute command: %s\n", err)
  } else {
    if err := cmd.Wait(); err != nil {
      log.Printf("Command finished with error: %v", err)
    }
  }
}

func createMenuItems(items []MenuItem) []*fyne.MenuItem {
  var menuItems []*fyne.MenuItem
  for _, item := range items {
    menuItem := fyne.NewMenuItem(item.Title, func() {
      executeCommand(item.Exec)
    })

    if item.Icon != "" {
      menuItem.Icon = iconNameToThemeIcon(item.Icon)
    }

    if len(item.Items) > 0 {
      subMenu := fyne.NewMenu(item.Title, createMenuItems(item.Items)...)
      menuItem.ChildMenu = subMenu
    }
    menuItems = append(menuItems, menuItem)
  }
  return menuItems
}

func loadAndSetupTray(a fyne.App, configPath string) error {
  config, err := loadConfig(configPath)
  if err != nil {
    return fmt.Errorf("failed to load config: %w", err)
  }
  setupTray(a, config)
  return nil
}

func setupTray(a fyne.App, config *TrayConfig) {
  if desk, ok := a.(desktop.App); ok {
    menuItems := createMenuItems(config.Items)
    menu := fyne.NewMenu(config.Title, menuItems...)

    desk.SetSystemTrayMenu(menu)

    icon, err := fyne.LoadResourceFromPath(config.Icon)
    if err != nil {
      log.Printf("Failed to load tray icon: %v", err)
      return
    }

    desk.SetSystemTrayIcon(icon)
  }
}

func watchConfigFile(a fyne.App, filePath string) {
  watcher, err := fsnotify.NewWatcher()
  if err != nil {
    log.Fatal(err)
  }

  go func() {
    for {
      select {
        case event, ok := <-watcher.Events:
          if !ok {
            return
          }
          if event.Op&fsnotify.Write == fsnotify.Write || event.Op&fsnotify.Rename == fsnotify.Rename || event.Op&fsnotify.Create == fsnotify.Create {
            log.Println("Configuration file modified, reloading...")
            if err := loadAndSetupTray(a, filePath); err != nil {
              log.Printf("Failed to reload tray configuration: %v", err)
            }
          }
        case err, ok := <-watcher.Errors:
          if !ok {
            return
          }
          log.Println("error:", err)
      }
    }
  }()
  abs, err := filepath.Abs(filePath)
  if err != nil {
    log.Fatal(err)
  }
  err = watcher.Add(abs)
  if err != nil {
    log.Fatal(err)
  }
}

func iconNameToThemeIcon(name string) fyne.Resource {
  iconMap := map[string]func() fyne.Resource{
    "AccountIcon": theme.AccountIcon,
    "BrokenImageIcon": theme.BrokenImageIcon,
    "CancelIcon": theme.CancelIcon,
    "CheckButtonIcon": theme.CheckButtonIcon,
    "CheckButtonCheckedIcon": theme.CheckButtonCheckedIcon,
    "ColorAchromaticIcon": theme.ColorAchromaticIcon,
    "ColorChromaticIcon": theme.ColorChromaticIcon,
    "ColorPaletteIcon": theme.ColorPaletteIcon,
    "ComputerIcon": theme.ComputerIcon,
    "ConfirmIcon": theme.ConfirmIcon,
    "ContentAddIcon": theme.ContentAddIcon,
    "ContentClearIcon": theme.ContentClearIcon,
    "ContentCopyIcon": theme.ContentCopyIcon,
    "ContentCutIcon": theme.ContentCutIcon,
    "ContentPasteIcon": theme.ContentPasteIcon,
    "ContentRedoIcon": theme.ContentRedoIcon,
    "ContentRemoveIcon": theme.ContentRemoveIcon,
    "ContentUndoIcon": theme.ContentUndoIcon,
    "DeleteIcon": theme.DeleteIcon,
    "DocumentCreateIcon": theme.DocumentCreateIcon,
    "DocumentIcon": theme.DocumentIcon,
    "DocumentPrintIcon": theme.DocumentPrintIcon,
    "DocumentSaveIcon": theme.DocumentSaveIcon,
    "DownloadIcon": theme.DownloadIcon,
    "ErrorIcon": theme.ErrorIcon,
    "FileApplicationIcon": theme.FileApplicationIcon,
    "FileAudioIcon": theme.FileAudioIcon,
    "FileIcon": theme.FileIcon,
    "FileImageIcon": theme.FileImageIcon,
    "FileTextIcon": theme.FileTextIcon,
    "FileVideoIcon": theme.FileVideoIcon,
    "FolderIcon": theme.FolderIcon,
    "FolderNewIcon": theme.FolderNewIcon,
    "FolderOpenIcon": theme.FolderOpenIcon,
    "GridIcon": theme.GridIcon,
    "HelpIcon": theme.HelpIcon,
    "HistoryIcon": theme.HistoryIcon,
    "HomeIcon": theme.HomeIcon,
    "InfoIcon": theme.InfoIcon,
    "ListIcon": theme.ListIcon,
    "LoginIcon": theme.LoginIcon,
    "LogoutIcon": theme.LogoutIcon,
    "MailAttachmentIcon": theme.MailAttachmentIcon,
    "MailComposeIcon": theme.MailComposeIcon,
    "MailForwardIcon": theme.MailForwardIcon,
    "MailReplyAllIcon": theme.MailReplyAllIcon,
    "MailReplyIcon": theme.MailReplyIcon,
    "MailSendIcon": theme.MailSendIcon,
    "MediaFastForwardIcon": theme.MediaFastForwardIcon,
    "MediaFastRewindIcon": theme.MediaFastRewindIcon,
    "MediaMusicIcon": theme.MediaMusicIcon,
    "MediaPauseIcon": theme.MediaPauseIcon,
    "MediaPhotoIcon": theme.MediaPhotoIcon,
    "MediaPlayIcon": theme.MediaPlayIcon,
    "MediaRecordIcon": theme.MediaRecordIcon,
    "MediaReplayIcon": theme.MediaReplayIcon,
    "MediaSkipNextIcon": theme.MediaSkipNextIcon,
    "MediaSkipPreviousIcon": theme.MediaSkipPreviousIcon,
    "MediaStopIcon": theme.MediaStopIcon,
    "MediaVideoIcon": theme.MediaVideoIcon,
    "MenuDropDownIcon": theme.MenuDropDownIcon,
    "MenuDropUpIcon": theme.MenuDropUpIcon,
    "MenuExpandIcon": theme.MenuExpandIcon,
    "MenuIcon": theme.MenuIcon,
    "MoreHorizontalIcon": theme.MoreHorizontalIcon,
    "MoreVerticalIcon": theme.MoreVerticalIcon,
    "MoveDownIcon": theme.MoveDownIcon,
    "MoveUpIcon": theme.MoveUpIcon,
    "NavigateBackIcon": theme.NavigateBackIcon,
    "NavigateNextIcon": theme.NavigateNextIcon,
    "QuestionIcon": theme.QuestionIcon,
    "RadioButtonCheckedIcon": theme.RadioButtonCheckedIcon,
    "RadioButtonIcon": theme.RadioButtonIcon,
    "SearchReplaceIcon": theme.SearchReplaceIcon,
    "SearchIcon": theme.SearchIcon,
    "SettingsIcon": theme.SettingsIcon,
    "StorageIcon": theme.StorageIcon,
    "UploadIcon": theme.UploadIcon,
    "ViewFullScreenIcon": theme.ViewFullScreenIcon,
    "ViewRefreshIcon": theme.ViewRefreshIcon,
    "ViewRestoreIcon": theme.ViewRestoreIcon,
    "VisibilityOffIcon": theme.VisibilityOffIcon,
    "VisibilityIcon": theme.VisibilityIcon,
    "VolumeDownIcon": theme.VolumeDownIcon,
    "VolumeMuteIcon": theme.VolumeMuteIcon,
    "VolumeUpIcon": theme.VolumeUpIcon,
    "WarningIcon": theme.WarningIcon,
    "ZoomFitIcon": theme.ZoomFitIcon,
    "ZoomInIcon": theme.ZoomInIcon,
    "ZoomOutIcon": theme.ZoomOutIcon,
  }
  if iconFunc, exists := iconMap[name]; exists {
    return iconFunc()
  }
  return nil
}

