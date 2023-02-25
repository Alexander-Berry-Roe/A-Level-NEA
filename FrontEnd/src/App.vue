<template>


    

    <transition name="fade">
        <login v-show="!logged_in" @loggedIn="setLogin()"></login>
    </transition>
    <transition name="fade">
        <popup v-show="logoutPrompt" >
            <logout @close="logoutPrompt = false"/>
        </popup>
    </transition>
    <transition name="fade">
        <popup v-show="accountSettingsShow">
            <menuContainer @close="accountSettingsShow=false">
                <accountSettings :username="this.username"/>
            </menuContainer>
        </popup>
    </transition>
    <transition name="fade">
       <popup v-show="cameraSettings">
            <menuContainer @close="cameraSettings = false">
                <cameraMenu/>
            </menuContainer>
        </popup>
    </transition>
    <transition name="fade">
       <popup v-show="viewRecordedPlayer">
            <menuContainer @close="viewRecordedPlayer = false">
                <recordedPlayer/>
            </menuContainer>
        </popup>
    </transition>
    <live-video/>
        
    <div id="app">
        <navbar :username="this.username" 
            @logoutOpen="logoutPrompt = true"
        >
        </navbar>
    </div>

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

import navbar from './components/navbar.vue';
import Login from './components/login.vue';
import axios from "axios";
import popup from "./components/popup.vue";
import logout from "./components/logout.vue";
import accountSettings from "./components/accountSettings.vue";
import menuContainer from "./components/menuContainer.vue";
import VModal from 'vue-js-modal';
import Modal from './components/modal.vue';
import { draggable } from 'vuedraggable';
import liveVideo from './components/livePlayer.vue';
import cameraMenu from './components/cameraMenu.vue';
import recordedPlayer from './components/recordedPlayer.vue'


export default {
    components: {
        navbar,
        Login,
        popup,
        logout,
        accountSettings,
        menuContainer,
        VModal,
        draggable,
        Modal,
        liveVideo, 
        cameraMenu,
        recordedPlayer
    },
    data() {
        return {
            logged_in: true,
            polling: null,
            username: "",
            logoutPrompt: false,
            accountMenuList: [],
            accountSettingsShow: false,
            cameraSettings: false,
            viewRecordedPlayer: false


        }
    },
    methods: {
        setLogin() {
            this.logged_in = true
            this.setUserInfo()


        },
        setUserInfo() {
            this.emitter.emit("loggedIn")
            this.refreshUserInfo()
        },
        refreshUserInfo() {
            axios
             .get("/api/getSelfUser")
             .then(response => {
                 this.username = response.data.id
             })
        },

        openAccountMenu() {
            this.accountSettingsShow = true
        },

        openCameraMenu() {
            this.cameraSettings = true
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
    
    },
    mounted() {
        this.emitter.on("openAccountMenu", this.openAccountMenu);
        this.emitter.on("reloadUserInfo", this.refreshUserInfo);
        this.emitter.on("openCameraMenu", this.openCameraMenu)
    }

}

</script>
