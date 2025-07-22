package utils

import (
	"fmt"
	"testing"
	"time"
)

type Action struct {
	EnableInterrupt string `json:"enable_interrupt" mapstructure:"enable_interrupt"`
	Args            string `json:"args" mapstructure:"args"`
	Texts           string `json:"texts" mapstructure:"texts"`
}

func TestAction(t *testing.T) {
	a := Action{Args: "111", EnableInterrupt: "111", Texts: "111"}
	ret := StructToStrStrMap(a)
	ret1 := StructToMap(a, false)

	ret["args"] = "123"
	ret["enable_interrupt"] = "123"
	ret["texts"] = "123"
	var b Action
	err := StrStrMapToStruct(ret, &b)
	if err != nil {
		t.Logf("StrStrMapToStruct failed:%v\n", ret)
	} else {
		t.Logf("------b=%#v\n", b)
	}
	var b1 Action
	err = MapToStruct(ret1, &b1)
	if err != nil {
		t.Logf("MapToStruct failed:%v\n", ret)
	} else {
		t.Logf("------b1=%#v\n", b1)
	}
}

type AsrStitistic struct {
	AsrAccuracy             float64 `json:"asr_accuracy" mapstructure:"asr_accuracy"`
	AsrCorrectAccuracy      float64 `json:"asr_correct_accuracy" mapstructure:"asr_correct_accuracy"`
	AsrNlpAccuracy          float64 `json:"asr_nlp_accuracy" mapstructure:"asr_nlp_accuracy"`
	AsrCorrectNlpAccuracy   float64 `json:"asr_correct_nlp_accuracy" mapstructure:"asr_correct_nlp_accuracy"`
	Total                   uint    `json:"total" mapstructure:"total"`
	AsrRightTotal           uint    `json:"asr_right_total" mapstructure:"asr_right_total"`
	AsrNlpRightTotal        uint    `json:"asr_nlp_right_total" mapstructure:"asr_nlp_right_total"`
	AsrCorrectRightTotal    uint    `json:"asr_correct_right_total" mapstructure:"asr_correct_right_total"`
	AsrCorrectNlpRightTotal uint    `json:"asr_correct_nlp_right_total" mapstructure:"asr_correct_nlp_right_total"`
}

func TestAsrStitistic(t *testing.T) {
	a := AsrStitistic{AsrAccuracy: 1.1, AsrCorrectAccuracy: 1.1, AsrNlpAccuracy: 1.1, AsrCorrectNlpAccuracy: 1.1, Total: 1, AsrRightTotal: 1, AsrNlpRightTotal: 1, AsrCorrectRightTotal: 1, AsrCorrectNlpRightTotal: 1}
	ret := StructToStrStrMap(a)
	ret1 := StructToMap(a, false)

	ret["asr_accuracy"] = "2.2"
	ret["asr_correct_accuracy"] = "2.2"
	ret["asr_nlp_accuracy"] = "2.2"
	var b AsrStitistic
	err := StrStrMapToStruct(ret, &b)
	if err != nil {
		fmt.Printf("StrStrMapToStruct failed:%v\n", ret)
	} else {
		fmt.Printf("------b=%#v\n", b)
	}
	var b1 AsrStitistic
	err = MapToStruct(ret1, &b1)
	if err != nil {
		fmt.Printf("MapToStruct failed:%v\n", ret)
	} else {
		fmt.Printf("------b1=%#v\n", b1)
	}
}

type Dialog struct {
	UnkownCount                int    `json:"unkown_count" mapstructure:"unkown_count"`
	User                       string `json:"user" mapstructure:"user"`
	PhoneNum                   string `json:"phone_num" mapstructure:"phone_num"`
	PhoneNumWaitConfirm        bool   `json:"phone_num_wait_confirm" mapstructure:"phone_num_wait_confirm"`
	WaitConfirmedPhoneNum      string `json:"wait_confirmed_phone_num" mapstructure:"wait_confirmed_phone_num"`
	PhoneNumStat               int    `json:"phone_num_stat" mapstructure:"phone_num_stat"`
	PhoneNumWrongTimes         int    `json:"phone_num_wrong_times" mapstructure:"phone_num_wrong_times"`
	LatestConfirmIntent        string `json:"latest_confirm_intent" mapstructure:"latest_confirm_intent"`
	LatestIntent               string `json:"latest_intent" mapstructure:"latest_intent"`
	LatestMostSimilar          string `json:"latest_most_similar" mapstructure:"latest_most_similar"`
	LatestInputMsg             string `json:"latest_input_msg" mapstructure:"latest_input_msg"`
	LatestQuestion             string `json:"latest_question" mapstructure:"latest_question"`
	ContinuousAskContinueCount int    `json:"continuous_ask_continue_count" mapstructure:"continuous_ask_continue_count"`
	ConfirmIntentIdx           int    `json:"confirm_intent_idx" mapstructure:"confirm_intent_idx"`
	ConfirmStat                int    `json:"confirm_stat" mapstructure:"confirm_stat"`
	IntentRanking              string `json:"intent_ranking" mapstructure:"intent_ranking"`
	CallID                     string `json:"call_id" mapstructure:"call_id"`
	SentencePathASR            string `json:"sentence_path_asr" mapstructure:"sentence_path_asr"`
	SentencePathTTS            string `json:"sentence_path_tts" mapstructure:"sentence_path_tts"`
	CallerNumber               string `json:"caller_number" mapstructure:"caller_number"`
	QueueCode                  string `json:"queue_code" mapstructure:"queue_code"`
	RecordFile                 string `json:"record_file" mapstructure:"record_file"`
	AnswerTimes                string `json:"answer_times" mapstructure:"answer_times"`
	PredictQuestions           string `json:"predict_questions" mapstructure:"predict_questions"`
	PredictSessions            string `json:"predict_sessions" mapstructure:"predict_sessions"`
}

