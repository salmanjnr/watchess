const round = document.getElementById("round-select")
const matchList = document.getElementById("matchList")

const gamesEndpoint = function(roundID) {
	return "/api/rounds/" + roundID
}

// TODO: find a better way for this
const getGamesDivID = function(matchID) {
	return "games" + matchID
}

const getMatchResultClass = function(result, side1, side2) {
	if (result) {
		if (result[side1] > result[side2]) {
			return "win"
		} else if (result[side1] < result[side2]) {
			return "loss"
		} else {
			return "draw"
		}
	}
	return ""
}

const getGameResultArr = function(result) {
	switch(result) {
		case 0:
			return [1, 0]
		case 1:
			return [0.5, 0.5]
		case 2:
			return [0, 1]
		default:
			return [0, 0]
	}
}

const getGameResultString = function(result) {
	switch(result) {
		case 0:
			return "1 - 0"
		case 1:
			return "0.5 - 0.5"
		case 2:
			return "0 - 1"
		default:
			return "-"
	}
}

const getMatchResult = function(match) {
	resInitial = {
		[match.side1]: 0,
		[match.side2]: 0,
	}
	if (match.games) {
		console.log(match.games)
		reducer = (acc, game) => {
			currentResult = getGameResultArr(game.result)
			acc[game.whiteMatchSide] += currentResult[0]
			acc[game.blackMatchSide] += currentResult[1]
			return acc
		}
		result = match.games.reduce(reducer, resInitial)
		return result
	}

	return resInitial
}

const matchHTML = function(match) {
	gamesDivID = getGamesDivID(match.id)
	matchResult = getMatchResult(match)
	matchResultClass = getMatchResultClass(matchResult, match.side1, match.side2)
	return `<a data-bs-toggle="collapse" href="#${gamesDivID}">
				<div class="row match-card ${matchResultClass}">
					<div class="col-4 text-start">
						<div class="text-nowrap overflow-hidden">
						${match.side1}
						</div>
					</div>
					<div class="col-4 text-center">
						${matchResult[match.side1]} - ${matchResult[match.side2]}
					</div>
					<div class="col-4 text-end">
						<div class="text-nowrap overflow-hidden">
						${match.side2}
						</div>
					</div>
				</div>
			</a>`
}

const gamesHTML = function(match) {
	gamesList = ""
	if (match.games) {
		match.games.forEach((game) => {
			matchColors = ["white", "black"]
			if (game.whiteMatchSide != match.side1) {
				matchColors = ["black", "white"]
			}
			gamesList += `
				<a>
					<div class="row game-card ${matchColors[0]}">
						<div class="col-4 text-start">
							<div class="text-nowrap overflow-hidden">
								${game[matchColors[0]]}
							</div>
						</div>
						<div class="col-4 text-center">
							${getGameResultString(game.result)}
						</div>
						<div class="col-4 text-end">
							<div class="text-nowrap overflow-hidden">
								${game[matchColors[1]]}
							</div>
						</div>
					</div>
				</a>
			`
		})
		return `
			<div class="row match-games collapse" id="${getGamesDivID(match.id)}">
				${gamesList}
			</div>
		`
	}
	return ""
}

const getRound = async function(roundID){
	const res = await fetch(gamesEndpoint(roundID))
	if (res.ok) {
		const games = await res.json()
		return games
	}
	return []
}

const populateRound = async function(roundID) {
	const round = await getRound(roundID)
	matchList.innerHTML = "";
	round.matches.forEach((match) => {
		matchList.innerHTML += matchHTML(match) + gamesHTML(match);
	})
}

round.addEventListener("change", function(_) {
	populateRound(round.value)
})

populateRound(round.value)
