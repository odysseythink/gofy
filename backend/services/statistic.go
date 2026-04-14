package services

import (
	"fmt"
	"strconv"
	"time"

	"github.com/odysseythink/mlog"
	"mlib.com/gofy/server/core/exceptions"
	dbengine "mlib.com/gofy/server/db_engine"
	"mlib.com/gofy/server/models"
	"mlib.com/gofy/server/utils"
)

var (
	statistic_url_method_relations = map[string]string{
		"/workflow/statistics/daily-conversations":      "WorkflowDailyRunsStatistic",
		"/workflow/statistics/daily-terminals":          "WorkflowDailyTerminalsStatistic",
		"/workflow/statistics/token-costs":              "WorkflowDailyTokenCostStatistic",
		"/workflow/statistics/average-app-interactions": "WorkflowAverageAppInteractionStatistic",
		"/statistics/daily-messages":                    "DailyMessageStatistic",
		"/statistics/daily-conversations":               "DailyConversationStatistic",
		"/statistics/daily-end-users":                   "DailyTerminalsStatistic",
		"/statistics/token-costs":                       "DailyTokenCostStatistic",
		"/statistics/average-session-interactions":      "AverageSessionInteractionStatistic",
		"/statistics/user-satisfaction-rate":            "UserSatisfactionRateStatistic",
		"/statistics/average-response-time":             "AverageResponseTimeStatistic",
		"/statistics/tokens-per-second":                 "TokensPerSecondStatistic",
	}
)

type StatisticService struct {
	method_mapping map[string]func(app_model *models.App, start time.Time, end time.Time) []map[string]any
}

func (s *StatisticService) Statistic(method string, app_model *models.App, start time.Time, end time.Time) []map[string]any {
	if len(s.method_mapping) == 0 {
		s.method_mapping = map[string]func(app_model *models.App, start time.Time, end time.Time) []map[string]any{
			"WorkflowDailyRunsStatistic":             s.WorkflowDailyRunsStatistic,
			"WorkflowDailyTerminalsStatistic":        s.WorkflowDailyTerminalsStatistic,
			"WorkflowDailyTokenCostStatistic":        s.WorkflowDailyTokenCostStatistic,
			"WorkflowAverageAppInteractionStatistic": s.WorkflowAverageAppInteractionStatistic,
			"DailyMessageStatistic":                  s.DailyMessageStatistic,
			"DailyConversationStatistic":             s.DailyConversationStatistic,
			"DailyTerminalsStatistic":                s.DailyTerminalsStatistic,
			"DailyTokenCostStatistic":                s.DailyTokenCostStatistic,
			"AverageSessionInteractionStatistic":     s.AverageSessionInteractionStatistic,
			"UserSatisfactionRateStatistic":          s.UserSatisfactionRateStatistic,
			"AverageResponseTimeStatistic":           s.AverageResponseTimeStatistic,
			"TokensPerSecondStatistic":               s.TokensPerSecondStatistic,
		}
	}
	if _, ok := s.method_mapping[method]; !ok {
		if _, ok := statistic_url_method_relations[method]; ok {
			method = statistic_url_method_relations[method]
		} else {
			mlog.Errorf("unsupported method=%s", method)
			panic(exceptions.NewValueError("unsupported method=" + method))
		}
	}
	return s.method_mapping[method](app_model, start, end)
}

