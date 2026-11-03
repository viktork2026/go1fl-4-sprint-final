package spentcalories

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

// Основные константы, необходимые для расчетов.
const (
	lenStep                    = 0.65 // средняя длина шага.
	mInKm                      = 1000 // количество метров в километре.
	minInH                     = 60   // количество минут в часе.
	stepLengthCoefficient      = 0.45 // коэффициент для расчета длины шага на основе роста.
	walkingCaloriesCoefficient = 0.5  // коэффициент для расчета калорий при ходьбе
)

func parseTraining(data string) (int, string, time.Duration, error) {
	parts := strings.Split(data, ",")
	if len(parts) != 3 {
		return 0, "", 0, fmt.Errorf(`invalid number of values: %d (%q), must be 3 ("steps,activity,duration")`, len(parts), data)
	}

	stepsRaw, activity, durationRaw := parts[0], parts[1], parts[2]

	steps, err := strconv.Atoi(stepsRaw)
	if err != nil {
		return 0, "", 0, err
	}
	if steps <= 0 {
		return 0, "", 0, fmt.Errorf("nonpositive number of steps: %d", steps)
	}

	duration, err := time.ParseDuration(durationRaw)
	if err != nil {
		return 0, "", 0, fmt.Errorf("invalid duration: %s", durationRaw)
	}
	if duration <= 0 {
		return 0, "", 0, fmt.Errorf("nonpositive duration: %v", duration)
	}

	return steps, activity, duration, nil
}

func distance(steps int, height float64) float64 {
	stepLength := height * stepLengthCoefficient

	return float64(steps) * stepLength / mInKm
}

func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	if duration <= 0 {
		return 0
	}

	dist := distance(steps, height)

	return dist / duration.Hours()
}

func TrainingInfo(data string, weight, height float64) (string, error) {
	steps, activity, duration, err := parseTraining(data)
	if err != nil {
		log.Println(err)
		return "", err
	}

	dist := distance(steps, height)
	mSpeed := meanSpeed(steps, height, duration)

	var cal float64

	switch activity {
	case "Ходьба":
		cal, err = WalkingSpentCalories(steps, weight, height, duration)
	case "Бег":
		cal, err = RunningSpentCalories(steps, weight, height, duration)
	default:
		err = fmt.Errorf("неизвестный тип тренировки")
	}

	if err != nil {
		log.Println(err)
		return "", err
	}

	str := `Тип тренировки: %s
Длительность: %.2f ч.
Дистанция: %.2f км.
Скорость: %.2f км/ч
Сожгли калорий: %.2f
`

	return fmt.Sprintf(str, activity, duration.Hours(), dist, mSpeed, cal), nil
}

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 {
		return 0, fmt.Errorf("invalid number of steps: %d", steps)
	}
	if weight <= 0 {
		return 0, fmt.Errorf("invalid weight: %.2f", weight)
	}
	if height <= 0 {
		return 0, fmt.Errorf("invalid height: %.2f", height)
	}
	if duration <= 0 {
		return 0, fmt.Errorf("invalid duration: %v", duration)
	}

	mSpeed := meanSpeed(steps, height, duration)

	return (weight * mSpeed * duration.Minutes()) / minInH, nil
}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	cal, err := RunningSpentCalories(steps, weight, height, duration)
	if err != nil {
		return 0, err
	}

	return cal * walkingCaloriesCoefficient, nil
}
