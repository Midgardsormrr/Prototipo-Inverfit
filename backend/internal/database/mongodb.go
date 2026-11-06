package database

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// Pinger define la interfaz para comprobar la conectividad con la base de datos.
type Pinger interface {
	Ping(ctx context.Context) error
}

// MongoDB encapsula el cliente de MongoDB y la referencia a la base de datos configurada.
type MongoDB struct {
	Client   *mongo.Client
	Database *mongo.Database
}

// Connect inicializa el cliente oficial de MongoDB, establece la conexión y valida la disponibilidad mediante Ping.
// Configura tiempos de espera prudentes para evitar bloqueos indefinidos.
func Connect(ctx context.Context, uri string, dbName string) (*MongoDB, error) {
	if strings.TrimSpace(uri) == "" {
		return nil, fmt.Errorf("la URI de conexión a MongoDB no puede estar vacía")
	}
	if strings.TrimSpace(dbName) == "" {
		return nil, fmt.Errorf("el nombre de la base de datos no puede estar vacío")
	}

	// Opciones del cliente con timeouts controlados
	clientOptions := options.Client().
		ApplyURI(uri).
		SetConnectTimeout(5 * time.Second).
		SetServerSelectionTimeout(5 * time.Second).
		SetSocketTimeout(10 * time.Second)

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("error al inicializar cliente de MongoDB: %w", err)
	}

	// Comprobación efectiva mediante Ping con timeout
	pingCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	if err := client.Ping(pingCtx, readpref.Primary()); err != nil {
		// En caso de fallo de conectividad, cerramos el cliente antes de retornar
		_ = client.Disconnect(context.Background())
		redactedURI := redactURI(uri)
		return nil, fmt.Errorf("no se pudo conectar con MongoDB en %s: %w", redactedURI, err)
	}

	return &MongoDB{
		Client:   client,
		Database: client.Database(dbName),
	}, nil
}

// Ping verifica la conectividad actual con el cluster de MongoDB.
func (m *MongoDB) Ping(ctx context.Context) error {
	if m == nil || m.Client == nil {
		return fmt.Errorf("cliente de MongoDB no inicializado")
	}
	return m.Client.Ping(ctx, readpref.Primary())
}

// Disconnect cierra ordenadamente el pool de conexiones de MongoDB.
func (m *MongoDB) Disconnect(ctx context.Context) error {
	if m == nil || m.Client == nil {
		return nil
	}
	return m.Client.Disconnect(ctx)
}

// redactURI oculta la información sensible (contraseñas) de la URI para mensajes de error seguros.
func redactURI(rawURI string) string {
	u, err := url.Parse(rawURI)
	if err != nil {
		parts := strings.Split(rawURI, "@")
		if len(parts) > 1 {
			return "mongodb://[REDACTED]@" + parts[len(parts)-1]
		}
		return "[REDACTED_URI]"
	}
	return u.Redacted()
}
