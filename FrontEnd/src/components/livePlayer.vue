<template>

    <div class="video-box" :grid="grid" :widgets="streams" :resizable="true">
      <template v-for="stream in streams" :key="stream.id" > 
        <div v-show="players" class="video-player" @mouseenter="stream.controls = true" @mouseleave="stream.controls = false">
          <video class="video-player" :ref="`video-${stream.id}`">
            <source :src="stream.url" :id="stream.id" type="application/x-mpegURL" />
          </video>
          <div v-show="stream.controls" class="camera-controls">
            {{stream.cameraName}}
          </div>
        </div>
      </template>
    </div>
</template>

<script>
import { defineComponent } from 'vue'
import { VueDraggableNext } from 'vue-draggable-next'
import Hls from 'hls.js'
import VueResizable from 'vue-resizable'
import Vue3DraggableResizable from 'vue3-draggable-resizable'
import { DraggableContainer } from 'vue3-draggable-resizable'
import axios from "axios";
import { ResponsiveDash, DashWidget } from "vue-responsive-dash";
import popup from "./popup.vue";


  export default defineComponent({
    components: {
      draggable: VueDraggableNext,
      vueResizable: VueResizable,
      Vue3DraggableResizable,
      DraggableContainer,
      ResponsiveDash,
      DashWidget,
      popup
    },
    data() {
      return {
        streams: [],
        grid: [2,2],
        draggable: true,
        initWidth: 0,
        isResizing: false,
        DraggableContainer,
        players: true
      }
      
    },
    methods: {
      log(event) {
        console.log(event)
      },
      playVideo(video) {
        const videoElement = this.$refs[`video-${video.id}`][0]
        //Checks if native or playback via HLS.js is available 
        if (videoElement.canPlayType('application/vnd.apple.mpegurl')) {
          // Native playback is available
          console.log("Native playback enabled")
          videoElement.src = video.url
          videoElement.addEventListener('loadedmetadata', () => {
            videoElement.play()
          })
        } else if (Hls.isSupported()) {
          // HLS is supported, use HLS.js
          const hlsInstance = new Hls()
          hlsInstance.loadSource(video.url)
          hlsInstance.attachMedia(videoElement)
          hlsInstance.on(Hls.Events.MANIFEST_PARSED, () => {
            videoElement.play()
          })
        } else {
          // Neither HLS nor native playback is available
          console.error('HLS and native playback are not supported')
      }

      },
      savePlayerLocation(stream) {
        axios({
        method: 'post',
        url: '/api/setStreamLocation',
        data: {
          cameraID: Number(stream.id),
          width: stream.width,
          height: stream.height,
          posx: stream.cordx,
          posy: stream.cordy
          }
        })
      },
      restartPlayback() {
        this.streams.forEach(video => {
          const player = this.$refs[`video-${video.id}`][0]

          if (player.seekable.length > 0) {
            player.currentTime = player.seekable.end(player.seekable.length - 1);
          } else {
            player.currentTime = 0;
          }
          player.playbackRate = 0.1;

          setTimeout(() => {
            player.playbackRate = 1;
            player.play();
          }, 500);

        })
      },
      stopStreams() {
        this.streams.forEach(video => {
          const player = this.$refs[`video-${video.id}`][0]

          player.pause();
          player.currentTime = 0;
        })
      },
      restartAllStreams() {
        this.stopStreams();
        
        this.streams = []
        axios
        .get("/api/getAllStreams").then((response) => {
          for (let i = 0; i < response.data.length; i++) {
            this.streams.push({id: response.data[i].id, url: response.data[i].live, controls: false, cameraName: response.data[i].Name})

        }
        this.clicked = true
      })
      } 
    },
    mounted() {
      axios
      .get("/api/getAllStreams").then((response) => {
        for (let i = 0; i < response.data.length; i++) {
          this.streams.push({id: response.data[i].id, url: response.data[i].live, controls: false, cameraName: response.data[i].Name})
        }
      })
      document.addEventListener('click', () => {
        if (!this.clicked) {
          this.streams.forEach(video => {
            this.playVideo(video)
            this.clicked = true
            
          })
        }
      })
      document.addEventListener('visibilitychange', this.restartPlayback);
    

    },
  
  })
</script>

<style scoped>
  .playerbox {
    top: 3rem;
    position: relative;
    width: 100%;
    z-index: -3;
    
  }

.video-player{
  float: left;
  position: relative;
  max-width:1200px;
  width: 100%;
  z-index: -3;
}

.camera-controls {
  position: absolute;
  bottom: 0;
  width: 100%;
  background-color: rgba(255, 255, 255, 0.5);
  margin: 0px;
  z-index: -3;
}
  .resize-button {
    position: absolute;
    right: 0;
    bottom: 0.1rem;
    z-index: -3;
  }

  .video-box {
    position: absolute;
    top: 3rem;
    left: 0;
    width: 100%;
    margin: 0;
      z-index: -3;
  }
</style>