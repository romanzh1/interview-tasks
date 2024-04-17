package mtc_digital

//

a := []int{1, 2, 3}
b := append(a[:1], 10)
fmt.Println(a, b) // [1, 10, 3] [1, 10]



//
func modS(s []int) {
	s[1] = 20
	s = append(s, 40)
}

func modS1(s *[]int) {
	(*s)[1] = 10
	*s = append(*s, 40) // [1, 10, 3]
}


s := []int{1, 2, 3}
modS(s)
modS1(&s)
fmt.Println(s) // [1, 10, 3, 40]





// [1,2,3,4,5] -> [1,3,5]
// [2,2,2] -> []
// [1,2,2] -> [1]
func filter(a []int) []int {
	var n int
	for i := 0; i < len(a); i++{
		if a[i] % 2 == 0{
			a[n] = a[i]
			n++
		}
	}

	return a[:n]
}