func TestDialog(t *testing.T) {
	a := Dialog{
		PhoneNumWaitConfirm: true,
		UnkownCount:         10,
		PhoneNum:            "123",
	}
	ret := StructToStrStrMap(a)
	fmt.Printf("-----------%#v\n", ret)
	ret1 := StructToMap(a, false)
	fmt.Printf("---------ret1=%#v\n", ret1)

	ret["unkown_count"] = "22"
	ret["phone_num_wait_confirm"] = "false"
	ret["phone_num"] = "4567"
	var b Dialog
	err := StrStrMapToStruct(ret, &b)
	if err != nil {
		fmt.Printf("StrStrMapToStruct failed:%v\n", ret)
	} else {
		fmt.Printf("------b=%#v\n", b)
	}
	var b1 Dialog
	err = MapToStruct(ret1, &b1)
	if err != nil {
		fmt.Printf("MapToStruct failed:%v\n", ret)
	} else {
		fmt.Printf("------b1=%#v\n", b1)
	}
}

type PredictSessionDetail struct {
	ID                      string     `mapstructure:"id" gorm:"primaryKey;column:id;type:varchar(36);not null" json:"id"`                                               // 主键ID
	CreatedAt               *time.Time `mapstructure:"created_at" gorm:"column:created_at;type:datetime;default:CURRENT_TIMESTAMP" json:"created_at"`                    // 创建时间
	UpdatedAt               *time.Time `mapstructure:"updated_at" gorm:"column:updated_at;type:datetime;default:CURRENT_TIMESTAMP" json:"updated_at"`                    // 更新时间
	DeletedAt               *time.Time `mapstructure:"deleted_at" gorm:"column:deleted_at;type:datetime(3)" json:"deleted_at"`                                           // 删除时间
	CreatedBy               int        `mapstructure:"created_by" gorm:"column:created_by;type:int(11)" json:"created_by"`                                               // 被谁创建
	UpdatedBy               int        `mapstructure:"updated_by" gorm:"column:updated_by;type:int(11)" json:"updated_by"`                                               // 被谁更新
	DeletedBy               int        `mapstructure:"deleted_by" gorm:"column:deleted_by;type:int(11)" json:"deleted_by"`                                               // 被谁删除
	UserID                  string     `mapstructure:"user_id" gorm:"column:user_id;type:varchar(255);not null" json:"user_id"`                                          // 接收id，表明消息来自哪里
	SessionStartAt          *time.Time `mapstructure:"session_start_at" gorm:"column:session_start_at;type:datetime;not null" json:"session_start_at"`                   // 会话开始时间
	SessionEndAt            *time.Time `mapstructure:"session_end_at" gorm:"column:session_end_at;type:datetime;not null" json:"session_end_at"`                         // 会话结束时间
	KnowledgeBaseName       string     `mapstructure:"knowledge_base_name" gorm:"column:knowledge_base_name;type:varchar(36);not null" json:"knowledge_base_name"`       // 知识库名
	InputMessage            string     `mapstructure:"input_message" gorm:"column:input_message;type:text;not null" json:"input_message"`                                // 客户输入信息
	InputMessageWithCorrect string     `mapstructure:"input_message_with_correct" gorm:"column:input_message_with_correct;type:text;" json:"input_message_with_correct"` // asr纠正结果
	CorrectMistakes         string     `mapstructure:"correct_mistakes" gorm:"column:correct_mistakes;type:text;" json:"correct_mistakes"`                               // 纠正详情
	InputRecording          string     `mapstructure:"input_recording" gorm:"column:input_recording;type:varchar(255)" json:"input_recording"`                           // 客户输入录音
	OutputMessage           string     `mapstructure:"output_message" gorm:"column:output_message;type:text;not null" json:"output_message"`                             // 机器人回复信息
	OutputRecording         string     `mapstructure:"output_recording" gorm:"column:output_recording;type:varchar(255)" json:"output_recording"`                        // 机器人回复录音
	PredictIntents          string     `mapstructure:"predict_intents" gorm:"column:predict_intents;type:varchar(255);not null" json:"predict_intents"`                  // 识别意图列表,json数组
	SesstionType            uint8      `mapstructure:"sesstion_type" gorm:"column:sesstion_type;type:tinyint(1);not null" json:"sesstion_type"`                          // 会话类型:
	//,0 -- 知识库知识;
	//,1 -- 闲聊;
	//,2 -- 流程
	CurConfirmIntentName string `mapstructure:"cur_confirm_intent_name" gorm:"column:cur_confirm_intent_name;type:varchar(255)" json:"cur_confirm_intent_name"` // 当前确认意图名
	CurConfirmIntent     string `mapstructure:"cur_confirm_intent" gorm:"column:cur_confirm_intent;type:text" json:"cur_confirm_intent"`                        // 当前确认意图
	ConfirmedTimes       uint32 `mapstructure:"confirmed_times" gorm:"column:confirmed_times;type:int(11) unsigned;not null" json:"confirmed_times"`            // 确认次数
	PhoneNum             string `mapstructure:"phone_num" gorm:"column:phone_num;type:varchar(50)" json:"phone_num"`                                            // 用户手机号
	Entities             string `mapstructure:"entities" gorm:"column:entities;type:varchar(255)" json:"entities"`                                              // 属性.json数组,格式为[{"属性名": "属性值"},{"属性名": "属性值"}]
	QuestionID           string `mapstructure:"question_id" gorm:"column:question_id;type:varchar(36);not null" json:"question_id"`                             // 问题ID
	ActionName           string `mapstructure:"action_name" gorm:"column:action_name;type:varchar(36)" json:"action_name"`                                      // 应答名
	ActionDetail         string `mapstructure:"action_detail" gorm:"column:action_detail;type:varchar(1024)" json:"action_detail"`                              // 应答详情
	IsNeededConfirm      uint8  `mapstructure:"is_needed_confirm" gorm:"column:is_needed_confirm;type:tinyint(1) unsigned zerofill" json:"is_needed_confirm"`   // 是否需要确认
	ConfirmResult        uint8  `mapstructure:"confirm_result" gorm:"column:confirm_result;type:tinyint(1)" json:"confirm_result"`                              // 意图确认结果: 1-- 是; 2 -- 否
	OperationType        uint8  `mapstructure:"operation_type" gorm:"column:operation_type;type:tinyint(1)" json:"operation_type"`                              // 后续操作类型（1-确认了意图；2-转人工，3-挂机）
	CallID               string `mapstructure:"call_id" gorm:"column:call_id;type:varchar(36)" json:"call_id"`                                                  // 通话ID
	CallerNumber         string `gorm:"column:caller_number;type:varchar(255)" json:"caller_number" mapstructure:"caller_number"`                               // 主叫号码
	CalleeNumber         string `gorm:"column:callee_number;type:varchar(255)" json:"callee_number" mapstructure:"callee_number"`                               // 被叫号码
	QueueCode            string `gorm:"column:queue_code;type:varchar(255)" json:"queue_code" mapstructure:"queue_code"`                                        // 队列号
	RecordFile           string `gorm:"column:record_file;type:varchar(255)" json:"record_file" mapstructure:"record_file"`                                     // 总录音
	IntentKind           uint8  `gorm:"column:intent_kind;" json:"intent_kind" mapstructure:"intent_kind"`
}

