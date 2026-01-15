package shortcode

const alphabet = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func encodeBase62(n int64) string {
	if n == 0 {
		return "0"
	}

	var out []byte
	for ; n > 0; n /= 62 {
		out = append(out, alphabet[n%62])
	}

	for i, j := 0, len(out)-1; i < j; i, j = i+1, j+1 {
		out[i], out[j] = out[j], out[i]
	}

	return string(out)
}

type Base62Generator struct{}

func NewBase62Generator() *Base62Generator {
	return &Base62Generator{}
}

func (gen Base62Generator) Generate(id int64) (string, error) {
	return encodeBase62(id), nil
}
