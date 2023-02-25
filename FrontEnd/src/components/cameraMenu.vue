<template>
    <div class="title">
        <h3>Camera settings</h3>
    </div>
    <div class="camera-selector">
        <div class="camera-selector-item" v-bind:style="{ backgroundColor: camera.selectionColor }" v-for="camera in cameras" :key="camera.CameraID" @click="onSelection(camera)">
            {{camera.Name}}
        </div>
        <div class="add-button" @click="addCamera()">
          Add Camera
        </div>
    </div>
    <div class="option-container">
        <div class="option">
            <div class="option-title">
                <t class="option-title">Camera name</t>
            </div>
            <input class="option-input" v-model="selctedCamera.Name"> 
            <div class="divider"/>
            <div class="option-title">
                <t class="option-title">Camera URL</t>
                <t class="error" v-if="showUsernameError"> Invalid URL </t>
            </div>
            <div class="divider"/>
            <input class="option-input" v-model="selctedCamera.CameraUrl" spellcheck="false">
            <div class="option-title">
              <t class="option-title">Recording retention time (days)</t>
            </div>
            <input class="option-input" type="number" v-model="recordingDays"/>
        </div>
      
    </div>
    <button class="save-button">
        Save
    </button>
    <button class="delete-button" @click="deleteCamera(selctedCamera)">
      Delete
    </button>
</template>

<script>
import axios from 'axios'
export default {
  data() {
    return {
        cameras: [],
        selctedCamera: {},
        loadedSettings: {},
        customBlue: 'rgb(0, 102, 255)'
    }
  },
  methods: {
    onSelection(camera) {
      this.unSelectAll()
      camera.selectionColor = this.customBlue
      this.selctedCamera = camera
    },
    unSelectAll() {
      for (let i = 0; i < this.cameras.length; i++) {
        this.cameras[i].selectionColor = ''

      }
    },
    deleteCamera(selectedCamera) {
      const index = this.cameras.indexOf(selectedCamera);
      if (index > -1) {
        this.cameras.splice(index, 1);
        this.selctedCamera = {};
      }
    },
    addCamera() {
      this.cameras.push({Name:'New Camera'})
    },
    save() {
      
    }
  },
  computed: {
    recordingDays: {
      get() {
        return this.selctedCamera.RecordingTime / 86400;
      },
      set(value) {
        this.selctedCamera.RecordingTime = value * 86400;
      },
    },
  },
  mounted() {
    axios
    .get("/api/getCameraSettings")
    .then(response => {
      for (let i = 0; i < response.data.length; i++) {
        this.cameras.push({CameraID: response.data[i].CameraID, Name: response.data[i].Name, CameraUrl: response.data[i].CameraUrl, RecordingTime: response.data[i].RecordingTime, selectionColor: 'none'})
      }
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
.delete-button {
  position: relative;
  margin-top: 0.5rem;
  background-color:rgb(255, 0, 0);
  bottom: 2rem;
  float: right;
  margin-right: 1rem;
  border-radius: 8px;
  outline: none;
  border: none;
  width: 3.15rem;
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
.option-container {
    position: relative;
    height: 35rem;
    width: calc(100% - 15rem);
    float:right;
}
.camera-selector-item {
    cursor: pointer;
}

.add-button {
  position: relative;
  cursor: pointer;

}
.divider {
  width: 100%;
  height: 1px;
  margin-top: 0.5rem;
  margin-bottom: 0.5rem;
  background-color: grey;
  
}
</style>