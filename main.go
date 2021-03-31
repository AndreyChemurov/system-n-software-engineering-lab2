package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type FSM struct {
	state int
}

func (fsm *FSM) update(input int) (err error) {
	/*

		STATES

		0: "engine_off"		Если двигатель не запущен, то нельзя переключить коробку
		1: "engine_on"		Если двигатель запущен, можно переключить коробку
		2: "gearbox_stop"	Если коробка не переключена, дифф-л не вращается, и нельзя давить на газ
		3: "gearbox_run"	Если коробка переключена, дифф-л готов вращаться
		4: "wheels_off"		Если колеса не крутятся, то машина стоит на месте
		5: "wheels_on"		Если колеса крутятся, то машина едет

	*/

	switch {
	case fsm.state == 0 && input == 0:
		// Двигатель не заведен -> Завести
		// Новое состояние: "engine_on"

		fsm.setState(1)
		fmt.Println("Engine on. Possible answers: 1, 2")
	case fsm.state == 1 && input == 1:
		// Двигатель заведен -> Заглушить
		// Новое состояние: "engine_off"

		fsm.setState(0)
		fmt.Println("Engine off. Possible answers: 0")
	case fsm.state == 1 && input == 2:
		// Двигатель заведен -> Повысить передачу
		// Новое состояние: "gearbox_run"

		fsm.setState(3)
		fmt.Println("Gearbox run. Possible answers: 3, 4")
	case fsm.state == 2 && input == 1:
		// Коробка не переключена -> Заглушить двигатель
		// Новое состояние: "engine_off"

		fsm.setState(0)
		fmt.Println("Engine off. Possible answers: 0")
	case fsm.state == 2 && input == 2:
		// Коробка не переключена -> Повысить передачу
		// Новое состояние: "gearbox_run"

		fsm.setState(3)
		fmt.Println("Gearbox run. Possible answers: 3, 4")
	case fsm.state == 3 && input == 3:
		// Коробка переключена -> Понизить передачу
		// Новое состояние: "gearbox_stop"

		fsm.setState(2)
		fmt.Println("Gearbox stop. Possible answers: 1, 2")
	case fsm.state == 3 && input == 4:
		// Коробка перелючена -> Нажать газ
		// Новое состояние: "wheels_on"

		fsm.setState(5)
		fmt.Println("Wheels on. Possible answers: 5")
	case fsm.state == 4 && input == 3:
		// Колеса не вращаются -> Понизить передачу
		// Новое состояние: "gearbox_stop"

		fsm.setState(2)
		fmt.Println("Gearbox stop. Possible answers: 1, 2")
	case fsm.state == 4 && input == 4:
		// Колеса не вращаются -> Нажать газ
		// Новое состояние: "wheels_on"

		fsm.setState(5)
		fmt.Println("Wheels on. Possible answers: 5")
	case fsm.state == 5 && input == 5:
		// Колеса вращаются -> Нажать на тормоз
		// Новое состояние: "wheels_off"

		fsm.setState(4)
		fmt.Println("Wheels off. Possible answers: 3, 4")
	default:
		return errors.New("wrong input")
	}

	return nil
}

func (fsm *FSM) setState(state int) {
	fsm.state = state
}

func validate(input string) (i int64, err error) {
	i, err = strconv.ParseInt(input, 10, 64)

	if err != nil {
		return -1, err
	}

	if i < 0 {
		return -1, errors.New("Enter number from 0-5")
	}

	return i, nil
}

func main() {
	var (
		fsm = new(FSM)

		err error
	)

	/*

		Действия пользователя (input):
			0: завести двигатель
			1: заглушить двигатель
			2: повысить передачу
			3: понизить передачу
			4: нажать на газ
			5: нажать на тормоз

		Для простоты:
			Действие 2 означает, что дифф-л может вращаться (включили передачу)
			Действие 3 означает, что дифф-л не может вращаться (нет передачи)
			Действие 4 означает, что машина едет (перевод состояния в "wheels_on")
			Действие 4 означает, что машина не едет (перевод состояния в "wheels_off")

	*/

	fsm.state = 0

	for {
		userInput := make(chan string)

		go func() {
			fmt.Print("Enter action: ")

			reader := bufio.NewReader(os.Stdin)
			text, _ := reader.ReadString('\n')
			text = strings.Replace(text, "\n", "", -1)

			userInput <- text
		}()

		select {
		case input := <-userInput:
			var i int64

			if i, err = validate(input); err != nil {
				fmt.Println(err.Error())
				continue
			}

			if err = fsm.update(int(i)); err != nil {
				fmt.Println(err.Error())
				continue
			}
		case <-time.After(100 * time.Second):
			os.Exit(0)
		}
	}
}
