<template>
    <div class="side-menu">
        <div class="camera-list-box">
            <div class="camera-list-item" v-for="camera in cameras" :key="camera">
                <img class = "thumbnail" :src="camera.Thumbnail"/>
                <div class="camera-list-text">
                    {{camera.Name}}
                </div>
            </div>
        </div>
    </div>
</template>
<script>
import { ref } from 'vue';
import axios from 'axios';


export default {
    data() {
        return {
            cameras: []
        }
    },
    methods: {
       
    },
    mounted() {
        axios
        .get("/api/getCameraListSideMenu").then((response) => {
            this.cameras = response.data
        })
    }
}
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