package daysteps

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Yandex-Practicum/tracker/internal/spentcalories"
)

const (
	// Длина одного шага в метрах
	stepLength = 0.65
	// Количество метров в одном километре
	mInKm = 1000
)

func parsePackage(data string) (int, time.Duration, error) {
	parts := strings.Split(data, ",")
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf(`invalid number of values: %d (%q), must be 2 ("steps,duration")`, len(parts), data)
	}

	stepsRaw, durationRaw := parts[0], parts[1]

	steps, err := strconv.Atoi(stepsRaw)
	if err != nil {
		return 0, 0, err
	}
	if steps <= 0 {
		return 0, 0, fmt.Errorf("nonpositive number of steps: %d", steps)
	}

	duration, err := time.ParseDuration(durationRaw)
	if err != nil {
		return 0, 0, fmt.Errorf("invalid duration: %s", durationRaw)
	}
	if duration <= 0 {
		return 0, 0, fmt.Errorf("nonpositive duration: %v", duration)
	}

	return steps, duration, nil
}

func DayActionInfo(data string, weight, height float64) string {
	steps, duration, err := parsePackage(data)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return ""
	}

	if steps <= 0 || weight <= 0 || height <= 0 {
		return ""
	}

	dist := float64(steps) * stepLength / mInKm
	cal, err := spentcalories.WalkingSpentCalories(steps, weight, height, duration)

	str := `Количество шагов: %d.
Дистанция составила %.2f км.
Вы сожгли %.2f ккал.
`
	return fmt.Sprintf(str, steps, dist, cal)
}
