<template>
    <popup v-show="player">
      <menuContainer @close="player = false" >
        <video class="video-player" ref="video" @click="start()" >
          <source :src="playlist" type="application/x-mpegURL" />
        </video>
      </menuContainer>
    </popup>
  
    <div class="side-menu">
      <div class="camera-list-box">
        <div class="camera-list-item" v-for="camera in cameras" :key="camera">
          <img class="thumbnail" :src="camera.thumbnail" />
          <div class="camera-list-text">{{ camera.name }}</div>
        </div>
      </div>
    </div>
  </template>
  
  <script>
  import { ref } from 'vue';
  import axios from 'axios';
  import modal from './modal.vue';
  import menuContainer from './menuContainer.vue';
  import Hls from 'hls.js';
  
  export default {
    components: {
      modal,
      menuContainer,
    },
    data() {
      return {
        cameras: [],
        player: true,
        playlist: '/stream/record/43/0/10000000000.m3u8',
      };
    },
    methods: {
      start() {
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
      },
    },
    mounted() {
      axios.get('/api/getCameraListSideMenu').then((response) => {
        this.cameras = response.data;
      });
    },
  };
  </script>

<style scoped>
.video-player {
    width:  90%;
    margin: 1rem;
}
div.side-menu {
    position: fixed;
    height: 100%;
    width: 20rem;
    background-color: white;
    z-index: 100;
    margin: 0;
    left: 0;
    top: 0;
    z-index: -2;
}

.camera-list-box {
    position: absolute;
    top: 3.1rem;
    width: 20rem;
    cursor: pointer;
    
}
.camera-list-box {
    width: 20rem;

}
.thumbnail {
    width: 10rem;
    height: 6rem;
    border-radius: 5px;
    object-fit: cover;
}
.camera-list-text {
    position: relative;
    float: right;
    
}
</style>