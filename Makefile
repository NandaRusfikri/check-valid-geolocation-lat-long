#================================
#== GOLANG ENVIRONMENT
#================================
GO := @go
GIN := @gin
SET := @set
ECHO := @echo

goinstall:
	${GO} get .

gobuild:
	${SET} GOOS=linux
	${SET} GOARCH=amd64
	${GO} build -o geolocation-linux main.go
	${ECHO} "Compiling done linux..."
	${SET} GOOS=windows
	${GO} build -o geolocation-win.exe main.go
	${ECHO} "Compiling done windows..."
	${SET} GOOS=darwin
	${GO} build -o geolocation-mac main.go
	${ECHO} "Compiling done mac os..."