func (s *StatisticService) WorkflowDailyRunsStatistic(app_model *models.App, start time.Time, end time.Time) []map[string]any {
	db := dbengine.Instance().DB.Model(&models.WorkflowRun{}).Select("DATE(created_at) AS date", "COUNT(id) AS runs ").Where("app_id = ? and triggered_from =?", app_model.ID, models.WorkflowRunTriggeredFrom_APP_RUN)
	if !start.IsZero() {
		db = db.Where("created_at >= ?", start)
	}
	if !end.IsZero() {
		db = db.Where("created_at < ?", end)
	}
	db = db.Group("date").Order("date DESC")
	type Data struct {
		Date string `gorm:"column:date"`
		Runs int    `gorm:"column:runs"`
	}

	datas := []Data{}
	err := db.Scan(&datas).Error
	if err != nil {
		mlog.Errorf("exec statistic sql failed:%v", err)
		return []map[string]any{}
	}
	response_datas := []map[string]any{}
	for _, v := range datas {
		mlog.Debugf("------v=%#v", v)
		response_datas = append(response_datas, map[string]any{
			"date": v.Date, "runs": v.Runs,
		})
	}

	return response_datas

}
func (s *StatisticService) WorkflowDailyTerminalsStatistic(app_model *models.App, start time.Time, end time.Time) []map[string]any {
	db := dbengine.Instance().DB.Model(&models.WorkflowRun{}).Select("DATE(created_at) AS date", "COUNT(DISTINCT workflow_runs.created_by) AS terminal_count ").Where("app_id = ? and triggered_from =?", app_model.ID, models.WorkflowRunTriggeredFrom_APP_RUN)
	if !start.IsZero() {
		db = db.Where("created_at >= ?", start)
	}
	if !end.IsZero() {
		db = db.Where("created_at < ?", end)
	}
	db = db.Group("date").Order("date DESC")
	type Data struct {
		Date          string `gorm:"column:date"`
		TerminalCount int    `gorm:"column:terminal_count"`
	}

	datas := []Data{}
	err := db.Scan(&datas).Error
	if err != nil {
		mlog.Errorf("exec statistic sql failed:%v", err)
		return []map[string]any{}
	}
	response_datas := []map[string]any{}
	for _, v := range datas {
		mlog.Debugf("------v=%#v", v)
		response_datas = append(response_datas, map[string]any{
			"date": v.Date, "terminal_count": v.TerminalCount,
		})
	}

	return response_datas
}

