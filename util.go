package main


func ReverseString(str string) string {
	newstr := ""
	for i := len(str)-1; i >= 0; i-- {
		newstr += string(str[i])
	}
	return newstr
}
