# IP Update

This mini app updates your ip address with Digital Ocean every 1 min. This is handy for doing a dynamic ip address on a home network.

# How to build

* `CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o ip-update`

# Install on MacOS

* Copy `ip-update.plist` to `/Library/LaunchDaemons/ip-update.plist`
* `sudo chown root:wheel /Library/LaunchDaemons/ip-update.plist`
* `sudo launchctl load /Library/LaunchDaemons/ip-update.plist`
* Reboot and verify output in `/tmp/ip-update.log`