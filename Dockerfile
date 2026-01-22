# Image de base légère Go
FROM golang:1.24-alpine AS builder

# Créer un répertoire de travail
WORKDIR /projetgo

# Copier les fichiers de modules
COPY go.mod go.sum ./

# Télécharger les dépendances
RUN go mod download

# Copier tout le code
COPY . .

# Compiler le binaire
RUN go build -o app

# --- Image finale légère ---
FROM alpine:latest
WORKDIR /app

# Copier le binaire compilé
COPY --from=builder /projetgo/app .

# Ajouter un utilisateur non-root pour la sécurité
RUN adduser -D appuser
USER appuser

# Exposer le port que ton service utilise
EXPOSE 8080

# Lancer le service
CMD ["./app"]
