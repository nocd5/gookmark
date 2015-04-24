package utils

import (
	"reflect"
	"testing"
)

func TestSplit0(t *testing.T) {
	source := "it is test string"
	dest, err := Split(source, ' ', '"')
	if err != nil {
		t.Errorf("Expected no error : %s", source)
	}

	expected := []string{"it", "is", "test", "string"}
	if !reflect.DeepEqual(dest, expected) {
		t.Errorf("Wrong result : expected:%s, dest:%s", expected, dest)
	}
}

func TestSplit1(t *testing.T) {
	source := "it is \"test string\""
	dest, err := Split(source, ' ', '"')
	if err != nil {
		t.Errorf("Expected no error : %s", source)
	}

	expected := []string{"it", "is", "\"test string\""}
	if !reflect.DeepEqual(dest, expected) {
		t.Errorf("Wrong result : expected:%s, dest:%s", expected, dest)
	}
}

func TestSplit2(t *testing.T) {
	source := "it    is    test     string"
	dest, err := Split(source, ' ', '"')
	if err != nil {
		t.Errorf("Expected no error : %s", source)
	}

	expected := []string{"it", "is", "test", "string"}
	if !reflect.DeepEqual(dest, expected) {
		t.Errorf("Wrong result : expected:%s, dest:%s", expected, dest)
	}
}

func TestSplit3(t *testing.T) {
	source := " it is test string "
	dest, err := Split(source, ' ', '"')
	if err != nil {
		t.Errorf("Expected no error : %s", source)
	}

	expected := []string{"it", "is", "test", "string"}
	if !reflect.DeepEqual(dest, expected) {
		t.Errorf("Wrong result : expected:%s, dest:%s", expected, dest)
	}
}

func TestSplit4(t *testing.T) {
	source := " it is t\"e\"st string "
	dest, err := Split(source, ' ', '"')
	if err != nil {
		t.Errorf("Expected no error : %s", source)
	}

	expected := []string{"it", "is", "t\"e\"st", "string"}
	if !reflect.DeepEqual(dest, expected) {
		t.Errorf("Wrong result : expected:%s, dest:%s", expected, dest)
	}
}

func TestSplit5(t *testing.T) {
	source := "これ は \"試験 文字列\" です"
	dest, err := Split(source, ' ', '"')
	if err != nil {
		t.Errorf("Expected no error : %s", source)
	}

	expected := []string{"これ", "は", "\"試験 文字列\"", "です"}
	if !reflect.DeepEqual(dest, expected) {
		t.Errorf("Wrong result : expected:%s, dest:%s", expected, dest)
	}
}

func TestSplitOtherChar0(t *testing.T) {
	source := "it,is,test,string"
	dest, err := Split(source, ',', '"')
	if err != nil {
		t.Errorf("Expected no error : %s", source)
	}

	expected := []string{"it", "is", "test", "string"}
	if !reflect.DeepEqual(dest, expected) {
		t.Errorf("Wrong result : expected:%s, dest:%s", expected, dest)
	}
}

func TestSplitOtherChar1(t *testing.T) {
	source := "it,is,\"test,string\""
	dest, err := Split(source, ',', '"')
	if err != nil {
		t.Errorf("Expected no error : %s", source)
	}

	expected := []string{"it", "is", "\"test,string\""}
	if !reflect.DeepEqual(dest, expected) {
		t.Errorf("Wrong result : expected:%s, dest:%s", expected, dest)
	}
}

func TestSplitOtherChar2(t *testing.T) {
	source := "it is 'test string'"
	dest, err := Split(source, ' ', '\'')
	if err != nil {
		t.Errorf("Expected no error : %s", source)
	}

	expected := []string{"it", "is", "'test string'"}
	if !reflect.DeepEqual(dest, expected) {
		t.Errorf("Wrong result : expected:%s, dest:%s", expected, dest)
	}
}

func TestSplitErrorExpect(t *testing.T) {
	source := "it is \"test s\"tring\""
	_, err := Split(source, ' ', '"')
	if err == nil {
		t.Errorf("Expect parse error : %s", source)
	}
}
