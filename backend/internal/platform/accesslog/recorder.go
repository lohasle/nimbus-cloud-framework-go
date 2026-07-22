package accesslog

import (
	"github.com/lohasle/nimbus-cloud-framework-go/internal/modules/infra"
	"github.com/lohasle/nimbus-cloud-framework-go/internal/platform/middleware"
	"gorm.io/gorm"
)

// Recorder persists authenticated request metadata in the shared infra log
// table while preserving the trace ID propagated through the gateway.
func Recorder(db *gorm.DB, applicationName string) middleware.RequestLogRecorder {
	return func(record middleware.RequestLog) {
		if record.TenantID == 0 {
			return
		}
		db.Create(&infra.APIAccessLog{
			TenantID: record.TenantID, TraceID: record.TraceID, ApplicationName: applicationName,
			UserID: record.UserID, Method: record.Method, Path: record.Path, Status: record.Status,
			Duration: record.Duration, IP: record.IP, UserAgent: record.UserAgent,
		})
	}
}
