<template>
    <div id="app">
        <navbar :username="this.username"
            @logoutOpen="logoutPrompt = true"
        >
        </navbar>
    </div>

    <transition name="fade">
        <login v-show="!logged_in" @loggedIn="setLogin()"></login>
    </transition>
    <transition name="fade">
        <popup v-show="logoutPrompt" >
            <logout @close="logoutPrompt = false"/>
        </popup>
    </transition>
</template>

<style scoped>
 .fade-enter-active, .fade-leave-active {
  transition: opacity .5s;
}
.fade-enter, .fade-leave-to {
  opacity: 0;
}
</style>

<script>
import navbar from './components/navbar.vue'
import Login from './components/login.vue'
import { ref } from 'vue'
import axios from "axios";
import popup from "./components/popup.vue"
import logout from "./components/logout.vue"

export default {
    components: {
        navbar,
        Login,
        popup,
        logout
    

    },
    data() {
        return {
            logged_in: true,
            polling: null,
            username: "",
            logoutPrompt: false

        }
    },
    methods: {
        setLogin() {
            this.logged_in = true
            this.setUserInfo()
        },
        setUserInfo() {
            axios
             .get("/api/getSelfUser")
             .then(response => {
                 this.username = response.data.id
             })
            


        }

    },
    created: function() {
     axios
      .get("/api/login/status")
      .then(response => {
        this.logged_in = response.data.auth

        if(this.logged_in) {
            this.setUserInfo()
        }
      })
    
    }

}

</script>