func (s *StatisticService) WorkflowDailyTokenCostStatistic(app_model *models.App, start time.Time, end time.Time) []map[string]any {
	db := dbengine.Instance().DB.Model(&models.WorkflowRun{}).Select("DATE(created_at) AS date", "SUM(total_tokens) AS token_count ").Where("app_id = ? and triggered_from =?", app_model.ID, models.WorkflowRunTriggeredFrom_APP_RUN)
	if !start.IsZero() {
		db = db.Where("created_at >= ?", start)
	}
	if !end.IsZero() {
		db = db.Where("created_at < ?", end)
	}
	db = db.Group("date").Order("date DESC")
	type Data struct {
		Date       string `gorm:"column:date"`
		TokenCount int    `gorm:"column:token_count"`
	}

	datas := []Data{}
	err := db.Scan(&datas).Error
	if err != nil {
		mlog.Errorf("exec statistic sql failed:%v", err)
		return []map[string]any{}
	}
	response_datas := []map[string]any{}
	for _, v := range datas {
		mlog.Debugf("------v=%#v", v)
		response_datas = append(response_datas, map[string]any{
			"date": v.Date, "token_count": v.TokenCount,
		})
	}

	return response_datas
}
func (s *StatisticService) WorkflowAverageAppInteractionStatistic(app_model *models.App, start time.Time, end time.Time) []map[string]any {
	sql := fmt.Sprintf(`
SELECT
    AVG(sub.interactions) AS interactions,
    sub.date
FROM
    (
        SELECT
            DATE(c.created_at) AS date,
            c.created_by,
            COUNT(c.id) AS interactions
        FROM
            workflow_runs c
        WHERE
            c.app_id = '%s'
            AND c.triggered_from = 'app-run'
`, app_model.ID)

	if !start.IsZero() {
		sql += fmt.Sprintf(`	AND created_at >= '%s'`, start.Format(time.DateTime))
	}
	if !end.IsZero() {
		sql += fmt.Sprintf(`	AND created_at < '%s'`, end.Format(time.DateTime))
	}
	sql += ` GROUP BY
	date, c.created_by
) sub
GROUP BY
sub.date;`

	type Data struct {
		Date         string  `gorm:"column:date"`
		Interactions float64 `gorm:"column:interactions"`
	}

	datas := []Data{}
	err := dbengine.Instance().DB.Raw(sql).Scan(&datas).Error
	if err != nil {
		mlog.Errorf("exec statistic sql failed:%v", err)
		return []map[string]any{}
	}
	response_datas := []map[string]any{}
	for _, v := range datas {
		mlog.Debugf("------v=%#v", v)
		response_datas = append(response_datas, map[string]any{
			"date": v.Date, "interactions": v.Interactions,
		})
	}

	return response_datas
}
func (s *StatisticService) DailyMessageStatistic(app_model *models.App, start, end time.Time) []map[string]any {
	sql_query := `SELECT
    DATE(created_at) AS date,
    COUNT(*) AS message_count
FROM
    messages
WHERE
    app_id = `
	sql_query += "'" + app_model.ID + "'"

	if !start.IsZero() {
		sql_query += " AND created_at >= '" + start.Format(time.DateTime) + "'"
	}
	if !end.IsZero() {
		sql_query += " AND created_at < '" + end.Format(time.DateTime) + "'"
	}
	sql_query += " GROUP BY date ORDER BY date;"
	type Response struct {
		Date         string `gorm:"column:date" json:"date"`
		MessageCount int64  `gorm:"column:message_count" json:"message_count"`
	}
	var response_datas []Response
	err := dbengine.Instance().DB.Raw(sql_query).Scan(&response_datas).Error
	if err != nil {
		mlog.Errorf("exec sql failed:%v", err)
		return nil
	}
	rsp := []map[string]any{}
	for _, response_data := range response_datas {
		rsp = append(rsp, map[string]any{"date": response_data.Date, "message_count": response_data.MessageCount})
	}
	return rsp
}
func (s *StatisticService) DailyConversationStatistic(app_model *models.App, start, end time.Time) []map[string]any {
	sql_query := `SELECT
    DATE(created_at) AS date,
    COUNT(DISTINCT messages.conversation_id) AS conversation_count
FROM
    messages
WHERE
    app_id = `
	sql_query += "'" + app_model.ID + "'"

	if !start.IsZero() {
		sql_query += " AND created_at >= '" + start.Format(time.DateTime) + "'"
	}
	if !end.IsZero() {
		sql_query += " AND created_at < '" + end.Format(time.DateTime) + "'"
	}
	sql_query += " GROUP BY date ORDER BY date;"
	type Response struct {
		Date              string `gorm:"column:date" json:"date"`
		ConversationCount int64  `gorm:"column:conversation_count" json:"conversation_count"`
	}
	var response_datas []Response
	err := dbengine.Instance().DB.Raw(sql_query).Scan(&response_datas).Error
	if err != nil {
		mlog.Errorf("exec sql failed:%v", err)
		return nil
	}
	rsp := []map[string]any{}
	for _, response_data := range response_datas {
		rsp = append(rsp, map[string]any{"date": response_data.Date, "conversation_count": response_data.ConversationCount})
	}
	return rsp
}
func (s *StatisticService) DailyTerminalsStatistic(app_model *models.App, start, end time.Time) []map[string]any {
	sql_query := `SELECT
    DATE(created_at) AS date,
    COUNT(DISTINCT messages.from_end_user_id) AS terminal_count
FROM
    messages
WHERE
    app_id = `
	sql_query += "'" + app_model.ID + "'"

	if !start.IsZero() {
		sql_query += " AND created_at >= '" + start.Format(time.DateTime) + "'"
	}
	if !end.IsZero() {
		sql_query += " AND created_at < '" + end.Format(time.DateTime) + "'"
	}
	sql_query += " GROUP BY date ORDER BY date;"
	type Response struct {
		Date          string `gorm:"column:date" json:"date"`
		TerminalCount int64  `gorm:"column:terminal_count" json:"terminal_count"`
	}
	var response_datas []Response
	err := dbengine.Instance().DB.Raw(sql_query).Scan(&response_datas).Error
	if err != nil {
		mlog.Errorf("exec sql failed:%v", err)
		return nil
	}
	rsp := []map[string]any{}
	for _, response_data := range response_datas {
		rsp = append(rsp, map[string]any{"date": response_data.Date, "terminal_count": response_data.TerminalCount})
	}
	return rsp
}
func (s *StatisticService) DailyTokenCostStatistic(app_model *models.App, start, end time.Time) []map[string]any {
	sql_query := `SELECT
    DATE(created_at) AS date,
    (SUM(messages.message_tokens) + SUM(messages.answer_tokens)) AS token_count,
    SUM(total_price) AS total_price
FROM
    messages
WHERE
    app_id = `
	sql_query += "'" + app_model.ID + "'"

	if !start.IsZero() {
		sql_query += " AND created_at >= '" + start.Format(time.DateTime) + "'"
	}
	if !end.IsZero() {
		sql_query += " AND created_at < '" + end.Format(time.DateTime) + "'"
	}
	sql_query += " GROUP BY date ORDER BY date;"
	type Response struct {
		Date       string `gorm:"column:date" json:"date"`
		TokenCount int64  `gorm:"column:token_count" json:"token_count"`
		TokenPrice int64  `gorm:"column:total_price" json:"total_price"`
	}
	var response_datas []Response
	err := dbengine.Instance().DB.Raw(sql_query).Scan(&response_datas).Error
	if err != nil {
		mlog.Errorf("exec sql failed:%v", err)
		return nil
	}
	rsp := []map[string]any{}
	for _, response_data := range response_datas {
		rsp = append(rsp, map[string]any{"date": response_data.Date, "token_count": response_data.TokenCount, "total_price": response_data.TokenPrice, "currency": "USD"})
	}
	return rsp
}
func (s *StatisticService) AverageSessionInteractionStatistic(app_model *models.App, start, end time.Time) []map[string]any {
	sql_query := `SELECT
    DATE(c.created_at) AS date,
    AVG(subquery.message_count) AS interactions
FROM
    (
        SELECT
            m.conversation_id,
            COUNT(m.id) AS message_count
        FROM
            conversations c
        JOIN
            messages m
            ON c.id = m.conversation_id
        WHERE
            c.app_id = `
	sql_query += "'" + app_model.ID + "'"

	if !start.IsZero() {
		sql_query += " AND c.created_at >= '" + start.Format(time.DateTime) + "'"
	}
	if !end.IsZero() {
		sql_query += " AND c.created_at < '" + end.Format(time.DateTime) + "'"
	}
	sql_query += `
        GROUP BY m.conversation_id
    ) subquery
LEFT JOIN
    conversations c
    ON c.id = subquery.conversation_id
GROUP BY
    date
ORDER BY
    date;`
	type Response struct {
		Date         string  `gorm:"column:date" json:"date"`
		Interactions float64 `gorm:"column:interactions" json:"interactions"`
	}
	var response_datas []Response
	err := dbengine.Instance().DB.Raw(sql_query).Scan(&response_datas).Error
	if err != nil {
		mlog.Errorf("exec sql failed:%v", err)
		return nil
	}
	rsp := []map[string]any{}
	for _, response_data := range response_datas {
		num, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", response_data.Interactions), 64)
		rsp = append(rsp, map[string]any{"date": response_data.Date, "interactions": num})
	}
	return rsp
}
func (s *StatisticService) UserSatisfactionRateStatistic(app_model *models.App, start, end time.Time) []map[string]any {
	sql_query := `SELECT
    DATE(m.created_at) AS date,
    COUNT(m.id) AS message_count,
    COUNT(mf.id) AS feedback_count
FROM
    messages m
LEFT JOIN
    message_feedbacks mf
    ON mf.message_id=m.id AND mf.rating='like'
WHERE
    m.app_id = `
	sql_query += "'" + app_model.ID + "'"

	if !start.IsZero() {
		sql_query += " AND m.created_at >= '" + start.Format(time.DateTime) + "'"
	}
	if !end.IsZero() {
		sql_query += " AND m.created_at < '" + end.Format(time.DateTime) + "'"
	}
	sql_query += " GROUP BY date ORDER BY date;"
	type Response struct {
		Date          string `gorm:"column:date" json:"date"`
		MessageCount  int64  `gorm:"column:message_count" json:"message_count"`
		FeedbackCount int64  `gorm:"column:interactions" json:"feedback_count"`
	}
	var response_datas []Response
	err := dbengine.Instance().DB.Raw(sql_query).Scan(&response_datas).Error
	if err != nil {
		mlog.Errorf("exec sql failed:%v", err)
		return nil
	}
	rsp := []map[string]any{}
	for _, response_data := range response_datas {
		rate := 0.0
		if response_data.MessageCount > 0 {
			rate = float64(response_data.FeedbackCount*1000) / float64(response_data.MessageCount)
		}

		rsp = append(rsp, map[string]any{"date": response_data.Date, "rate": utils.RoundWithPrecision(rate, 2)})
	}
	return rsp
}
func (s *StatisticService) AverageResponseTimeStatistic(app_model *models.App, start, end time.Time) []map[string]any {
	sql_query := `SELECT
    DATE(created_at) AS date,
    AVG(provider_response_latency) AS latency
FROM
    messages
WHERE
    app_id = `
	sql_query += "'" + app_model.ID + "'"

	if !start.IsZero() {
		sql_query += " AND created_at >= '" + start.Format(time.DateTime) + "'"
	}
	if !end.IsZero() {
		sql_query += " AND created_at < '" + end.Format(time.DateTime) + "'"
	}
	sql_query += " GROUP BY date ORDER BY date;"
	type Response struct {
		Date    string  `gorm:"column:date" json:"date"`
		Latency float64 `gorm:"column:latency" json:"latency"`
	}
	var response_datas []Response
	err := dbengine.Instance().DB.Raw(sql_query).Scan(&response_datas).Error
	if err != nil {
		mlog.Errorf("exec sql failed:%v", err)
		return nil
	}
	rsp := []map[string]any{}
	for _, response_data := range response_datas {
		rsp = append(rsp, map[string]any{"date": response_data.Date, "latency": utils.RoundWithPrecision(response_data.Latency*float64(1000), 4)})
	}
	return rsp
}
func (s *StatisticService) TokensPerSecondStatistic(app_model *models.App, start, end time.Time) []map[string]any {
	sql_query := `SELECT
    DATE(created_at) AS date,
    CASE
        WHEN SUM(provider_response_latency) = 0 THEN 0
        ELSE (SUM(answer_tokens) / SUM(provider_response_latency))
    END as tokens_per_second
FROM
    messages
WHERE
    app_id = `
	sql_query += "'" + app_model.ID + "'"

	if !start.IsZero() {
		sql_query += " AND created_at >= '" + start.Format(time.DateTime) + "'"
	}
	if !end.IsZero() {
		sql_query += " AND created_at < '" + end.Format(time.DateTime) + "'"
	}
	sql_query += " GROUP BY date ORDER BY date"
	type Response struct {
		Date            string  `gorm:"column:date" json:"date"`
		TokensPerSecond float64 `gorm:"column:tokens_per_second" json:"tokens_per_second"`
	}
	var response_datas []Response
	err := dbengine.Instance().DB.Raw(sql_query).Scan(&response_datas).Error
	if err != nil {
		mlog.Errorf("exec sql failed:%v", err)
		return nil
	}
	rsp := []map[string]any{}
	for _, response_data := range response_datas {
		rsp = append(rsp, map[string]any{"date": response_data.Date, "tps": utils.RoundWithPrecision(response_data.TokensPerSecond, 4)})
	}
	return rsp
}
