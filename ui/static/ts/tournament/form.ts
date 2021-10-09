const endDate = <HTMLInputElement>document.getElementById("end-date");
const startDate = <HTMLInputElement>document.getElementById("start-date");

const isValidDatePair = function (start: string, end: string): boolean {
	const startTime = Date.parse(start);
	const endTime = Date.parse(end);
	if (isNaN(startTime) || isNaN(endTime)) {
		return true;
	} else if (startTime > endTime) {
		return false;
	} else {
		return true;
	}
}

const validateDatePair = function (start: HTMLInputElement, end: HTMLInputElement) {
	return function (_: Event) {
		if (isValidDatePair(start.value, end.value)) {
			start.setCustomValidity("");
			end.setCustomValidity("");
		} else {
			start.setCustomValidity("Start date can't be after end date");
			end.setCustomValidity("End date can't be before start date");
		}
	}
}

endDate.addEventListener("input", validateDatePair(startDate, endDate));
startDate.addEventListener("input", validateDatePair(startDate, endDate));
