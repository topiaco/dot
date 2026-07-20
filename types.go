package dot

// ImportResult 通用运维任务导入结果
type ImportResult struct {
	TaskName  string   `json:"task_name"`
	TotalRows int      `json:"total_rows"`
	Inserted  int      `json:"inserted"`
	Updated   int      `json:"updated"`
	Skipped   int      `json:"skipped"`
	Failed    int      `json:"failed"`
	Errors    []string `json:"errors,omitempty"`
	Message   string   `json:"message,omitempty"`
}

// TaskArgs 运维/数据迁移任务触发参数
type TaskArgs struct {
	TaskName string         `json:"task_name" binding:"required"`
	DryRun   bool           `json:"dry_run"`
	Operator string         `json:"operator"`
	Params   map[string]any `json:"params"`
}

// Mode 获取当前运行模式的中文描述
func (args *TaskArgs) Mode() string {
	if args.DryRun {
		return "演示"
	}
	return "正式"
}

// GetStringParam 获取字符串类型参数
func (args *TaskArgs) GetStringParam(key, defaultValue string) string {
	if args.Params == nil {
		return defaultValue
	}
	if val, ok := args.Params[key].(string); ok {
		return val
	}
	return defaultValue
}

// GetBoolParam 获取布尔类型参数（支持 bool 和 string "true"/"1"）
func (args *TaskArgs) GetBoolParam(key string, defaultValue bool) bool {
	if args.Params == nil {
		return defaultValue
	}
	if val, ok := args.Params[key].(bool); ok {
		return val
	}
	if val, ok := args.Params[key].(string); ok {
		return val == "true" || val == "1"
	}
	return defaultValue
}