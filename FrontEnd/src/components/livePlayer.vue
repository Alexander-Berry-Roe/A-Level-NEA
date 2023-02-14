<template>

    <div class="video-box" :grid="grid" :widgets="streams" :resizable="true">
      <template v-for="stream in streams" :key="stream.id" > 
        <div class="video-player" @mouseenter="stream.controls = true" @mouseleave="stream.controls = false">
          <video class="video-player" :ref="`video-${stream.id}`" :style="{ width: stream.width + 'px' }">
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
        DraggableContainer 
      }
      
    },
    methods: {
      log(event) {
        console.log(event)
      },
      playVideo(video) {
        if (Hls.isSupported()) {
          const hlsInstance = new Hls()
          hlsInstance.loadSource(video.url)
          hlsInstance.attachMedia(this.$refs[`video-${video.id}`][0])
          hlsInstance.on(Hls.Events.MANIFEST_PARSED, () => {
            this.$refs[`video-${video.id}`][0].play()
          })
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
      } 
    },
    mounted() {
      axios
      .get("/api/getAllStreams").then((response) => {
        for (let i = 0; i < response.data.length; i++) {
          this.streams.push({id: response.data[i].id, url: response.data[i].live, controls: false, cameraName: response.data[i].Name, cordx: response.data[i].PosX, cordy: response.data[i].PosY, height: response.data[i].Height, width: response.data[i].Width})
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
    
  }

.video-player{
  float: left;
  position: relative;
}

.camera-controls {
  position: absolute;
  bottom: 0;
  width: 100%;
  background-color: rgba(255, 255, 255, 0.5);
  margin: 0px;
}
  .resize-button {
    position: absolute;
    right: 0;
    bottom: 0.1rem;
  }

  .video-box {
    position: absolute;
    top: 3rem;
    left: 0;
    width: 100%;
    margin: 0;
  }
</style>