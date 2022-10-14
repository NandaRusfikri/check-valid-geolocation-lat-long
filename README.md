## Check Valid Geolocation With Latitudes and Longitudes
#### By NandaRusfikri - 14/10/2022

![](img/img1.jpg)
![](img/img2.jpg)

- ### Note
  - Input Datasource with CSV
  - Output Geolocation Valid or Not Contains CSV
  - Checking use API https://www.gps-coordinates.net/
  

  ```sh
  $ go get . || go mod 
  ```

  - Start application

  ```sh
  $ go run main.go  
  ```
  - Build application Windows

  ```sh
  $ set GOOS=windows
  $ go build -o geolocation-win.exe main.go 
  ```

  - Build application Linux

  ```sh
  $ set GOOS=linux
  $ go build -o geolocation-linux main.go 
  ```

  - Build application Mac

  ```sh
  $ set GOOS=darwin
  $ go build -o geolocation-mac main.go 
  ```



