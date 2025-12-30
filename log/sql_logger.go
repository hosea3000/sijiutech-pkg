package repository

import (
	"context"
	"database/sql"
	"time"

	"sijiutech-material-api/pkg/log"

	"go.uber.org/zap"
)

// sqlLogger 实现了 DBTX 接口，用于记录SQL执行日志
type sqlLogger struct {
	db     *sql.DB
	logger *log.Logger
}

// NewSQLLogger 创建一个带日志记录的数据库包装器
func NewSQLLogger(db *sql.DB, logger *log.Logger) *sqlLogger {
	return &sqlLogger{
		db:     db,
		logger: logger,
	}
}

// ExecContext 执行SQL并记录日志
func (s *sqlLogger) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	start := time.Now()
	result, err := s.db.ExecContext(ctx, query, args...)
	duration := time.Since(start)

	logFields := []zap.Field{
		zap.String("sql", query),
		zap.Any("args", args),
		zap.Duration("duration", duration),
	}
	if err != nil {
		logFields = append(logFields, zap.Error(err))
		s.logger.WithContext(ctx).Error("SQL Exec", logFields...)
	} else {
		s.logger.WithContext(ctx).Debug("SQL Exec", logFields...)
	}

	return result, err
}

// PrepareContext 准备SQL语句并记录日志
func (s *sqlLogger) PrepareContext(ctx context.Context, query string) (*sql.Stmt, error) {
	start := time.Now()
	stmt, err := s.db.PrepareContext(ctx, query)
	duration := time.Since(start)

	logFields := []zap.Field{
		zap.String("sql", query),
		zap.Duration("duration", duration),
	}
	if err != nil {
		logFields = append(logFields, zap.Error(err))
		s.logger.WithContext(ctx).Error("SQL Prepare", logFields...)
	} else {
		s.logger.WithContext(ctx).Debug("SQL Prepare", logFields...)
	}

	// 注意：返回原始的*sql.Stmt以符合DBTX接口
	// 如果需要记录Prepared Statement的执行日志，可以考虑使用driver wrapper
	return stmt, err
}

// QueryContext 查询SQL并记录日志
func (s *sqlLogger) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	start := time.Now()
	rows, err := s.db.QueryContext(ctx, query, args...)
	duration := time.Since(start)

	logFields := []zap.Field{
		zap.String("sql", query),
		zap.Any("args", args),
		zap.Duration("duration", duration),
	}
	if err != nil {
		logFields = append(logFields, zap.Error(err))
		s.logger.WithContext(ctx).Error("SQL Query", logFields...)
	} else {
		s.logger.WithContext(ctx).Debug("SQL Query", logFields...)
	}

	return rows, err
}

// QueryRowContext 查询单行SQL并记录日志
func (s *sqlLogger) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	start := time.Now()
	row := s.db.QueryRowContext(ctx, query, args...)
	duration := time.Since(start)

	s.logger.WithContext(ctx).Debug("SQL QueryRow",
		zap.String("sql", query),
		zap.Any("args", args),
		zap.Duration("duration", duration),
	)

	return row
}
