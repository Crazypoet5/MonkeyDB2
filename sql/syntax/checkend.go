package syntax

func checkend(tr *TokenReader) bool {
	t := tr.Read()
	if t.Kind == "split" {
		t = tr.Read()
	}
	if t.Kind == "structs" && string(t.Raw) == ";" {
		t = tr.Read()
		for t.Kind == "split" {
			t = tr.Read()
		}
		if !tr.Empty() {
			return false
		}
		return true
	}
	return false
}
