type GameResult = 0 | 1 | 2 | undefined;
type RoundID = string;
type MatchID = string;
type GameID = string;

interface Result {
	[key: string]: number;
}

interface Game {
	id: GameID;
	white: string;
	black: string;
	result: GameResult;
	whiteMatchSide: string;
	blackMatchSide: string;
	pgn: string;
	matchID: MatchID;
	roundID: RoundID;
}

interface Match {
	id: MatchID;
	side1: string;
	side2: string;
	roundID: RoundID;
	games: Game[];
}

const roundsCache = Object.create({})
const round = <HTMLInputElement>document.getElementById("round-select")
const matchList = document.getElementById("matchList")

const gamesEndpoint = function(roundID: RoundID) {
	return "/api/rounds/" + roundID
}

const tournamentEndpoint = function(tournamentID: string) {
	return "/api/tournaments/" + tournamentID
}

const getTournamentID = function(): string | null {
	const tournamentIDRe = new RegExp('(?<=tournaments/)[0-9]+')
	const matches = tournamentIDRe.exec(window.location.pathname)
	if ((matches) && (matches.length >= 1)) {
		return matches[0]
	}
	return null
}

// TODO: find a better way for this
const getGamesDivID = function(matchID: MatchID) {
	return "games" + matchID
}

const getMatchResultClass = function(result: Result, side1: string, side2: string) {
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

const getGameResultArr = function(result: GameResult) {
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

const getGameResultString = function(result: GameResult) {
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

const getMatchResult = function(match: Match) {
	const resInitial: Result = {
		[match.side1]: 0,
		[match.side2]: 0,
	}

	if (match.games) {
		const reducer = (acc: Result, game: Game) => {
			const currentResult = getGameResultArr(game.result)
			acc[game.whiteMatchSide] += currentResult[0]
			acc[game.blackMatchSide] += currentResult[1]
			return acc
		}
		const result = match.games.reduce(reducer, resInitial)
		return result
	}

	return resInitial
}

const matchHTML = function(match: Match) {
	const gamesDivID = getGamesDivID(match.id)
	const matchResult = getMatchResult(match)
	const matchResultClass = getMatchResultClass(matchResult, match.side1, match.side2)
	return `
		<a data-bs-toggle="collapse" href="#${gamesDivID}">
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
		</a>
		`
}

const gamesHTML = function(match: Match) {
	var gamesList = ""
	if (match.games) {
		match.games.forEach((game) => {
			var side1 = game.white
			var side2 = game.black
			var colorClass = "white"
			if (game.whiteMatchSide != match.side1) {
				[side1, side2] = [side2, side1]
				colorClass = "black"
			}
			gamesList += `
				<a>
					<div class="row game-card ${colorClass}">
						<div class="col-4 text-start">
							<div class="text-nowrap overflow-hidden">
								${side1}
							</div>
						</div>
						<div class="col-4 text-center">
							${getGameResultString(game.result)}
						</div>
						<div class="col-4 text-end">
							<div class="text-nowrap overflow-hidden">
								${side2}
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

const getRound = async function(roundID: RoundID){
	if (roundID in roundsCache) {
		return roundsCache[roundID]
	}
	const res = await fetch(gamesEndpoint(roundID))
	if (res.ok) {
		const games = await res.json()
		roundsCache[roundID] = games
		return games
	}
	return []
}

const populateRound = async function(roundID: RoundID) {
	const round = await getRound(roundID)
	if (!matchList) {
		return
	}
	matchList.innerHTML = "";
	round.matches.forEach((match: Match) => {
		matchList.innerHTML += matchHTML(match) + gamesHTML(match);
	})
}

if (round) {
	round.addEventListener("change", function(_: Event) {
		populateRound(round.value)
	})
}

populateRound(round.value)
