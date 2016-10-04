import Vue from "vue"
import VueRouter from "vue-router"
Vue.use(VueRouter);

import App from "./components/app.vue"
import Home from "./components/home.vue"
import Settings from "./components/settings.vue"

const router = new VueRouter({
	mode: "history",
	base: __dirname,
	routes: [
		{ path: "/", component: Home },
		{ path: "/settings", component: Settings}
	]
});

new Vue({
	router,
	render: h => h(App)
}).$mount("#app");

// Is this needed?
router.push("/");

const local_zifd = "http://localhost:8080/";

// setup API routes
$.fn.api.settings.api = {
	"bootstrap": local_zifd + "self/bootstrap/{address}/"
};
