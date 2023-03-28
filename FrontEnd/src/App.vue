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
    <popup v-show="player">
      <menuContainer @close="player = false" >
        <video class="video-player" ref="video" @click="startPlayer()" >
          <source :src="playlist" type="application/x-mpegURL" />
        </video>
      </menuContainer>
    </popup>

</template>

<style scoped>
 .fade-enter-active, .fade-leave-active {
  transition: opacity .5s;
}
.fade-enter, .fade-leave-to {
  opacity: 0;
}
.video-player {
    width:  90%;
    margin: 1rem;
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
import liveVideo from './components/livePlayer.vue';
import cameraMenu from './components/cameraMenu.vue';
import recordedPlayer from './components/recordedPlayer.vue'
import Hls from 'hls.js';


export default {
    components: {
        navbar,
        Login,
        popup,
        logout,
        accountSettings,
        menuContainer,
        VModal,
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
            viewRecordedPlayer: false,
            player: false,
            playlist: ""


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
        },
        startPlayer() {
        const videoElement = this.$refs.video;
        if (videoElement.paused || videoElement.ended) {
          if (videoElement.canPlayType('application/vnd.apple.mpegurl')) {
            // Native playback is available
            console.log('Native playback enabled');
            videoElement.src = this.playlist;
            videoElement.controls = true;
            videoElement.addEventListener('loadedmetadata', () => {
              videoElement.play();

            });
          } else if (Hls.isSupported()) {
            // HLS is supported, use HLS.js
            const hlsInstance = new Hls();
            hlsInstance.attachMedia(videoElement);
            hlsInstance.loadSource(this.playlist);
            videoElement.controls = true;
            hlsInstance.on(Hls.Events.MANIFEST_PARSED, () => {
              videoElement.play();
            });
          } else {
            // Neither HLS nor native playback is available
            console.error('HLS and native playback are not supported');
          }
        } else {
          videoElement.pause();
        }
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
        //Opens recorded video player and plays selected camera.
        this.emitter.on("playRecorded", (id) => {
            const unixtime = Math.floor(Date.now() / 1000);
            this.player = true
            this.playlist = "/stream/record/" + id + "/0/" + unixtime +".m3u8"
            setTimeout(() => this.startPlayer, 250);
            setTimeout(() => this.startPlayer(), 500); 
        });
    }

}

</script>
