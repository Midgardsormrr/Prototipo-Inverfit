package database

import (
	"context"
	"strings"
	"testing"
	"time"
)

func TestConnectEmptyURI(t *testing.T) {
	ctx := context.Background()
	_, err := Connect(ctx, "", "test_db")
	if err == nil {
		t.Fatal("esperaba error por URI vacía, pero fue nil")
	}
}

func TestConnectEmptyDB(t *testing.T) {
	ctx := context.Background()
	_, err := Connect(ctx, "mongodb://localhost:27017", "")
	if err == nil {
		t.Fatal("esperaba error por nombre de base de datos vacío, pero fue nil")
	}
}

func TestConnectUnreachable(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// URI con credenciales ficticias hacia un puerto inexistente
	fakeURI := "mongodb://testuser:supersecretpass@127.0.0.1:29999/inverfit_test"
	_, err := Connect(ctx, fakeURI, "inverfit_test")
	if err == nil {
		t.Fatal("esperaba error de conexión al puerto 29999, pero fue nil")
	}

	errMsg := err.Error()
	if strings.Contains(errMsg, "supersecretpass") {
		t.Errorf("el mensaje de error no debe revelar la contraseña: %s", errMsg)
	}
}

func TestNilMongoDBPingAndDisconnect(t *testing.T) {
	var m *MongoDB
	err := m.Ping(context.Background())
	if err == nil {
		t.Fatal("esperaba error de Ping con instancia nil")
	}

	err = m.Disconnect(context.Background())
	if err != nil {
		t.Errorf("Disconnect con instancia nil debe ser seguro y retornar nil, obtenido: %v", err)
	}
}
