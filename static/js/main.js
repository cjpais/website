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
			datestring: null,
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
				this.createYears(response.data);
			}
			this.datestring = d.toLocaleString(luxon.DateTime.DATE_HUGE);
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
			Object.values(response).map((data) => {
				list.push(data);
			})
			this.data = list;
			year = this.data[0].Int;
			month = this.data[0].Months[0].Int;
			day = this.data[0].Months[0].Days[0].Int;
			if (d.year == year && d.month == month && d.day == day) {
				this.latestIsToday = true;
			}
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
			this.newtime = date.toLocaleString(luxon.DateTime.TIME_WITH_SECONDS);
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
						//console.log("post good");
					} else if (req.status == 201) {
						resp = JSON.parse(req.response)
						app.createYears(resp)
						// TODO clean this up bad
						// that is that we need to always be selecting by #
						e.target.querySelector("#newarea").classList.add("hidden")
						this.newbool = false;
					} else {
						//console.log("post failed");
					}
				};

				req.send(data);
			}, false);
		},
		cancelAdd: function (e) {
			console.log(e)
			e.target.parentElement.classList.add("hidden")
			this.newbool = false;
			e.preventDefault();
		},
		removeMoment: function (e) {
			var c = confirm("are you sure you want to delete this moment?")
			if (c == true) {
				form = e.target.parentElement.parentElement.parentElement.parentElement;
				form.addEventListener("submit", function(e) {
					e.preventDefault()
					var data = new FormData(form)
					var req = new XMLHttpRequest();
					req.open("POST", "/api/remove", true);
					req.onload = function(ev) {
						if (req.status == 201) {
							resp = JSON.parse(req.response)
							app.createYears(resp)
						}
					};

					req.send(data);
				}, false)
			} else {
				e.preventDefault();
			}
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
	props: {
		post: Object,
		n: Boolean
	},
	template: "#photo-component",
	methods: {
		autosize: function (e){
			var el = e.target;
			setTimeout(function(){
				el.style.cssText = 'height:' + el.scrollHeight + 'px';
			},0);
		}
	}
})

