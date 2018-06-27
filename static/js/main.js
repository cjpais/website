window.onclick = function(event) {
	console.log("hit")
	document.getElementById("dropdown").setAttribute("hidden");
}

// Define a new component called button-counter
Vue.component('post', {
	props: ['post'],
	template: "#post-component"
})

Vue.component('photo', {
	props: ['post'],
	template: "#photo-component"
})

var app = new Vue({
	el: '#app',
	data () {
		return {
			data: null,
			auth: false,
			latestIsToday: false,
			date: null,
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
		newPicker: function () {
			document.getElementById("dropdown").removeAttribute("hidden")
			console.log(this)
		}
	}
})
