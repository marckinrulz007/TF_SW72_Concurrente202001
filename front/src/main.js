import Vue from 'vue'
import './plugins/fontawesome'
import App from './App.vue'
import vuetify from './plugins/vuetify'
import router from './router'
import { FontAwesomeIcon } from '@fortawesome/vue-fontawesome'

Vue.config.productionTip = false
Vue.component('font-awesome-icon', FontAwesomeIcon)

new Vue({
  vuetify,
  router,
  render: h => h(App)
}).$mount('#app')
