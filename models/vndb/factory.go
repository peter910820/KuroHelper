package vndb

// vndb request factory
func VndbCreate() *BasicRequest {
	results := 100
	return &BasicRequest{
		Results: &results,
	}
}
