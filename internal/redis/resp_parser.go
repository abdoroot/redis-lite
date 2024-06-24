package redis

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"
)

// marshal
// unmarshal
const newLine = "\r\n"

var dataTypeMap = map[string]string{
	"+": "sstring",
	"-": "error",
	":": "integers",
	"$": "bstring",
	"*": "array",
}

func dataType(s string) string {
	if t, exist := dataTypeMap[s[:1]]; exist {
		return t
	}
	return "unknown"
}

func Marshal(s string) ([]byte, error) {
	var serializedStr strings.Builder
	//set name ahmed \r\n = new line
	//*3\r\n$3\r\nSET\r\n$4\r\nname\r\n$5\r\nahmed\r\n
	//send array of bulk string
	arr := strings.Split(s, " ")
	arrlen := len(arr)
	serializedStr.WriteString("*")                  //array symbol
	serializedStr.WriteString(strconv.Itoa(arrlen)) //array length
	serializedStr.WriteString(newLine)              //new line
	for i := 0; i < arrlen; i++ {
		/*
			Write each word
			Bulk of string $<length>\r\n<data>\r\n
		*/
		w := arr[i]
		wlen := strconv.Itoa(len(w))          //word length string
		serializedStr.WriteString("$" + wlen) //symbol +  length
		serializedStr.WriteString(newLine)    //new line \r\n
		serializedStr.WriteString(w)          //Data
		serializedStr.WriteString(newLine)    //new line \r\n
	}

	return []byte(serializedStr.String()), nil
}

func unmarshalSlice(s string) ([]any, error) {
	var (
		mtypes          string //main type
		c               int    //counter
		nextLineCounter int    //use to detemine the line num that have string
		arr             []any
	)
	
	sc := bufio.NewScanner(strings.NewReader(s))
	for sc.Scan() {
		l := sc.Text()
		if c == 0 {
			mtypes = dataType(l[:1]) //index 0 represent main data type
			if mtypes == "array" {
				llen := l[1:2] //array length
				_, err := strconv.Atoi(llen)
				if err != nil {
					fmt.Println("error converting array length to int", err)
					return arr, fmt.Errorf("error converting array length to int %v", err)
				}
			}
		} else {
			// array element
			lineType := dataType(l[:1])
			if lineType == "bstring" {
				//next line is the string
				nextLineCounter = c + 1
			}
		}

		if nextLineCounter != 0 && c == nextLineCounter {
			//add the string
			arr = append(arr, l)
		}
		c++
	}
	return arr, nil
}

func Unmarshal(s string) (any, error) {
	dtype := dataType(s[:1])
	switch dtype {
	case "array":
		return unmarshalSlice(s)
	default:
		return "", fmt.Errorf("unsupported type")
	}
}
