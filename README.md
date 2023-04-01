# **Gown**
**Gown** (Go download) is a download manager written in Go (backend) and flutter (frontend). This application will be targeted as desktop app for Internet Download Manager alternative

## **Getting Started**
### Prerequisites
- Go
- Flutter

### Installation (Fedora Linux)
1. Install Go
2. Download the flutter tar from [here](https://docs.flutter.dev/get-started/install/linux)
3. Install the flutter
   
```bash
# extract the tar
tar xf ~/Downloads/flutter_linux_3.7.9-stable.tar.xz

# copy to .zshrc
export PATH=$PATH:~/flutter/bin

# refresh the shell
source ~/.zshrc

# check if flutter successfully installed 
flutter --version
```

4. Install the required dependency for desktop app development
```bash
sudo dnf install clang cmake ninja-build pkgconf gtk3 gtk3-devel
```
5. Install [Hover](https://hover.build)
```bash
GO111MODULE=on go get -u -a github.com/go-flutter-desktop/hover@latest
```
6. Or just run `go mod tidy`


## **State**
The current state of this project is prototype