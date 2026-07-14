package dot

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/xuri/excelize/v2"
)

// ExcelImportResult 封装通用的导入计数与错误留痕结果
type ExcelImportResult struct {
	SheetName string   `json:"sheet_name"`
	TotalRows int      `json:"total_rows"`
	Inserted  int      `json:"inserted"`
	Updated   int      `json:"updated"`
	Skipped   int      `json:"skipped"`
	Failed    int      `json:"failed"`
	Errors    []string `json:"errors"`
}

// ExcelImporter 封装通用 Excel 导入器结构体
type ExcelImporter struct {
	filePath       string
	excelFile      *excelize.File
	currentSheet   string
	currentRows    [][]string
	headerMap      map[string]int
	headerRowIndex int
}

// NewExcelImporter 仅负责打开 Excel 文件描述符（支持本地路径与 HTTP/HTTPS 远端链接）
func NewExcelImporter(filePath string) (*ExcelImporter, error) {
	var f *excelize.File
	var err error

	filePathTrimmed := strings.TrimSpace(filePath)
	if strings.HasPrefix(filePathTrimmed, "http://") || strings.HasPrefix(filePathTrimmed, "https://") {
		resp, httpErr := http.Get(filePathTrimmed)
		if httpErr != nil {
			return nil, fmt.Errorf("下载远端 Excel 文件失败: %w", httpErr)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("下载远端 Excel 文件返回状态码异常: %d", resp.StatusCode)
		}

		f, err = excelize.OpenReader(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("解析远端 Excel 数据流失败: %w", err)
		}
	} else {
		f, err = excelize.OpenFile(filePathTrimmed)
		if err != nil {
			return nil, fmt.Errorf("打开本地 Excel 文件失败: %w", err)
		}
	}

	return &ExcelImporter{
		filePath:  filePath,
		excelFile: f,
	}, nil
}

// GetSheetList 返回当前 Excel 包含的所有 Sheet 名字列表
func (imp *ExcelImporter) GetSheetList() []string {
	if imp.excelFile == nil {
		return nil
	}
	return imp.excelFile.GetSheetList()
}

// SelectSheetByName 切换至指定名称的 Sheet，并加载 rows，智能识别并建立表头映射
func (imp *ExcelImporter) SelectSheetByName(name string) error {
	if imp.excelFile == nil {
		return errors.New("excel 句柄已关闭或未初始化")
	}

	rows, err := imp.excelFile.GetRows(name)
	if err != nil {
		return fmt.Errorf("读取工作表 %s 数据失败: %w", name, err)
	}

	imp.currentSheet = name
	imp.currentRows = rows
	imp.headerMap = nil
	imp.headerRowIndex = -1

	if len(rows) > 0 {
		imp.searchAndBuildHeaderMap()
	}
	return nil
}

// SelectSheetByIndex 切换至指定物理索引的 Sheet，加载 rows 并建立表头映射
func (imp *ExcelImporter) SelectSheetByIndex(index int) error {
	sheetList := imp.GetSheetList()
	if index < 0 || index >= len(sheetList) {
		return fmt.Errorf("工作表索引 %d 越界, 总数: %d", index, len(sheetList))
	}
	return imp.SelectSheetByName(sheetList[index])
}

// GetCurrentSheetName 获取当前处于活动状态 of Sheet 名字
func (imp *ExcelImporter) GetCurrentSheetName() string {
	return imp.currentSheet
}

// GetFile 返回底层的 excelize.File 句柄对象
func (imp *ExcelImporter) GetFile() *excelize.File {
	return imp.excelFile
}

// Close 关闭释放文件句柄
func (imp *ExcelImporter) Close() error {
	if imp.excelFile != nil {
		err := imp.excelFile.Close()
		imp.excelFile = nil
		return err
	}
	return nil
}

// GetCellValue 获取列对应的值
func (imp *ExcelImporter) GetCellValue(row []string, colName string) string {
	if imp.headerMap == nil {
		return ""
	}
	colKey := imp.normalizeHeaderKey(colName)
	idx, ok := imp.headerMap[colKey]
	if !ok || idx >= len(row) {
		return ""
	}
	val := strings.TrimSpace(row[idx])
	if strings.EqualFold(val, "NULL") || val == "/" || val == "\\" || strings.EqualFold(val, "#N/A") {
		return ""
	}
	return val
}

// GetCellValueFallback 根据别名匹配获取值
func (imp *ExcelImporter) GetCellValueFallback(row []string, colNames ...string) string {
	for _, name := range colNames {
		if val := imp.GetCellValue(row, name); val != "" {
			return val
		}
	}
	return ""
}

// GetCellValueByOffset 根据表头名称加上物理列偏移获取单元格值
func (imp *ExcelImporter) GetCellValueByOffset(row []string, colName string, offset int) string {
	if imp.headerMap == nil {
		return ""
	}
	colKey := imp.normalizeHeaderKey(colName)
	baseIdx, ok := imp.headerMap[colKey]
	if !ok {
		return ""
	}

	targetIdx := baseIdx + offset
	if targetIdx < 0 || targetIdx >= len(row) {
		return ""
	}

	val := strings.TrimSpace(row[targetIdx])
	if strings.EqualFold(val, "NULL") || val == "/" || val == "\\" || strings.EqualFold(val, "#N/A") {
		return ""
	}
	return val
}

