<template>
    <div class="account-menu">
        <div class="account-text" v-for="option in options" :key="option.title">
            <button class="account-button" @click="buttonClick(option)">
                <h1 class="account-text">{{option.title}}</h1>  
            </button>
            <div class="line" v-if="option.title != 'Logout'"/> 
        </div>
    </div>
</template>

<style>
.account-menu{
    margin-top: 3rem;
    background-color: white;
    border: 0rem;
    margin-left: auto;
    width: 10rem;
    border: 0rem;
    border-bottom-left-radius: 0.2rem;
    border-bottom-right-radius: 0.2rem;
    overflow: hidden;

}
.account-text {
    margin: 0rem;
    text-align: left;
    font-size: 1rem;

}

.account-button {
    width: 10rem;
    border: 0rem;
    background-color: transparent;
    box-shadow: inset 0 0 0 0 transparent;
    outline: none;
    transition: ease-out 0.3s;
    padding: auto;
    margin-top: 0.1rem;
    margin-bottom: 0.1rem;
    cursor: pointer;

}

.account-button:hover {
    box-shadow: inset 10rem 0 0 rgb(0, 102, 255);
    cursor: pointer;
    color: black;
    
}

.line {
    padding: 0rem;
    size: 2;
    margin: 0;
    width: 9rem;
    height: 1px;
    margin-right: auto;
    margin-left: auto;
    background: grey;
}

</style>
<script>
import { ref } from 'vue';
export default {
    data() {
        return {
            options: [{title: 'Settings', hover: false}, {title: 'Logout', hover: false}]
        }
    },
    methods: {
        buttonClick(option) {
            if (option.title == "Logout") {
                setTimeout(() => this.$emit("logoutOpen"), 200);
            } else {
                console.log("Unkown option")
            }
        },
        logOut() {

        }
    }
    
}
</script>