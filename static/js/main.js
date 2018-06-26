function createYears(response) {
	var list = [];
	Object.values(response).map((data) => {
		list.push(data);
	})
	return list;
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
			auth: false
		}
	},
	mounted () {
		axios
		.get("/api/days")
		.then(response => {
			this.data = createYears(response.data);
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
})
