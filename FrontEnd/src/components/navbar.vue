<template>
    <div> 
        <transition name="slide-fade">
            <account-menu class="account-menu" v-show="accountMenuOpen" @logoutOpen="openLogout()" @close="accountMenuOpen=false"/>
        </transition> 
        <sideMenu v-show="sideMenuShown"/>


        <div class="navbar">
                <button class="side-menu-button" @click="sideMenuShown = !sideMenuShown">
                    <div class="side-menu-button"></div>
                    <div class="side-menu-button"></div>
                    <div class="side-menu-button"></div>
                </button>
            <div class="navbar-account">

                <button class="navbar-account" @click="accountMenuOpen = !accountMenuOpen">
                     <h3>
                        Settings 

                     </h3>
                </button>

            </div>

        </div>

    </div>
</template>



<script>
import { ref } from 'vue';
import accountMenu from './accountMenu.vue'
import sideMenu from './sidemenu.vue'
export default {
    props: ['username'],
    components: {
        accountMenu,
        sideMenu
    },
    data() {
        return {
            accountMenuOpen: false,
            sideMenuShown: false   
        }
    },
    methods: {
        openLogout() {
            this.$emit("logoutOpen")
        }
    }
}

</script>

<style scoped>
body{
    background-color: rgb(255, 255, 255);
}

.account-menu{
    top: 3rem;
    z-index: -2;
}
.navbar-account{
    display: flex;
    justify-content: right;
    align-items: center;  
    
}

div.navbar-account{
    position: fixed;
    right: 0;
    top: 0;
    
}

img.navbar-account-icon{
    background-color: rgb(165, 162, 162);
    overflow: hidden;
    border-radius: 50px;

}

button.navbar-account{
    color: white;
    height: 3rem;
    outline-style: none;
    background: transparent;
    border: 0px;
    cursor: pointer;
}



.navbar{

    position: fixed;
    background-color: rgb(44, 44, 44);
    height: 3rem;
    width: 100%;
    top:0;
    left: 0;
    z-index: -1;
}
button.side-menu-button {
    position: fixed;
    background: none;
    height:3rem;
    width: 2.5rem;
    color: white;
    border: none;
    cursor: pointer;
    
    
}

div.side-menu-button {
  height: 2px;
  margin: 4px;
  background-color: white;
  margin: 6px 0;
  
}

</style>
