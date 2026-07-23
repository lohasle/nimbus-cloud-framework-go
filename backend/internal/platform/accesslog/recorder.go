package accesslog

import (
	"net/http"
	"time"

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
		_ = db.Create(&infra.APIAccessLog{
			TenantID: record.TenantID, TraceID: record.TraceID, ApplicationName: applicationName,
			UserID: record.UserID, Method: record.Method, Path: record.Path, Status: record.Status,
			Duration: record.Duration, IP: record.IP, UserAgent: record.UserAgent,
		}).Error
		if record.Status >= http.StatusInternalServerError {
			_ = db.Create(&infra.APIErrorLog{
				TenantID: record.TenantID, TraceID: record.TraceID, UserID: record.UserID,
				UserType: 2, ApplicationName: applicationName, RequestMethod: record.Method,
				RequestURL: record.Path, UserIP: record.IP, UserAgent: record.UserAgent,
				ExceptionTime: time.Now(), ExceptionName: http.StatusText(record.Status),
				ExceptionMessage: "HTTP request failed",
				ExceptionRootCauseMessage: "HTTP request returned server error",
				ResultCode:                record.Status,
			}).Error
		}
	}
}
