const applyValidation = function (fields, ...args) {
	return function (_) {
		res = "";
		for (i in args) {
			v = args[i]();
			console.log(v);
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

const validateEmail = function (field) {
	return function () {
		if (field.validity.typeMismatch) {
			return "Invalid email address";
		} else {
			return "";
		}
	}
}

const validateMatchingPair = function (field1, field2) {
	return function (_) {
		if (field1.value == field2.value) {
			return "";
		} else {
			return `${field1.name} and ${field2.name} don't match`;
		}
	}
}

const validateMinLength = function (field, minLength) {
	return function () {
		if (field.value.length < minLength) {
			return `This field is too short (minimum is ${minLength})`;
		} else {
			return "";
		}
	}
}

const email = document.getElementById("email");
const pass = document.getElementById("password");
const passConfirmation = document.getElementById("confirm-password");

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
