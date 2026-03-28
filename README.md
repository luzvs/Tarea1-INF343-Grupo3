Tener instalado go 1.2+ recomendable
Esto se hizo con WSL linux en Visual Studio Code

Para el main.go:
Asegurarse de estar en:
cd parte1/api

En la terminal:
go mod init parte1/api
go get github.com/gin-gonic/gin
go mod tidy

Ejecutar:
go run main.go

En otra terminal sin cerrar la anterior:

Asegurarse de estar en:
cd parte1/cliente

Ejecutar:
go run cliente.go
