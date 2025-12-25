package variables

import (
	"fmt"

	variableenumtypes "mlib.com/gofy/server/enum_types/variable"
)

// Example usage demonstrating the segment implementation
func ExampleUsage() {
	// Create different types of segments
	stringSeg := NewStringSegment("Hello World")
	intSeg := NewIntegerSegment(42)
	floatSeg := NewFloatSegment(3.14159)
	noneSeg := NewNoneSegment()

	// Create an object segment
	objData := map[string]any{
		"name": "John",
		"age":  30,
	}
	objSeg := NewObjectSegment(objData)

	// Create array segments
	strArr := []string{"apple", "banana", "cherry"}
	arrStrSeg := NewArrayStringSegment(strArr)

	// All segments implement the Segmenter interface
	var segments []Segmenter = []Segmenter{
		stringSeg,
		intSeg,
		floatSeg,
		noneSeg,
		objSeg,
		arrStrSeg,
	}

	// Demonstrate polymorphism through the interface
	for _, segment := range segments {
		fmt.Printf("Type: %s, Value: %s, Size: %d bytes\n",
			segment.ValueType(),
			segment.Text(),
			segment.Size())
	}

	// Demonstrate discriminator function
	discriminator := GetSegmentDiscriminator(stringSeg)
	if discriminator != variableenumtypes.Segment_STRING {
		fmt.Printf("Discriminator for string segment: %s\n", discriminator)
	}
}
