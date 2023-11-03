package cast

import "testing"

func TestLowerSnakeCase(t *testing.T) {

	if StringLowerSnakeCase("HelloWorld") != "hello_world" {
		t.Errorf("Error")
	}

	if StringLowerSnakeCase("HELLOWORLD") != "helloworld" {
		t.Errorf("Error")
	}

	if StringLowerSnakeCase("CAMELStuff") != "camel_stuff" {
		t.Errorf("Error")
	}

	if StringLowerSnakeCase("StuffCAMEL") != "stuff_camel" {
		t.Logf("%v\n", StringLowerSnakeCase("StuffUUID"))
		t.Errorf("Error")
	}
}
