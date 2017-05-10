package main

import (
	"flag"
	"encoding/json"
	"errors"
	"os"
	"fmt"
	"math"
)

/**
example 1
./main -a="[[3,-9,3],[2,-4,4],[1,8,-18]]" -b="[-18,-10,35]"

example 2
./main -a="[[1,-3,1],[0,0,-2],[0,11,-20]]" -b="[-6,2,42]"
 */

type matrix [][]float64
type vector []float64

func main() {

	// парсинг входных данных
	err, a, b := parseInput();

	// проверка ошибок парсинга
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// индекс, определяет порядок колонок в матрице
	index := make([]int, len(a))
	for i := range index {
		index[i] = i
	}

	// отображаем исходные данные
	fmt.Println("matrix A")
	a.dump(index)
	fmt.Println("vector B")
	b.dump()

	// прямой ход
	for i := 0; i < len(a); i++ {

		// главный элемент, значение по умолчанию
		r := a[i][index[i]]

		// если главный элемент равен нулю, нужно найти другой
		// методом перестановки колонок в матрице
		if r == 0 {
			var kk int

			// двигаемся вправо от диаганаотного элемента, для поиска максимального по модулю элемента
			for k := i; k < len(a); k++ {
				if math.Abs(a[i][index[k]]) > r {
					kk = k
				}
			}

			// если удалось найти главный элемент
			if kk > 0 {
				// меняем местами колонки, так чтобы главный элемент встал в диагональ матрицы
				index[i], index[kk] = index[kk], index[i]
			}

			// получаем главный элемента, текущей строки из диагонали
			r = a[i][index[i]];
		}

		// если главный элемент строки равен 0, метод гаусса не работает
		if r == 0 {
			if b[i] == 0 {
				fmt.Println("система имеет множество решений")
			} else {
				fmt.Println("система не имеет решений")
			}
			os.Exit(1)
		}

		// деление элементов текущей строки, на главный элемент
		for j := 0; j < len(a[i]); j++ {
			a[i][index[j]] /= r
		}
		b[i] /= r

		// вычитание текущей строки из всех ниже расположенных строк с занулением I - ого элемента в каждой из них
		for k := i + 1; k < len(a); k++ {
			r = a[k][index[i]];
			for j := 0; j < len(a[i]); j++ {
				a[k][index[j]] = a[k][index[j]] - a[i][index[j]]*r
			}
			b[k] = b[k] - b[i]*r
		}

		// отображаем дамп матрицы A и вектора B
		fmt.Println("++++++++++++\n")
		fmt.Println("matrix A")
		a.dump(index)
		fmt.Println("vector B")
		b.dump()
	}

	var x vector = make(vector, len(b))

	// обратный ход
	for i := len(a) - 1; i >= 0; i-- {
		// Задается начальное значение элемента x[I].
		x[i] = b[i]

		// Корректируется искомое значение x[I].
		// В цикле по J от I+1 до N (в случае, когда I=N, этот шаг не выполняется) производятся вычисления x[I]:=  x[I] - x[J]* A[I, J].
		for j := i + 1; j < len(a); j++ {
			x[i] = x[i] - (x[j] * a[i][index[j]]);
		}
	}

	fmt.Println("++++++++++++\n")
	fmt.Println("vector X")
	for i := 0; i < len(x); i++ {
		fmt.Printf("[%v] ", x[index[i]])
	}
	fmt.Println()
}

// отображение дампа матрицы
func (a matrix) dump(index []int) {
	for i := range a {
		for j := range a[i] {
			if (a[i][index[j]] == 0) {
				// необходимо чтобы избавиться от -0
				fmt.Printf("[0] ")
			} else {
				fmt.Printf("[%v] ", a[i][index[j]])
			}
		}
		fmt.Println()
	}
	fmt.Println()
}

// отображение дампа вектора
func (b vector) dump() {

	for i := 0; i < len(b); i++ {
		fmt.Printf("[%v] ", b[i])
	}

	fmt.Println()
	fmt.Println()
}

// Парсинг входных данных
func parseInput() (error, matrix, vector) {

	// квадратная матрица размером N на N
	var a matrix

	// числовой вектор-столбец размером N
	var b vector

	// описание флагов командной строки
	aJson := flag.String("a", "[3,-9,3],[2,-4,4],[1,8,-18]", "квадратная матрица размером N на N")
	bJson := flag.String("b", "[-18,-10,35]", "числовой вектор-столбец размером N")

	// парсинг флагов командной строки
	flag.Parse()

	// парсинг матрицы a из Json
	if err1 := json.Unmarshal([]byte(*aJson), &a); err1 != nil {
		return err1, nil, nil
	}

	// парсинг вектора b из Json
	if err2 := json.Unmarshal([]byte(*bJson), &b); err2 != nil {
		return err2, nil, nil
	}

	// вылидация данных
	if len(a) < 2 || len(a) != len(b) {
		return errors.New("Не верный формат данных"), nil, nil
	}

	// вылидация данных
	for i := 0; i < len(a); i++ {
		if (len(a[i]) != len(b)) {
			return errors.New("Не верный формат данных"), nil, nil
		}
	}

	return nil, a, b;
}
