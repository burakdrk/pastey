# pastey

### API

Required Dependencies:

1. sqlc, `$ go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest`
2. golang-migrate, `$ brew install golang-migrate` or https://github.com/golang-migrate/migrate/releases
3. GNU make
4. Docker Postgres, `$ docker pull postgres:16-alpine`

### iOS App

<img height=500 src="https://raw.githubusercontent.com/burakdrk/pastey/main/screenshots/iosapp.jpg"/>

You can build the .ipa file with Xcode and sideload with AltStore or a similar tool.

### Desktop App

<img height=500 src="https://raw.githubusercontent.com/burakdrk/pastey/main/screenshots/wailsapp.png"/>

Building:
1. From project root directory, `$ cd pastey-wails`
2. Build for your operating system (Not tested on Mac/Linux) `$ wails build`
3. The executable can be found in build/bin
