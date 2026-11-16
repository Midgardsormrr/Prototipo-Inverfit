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
	"inverfit-backend/internal/modules/payments"
	"inverfit-backend/internal/modules/schedule"
	"inverfit-backend/internal/modules/students"
	"inverfit-backend/internal/router"
)

// studentReaderAdapter conecta el servicio de alumnos con la interfaz requerida por Agenda,
// garantizando la independencia modular e imposibilitando el acceso directo a la colección students.
type studentReaderAdapter struct {
	service students.Service
}

func (a *studentReaderAdapter) GetStudentsSummary(ctx context.Context, ids []string) ([]schedule.StudentSummary, error) {
	var summaries []schedule.StudentSummary
	for _, id := range ids {
		st, err := a.service.GetByID(ctx, id)
		if err != nil {
			return nil, err
		}
		summaries = append(summaries, schedule.StudentSummary{
			ID:       st.ID,
			FullName: st.FullName,
			Archived: false,
		})
	}
	return summaries, nil
}

// paymentsStudentReaderAdapter conecta el servicio de alumnos con la interfaz requerida por Pagos,
// permitiendo consultar datos del alumno y activar membresías "Por revisar" sin acoplamiento a MongoDB.
type paymentsStudentReaderAdapter struct {
	service students.Service
}

func (a *paymentsStudentReaderAdapter) GetStudent(ctx context.Context, id string) (*payments.StudentInfo, error) {
	st, err := a.service.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return &payments.StudentInfo{
		ID:               st.ID,
		FullName:         st.FullName,
		PlanName:         st.PlanName,
		Frequency:        st.Frequency,
		MembershipStatus: st.MembershipStatus,
		Archived:         false,
	}, nil
}

func (a *paymentsStudentReaderAdapter) ActivateMembershipIfReviewing(ctx context.Context, studentID string) error {
	st, err := a.service.GetByID(ctx, studentID)
	if err != nil {
		return err
	}
	if st.MembershipStatus == "Por revisar" {
		statusActive := "Activo"
		_, err = a.service.Update(ctx, studentID, students.UpdateStudentRequest{
			MembershipStatus: &statusActive,
		})
		return err
	}
	return nil
}

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

	// 4. Configurar Módulo de Alumnos
	var studentsRepo students.Repository
	if mongoDB != nil && mongoDB.Database != nil {
		studentsRepo = students.NewRepository(mongoDB.Database)
	} else {
		studentsRepo = students.NewRepository(nil)
	}
	studentsService := students.NewService(studentsRepo)
	studentsHandler := students.NewHandler(studentsService)

	// 5. Configurar Módulo de Agenda (Schedule)
	var scheduleRepo schedule.Repository
	if mongoDB != nil && mongoDB.Database != nil {
		scheduleRepo = schedule.NewRepository(mongoDB.Database)
		// Asegurar índices de colecciones de agenda
		idxCtx, idxCancel := context.WithTimeout(context.Background(), 5*time.Second)
		if err := scheduleRepo.EnsureIndexes(idxCtx); err != nil {
			logger.Warn("No se pudieron verificar los índices de agenda", slog.String("error", err.Error()))
		}
		idxCancel()
	} else {
		scheduleRepo = schedule.NewRepository(nil)
	}
	studentReader := &studentReaderAdapter{service: studentsService}
	scheduleService := schedule.NewService(scheduleRepo, studentReader)
	scheduleHandler := schedule.NewHandler(scheduleService)

	// 6. Configurar Módulo de Pagos
	var paymentRepo payments.Repository
	if mongoDB != nil && mongoDB.Database != nil {
		paymentRepo = payments.NewRepository(mongoDB.Database)
		idxCtx, idxCancel := context.WithTimeout(context.Background(), 5*time.Second)
		if err := paymentRepo.EnsureIndexes(idxCtx); err != nil {
			logger.Warn("No se pudieron verificar los índices de pagos", slog.String("error", err.Error()))
		}
		idxCancel()
	} else {
		paymentRepo = payments.NewRepository(nil)
	}
	paymentStudentReader := &paymentsStudentReaderAdapter{service: studentsService}
	paymentService := payments.NewService(paymentRepo, paymentStudentReader)
	paymentHandler := payments.NewHandler(paymentService)

	// 7. Configurar Router de Gin
	var pinger database.Pinger
	if mongoDB != nil {
		pinger = mongoDB
	}
	r := router.Setup(cfg, pinger, studentsHandler, scheduleHandler, paymentHandler, logger)

	// 8. Configurar Servidor HTTP con timeouts defensivos
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
