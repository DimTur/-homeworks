package main

import "testing"

func BenchmarkInefficientLogger(b *testing.B) {
	ilogger := &InefficientLogger{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ilogger.Info1("InefficientLogger")
	}
}

func BenchmarkEfficientLogger(b *testing.B) {
	elogger := NewEfficientLogger()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		elogger.Info2("EfficientLogger")
	}
}
