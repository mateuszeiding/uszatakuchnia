package services

import (
	"regexp"
	"strconv"
	"strings"
)

// Quantity — sparsowana ilość składnika. Dla pojedynczej liczby Min == Max.
type Quantity struct {
	Empty bool
	Min   float64
	Max   float64
}

var unicodeFractions = map[string]float64{"½": 0.5, "¼": 0.25, "¾": 0.75, "⅓": 1.0 / 3, "⅔": 2.0 / 3}

var (
	reWholeUnicode = regexp.MustCompile(`^(\d+)\s*([½¼¾⅓⅔])$`)
	reMixed        = regexp.MustCompile(`^(\d+)\s+(\d+)/(\d+)$`)
	reFraction     = regexp.MustCompile(`^(\d+)/(\d+)$`)
	reDecimal      = regexp.MustCompile(`^\d+([.,]\d+)?$`)
	reRangeSep     = regexp.MustCompile(`\s*[-–—]\s*`)
)

func parseNumber(s string) (float64, bool) {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0, false
	}
	if v, ok := unicodeFractions[s]; ok {
		return v, true
	}
	atoi := func(x string) float64 { n, _ := strconv.Atoi(x); return float64(n) }
	if m := reWholeUnicode.FindStringSubmatch(s); m != nil {
		return atoi(m[1]) + unicodeFractions[m[2]], true
	}
	if m := reMixed.FindStringSubmatch(s); m != nil && atoi(m[3]) != 0 {
		return atoi(m[1]) + atoi(m[2])/atoi(m[3]), true
	}
	if m := reFraction.FindStringSubmatch(s); m != nil && atoi(m[2]) != 0 {
		return atoi(m[1]) / atoi(m[2]), true
	}
	if reDecimal.MatchString(s) {
		v, err := strconv.ParseFloat(strings.Replace(s, ",", ".", 1), 64)
		return v, err == nil
	}
	return 0, false
}

// ParseQuantity — liczby, ułamki (½, 1/4, 1 1/2, 1,5) i zakresy (2–3). Port parseQty z prototypu (cr2-ui.jsx).
func ParseQuantity(raw string) (Quantity, bool) {
	s := strings.TrimSpace(raw)
	if s == "" {
		return Quantity{Empty: true}, true
	}
	parts := reRangeSep.Split(s, -1)
	if len(parts) == 2 {
		a, okA := parseNumber(parts[0])
		b, okB := parseNumber(parts[1])
		if okA && okB && b > a {
			return Quantity{Min: a, Max: b}, true
		}
		return Quantity{}, false
	}
	v, ok := parseNumber(s)
	if !ok || v <= 0 {
		return Quantity{}, false
	}
	return Quantity{Min: v, Max: v}, true
}
