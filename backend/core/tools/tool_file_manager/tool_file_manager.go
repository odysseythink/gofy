package toolfilemanager

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"mime"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/odysseythink/confy"
	"github.com/odysseythink/mlog"
	uuid "github.com/satori/go.uuid"
	"mlib.com/gofy/server/core/exceptions"
	dbengine "mlib.com/gofy/server/db_engine"
	"mlib.com/gofy/server/models"
	"mlib.com/gofy/server/utils"
	"mlib.com/gofy/server/utils/mimetypes"
)

type ToolFileManager struct{}

func (mgr *ToolFileManager) SignFile(tool_file_id string, extension string) string {
	base_url := confy.Get[string]("FILES_URL")
	file_preview_url := fmt.Sprintf("%s/files/tools/%s%s", base_url, tool_file_id, extension)

	timestamp := strconv.Itoa(int(time.Now().Unix()))
	nonce := utils.GenerateRandomHex(16)
	data_to_sign := fmt.Sprintf("file-preview|%s|%s|%s", tool_file_id, timestamp, nonce)
	secretKey := []byte(confy.Get[string]("SECRET_KEY")) // 替换为实际的密钥

	// 创建 HMAC-SHA256 签名
	sign := hmac.New(sha256.New, secretKey)
	sign.Write([]byte(data_to_sign))
	sign_bytes := sign.Sum(nil)
	encoded_sign := base64.URLEncoding.EncodeToString(sign_bytes)

	return fmt.Sprintf("%s?timestamp=%s&nonce=%s&sign=%s", file_preview_url, url.QueryEscape(timestamp), url.QueryEscape(nonce), url.QueryEscape(encoded_sign))
}
func (mgr *ToolFileManager) VerifyFile(file_id string, timestamp string, nonce string, sign string) bool {
	/*
	   verify signature
	*/
	timestamp_num, err := strconv.Atoi(timestamp)
	if err != nil {
		mlog.Errorf("invalid timestamp")
		panic(exceptions.NewValueError("invalid timestatmp"))
	}
	data_to_sign := fmt.Sprintf("file-preview|%s|%s|%s", file_id, timestamp, nonce)
	secretKey := []byte(confy.Get[string]("SECRET_KEY"))
	recalculated_sign := hmac.New(sha256.New, secretKey)
	recalculated_sign.Write([]byte(data_to_sign))
	sign_bytes := recalculated_sign.Sum(nil)
	recalculated_encoded_sign := base64.URLEncoding.EncodeToString(sign_bytes)

	// verify signature
	if sign != recalculated_encoded_sign {
		return false
	}
	current_time := time.Now().Unix()
	return int(current_time)-timestamp_num <= confy.GetWithDefault[int]("FILES_ACCESS_TIMEOUT", 300)

}
func (mgr *ToolFileManager) CreateFileByRaw(
	user_id string,
	tenant_id string,
	conversation_id string,
	file_binary []byte,
	mimetype string,
) *models.ToolFile {
	extensions, err := mime.ExtensionsByType(mimetype)
	if err != nil || len(extensions) < 1 {
		mlog.Warningf("guss extension by mimetype=%s failed:%v", mimetype, err)
		extensions = []string{".bin"}
	}
	unique_name := uuid.NewV4().String()
	filename := fmt.Sprintf("%s%s", unique_name, extensions[0])
	filepath := fmt.Sprintf("tools/%s/%s", tenant_id, filename)
	err = os.WriteFile(filepath, file_binary, 0644)
	if err != nil {
		panic(exceptions.NewValueError(fmt.Sprintf("write file(%s) failed:%v", filepath, err)))
	}

	tool_file := models.NewToolFile(
		user_id,
		tenant_id,
		conversation_id,
		filepath,
		mimetype,
		"",
		filename,
		len(file_binary),
	)
	dbengine.Instance().DB.Create(tool_file)

	return tool_file

}
func (mgr *ToolFileManager) GetFileBinary(id string) ([]byte, string) {
	/*
	   get file binary

	   :param id: the id of the file

	   :return: the binary of the file, mime type
	*/
	tool_file := new(models.ToolFile)
	err := dbengine.Instance().DB.Model(&models.ToolFile{}).Where("id = ?", id).First(tool_file).Error
	if err != nil {
		mlog.Errorf("get ToolFile failed:%v", err)
		tool_file = nil
	}

	if tool_file == nil {
		return nil, ""
	}
	blob, err := os.ReadFile(tool_file.FileKey)
	if err != nil {
		mlog.Errorf("read file(%s) failed:%v", tool_file.FileKey, err)
		return nil, ""
	}

	return blob, tool_file.Mimetype

}

