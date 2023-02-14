<template id="login-form">
    <div class="login">
        <form>
            

            <div class="floating-box">
                <h3 class="title" >Sign In</h3>
                <div class="error-box">
                    <div class="input-group">
                        <transition name="bounce">
                            <p class="error" v-if="error">Error incorrect password/username</p>
                        </transition>
                    </div>
                </div>
                <div class="input-group">
                    <input type="text" id="username" class="login-input" placeholder="username" v-model="username"/>
                </div>

                <div class="input-group">
                    <input type="password" id="password" class="login-input" placeholder="password" v-model="password"/>
                </div>
                
                <button class="input-group-button" @click="login()">Log in</button>


            </div>
 
        </form>


    </div>
</template>
 
<script>
import axios from "axios";
import { ref } from 'vue';
import modal from './modal.vue';
export default {
  components: { modal },
    name: "LoginForm",
    template: "#login-form",
    data() {
        return {
            username: "",
            password: "",
            loggedin: false,
            error: false
        };
    },
    methods: {
        login() {

            this.loading = true

            axios({
                method: 'post',
                url: '/api/login/auth',
                data: {
                    id: this.username,
                    passwd: this.password
                }
            })
            .then(response => {
		        this.loggedin = response.data.auth
                
                if (this.loggedin) {
                    this.$emit('loggedIn')
                 } else {
                    this.error = true
                    this.loading = false
                    this.password = "";
                }
             
	        })




        }
    }
}
</script>
<style scoped>
div.login {
  font-family: 'Courier New', Courier, monospace;
  font-size: 24px;


  display: flex;
  justify-content: center;
  
}

div.floating-box {
  width: 300px;
  /* height: 270px; */
  padding: 0.5rem;
  margin-top: 4rem;
  z-index: 1;
  
  background-color: rgb(255, 255, 255);
  text-align: center;
  border-radius: 8px;
}

div.input-group{
    width: 300px
  
}

h3.title {
  padding: 15px;
}

.login-input {
  font-family: 'Courier New', Courier, monospace;
  border-style:solid;
  border-width: 1px;
  border-color:rgb(178, 178, 182);
  margin: 0.6rem;
  height: 2rem;
  border-radius: 5px;
  margin-bottom: 1px;
  width: 190px;
}

.input-group-button{
  font-family: 'Courier New', Courier, monospace;
  border-color:rgb(178, 178, 182);
  background-color: rgb(0, 102, 255);
  border-style: solid;
  border-width: 1px;
  border-radius: 5px;
  margin-top: 10px;
  margin-right: auto;
  margin-left: auto;
  width: 200px;
  height: 2rem;
  color: white
}

.error{
    color: red;
    font-size: 1rem;
}

.error-box{
    height: 2rem;
}

.bounce-enter-active {
  animation: bounce-in .5s;
}
.bounce-leave-active {
  animation: bounce-in .5s reverse;
}
@keyframes bounce-in {
  0% {
    transform: scale(0);
  }
  50% {
    transform: scale(1.5);
  }
  100% {
    transform: scale(1);
  }
}
.login {
  position: absolute;
  margin: 0px;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  z-index: 99;

  
}
.loading-icon {
  height: 2rem;
}

</style>