package variables

import (
	"testing"

	variableenumtypes "mlib.com/gofy/server/enum_types/variable"
)

func TestSegmentImplementation(t *testing.T) {
	// Test StringSegment
	strSeg := NewStringSegment("hello")
	if strSeg.ValueType() != variableenumtypes.Segment_STRING {
		t.Errorf("Expected STRING segment type, got %v", strSeg.ValueType())
	}
	if strSeg.Text() != "hello" {
		t.Errorf("Expected 'hello', got %v", strSeg.Text())
	}

	// Test IntegerSegment
	intSeg := NewIntegerSegment(42)
	if intSeg.ValueType() != variableenumtypes.Segment_INTEGER {
		t.Errorf("Expected INTEGER segment type, got %v", intSeg.ValueType())
	}
	if intSeg.Text() != "42" {
		t.Errorf("Expected '42', got %v", intSeg.Text())
	}

	// Test FloatSegment
	floatSeg := NewFloatSegment(3.14)
	if floatSeg.ValueType() != variableenumtypes.Segment_FLOAT {
		t.Errorf("Expected FLOAT segment type, got %v", floatSeg.ValueType())
	}

	// Test NoneSegment
	noneSeg := NewNoneSegment()
	if noneSeg.ValueType() != variableenumtypes.Segment_NONE {
		t.Errorf("Expected NONE segment type, got %v", noneSeg.ValueType())
	}
	if noneSeg.Text() != "" {
		t.Errorf("Expected empty string, got %v", noneSeg.Text())
	}

	// Test ObjectSegment
	objData := map[string]any{"key": "value"}
	objSeg := NewObjectSegment(objData)
	if objSeg.ValueType() != variableenumtypes.Segment_OBJECT {
		t.Errorf("Expected OBJECT segment type, got %v", objSeg.ValueType())
	}

	// Test ArrayStringSegment
	strArr := []string{"a", "b", "c"}
	arrStrSeg := NewArrayStringSegment(strArr)
	if arrStrSeg.ValueType() != variableenumtypes.Segment_ARRAY_STRING {
		t.Errorf("Expected ARRAY_STRING segment type, got %v", arrStrSeg.ValueType())
	}

	// Test Segmenter interface compliance
	var _ Segmenter = strSeg
	var _ Segmenter = intSeg
	var _ Segmenter = floatSeg
	var _ Segmenter = noneSeg
	var _ Segmenter = objSeg
	var _ Segmenter = arrStrSeg
}

func TestGetSegmentDiscriminator(t *testing.T) {
	strSeg := NewStringSegment("test")
	discriminator := GetSegmentDiscriminator(strSeg)

	if discriminator != variableenumtypes.Segment_STRING {
		t.Errorf("Expected STRING discriminator, got %v", discriminator)
	}

	// Test with map
	mapData := map[string]any{"value_type": "string"}
	discriminator = GetSegmentDiscriminator(mapData)

	if discriminator != variableenumtypes.Segment_STRING {
		t.Errorf("Expected STRING discriminator from map, got %v", discriminator)
	}

	// Test with invalid data
	invalidData := "not a segment or map"
	discriminator = GetSegmentDiscriminator(invalidData)
	if string(discriminator) != "" {
		t.Error("Expected discriminator to be nil for invalid data")
	}
}