func (mgr *ToolFileManager) GetFileBinaryByMessageFileID(id string) ([]byte, string) {
	/*
	   get file binary

	   :param id: the id of the file

	   :return: the binary of the file, mime type
	*/

	message_file := new(models.MessageFile)
	err := dbengine.Instance().DB.Model(&models.MessageFile{}).Where("id = ?", id).First(message_file).Error
	if err != nil {
		mlog.Errorf("get MessageFile failed:%v", err)
		message_file = nil
	}
	tool_file_id := ""
	// Check if message_file is not None
	if message_file != nil {
		// get tool file id
		if message_file.URL != "" {
			tmplist := strings.Split(message_file.URL, "/")
			tool_file_id = tmplist[len(tmplist)-1]
			// trim extension
			tool_file_id = strings.Split(tool_file_id, ".")[0]
		} else {
			tool_file_id = ""
		}
	} else {
		tool_file_id = ""
	}
	tool_file := new(models.ToolFile)
	err = dbengine.Instance().DB.Model(&models.ToolFile{}).Where("id = ?", tool_file_id).First(tool_file).Error
	if err != nil {
		mlog.Errorf("get ToolFile failed:%v", err)
		tool_file = nil
	}

	if tool_file == nil {
		return nil, ""
	}
	blob, err := os.ReadFile(tool_file.FileKey)
	if err != nil {
		mlog.Errorf("read file(%s) failed:%v", tool_file.FileKey, err)
		return nil, ""
	}

	return blob, tool_file.Mimetype

}

// func(mgr *ToolFileManager) get_file_generator_by_tool_file_id(tool_file_id string){
//         /*
//         get file binary

//         :param tool_file_id: the id of the tool file

//         :return: the binary of the file, mime type
//         */
//         tool_file = (
//             db.session.query(ToolFile)
//             .filter(
//                 ToolFile.id == tool_file_id,
//             )
//             .first()
//         )

//         if not tool_file:
//             return None, None

//         stream = storage.load_stream(tool_file.file_key)

//	        return stream, tool_file
//	}
func (mgr *ToolFileManager) CreateFileByURL(
	user_id string,
	tenant_id string,
	conversation_id string,
	file_url string,
) *models.ToolFile {
	// try to download image
	resp, err := http.Get(file_url)
	if err != nil {
		mlog.Errorf("http get failed:%v", err)
		panic(exceptions.NewValueError(fmt.Sprintf("http get failed:%v", err)))
	}
	mimetype := mimetypes.GuessType(file_url, true)
	if mimetype == "" {
		mimetype = "octet/stream"
	}
	extensions, err := mime.ExtensionsByType(mimetype)
	extension := extensions[0]
	if extension == "" {
		extension = ".bin"
	}

	unique_name := uuid.NewV4().String()
	filename := fmt.Sprintf("%s%s", unique_name, extension)
	filepath := fmt.Sprintf("tools/%s/%s", tenant_id, filename)
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	os.WriteFile(filepath, data, 0644)
	tool_file := models.NewToolFile(
		user_id,
		tenant_id,
		conversation_id,
		filepath,
		mimetype,
		file_url,
		"",
		len(data),
	)
	dbengine.Instance().DB.Create(tool_file)

	return tool_file

}

// // init tool_file_parser
// from core.file.tool_file_parser import tool_file_manager

// tool_file_manager["manager"] = ToolFileManager
