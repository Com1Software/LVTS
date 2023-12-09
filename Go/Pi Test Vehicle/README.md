# Pi Test Vehicle

## Installation Build

```shell
https://github.com/Com1Software/LVTS.git
cd 'LVTS/Gp/Pi Test Vehicle"
go mod init test
go mod tidy
go build

```

## Auto Start Setup
To make the vehicle run automaticly when the Rasberry Pi is first turned on,
you can add the cammand to the .bashrc. To do this at the command line enter

sudo vim .bashrc

At the very end of file add:
```shell
sudo nmcli device wifi hotspot ssid PiLVTS password yourpassword
'./LVTS/Go/Pi Test Vevicle/test'

```
Next run sudo raspi-config and in the system change the boot option to CLI user


## Useful Links

[Test-Magnetometer](https://github.com/Com1Software/Test-Magnetometer)
