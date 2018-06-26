//window.onload = function() {
Date.prototype.toDateInputValue = (function() {
    var local = new Date(this);
    local.setMinutes(this.getMinutes() - this.getTimezoneOffset());
    return local.toJSON().slice(0,10);
});

console.log("loaded")
d = new Date()
document.getElementById("date").value = d.toDateInputValue();
document.getElementById("time").value = d.toTimeString().slice(0,5);
document.getElementById("tz").value = Intl.DateTimeFormat().resolvedOptions().timeZone;
//}


var form = document.forms.namedItem("newmoment");
form.addEventListener("submit", function(e) {
	console.log("SUBMIT")
	var data = new FormData(form);
	var req = new XMLHttpRequest();
	req.open("POST", "/api/new", true);
	req.onload = function(ev) {
		if (req.status == 200) {
			console.log("upload good");
		} else {
			console.log("upload failed");
		}
	};

	req.send(data);
	e.preventDefault();
}, false);
