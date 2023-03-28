<template>

  
    <div class="side-menu">
      <div class="camera-list-box">
        <div class="camera-list-item" v-for="camera in cameras" :key="camera" @click="openPlayer(camera.cameraID)">
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
      openPlayer(id) {
          this.emitter.emit("playRecorded", id)
        }
    },
    mounted() {
      axios.get('/api/getCameraListSideMenu').then((response) => {
        this.cameras = response.data;
      });
    },
  };
  </script>

<style scoped>

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