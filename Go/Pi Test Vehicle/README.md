# Pi Test Vehicle

## Installation and Build

```shell
git clone https://github.com/Com1Software/LVTS.git
cd 'LVTS/Gp/Pi Test Vehicle"
go mod init test
go mod tidy
go build
```

## Auto Start Setup With WiFi Hotspot
To make the vehicle run automaticly when the Rasberry Pi is first turned on,
you can add the command to the .bashrc. To do this at the command line enter

```shell
sudo vim .bashrc
```

At the very end of file add:
```shell
sudo nmcli device wifi hotspot ssid Pi-LVTS password yourpassword
'./LVTS/Go/Pi Test Vevicle/test'
```
Next run sudo raspi-config and in the system change the boot option to CLI user


## Useful Links

[Test-Magnetometer](https://github.com/Com1Software/Test-Magnetometer)
