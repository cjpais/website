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


// form submission
var form = document.forms.namedItem("newmoment");
form.addEventListener("submit", function(e) {
	e.preventDefault();

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
}, false);

function dropHandler(e) {
	e.preventDefault();
	dt = e.dataTransfer;
	files = dt.files;
	handleFiles(files)
}

function handleFiles(fs) {
	files = [...files]
	files.forEach(uploadFile)
	previewFile(files[0])
}

function previewFile(file) {
	let reader = new FileReader()
	reader.readAsDataURL(file)
	reader.onloadend = function() {
		let img = document.createElement('img')
		img.src = reader.result
		// remove text
		var msg = document.getElementById("dmsg");
		msg.parentNode.removeChild(msg);
		document.getElementById('dropzone').appendChild(img)
	}
}

function uploadFile(file) {
	console.log(file)
}

function clickHandler(e) {
	console.log("hi");
}

function dragOverHandler(e) {
	e.preventDefault();
	// change the text in the dropzone to be a +
	message = document.getElementById("dmsg");
	message.innerHTML = '+';
}
