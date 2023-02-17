<template>
    <div class="title">
        <h3>Camera settings</h3>
    </div>
    <div class="camera-selector">
        <div class="camera-selector-item" v-for="camera in cameras" :key="camera.CameraID">
            {{camera.Name}}
        </div>
    </div>
    <div class="option-contianer">
        <div class="option">
            <div class="option-title">
                <t class="option-title">Camera name</t>
            </div>
            <input class="option-input"> 
        </div>
        <div class="option">
            <div class="option-title">
                <t class="option-title">Camera URL</t>
                <t class="error" v-if="showUsernameError"> Invalid URL </t>
            </div>
            <input class="option-input">
        </div>
    </div>
    <button class="save-button">
        Save
    </button>
</template>

<script>
import axios from 'axios'
export default {
  data() {
    return {
        cameras: [],
        loadedSettings: {},
        customBlue: 'rgb(0, 102, 255)' 
    }
  },
  methods: {
  },
  mounted() {
    axios
    .get("/api/getCameraSettings")
    .then(response => {
        this.cameras = response.data
    })

}

}
</script>

<style scoped>
.title {
  margin: 1rem;
}

.option {
  margin-top: 0px;
  margin-bottom: 1rem;
  margin-right: 2rem;
  margin-left: 0rem;
  padding: 0.5rem;
  border-radius: 8px;
  background-color: transparent;
  border: 0.1rem;
  border-style: solid;
  border-color: rgb(0, 102, 255);
  float: left;
  width: calc(100% - 1rem);
  box-sizing: border-box;
}
input.option-input {
  margin-left: 0rem;
  margin-right: 0rem;
  width:99%;
  font-family: 'Courier New', Courier, monospace;
  border-style:solid;
  border-width: 1px;
  border-color:rgb(178, 178, 182);
  height: 1.5rem;
  border-radius: 5px;
  margin-bottom: 1px;
}
.option-title {
  font-weight: bold;
  width: auto;
  margin-bottom: 0.5rem;
}
.save-button {
  position: relative;
  margin-top: 0.5rem;
  background-color:rgb(0, 102, 255);
  bottom: 2rem;
  float: right;
  margin-right: 1rem;
  border-radius: 8px;
  outline: none;
  border: none;
  width: 3rem;
  height: 2rem;
  color: white;
  cursor: pointer;
}
.error {
  float:right;
  margin-top: 0.5rem;
  color: red;

}

.camera-selector {
    height: 35rem;
    width: 12.8rem;
    margin-left: 1rem;
    float: left;
    background-color: transparent;
    border: 0.1rem;
    border-style: solid;
    border-color: rgb(0, 102, 255);
    border-radius: 8px;
    padding-top: 0.4rem;
}
.option-contianer {
    position: relative;
    height: 35rem;
    width: calc(100% - 15rem);
    float:right;
}
.camera-selector-item {
    cursor: pointer;
}
</style>