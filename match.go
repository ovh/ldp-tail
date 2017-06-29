package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
)

var supportedMatchOperatorsMap = map[string]struct {
	description     string
	description_not string
}{
	"present": {"Field %[1]q is present", "Field %[1]q is not present"},
	"begin":   {"Field %[1]q begins with %[2]v", "Field %[1]q doest not begins with %[2]v"},
	"contain": {"Field %[1]q contains %[2]v", "Field %[1]q doest not contains %[2]v"},
	"lt":      {"Field %[1]q is less than %[2]v", "Field %[1]q is greater than or equal to %[2]v"},
	"le":      {"Field %[1]q is less than or equal to %[2]v", "Field %[1]q is greater than %[2]v"},
	"eq":      {"Field %[1]q is equal to %[2]v", "Field %[1]q is not equal to %[2]v"},
	"ge":      {"Field %[1]q is greater than or equal to %[2]v", "Field %[1]q is less than %[2]v"},
	"gt":      {"Field %[1]q is greater than %[2]v", "Field %[1]q is less than or equal to %[2]v"},
}

var supportedMatchOperators = func() []string {
	a := make([]string, 0, len(supportedMatchOperatorsMap))

	for k := range supportedMatchOperatorsMap {
		a = append(a, k)
	}

	return a
}()

type matchCriterion struct {
	Key      string
	Operator string
	Value    interface{}
	Not      bool
}

func isValidMatchCriteria(criteria []matchCriterion) (bool, error) {
	for k, m := range criteria {
		if _, ok := supportedMatchOperatorsMap[m.Operator]; !ok {
			return false, fmt.Errorf("invalid operator %q in %+v", m.Operator, m)
		}
		if m.Operator == "begin" {
			if _, ok := m.Value.(string); !ok {
				return false, fmt.Errorf("invalid value for operator 'begin' in %+v", m)
			}
		}
		if m.Operator == "lt" || m.Operator == "le" || m.Operator == "ge" || m.Operator == "gt" {
			f, err := toNumber(m.Value)
			if err != nil {
				return false, fmt.Errorf("invalid value for operator '%s' in %+v: %s", m.Operator, m, err.Error())
			}
			// Overwrite value
			criteria[k].Value = f
		}
	}

	return true, nil
}

func match(value map[string]interface{}, criteria []matchCriterion) bool {
	for _, m := range criteria {
		v, ok := value[m.Key]

		// Key is not present, don't match
		if !ok && m.Operator != "present" {
			return false
		}

		switch m.Operator {
		case "present":
			if !ok != m.Not {
				return false
			}
		case "eq":
			if (v != m.Value) != m.Not {
				return false
			}
		case "begin":
			vString, ok := v.(string)
			if !ok {
				// Stringify value
				vString = fmt.Sprintf("%v", v)
			}
			if (!strings.HasPrefix(vString, m.Value.(string))) != m.Not {
				return false
			}
		case "contain":
			vString, ok := v.(string)
			if !ok {
				// Stringify value
				vString = fmt.Sprintf("%v", v)
			}
			if (!strings.Contains(vString, m.Value.(string))) != m.Not {
				return false
			}
		case "lt":
			vFloat, err := toNumber(v)
			if err != nil {
				log.Printf("lt: incompatible value for comparison: %s", err.Error())
			}
			if (vFloat >= m.Value.(float64)) != m.Not {
				return false
			}
		case "le":
			vFloat, err := toNumber(v)
			if err != nil {
				log.Printf("lt: incompatible value for comparison: %s", err.Error())
			}
			if (vFloat > m.Value.(float64)) != m.Not {
				return false
			}
		case "ge":
			vFloat, err := toNumber(v)
			if err != nil {
				log.Printf("lt: incompatible value for comparison: %s", err.Error())
			}
			if (vFloat < m.Value.(float64)) != m.Not {
				return false
			}
		case "gt":
			vFloat, err := toNumber(v)
			if err != nil {
				log.Printf("lt: incompatible value for comparison: %s", err.Error())
			}
			if (vFloat <= m.Value.(float64)) != m.Not {
				return false
			}
		default:
			panic("Unhandled operator")
		}
	}
	return true
}

func toNumber(v interface{}) (float64, error) {
	switch value := v.(type) {
	case string:
		// Try to parse value as float64
		f, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return 0, fmt.Errorf("'%v' can't be parsed as a number", v)
		}
		return f, nil
	case uint8:
		return float64(value), nil
	case uint16:
		return float64(value), nil
	case uint32:
		return float64(value), nil
	case uint64:
		return float64(value), nil
	case int8:
		return float64(value), nil
	case int16:
		return float64(value), nil
	case int32:
		return float64(value), nil
	case int64:
		return float64(value), nil
	case float32:
		return float64(value), nil
	case float64:
		return value, nil
	default:
		return 0, fmt.Errorf("can't parse type %T as a number", v)
	}
}
