package main

import (
	"fmt"
)

// Human — Родительский класс(продавец)
type Human struct {
	Name string
	Age  int
	Role string
}

func (h Human) Introduce() {
	fmt.Printf("Меня зовут %s, и мне %d лет, моя роль: %s\n", h.Name, h.Age, h.Role)
}

func (h *Human) Birthday() {
	h.Age++
	fmt.Printf("%s отпраздновал день рождения, теперь ему %d лет\n", h.Name, h.Age)
}

// Action — Встраиваю Human
type Action struct {
	Human
	Product string
	Price   int
}

func (a *Action) AddProduct(product string, price int) {
	a.Product = product
	a.Price = price
	fmt.Printf("%s добавил товар: %s (цена: %d)\n", a.Name, a.Product, a.Price)
}

func (a Action) ShowProduct() {
	fmt.Printf("Товар: %s, цена: %d\n", a.Product, a.Price)
}

func main() {
	h := Human{Name: "Дмитрий", Age: 30, Role: "продавец"}

	a := Action{Human: h}

	a.Introduce()
	a.Birthday()

	a.AddProduct("Смартфон", 77777)
	a.ShowProduct()
}
