package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"inverfit-backend/internal/config"
	"inverfit-backend/internal/database"
	"inverfit-backend/internal/router"
)

func main() {
	// 1. Cargar configuración desde variables de entorno y archivo .env opcional
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error crítico de configuración: %v\n", err)
		os.Exit(1)
	}

	// 2. Inicializar logger estructurado según el nivel configurado
	var logLevel slog.Level
	switch cfg.LogLevel {
	case "debug":
		logLevel = slog.LevelDebug
	case "warn":
		logLevel = slog.LevelWarn
	case "error":
		logLevel = slog.LevelError
	default:
		logLevel = slog.LevelInfo
	}

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: logLevel,
	}))
	slog.SetDefault(logger)

	logger.Info("Iniciando servicio backend INVERFIT",
		slog.String("app_env", cfg.AppEnv),
		slog.String("http_port", cfg.HTTPPort),
		slog.String("log_level", cfg.LogLevel),
		slog.String("mongodb_uri", cfg.RedactedMongoURI()),
		slog.String("mongodb_database", cfg.MongoDBDatabase),
	)

	// 3. Conectar a MongoDB con timeout controlado
	dbCtx, dbCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer dbCancel()

	var mongoDB *database.MongoDB
	mongoDB, err = database.Connect(dbCtx, cfg.MongoDBURI, cfg.MongoDBDatabase)
	if err != nil {
		// Se registra advertencia: el proceso inicia pero /health/ready reportará 503 hasta que la BD esté disponible
		logger.Warn("No se pudo establecer conexión inmediata con MongoDB al iniciar. El servidor operará en modo degradado (/health/ready responderá 503)",
			slog.String("error", err.Error()),
		)
	} else {
		logger.Info("Conexión con MongoDB verificada exitosamente",
			slog.String("database", cfg.MongoDBDatabase),
		)
	}

	// 4. Configurar Router de Gin
	var pinger database.Pinger
	if mongoDB != nil {
		pinger = mongoDB
	}
	r := router.Setup(cfg, pinger, logger)

	// 5. Configurar Servidor HTTP con timeouts defensivos
	srv := &http.Server{
		Addr:         ":" + cfg.HTTPPort,
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Canal para escuchar señales de terminación del sistema operativo
	shutdownChan := make(chan os.Signal, 1)
	signal.Notify(shutdownChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Iniciar servidor HTTP en una goroutine independiente
	go func() {
		logger.Info(fmt.Sprintf("Servidor HTTP escuchando en puerto :%s", cfg.HTTPPort))
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("Fallo inesperado en el servidor HTTP", slog.String("error", err.Error()))
			os.Exit(1)
		}
	}()

	// 6. Esperar señal de terminación para ejecutar cierre ordenado (Graceful Shutdown)
	sig := <-shutdownChan
	logger.Info("Señal de apagado recibida, iniciando cierre ordenado del servidor",
		slog.String("signal", sig.String()),
	)

	// Contexto con timeout de 10 segundos para el drenado de peticiones activas
	shutdownCtx, cancelShutdown := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelShutdown()

	// Detener servidor HTTP
	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Error("Error durante el cierre ordenado del servidor HTTP", slog.String("error", err.Error()))
	} else {
		logger.Info("Servidor HTTP detenido correctamente sin interrumpir peticiones en curso")
	}

	// Desconectar cliente de MongoDB si estaba activo
	if mongoDB != nil {
		if err := mongoDB.Disconnect(shutdownCtx); err != nil {
			logger.Error("Error al cerrar conexión con MongoDB", slog.String("error", err.Error()))
		} else {
			logger.Info("Pool de conexiones de MongoDB cerrado limpiamente")
		}
	}

	logger.Info("Apagado ordenado de INVERFIT Backend completado exitosamente")
}
