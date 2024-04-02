# Tray for golang

Tray provides a GUI system tray for your application, dynamically generated and updated from a JSON file.

## Screenshot

[!Tray](https://raw.githubusercontent.com/ipv6rslimited/tray/main/screenshot.png)

## Features

- Custom System Tray icon in Mac, windows and Linux generated from a JSON

- Updates the tray menu based on the JSON file updates (watches the JSON file)

- Super small binary when compiled.

## Use Case

- We wanted to provide a super simple way for people to run shell commands. This + [Configurator](https://github.com/ipv6rslimited/configurator) does the trick.

- This was to help people who use [IPv6rs](https://ipv6.rs) as well as anyone who needs a system tray, period.

## Example

- See the test folder to see how to use the Tray

- Build the tray-test by typing:
```
cd test
go run test.go config.json
```

## License

Distributed under the COOL License.

Copyright (c) 2024 IPv6.rs <https://ipv6.rs>
All Rights Reserved