func TestPredictSessionDetail(t *testing.T) {
	nowtime := time.Now()
	a := PredictSessionDetail{
		ID:        "111111111",
		CreatedAt: &nowtime,
	}
	ret := StructToStrStrMap(a)
	fmt.Printf("-----------%#v\n", ret)
	ret1 := StructToMap(a, false)
	fmt.Printf("---------ret1=%#v\n", ret1)
	ret2 := StructToMap(&a, false)
	fmt.Printf("---------ret2=%#v\n", ret2)

	ret["id"] = "22"
	var b PredictSessionDetail
	err := StrStrMapToStruct(ret, &b)
	if err != nil {
		fmt.Printf("StrStrMapToStruct failed:%v\n", err)
	} else {
		fmt.Printf("------b=%#v\n", b)
	}
	var b1 PredictSessionDetail
	err = MapToStruct(ret1, &b1)
	if err != nil {
		fmt.Printf("MapToStruct failed:%v\n", err)
	} else {
		fmt.Printf("------b1=%#v\n", b1)
	}

	var b2 PredictSessionDetail
	err = MapToStruct(ret2, &b2)
	if err != nil {
		fmt.Printf("MapToStruct failed:%v\n", err)
	} else {
		fmt.Printf("------b2=%#v\n", b2)
	}
}
