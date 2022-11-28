<template>
    <div> 
        <transition name="slide-fade">
            <account-menu class="account-menu" v-show="accountMenuOpen" @logoutOpen="openLogout()"/>
        </transition> 
        <div class="navbar">
            <div class="navbar-account">
                <button class="navbar-account" @click="accountMenuOpen = !accountMenuOpen">
                     <h3>
                        {{username}}
                        <img src="./icons/arrow-down-sign-to-navigate.png"/>
                     </h3>
                </button>

            </div>

        </div>

    </div>
</template>


<style scoped>
body{
    background-color: rgb(255, 255, 255);
}

.navbar-account{
    display: flex;
    justify-content: right;
    align-items: center;  
    
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
}

</style>

<script>
import { ref } from 'vue';
import accountMenu from './accountMenu.vue'
export default {
    props: ['username'],
    components: {
        accountMenu
    },
    data() {
        return {
            accountMenuOpen: false   
        }
    },
    methods: {
        openLogout() {
            this.$emit("logoutOpen")
        }
    }
}

</script>