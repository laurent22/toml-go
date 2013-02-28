package toml

import (
	"time"
)

func (this Document) GetArray(name string, defaultValue...[]Value) []Value {
	v, ok := this.GetValue(name)
	if !ok {
		if len(defaultValue) >= 1 {
			return defaultValue[0]
		} else {
			return make([]Value, 0)
		}
	}
	return v.AsArray()
}

func (this Document) GetString(name string, defaultValue...string) string {
	v, ok := this.GetValue(name)
	if !ok {
		if len(defaultValue) >= 1 {
			return defaultValue[0]
		} else {
			return ""
		}
	}
	return v.AsString()
}

func (this Document) GetInt(name string, defaultValue...int) int {
	v, ok := this.GetValue(name)
	if !ok {
		if len(defaultValue) >= 1 {
			return defaultValue[0]
		} else {
			return 0
		}
	}
	return v.AsInt()
}

func (this Document) GetInt8(name string, defaultValue...int8) int8 {
	v, ok := this.GetValue(name)
	if !ok {
		if len(defaultValue) >= 1 {
			return defaultValue[0]
		} else {
			return 0
		}
	}
	return v.AsInt8()
}

func (this Document) GetInt16(name string, defaultValue...int16) int16 {
	v, ok := this.GetValue(name)
	if !ok {
		if len(defaultValue) >= 1 {
			return defaultValue[0]
		} else {
			return 0
		}
	}
	return v.AsInt16()
}

func (this Document) GetInt32(name string, defaultValue...int32) int32 {
	v, ok := this.GetValue(name)
	if !ok {
		if len(defaultValue) >= 1 {
			return defaultValue[0]
		} else {
			return 0
		}
	}
	return v.AsInt32()
}

func (this Document) GetInt64(name string, defaultValue...int64) int64 {
	v, ok := this.GetValue(name)
	if !ok {
		if len(defaultValue) >= 1 {
			return defaultValue[0]
		} else {
			return 0
		}
	}
	return v.AsInt64()
}

func (this Document) GetFloat(name string, defaultValue...float64) float64 {
	v, ok := this.GetValue(name)
	if !ok {
		if len(defaultValue) >= 1 {
			return defaultValue[0]
		} else {
			return 0.0
		}
	}
	return v.AsFloat()
}

func (this Document) GetFloat32(name string, defaultValue...float32) float32 {
	v, ok := this.GetValue(name)
	if !ok {
		if len(defaultValue) >= 1 {
			return defaultValue[0]
		} else {
			return 0.0
		}
	}
	return v.AsFloat32()
}

func (this Document) GetFloat64(name string, defaultValue...float64) float64 {
	v, ok := this.GetValue(name)
	if !ok {
		if len(defaultValue) >= 1 {
			return defaultValue[0]
		} else {
			return 0.0
		}
	}
	return v.AsFloat64()
}

func (this Document) GetBool(name string, defaultValue...bool) bool {
	v, ok := this.GetValue(name)
	if !ok {
		if len(defaultValue) >= 1 {
			return defaultValue[0]
		} else {
			return false
		}
	}
	return v.AsBool()
}

func (this Document) GetDate(name string, defaultValue...time.Time) time.Time {
	v, ok := this.GetValue(name)
	if !ok {
		if len(defaultValue) >= 1 {
			return defaultValue[0]
		} else {
			return time.Now()
		}
	}
	return v.AsDate()
}

