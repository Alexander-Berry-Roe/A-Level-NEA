<template>
  <div class="title">
    <h3>Account settings</h3>
  </div>
  <div class="option">
    <div class="option-title">
      <t class="option-title">Username</t>
      <t class="error" v-if="showUsernameError"> Username already taken </t>
    </div>
    <input class="option-input" v-model="newUsername"> 
    <button class="save-button" @click="changeUsername()">
        Update
    </button>
  </div>
  <div class="option">
    <div class="option-title">
      <t class="option-title">Change password</t>
    </div>
    <t>New password</t>
    <input class="option-input">
    <t>Confirm password</t>
    <input class="option-input">
    <button class="save-button">
        Update
    </button>
  </div>

</template>

<script>
import axios from 'axios'
export default {
  props: ['username'],
  data() {
    return {
      newUsername: "",
      showUsernameError: false
    }
  },
  methods: {
    refresh() {
      this.newUsername = this.username
    },
    changeUsername() {
      axios.post('/api/setOwnUsername', {
        username: this.newUsername
      }).then(response => {
        if (response.data.success) {
          this.emitter.emit("reloadUserInfo")
          this.showUsernameError=false
        } else {
          this.showUsernameError=true
        }
      } )
    }
  },
  mounted() {
        this.emitter.on("openAccountMenu", this.refresh);
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
  margin-left: 2rem;
  padding: 0.5rem;
  border-radius: 8px;
  background-color: transparent;
  border: 0.1rem;
  border-style: solid;
  border-color: rgb(0, 102, 255);
  
}
input.option-input {
  margin-left: 0rem;
  margin-right: 0rem;
  width: 99%;
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
  margin-top: 0.5rem;
}
.error {
  float:right;
  margin-top: 0.5rem;
  color: red;
}
.save-button {
  position: relative;
  margin-top: 0.5rem;
  background-color:rgb(0, 102, 255);
  margin-right: 1rem;
  border-radius: 8px;
  outline: none;
  border: none;
  width: 4rem;
  height: 2rem;
  color: white;
  cursor: pointer;
}

</style>