// GetCellValueFallbackByOffset 根据候选表头别名及物理列偏移获取单元格值
func (imp *ExcelImporter) GetCellValueFallbackByOffset(row []string, offset int, colNames ...string) string {
	for _, name := range colNames {
		if val := imp.GetCellValueByOffset(row, name, offset); val != "" {
			return val
		}
	}
	return ""
}


// GetCellValueAsInt 获取整型单元格值，转换失败返回默认值
func (imp *ExcelImporter) GetCellValueAsInt(row []string, colName string, defaultVal int) int {
	val := imp.GetCellValue(row, colName)
	if val == "" {
		return defaultVal
	}
	i, err := strconv.Atoi(val)
	if err != nil {
		if f, err2 := strconv.ParseFloat(val, 64); err2 == nil {
			return int(f)
		}
		return defaultVal
	}
	return i
}

// GetCellValueAsFloat 获取浮点单元格值，转换失败返回默认值
func (imp *ExcelImporter) GetCellValueAsFloat(row []string, colName string, defaultVal float64) float64 {
	val := imp.GetCellValue(row, colName)
	if val == "" {
		return defaultVal
	}
	f, err := strconv.ParseFloat(val, 64)
	if err != nil {
		return defaultVal
	}
	return f
}

// GetPictures 提取当前行指定单元格的图片
func (imp *ExcelImporter) GetPictures(rowIndex int, colLetter string) ([]excelize.Picture, error) {
	if imp.excelFile == nil || imp.currentSheet == "" {
		return nil, errors.New("excel 句柄不可用")
	}
	cellName := fmt.Sprintf("%s%d", strings.ToUpper(colLetter), rowIndex+1)
	return imp.excelFile.GetPictures(imp.currentSheet, cellName)
}

// Import 通用导入驱动流框架
func (imp *ExcelImporter) Import(ctx context.Context, handler func(rowIndex int, row []string) (string, error)) (*ExcelImportResult, error) {
	if imp.currentSheet == "" {
		return nil, errors.New("未选择任何活动的工作表")
	}

	result := &ExcelImportResult{
		SheetName: imp.currentSheet,
	}

	if len(imp.currentRows) <= imp.headerRowIndex+1 {
		return result, nil
	}

	startRow := imp.headerRowIndex + 1
	for rowIndex := startRow; rowIndex < len(imp.currentRows); rowIndex++ {
		if ctx != nil {
			if err := ctx.Err(); err != nil {
				return nil, err
			}
		}

		row := imp.currentRows[rowIndex]
		if imp.isEmptyRow(row) {
			result.Skipped++
			continue
		}

		action, err := handler(rowIndex, row)
		if err != nil {
			result.Failed++
			result.Errors = append(result.Errors, fmt.Sprintf("第%d行导入失败: %v", rowIndex+1, err))
			continue
		}

		switch action {
		case "insert":
			result.TotalRows++
			result.Inserted++
		case "update":
			result.TotalRows++
			result.Updated++
		case "skip":
			result.Skipped++
		default:
			result.TotalRows++
		}
	}

	return result, nil
}

// BuildExcelHeaderMap 建立列名 → 列索引的映射
func BuildExcelHeaderMap(headerRow []string) map[string]int {
	m := make(map[string]int, len(headerRow))
	for i, h := range headerRow {
		key := strings.TrimSpace(h)
		key = strings.ReplaceAll(key, "\n", "")
		key = strings.ReplaceAll(key, " ", "")
		key = strings.ReplaceAll(key, "（", "(")
		key = strings.ReplaceAll(key, "）", ")")
		if key != "" {
			m[key] = i
		}
	}
	return m
}

// searchAndBuildHeaderMap 扫描前 5 行进行智能表头寻址
func (imp *ExcelImporter) searchAndBuildHeaderMap() {
	keywords := []string{"货号", "公司货号", "产品货号", "item_no", "设计号", "design_number", "姓名", "客户名", "客户"}
	maxScanRows := 5
	if len(imp.currentRows) < maxScanRows {
		maxScanRows = len(imp.currentRows)
	}

	for i := 0; i < maxScanRows; i++ {
		m := BuildExcelHeaderMap(imp.currentRows[i])
		for _, kw := range keywords {
			colKey := imp.normalizeHeaderKey(kw)
			if _, ok := m[colKey]; ok {
				imp.headerMap = m
				imp.headerRowIndex = i
				return
			}
		}
	}

	defaultHeaderIndex := 1
	if len(imp.currentRows) <= defaultHeaderIndex {
		defaultHeaderIndex = 0
	}
	imp.headerMap = BuildExcelHeaderMap(imp.currentRows[defaultHeaderIndex])
	imp.headerRowIndex = defaultHeaderIndex
}

// normalizeHeaderKey 格式化表头 key 剔除多余格式
func (imp *ExcelImporter) normalizeHeaderKey(name string) string {
	key := strings.TrimSpace(name)
	key = strings.ReplaceAll(key, " ", "")
	key = strings.ReplaceAll(key, "（", "(")
	key = strings.ReplaceAll(key, "）", ")")
	return key
}

// isEmptyRow 判断一行是否全部为空字符串
func (imp *ExcelImporter) isEmptyRow(row []string) bool {
	if len(row) == 0 {
		return true
	}
	for _, val := range row {
		if strings.TrimSpace(val) != "" {
			return false
		}
	}
	return true
}
