var app = new Vue({
	el: '#app',
	data () {
		return {
			data: null,
			auth: false,
			latestIsToday: false,
			newbool: false,
			newtime: null,
			newtype: null,
			tz: null,
			date: null,
			inc: 0,
			dropdown: false,
		}
	},
	mounted () {
		axios
		.get("/api/days")
		.then(response => {
			d = luxon.DateTime.local();
			if (response.data != null) {
				this.data = this.createYears(response.data);
				year = this.data[0].Int;
				month = this.data[0].Months[0].Int;
				day = this.data[0].Months[0].Days[0].Int;
				if (d.year == year && d.month == month && d.day == day) {
					console.log("latest is today")
					this.latestIsToday = true;
				}
			}
			this.date = d.toLocaleString(luxon.DateTime.DATE_HUGE);
		})
		axios.get("/api/auth").then(response => this.auth = response.data)
	},
	filters: {
		dateformat: function(year, month, day) {
			d = luxon.DateTime.fromObject({year: year, month: month, day: day})
			return d.toLocaleString(luxon.DateTime.DATE_HUGE)
		},
		timeformat: function(time) {
			// we are treating everything as utc time even though that's false.
			// it is the time I am experiencing at any moment
			d = luxon.DateTime.fromISO(time, {zone: 'utc'})
			return d.toLocaleString(luxon.DateTime.TIME_SIMPLE)
		}
	},
	methods: {
		createYears: function (response) {
			var list = [];
			console.log("create years")
			Object.values(response).map((data) => {
				list.push(data);
			})
			return list;
		},
		dropclick: function (e) {
			if (e.target.tagName == "A") {
				// set child invis
				e.target.parentElement.parentElement.querySelector("#dropdown").classList.add("hidden")
			} else {
				// set child vis
				e.target.parentElement.querySelector("#dropdown").classList.remove("hidden")
				e.preventDefault();
			}
			if (e.target.href == "") {
				console.log("no href")
				e.preventDefault();
			}
		},
		addMoment: function (e, type) {
			var form = e.target.parentElement.parentElement.parentElement.parentElement
			form.querySelector("#newarea").classList.remove("hidden")
			var date = luxon.DateTime.local()
			this.dropdown = false;
			this.newtype = type;
			this.newbool = true;
			this.newtime = date.toLocaleString(luxon.DateTime.TIME_SIMPLE);
			this.tz = date.zoneName;
			// remove the pre class 
			var pre = document.getElementsByClassName("pre");
			if (pre.length > 0) {
				pre[0].classList.remove("pre")
			}
		},
		submitMoment: function (e) {
			form = e.target.parentElement.parentElement;
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
			// TODO refresh the page with updated data or something
			// TODO that is hide and clear everything and then add new vue element to array
		}
	}
})

Vue.component('post', {
	props: {
		post: Object,
		n: Boolean
	},
	template: "#post-component",
	methods: {
		autosize: function (e){
			var el = e.target;
			setTimeout(function(){
				el.style.cssText = 'height:' + el.scrollHeight + 'px';
			},0);
		}
	}
})

Vue.component('photo', {
	props: ['post'],
	template: "#photo-component"
})

