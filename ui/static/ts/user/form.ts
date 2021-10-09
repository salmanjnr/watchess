const applyValidation = function (fields: HTMLInputElement[], ...args: ((e: Event) => string)[]) {
	return function (e: Event) {
		var res = "";
		for (var i in args) {
			var v = args[i](e);
			if (v != "") {
				res = v;
				break;
			}
		}
		for (i in fields) {
			fields[i].setCustomValidity(res);
		}
	}
}

const validateEmail = function (field: HTMLInputElement) {
	return function (_: Event) {
		if (field.validity.typeMismatch) {
			return "Invalid email address";
		} else {
			return "";
		}
	}
}

const validateMatchingPair = function (field1: HTMLInputElement, field2: HTMLInputElement) {
	return function (_: Event) {
		if (field1.value == field2.value) {
			return "";
		} else {
			return `${field1.name} and ${field2.name} don't match`;
		}
	}
}

const validateMinLength = function (field: HTMLInputElement, minLength: number) {
	return function (_: Event) {
		if (field.value.length < minLength) {
			return `This field is too short (minimum is ${minLength})`;
		} else {
			return "";
		}
	}
}

const email = <HTMLInputElement>document.getElementById("email");
const pass = <HTMLInputElement>document.getElementById("password");
const passConfirmation = <HTMLInputElement>document.getElementById("confirm-password");

email.addEventListener("input", applyValidation([email], validateEmail(email)));
pass.addEventListener("input", applyValidation(
 	[pass, passConfirmation],
	validateMatchingPair(pass, passConfirmation),
	validateMinLength(pass, 8)
));
passConfirmation.addEventListener("input", applyValidation(
	[pass, passConfirmation],
	validateMatchingPair(pass, passConfirmation),
	validateMinLength(passConfirmation, 8)
));
