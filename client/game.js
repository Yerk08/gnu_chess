function copy(object) {
	return JSON.parse(JSON.stringify(object))
}


const queryString = window.location.search;
const urlParams = new URLSearchParams(queryString);

const token = urlParams.get('token')
const board_element = document.getElementById("board")
const board_ctx = board_element.getContext('2d');

image_from_code = {
	"B": "BlackBishop.svg",
	"C": "BlackCamel.svg",
	"G": "BlackGnu.svg",
	"K": "BlackKing.svg",
	"H": "BlackKnight.svg",
	"P": "BlackPawn.svg",
	"Q": "BlackQueen.svg",
	"R": "BlackRook.svg",

	"b": "WhiteBishop.svg",
	"c": "WhiteCamel.svg",
	"g": "WhiteGnu.svg",
	"k": "WhiteKing.svg",
	"h": "WhiteKnight.svg",
	"p": "WhitePawn.svg",
	"q": "WhiteQueen.svg",
	"r": "WhiteRook.svg",
	".": "Empty.svg",
	
	"<": "Backward.svg",
	">": "Forward.svg",
	"@": "RotateBoard.svg",
}

boards = {
	chess: {
		rows: 8,
		cols: 8,
		field: [
			"RHBQKBHR",
			"PPPPPPPP",
			"........",
			"........",
			"........",
			"........",
			"pppppppp",
			"rhbqkbhr",
		],
	},	
	almost_wildebeest: {
		rows: 10,
		cols: 11,
		field: [
			"RHCCGKQBBHR",
			"PPPPPPPPPPP",
			"...........",
			"...........",
			"...........",
			"...........",
			"...........",
			"...........",
			"ppppppppppp",
			"rhbbqkgcchr",
		],
	},
}

var rows, cols, board_elms, saved_data, tilesize = 30, svr, svc
var is_rotated = false
function redraw_board() {
	board_elms = copy(boards[saved_data.gamename].field);
	for (var i = 0; i < saved_data.current; ++i) {
		mv = saved_data.moves[i]
		i1 = mv.charCodeAt(0) - 97
		j1 = mv.charCodeAt(1) - 97
		i2 = mv.charCodeAt(2) - 97
		j2 = mv.charCodeAt(3) - 97
		board_elms[i2] = board_elms[i2].slice(0, j2) + board_elms[i1][j1] + board_elms[i2].slice(j2 + 1, cols)
		board_elms[i1] = board_elms[i1].slice(0, j1) + '.' + board_elms[i1].slice(j1 + 1, cols)
	}

	board_element.width = window.innerWidth;
	board_element.height = window.innerHeight;
	board_ctx.clearRect(0, 0, board_element.width, board_element.height);
	board_ctx.font = Math.ceil(tilesize.toString()) + "px serif"
	tilesize = Math.min(board_element.width / (cols + 2), board_element.height / rows)
	for (let i = 0; i < rows; ++i) {
		for (let j = 0; j < cols; ++j) {
			cell = board_elms[i][j]
			selected = i == svr && j == svc
			color = (i + j) % 2 == 0
			if (is_rotated) {
				i = rows - i - 1
				j = cols - j - 1
			}
			if (color) {
				board_ctx.fillStyle = "white"
				board_ctx.fillRect(j * tilesize, i * tilesize, tilesize, tilesize)
				board_ctx.fillStyle = "black"
			} else {
				board_ctx.fillRect(j * tilesize, i * tilesize, tilesize, tilesize)
			}
			if (selected) {
				board_ctx.font = Math.ceil(tilesize.toString()) + "px serif bold"
				board_ctx.fillStyle = "blue"
				board_ctx.fillText(cell, j * tilesize + 8, i * tilesize + tilesize * 0.8)
				board_ctx.fillStyle = "black"
				board_ctx.font = Math.ceil(tilesize.toString()) + "px serif"
			} else {
				if (color) {
					board_ctx.fillText(cell, j * tilesize + 8, i * tilesize + tilesize * 0.8)
				} else {
					board_ctx.fillStyle = "white"
					board_ctx.fillText(cell, j * tilesize + 8, i * tilesize + tilesize * 0.8)
					board_ctx.fillStyle = "black"
				}
			}
			if (is_rotated) {
				i = rows - i - 1
				j = cols - j - 1
			}
		}
	}
	board_ctx.fillText('<', (cols + 1) * tilesize + 8, tilesize + tilesize * 0.8)
	board_ctx.fillText('>', (cols + 1) * tilesize + 8, 2 * tilesize + tilesize * 0.8)
	board_ctx.fillText('@', (cols + 1) * tilesize + 8, 3 * tilesize + tilesize * 0.8)
}

function update_board() {
	fetch('/api/board/get?token=' + token, {
		method: 'GET'
	})
	.then(function(response) {
		return response.json();
	})
	.then(function(data) {
		if (data.error == "unknown token") {
			fetch("/api/board/createnew?gamename=almost_wildebeest")
			.then(function(req) {return req.json()})
			.then(function(data) {
				window.location.replace("/?token=" + data.token)
			})
			return
		}
		saved_data = data
		saved_data.token = token
		rows = boards[saved_data.gamename].rows
		cols = boards[saved_data.gamename].cols
		redraw_board()
	}).catch(function(err) {
		board_ctx.font = Math.ceil(tilesize.toString()) + "px serif bold"
		board_ctx.fillStyle = "purple"
		board_ctx.fillText('error connection', 8, tilesize * 0.8 + tilesize / 2)
		board_ctx.fillStyle = "black"
	})
}

setInterval(update_board, 300);

addEventListener("click", (event) => {
	var r = event.pageY
	var c = event.pageX
	r /= tilesize
	c /= tilesize
	r = Math.floor(r)
	c = Math.floor(c)
	var is_diff = false
	if (c >= cols) {
		if (r == 1) {
			saved_data.current -= 1
			if (saved_data.current < 0) {
				saved_data.current = 0
			}
			is_diff = true
		} else if (r == 2) {
			saved_data.current += 1
			if (saved_data.current > saved_data.moves.length) {
				saved_data.current = saved_data.moves.length
			}
			is_diff = true
		} else if (r == 3) {
			is_rotated = !is_rotated
		}
	} else {
		if (r < 0 || r >= rows || c < 0 || c >= cols) {
			return
		}
		if (is_rotated) {
			r = rows - r - 1
			c = cols - c - 1
		}
		if (svr) {
			if (svr == r && svc == c) {
				svr = undefined
				svc = undefined
				return
			}
			saved_data.moves = saved_data.moves.slice(0, saved_data.current)
			saved_data.current += 1
			mv = String.fromCharCode(svr + 97) + String.fromCharCode(svc + 97)
			mv = mv + String.fromCharCode(r + 97) + String.fromCharCode(c + 97)
			saved_data.moves.push(mv)
			is_diff = true
			svr = null
			svc = null
		} else {
			if (board_elms[r][c] != '.') {
				svr = r
				svc = c
			}
		}
	}
	if (is_diff) {
		fetch("/api/board/set", {
			method: "POST",
			body: JSON.stringify(saved_data)
		});
	}
	redraw_board()
});
