package bento

import "gopkg.in/yaml.v3"

// Marshal encodes b to YAML bytes. No filesystem I/O is performed.
//
// The output is suitable for writing to a bento file or sending to another
// user. Use Validate before marshaling to catch incomplete bentos early.
func Marshal(b Bento) ([]byte, error) {
	return yaml.Marshal(b)
}

// Unmarshal decodes YAML data into a Bento. No filesystem I/O is performed.
//
// Unmarshal does not validate the decoded value; call Validate afterward when
// consuming untrusted or user-supplied input.
func Unmarshal(data []byte) (Bento, error) {
	var b Bento
	if err := yaml.Unmarshal(data, &b); err != nil {
		return Bento{}, err
	}
	return b, nil
}
