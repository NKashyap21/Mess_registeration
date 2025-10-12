package logs

import "github.com/LambdaIITH/mess_registration/services"

// InitLogsController initializes and returns a new LogsController
func InitLogsController() *LogsController {
	return &LogsController{
		loggerService: services.GetLoggerService(),
	}
}
