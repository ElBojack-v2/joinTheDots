package main

import ("fmt"
 "os"
 "bufio"
 "sort"
 "strings"
)

type position struct {
	x, y int
	char rune
}

func (p position) compare(pp position) bool {
	if p.x == pp.x && p.y == pp.y {
		return true
	} else {
		return false	
	}
}

func (p position) onSameLine(pp position) bool {
	if p.x == pp.x || p.y == pp.y {
		return true
	} else {
		return false	
	}
}

func generatePoints(p1 position, p2 position, r_pos *[]position) {
	x := p2.x - p1.x
	y := p2.y - p1.y
	x_isneg := x < 0
	y_isneg := y < 0
	xbal := 1
	ybal := 1

	if (x_isneg) {
		x *= -1
		xbal = -1
	}
	if (y_isneg) {
		y *= -1
		ybal = -1
	}

	if (x == 0) {// Same vertical line
		for i:= 1; i < y; i++ {
			*r_pos = append(*r_pos, position{p1.x, p1.y + (i * ybal), '|'})
		}
	} else if (y == 0) {// Same Horizontal line
		for i:= 1; i < x; i++ {
			*r_pos = append(*r_pos, position{p1.x + (i * xbal), p1.y, '-'})
		}
	} else { // Diagonal
		for i:= 1; i < y; i++ {
			if (p1.x > p2.x && p1.y < p2.y || p1.x < p2.x && p1.y > p2.y) { // slash
				*r_pos = append(*r_pos, position{p1.x + (i * xbal), p1.y + (i * ybal), '/'})
			} else { // backslash
				*r_pos = append(*r_pos, position{p1.x + (i * xbal), p1.y + (i * ybal), '\\'})
			}
		}
	}
}

func replacePoints(board *[]string, replace *[]position) {
	for _, rep := range *replace {
		for y, p := range *board {
			for x, c := range p {
				if (x == rep.x && y == rep.y) {
					if ((*board)[y][x] == '.') {
						(*board)[y] = (*board)[y][:x] + string(rep.char) + (*board)[y][x+1:]
					} else if (((*board)[y][x] == '/' && rep.char == '\\') || ((*board)[y][x] == '\\' && rep.char == '/')) {
						(*board)[y] = (*board)[y][:x] + string('X') + (*board)[y][x+1:]
					} else if (((*board)[y][x] == '|' && rep.char == '-') || ((*board)[y][x] == '-' && rep.char == '|')) {
						(*board)[y] = (*board)[y][:x] + string('+') + (*board)[y][x+1:]
					}
					break
				}
				_ = c
			}
		}
	}
}

func replacePinPoints(board *[]string, spePos *[]position) {
	for _, sp := range *spePos {
		(*board)[sp.y] = (*board)[sp.y][:sp.x] + string('o') + (*board)[sp.y][sp.x+1:]	
	}

	for y, line := range *board {
		tmp := strings.Replace(line, ".", " ", -1)
		(*board)[y] = strings.TrimRight(tmp, " ")
	}
}

func replaceStar(board *[]string, spePos *[]position) {
	y := len(*board)
	x := len((*board)[0])

	for i := 0; i < y; i++ {
		for j := 0; j < x; j++ {
			count := 0
			for _, b := range *spePos {
				if b.x == j && b.y == i {
					count++
				}

				if count > 2 {
					(*board)[i] = (*board)[i][:j] + string('*') + (*board)[i][j+1:]	
					break;
				}
			}
		}

	}
}

func main() {
    scanner := bufio.NewScanner(os.Stdin)
    scanner.Buffer(make([]byte, 1000000), 1000000)

	//init Values
    var H, W int
    scanner.Scan()
    fmt.Sscan(scanner.Text(),&H, &W)
    spePos := make([]position, 0)
	board := make([]string,0)

	// Get line and joints
    for i := 0; i < H; i++ {
        scanner.Scan()
        row := scanner.Text()
		for j, char := range row {
			switch char {
			case '.':
			default:
				spePos = append(spePos, position{j, i, char})
			}
		}
		board = append(board, row)
    }

	if len(spePos) < 2 {
		return;
	}

	// sort joints
	sort.Slice(spePos, func(i, j int) bool {
		return spePos[i].char < spePos[j].char
	})

	//Compute directions + store join lines
	current := spePos[0]
	next := spePos[1]
	newSymbol := make([]position, 0)

	for i:=2; i <= len(spePos);i++ { // When spePos Over
		generatePoints(current, next, &newSymbol) // Get Points between
		current = next
		if (i < len(spePos)) {
			next = spePos[i]
		}
	}
	replacePoints(&board, &newSymbol)// Replace point (mixing cross rule)
	replaceStar(&board, &newSymbol)// Replace Star Point (+3 time Cross)
	replacePinPoints(&board, &spePos)// Replace Letter and Digits

    // fmt.Fprintln(os.Stderr, "Debug messages...")
	for _,l := range board {
		fmt.Fprintln(os.Stdout, l)
	}
}
