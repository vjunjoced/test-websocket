# Etapa de construcción
FROM golang:1.19 AS builder

# Directorio de trabajo
WORKDIR /app

# Copiar los archivos del código fuente al contenedor
COPY . .

# Descargar las dependencias
RUN go mod download

# Compilar la aplicación
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Etapa de ejecución
FROM alpine:latest

# Directorio de trabajo
WORKDIR /root/

# Copiar el ejecutable desde la etapa de construcción
COPY --from=builder /app/main .

# Exponer el puerto en el que se ejecutará la aplicación
EXPOSE 3000

# Comando para ejecutar la aplicación
CMD ["./main"]
