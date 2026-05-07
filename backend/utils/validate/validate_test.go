package validate

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/odysseythink/gofy/backend/models/request"
)

type PageInfoTest struct {
	PageInfo request.PageInfo
	Name     string
}

func TestVerify(t *testing.T) {
	PageInfoVerify := Rules{"Page": {NotEmpty()}, "PageSize": {NotEmpty()}, "Name": {NotEmpty()}}
	var testInfo PageInfoTest
	testInfo.Name = "test"
	testInfo.PageInfo.Page = 0
	testInfo.PageInfo.PageSize = 0
	err := Verify(testInfo, PageInfoVerify)
	if err == nil {
		t.Error("校验失败，未能捕捉0值")
	}
	testInfo.Name = ""
	testInfo.PageInfo.Page = 1
	testInfo.PageInfo.PageSize = 10
	err = Verify(testInfo, PageInfoVerify)
	if err == nil {
		t.Error("校验失败，未能正常检测name为空")
	}
	// testInfo.Name = "test"
	// testInfo.PageInfo.Page = 1
	// testInfo.PageInfo.PageSize = 10
	// err = Verify(testInfo, PageInfoVerify)
	// if err != nil {
	// 	t.Error("校验失败，未能正常通过检测")
	// }

	testmap := make(map[string]any)
	testmap["user"] = "test"
	testmap["msg"] = "test"
	testmap["empty"] = ""
	testmap["testslice"] = []string{"1", "2"}
	if !InterfaceTypeVerify(testmap, reflect.Slice) {
		t.Error("校验失败，testmap未能正常通过检测")
	}
	if err := StringMapTypeVerify(testmap, Rules{
		"user":      {RuleTypeOfField(reflect.String), NotEmpty()},
		"msg":       {RuleTypeOfField(reflect.String), NotEmpty()},
		"empty":     {RuleTypeOfField(reflect.String), NotEmpty()},
		"testslice": {RuleTypeOfField(reflect.Slice), NotEmpty()},
	}); err != nil {
		t.Error("校验失败，testmap未能正常通过检测:", err)
	}
	fmt.Println("----------", testmap["testslice"].([]string))
}
