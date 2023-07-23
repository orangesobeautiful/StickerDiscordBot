package ginext

type EmptyResp struct{}

func isEmptyResp(v any) bool {
	switch v.(type) {
	case *EmptyResp, EmptyResp:
		return true
	}

	return false
}
