package scheduler

import (
	"context"
	"os"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
)

// JobFunc representa a função executada pelo scheduler.
type JobFunc func(ctx context.Context) error

// JobConfig configura um job individual.
type JobConfig struct {
	Name         string
	Schedule     string
	Enabled      bool
	FeatureFlag  string
	Job          JobFunc
	Tenants      []string
	TenantRunner func(ctx context.Context, tenantID string) error
}

// Scheduler encapsula o cron com métricas e logging.
type Scheduler struct {
	cron     *cron.Cron
	logger   *zap.Logger
	metrics  *Metrics
	runLogFn func(entryID cron.EntryID, job JobConfig, start time.Time, err error)
}

// Metrics agrupa métricas Prometheus para os jobs.
type Metrics struct {
	duration *prometheus.SummaryVec
	errors   *prometheus.CounterVec
}

// New cria um scheduler com cron spec padrão (segundos habilitados).
func New(logger *zap.Logger) *Scheduler {
	c := cron.New(cron.WithSeconds())
	return &Scheduler{
		cron:    c,
		logger:  logger,
		metrics: newMetrics(),
		runLogFn: func(entryID cron.EntryID, job JobConfig, start time.Time, err error) {
			// Placeholder para futura persistência em cron_run_logs (sqlc)
			_ = entryID
			_ = job
			_ = start
			_ = err
		},
	}
}

func newMetrics() *Metrics {
	return &Metrics{
		duration: promauto.NewSummaryVec(prometheus.SummaryOpts{
			Name:       "bap_scheduler_job_duration_seconds",
			Help:       "Duração de cada execução de job",
			Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
		}, []string{"job"}),
		errors: promauto.NewCounterVec(prometheus.CounterOpts{
			Name: "bap_scheduler_job_errors_total",
			Help: "Total de erros por job",
		}, []string{"job"}),
	}
}

// AddJob registra um job se estiver habilitado e com função válida.
func (s *Scheduler) AddJob(cfg JobConfig) error {
	if !cfg.Enabled {
		s.logger.Info("Job desabilitado (env)", zap.String("job", cfg.Name))
		return nil
	}
	if cfg.FeatureFlag != "" && strings.EqualFold(os.Getenv(cfg.FeatureFlag), "false") {
		s.logger.Info("Job bloqueado por feature flag", zap.String("job", cfg.Name), zap.String("flag", cfg.FeatureFlag))
		return nil
	}

	run := s.wrap(cfg)
	entryID, err := s.cron.AddFunc(cfg.Schedule, run)
	if err != nil {
		return err
	}

	s.logger.Info("Job registrado",
		zap.String("job", cfg.Name),
		zap.String("schedule", cfg.Schedule),
		zap.Int64("entry_id", int64(entryID)),
		zap.Bool("multi_tenant", len(cfg.Tenants) > 0),
	)
	return nil
}

func (s *Scheduler) wrap(cfg JobConfig) func() {
	return func() {
		start := time.Now()
		ctx := context.Background()

		var err error
		if len(cfg.Tenants) > 0 && cfg.TenantRunner != nil {
			for _, tenantID := range cfg.Tenants {
				if tenantID == "" {
					continue
				}
				if e := cfg.TenantRunner(ctx, tenantID); e != nil {
					err = e
					s.metrics.errors.WithLabelValues(cfg.Name).Inc()
					s.logger.Warn("Erro ao executar job para tenant",
						zap.String("job", cfg.Name),
						zap.String("tenant_id", tenantID),
						zap.Error(e),
					)
				}
			}
		} else if cfg.Job != nil {
			err = cfg.Job(ctx)
			if err != nil {
				s.metrics.errors.WithLabelValues(cfg.Name).Inc()
				s.logger.Warn("Erro ao executar job",
					zap.String("job", cfg.Name),
					zap.Error(err),
				)
			}
		} else {
			s.logger.Warn("Job sem função associada", zap.String("job", cfg.Name))
		}

		elapsed := time.Since(start)
		s.metrics.duration.WithLabelValues(cfg.Name).Observe(elapsed.Seconds())
		s.runLogFn(0, cfg, start, err)
	}
}

// Start inicia o cron em background.
func (s *Scheduler) Start() {
	s.cron.Start()
	s.logger.Info("Scheduler iniciado")
}

// Stop encerra o cron e espera execuções atuais terminarem.
func (s *Scheduler) Stop(ctx context.Context) {
	stopCtx := s.cron.Stop()
	select {
	case <-stopCtx.Done():
	case <-ctx.Done():
	}
	s.logger.Info("Scheduler encerrado")
}

// ParseTenantEnv lê a env e devolve lista de tenants.
func ParseTenantEnv(envName string) []string {
	val := os.Getenv(envName)
	if val == "" {
		return nil
	}
	parts := strings.Split(val, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if trimmed := strings.TrimSpace(p); trimmed != "" {
			out = append(out, trimmed)
		}
	}
	return out
}
